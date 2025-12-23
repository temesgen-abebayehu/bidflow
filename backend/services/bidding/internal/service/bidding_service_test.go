package service

import (
	"context"
	"errors"
	"testing"

	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/domain"
)

// Mocks
type MockBidRepo struct {
	CreateFunc          func(ctx context.Context, bid *domain.Bid) error
	GetByIDFunc         func(ctx context.Context, id string) (*domain.Bid, error)
	ListByAuctionIDFunc func(ctx context.Context, auctionID string) ([]domain.Bid, error)
	GetHighestBidFunc   func(ctx context.Context, auctionID string) (*domain.Bid, error)
}

func (m *MockBidRepo) Create(ctx context.Context, bid *domain.Bid) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, bid)
	}
	return nil
}
func (m *MockBidRepo) GetByID(ctx context.Context, id string) (*domain.Bid, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockBidRepo) ListByAuctionID(ctx context.Context, auctionID string) ([]domain.Bid, error) {
	if m.ListByAuctionIDFunc != nil {
		return m.ListByAuctionIDFunc(ctx, auctionID)
	}
	return nil, nil
}
func (m *MockBidRepo) GetHighestBid(ctx context.Context, auctionID string) (*domain.Bid, error) {
	if m.GetHighestBidFunc != nil {
		return m.GetHighestBidFunc(ctx, auctionID)
	}
	return nil, nil
}

type MockEventProducer struct {
	PublishBidPlacedFunc func(ctx context.Context, bid *domain.Bid) error
}

func (m *MockEventProducer) PublishBidPlaced(ctx context.Context, bid *domain.Bid) error {
	if m.PublishBidPlacedFunc != nil {
		return m.PublishBidPlacedFunc(ctx, bid)
	}
	return nil
}

type MockAuctionClient struct {
	ValidateBidFunc        func(ctx context.Context, auctionID string, amount float64, bidderID string) (bool, string, error)
	UpdateAuctionPriceFunc func(ctx context.Context, auctionID string, amount float64) error
}

func (m *MockAuctionClient) ValidateBid(ctx context.Context, auctionID string, amount float64, bidderID string) (bool, string, error) {
	if m.ValidateBidFunc != nil {
		return m.ValidateBidFunc(ctx, auctionID, amount, bidderID)
	}
	return true, "", nil
}

func (m *MockAuctionClient) UpdateAuctionPrice(ctx context.Context, auctionID string, amount float64) error {
	if m.UpdateAuctionPriceFunc != nil {
		return m.UpdateAuctionPriceFunc(ctx, auctionID, amount)
	}
	return nil
}

// Tests
func TestPlaceBid(t *testing.T) {
	tests := []struct {
		name          string
		auctionID     string
		bidderID      string
		amount        float64
		mockSetup     func(*MockBidRepo, *MockEventProducer, *MockAuctionClient)
		expectedError bool
	}{
		{
			name:      "Success",
			auctionID: "auction-1",
			bidderID:  "user-1",
			amount:    100.0,
			mockSetup: func(r *MockBidRepo, e *MockEventProducer, c *MockAuctionClient) {
				c.ValidateBidFunc = func(ctx context.Context, auctionID string, amount float64, bidderID string) (bool, string, error) {
					return true, "valid", nil
				}
				r.CreateFunc = func(ctx context.Context, bid *domain.Bid) error {
					return nil
				}
				c.UpdateAuctionPriceFunc = func(ctx context.Context, auctionID string, amount float64) error {
					return nil
				}
				e.PublishBidPlacedFunc = func(ctx context.Context, bid *domain.Bid) error {
					return nil
				}
			},
			expectedError: false,
		},
		{
			name:      "Invalid Bid",
			auctionID: "auction-1",
			bidderID:  "user-1",
			amount:    50.0,
			mockSetup: func(r *MockBidRepo, e *MockEventProducer, c *MockAuctionClient) {
				c.ValidateBidFunc = func(ctx context.Context, auctionID string, amount float64, bidderID string) (bool, string, error) {
					return false, "too low", nil
				}
			},
			expectedError: true,
		},
		{
			name:      "Repo Error",
			auctionID: "auction-1",
			bidderID:  "user-1",
			amount:    100.0,
			mockSetup: func(r *MockBidRepo, e *MockEventProducer, c *MockAuctionClient) {
				c.ValidateBidFunc = func(ctx context.Context, auctionID string, amount float64, bidderID string) (bool, string, error) {
					return true, "valid", nil
				}
				r.CreateFunc = func(ctx context.Context, bid *domain.Bid) error {
					return errors.New("db error")
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockBidRepo{}
			producer := &MockEventProducer{}
			client := &MockAuctionClient{}

			if tt.mockSetup != nil {
				tt.mockSetup(repo, producer, client)
			}

			svc := NewBiddingService(repo, producer, client)
			_, err := svc.PlaceBid(context.Background(), tt.auctionID, tt.bidderID, tt.amount)

			if (err != nil) != tt.expectedError {
				t.Errorf("PlaceBid() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestGetBidsByAuction(t *testing.T) {
	repo := &MockBidRepo{
		ListByAuctionIDFunc: func(ctx context.Context, auctionID string) ([]domain.Bid, error) {
			return []domain.Bid{
				{ID: "1", Amount: 100},
				{ID: "2", Amount: 90},
			}, nil
		},
	}
	svc := NewBiddingService(repo, &MockEventProducer{}, &MockAuctionClient{})

	bids, err := svc.GetBidsByAuction(context.Background(), "auction-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(bids) != 2 {
		t.Errorf("expected 2 bids, got %d", len(bids))
	}
}
