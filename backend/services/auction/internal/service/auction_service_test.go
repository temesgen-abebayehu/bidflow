package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/domain"
	"go.uber.org/zap"
)

type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...zap.Field)  {}
func (m *MockLogger) Info(msg string, fields ...zap.Field)   {}
func (m *MockLogger) Warn(msg string, fields ...zap.Field)   {}
func (m *MockLogger) Error(msg string, fields ...zap.Field)  {}
func (m *MockLogger) Fatal(msg string, fields ...zap.Field)  {}
func (m *MockLogger) With(fields ...zap.Field) logger.Logger { return m }
func (m *MockLogger) Sync() error                            { return nil }

type MockAuctionRepo struct {
	CreateFunc  func(ctx context.Context, auction *domain.Auction) error
	GetByIDFunc func(ctx context.Context, id string) (*domain.Auction, error)
	UpdateFunc  func(ctx context.Context, auction *domain.Auction) error
	DeleteFunc  func(ctx context.Context, id string) error
	ListFunc    func(ctx context.Context, page, limit int, status domain.AuctionStatus, category string) ([]domain.Auction, int64, error)
}

func (m *MockAuctionRepo) Create(ctx context.Context, auction *domain.Auction) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, auction)
	}
	return nil
}

func (m *MockAuctionRepo) GetByID(ctx context.Context, id string) (*domain.Auction, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockAuctionRepo) Update(ctx context.Context, auction *domain.Auction) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, auction)
	}
	return nil
}

func (m *MockAuctionRepo) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockAuctionRepo) List(ctx context.Context, page, limit int, status domain.AuctionStatus, category string) ([]domain.Auction, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, page, limit, status, category)
	}
	return nil, 0, nil
}

type MockEventProducer struct {
	PublishAuctionCreatedFunc func(ctx context.Context, auction *domain.Auction) error
	PublishAuctionUpdatedFunc func(ctx context.Context, auction *domain.Auction) error
	PublishAuctionClosedFunc  func(ctx context.Context, auction *domain.Auction, winnerID string) error
}

func (m *MockEventProducer) PublishAuctionCreated(ctx context.Context, auction *domain.Auction) error {
	if m.PublishAuctionCreatedFunc != nil {
		return m.PublishAuctionCreatedFunc(ctx, auction)
	}
	return nil
}

func (m *MockEventProducer) PublishAuctionUpdated(ctx context.Context, auction *domain.Auction) error {
	if m.PublishAuctionUpdatedFunc != nil {
		return m.PublishAuctionUpdatedFunc(ctx, auction)
	}
	return nil
}

func (m *MockEventProducer) PublishAuctionClosed(ctx context.Context, auction *domain.Auction, winnerID string) error {
	if m.PublishAuctionClosedFunc != nil {
		return m.PublishAuctionClosedFunc(ctx, auction, winnerID)
	}
	return nil
}

func TestCreateAuction(t *testing.T) {
	tests := []struct {
		name        string
		sellerID    string
		title       string
		description string
		startPrice  float64
		startTime   time.Time
		endTime     time.Time
		category    string
		imageURL    string
		mockRepo    func() *MockAuctionRepo
		mockProd    func() *MockEventProducer
		wantErr     bool
	}{
		{
			name:        "Success",
			sellerID:    "seller-1",
			title:       "Test Auction",
			description: "Description",
			startPrice:  10.0,
			startTime:   time.Now().Add(1 * time.Hour),
			endTime:     time.Now().Add(2 * time.Hour),
			category:    "Electronics",
			imageURL:    "http://image.com",
			mockRepo: func() *MockAuctionRepo {
				return &MockAuctionRepo{
					CreateFunc: func(ctx context.Context, auction *domain.Auction) error {
						if auction.SellerID != "seller-1" {
							return errors.New("unexpected seller id")
						}
						return nil
					},
				}
			},
			mockProd: func() *MockEventProducer {
				return &MockEventProducer{
					PublishAuctionCreatedFunc: func(ctx context.Context, auction *domain.Auction) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "Invalid Time Range",
			sellerID:    "seller-1",
			title:       "Test Auction",
			description: "Description",
			startPrice:  10.0,
			startTime:   time.Now().Add(2 * time.Hour),
			endTime:     time.Now().Add(1 * time.Hour),
			category:    "Electronics",
			imageURL:    "http://image.com",
			mockRepo: func() *MockAuctionRepo {
				return &MockAuctionRepo{}
			},
			mockProd: func() *MockEventProducer {
				return &MockEventProducer{}
			},
			wantErr: true,
		},
		{
			name:        "Negative Price",
			sellerID:    "seller-1",
			title:       "Test Auction",
			description: "Description",
			startPrice:  -10.0,
			startTime:   time.Now().Add(1 * time.Hour),
			endTime:     time.Now().Add(2 * time.Hour),
			category:    "Electronics",
			imageURL:    "http://image.com",
			mockRepo: func() *MockAuctionRepo {
				return &MockAuctionRepo{}
			},
			mockProd: func() *MockEventProducer {
				return &MockEventProducer{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockRepo()
			prod := tt.mockProd()
			svc := NewAuctionService(repo, prod, &MockLogger{})

			_, err := svc.CreateAuction(context.Background(), tt.sellerID, tt.title, tt.description, tt.startPrice, tt.startTime, tt.endTime, tt.category, tt.imageURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAuction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAuction(t *testing.T) {
	mockRepo := &MockAuctionRepo{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.Auction, error) {
			if id == "found" {
				return &domain.Auction{ID: "found"}, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := NewAuctionService(mockRepo, &MockEventProducer{}, &MockLogger{})

	t.Run("Found", func(t *testing.T) {
		auction, err := svc.GetAuction(context.Background(), "found")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if auction.ID != "found" {
			t.Errorf("expected id found, got %s", auction.ID)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		_, err := svc.GetAuction(context.Background(), "missing")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUpdateAuction(t *testing.T) {
	mockRepo := &MockAuctionRepo{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.Auction, error) {
			if id == "active" {
				return &domain.Auction{ID: "active", Status: domain.AuctionStatusActive}, nil
			}
			if id == "closed" {
				return &domain.Auction{ID: "closed", Status: domain.AuctionStatusClosed}, nil
			}
			return nil, errors.New("not found")
		},
		UpdateFunc: func(ctx context.Context, auction *domain.Auction) error {
			return nil
		},
	}
	mockProd := &MockEventProducer{
		PublishAuctionUpdatedFunc: func(ctx context.Context, auction *domain.Auction) error {
			return nil
		},
	}
	svc := NewAuctionService(mockRepo, mockProd, &MockLogger{})

	t.Run("Success", func(t *testing.T) {
		_, err := svc.UpdateAuction(context.Background(), "active", "New Title", "", "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Closed Auction", func(t *testing.T) {
		_, err := svc.UpdateAuction(context.Background(), "closed", "New Title", "", "")
		if err == nil {
			t.Error("expected error for closed auction, got nil")
		}
	})
}

func TestCloseAuction(t *testing.T) {
	mockRepo := &MockAuctionRepo{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.Auction, error) {
			return &domain.Auction{ID: id, Status: domain.AuctionStatusActive}, nil
		},
		UpdateFunc: func(ctx context.Context, auction *domain.Auction) error {
			if auction.Status != domain.AuctionStatusClosed {
				return errors.New("status not updated")
			}
			return nil
		},
	}
	mockProd := &MockEventProducer{
		PublishAuctionClosedFunc: func(ctx context.Context, auction *domain.Auction, winnerID string) error {
			return nil
		},
	}
	svc := NewAuctionService(mockRepo, mockProd, &MockLogger{})

	err := svc.CloseAuction(context.Background(), "1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateBid(t *testing.T) {
	now := time.Now()
	mockRepo := &MockAuctionRepo{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.Auction, error) {
			if id == "active" {
				return &domain.Auction{
					ID:           "active",
					Status:       domain.AuctionStatusActive,
					CurrentPrice: 100,
					EndTime:      now.Add(1 * time.Hour),
				}, nil
			}
			if id == "ended" {
				return &domain.Auction{
					ID:           "ended",
					Status:       domain.AuctionStatusActive,
					CurrentPrice: 100,
					EndTime:      now.Add(-1 * time.Hour),
				}, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := NewAuctionService(mockRepo, &MockEventProducer{}, &MockLogger{})

	t.Run("Valid Bid", func(t *testing.T) {
		valid, msg, err := svc.ValidateBid(context.Background(), "active", 150)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !valid {
			t.Errorf("expected valid bid, got invalid: %s", msg)
		}
	})

	t.Run("Low Bid", func(t *testing.T) {
		valid, _, err := svc.ValidateBid(context.Background(), "active", 50)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if valid {
			t.Error("expected invalid bid, got valid")
		}
	})

	t.Run("Ended Auction", func(t *testing.T) {
		valid, _, err := svc.ValidateBid(context.Background(), "ended", 150)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if valid {
			t.Error("expected invalid bid for ended auction, got valid")
		}
	})
}

func TestUpdateCurrentPrice(t *testing.T) {
	mockRepo := &MockAuctionRepo{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.Auction, error) {
			return &domain.Auction{ID: id, CurrentPrice: 100}, nil
		},
		UpdateFunc: func(ctx context.Context, auction *domain.Auction) error {
			if auction.CurrentPrice != 200 {
				return errors.New("price not updated")
			}
			return nil
		},
	}
	mockProd := &MockEventProducer{
		PublishAuctionUpdatedFunc: func(ctx context.Context, auction *domain.Auction) error {
			return nil
		},
	}
	svc := NewAuctionService(mockRepo, mockProd, &MockLogger{})

	err := svc.UpdateCurrentPrice(context.Background(), "1", 200)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestListAuctions(t *testing.T) {
	mockRepo := &MockAuctionRepo{
		ListFunc: func(ctx context.Context, page, limit int, status domain.AuctionStatus, category string) ([]domain.Auction, int64, error) {
			if page == 1 && limit == 10 {
				return []domain.Auction{{ID: "1"}}, 1, nil
			}
			return nil, 0, errors.New("invalid params")
		},
	}
	svc := NewAuctionService(mockRepo, &MockEventProducer{}, &MockLogger{})

	t.Run("Success", func(t *testing.T) {
		auctions, count, err := svc.ListAuctions(context.Background(), 1, 10, "", "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("expected count 1, got %d", count)
		}
		if len(auctions) != 1 {
			t.Errorf("expected 1 auction, got %d", len(auctions))
		}
	})

	t.Run("Default Params", func(t *testing.T) {
		// Test that 0, 0 defaults to 1, 10
		auctions, _, err := svc.ListAuctions(context.Background(), 0, 0, "", "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(auctions) != 1 {
			t.Errorf("expected 1 auction, got %d", len(auctions))
		}
	})
}
