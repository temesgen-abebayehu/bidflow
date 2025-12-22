package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/temesgen-abebayehu/bidflow/backend/proto/pb"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockAuctionService is a mock implementation of domain.AuctionService
type MockAuctionService struct {
	CreateAuctionFunc      func(ctx context.Context, sellerID, title, description string, startPrice float64, startTime, endTime time.Time, category, imageURL string) (*domain.Auction, error)
	GetAuctionFunc         func(ctx context.Context, id string) (*domain.Auction, error)
	ListAuctionsFunc       func(ctx context.Context, page, limit int, status string, category string) ([]domain.Auction, int64, error)
	UpdateAuctionFunc      func(ctx context.Context, id string, title, description, imageURL string) (*domain.Auction, error)
	CloseAuctionFunc       func(ctx context.Context, id string) error
	ValidateBidFunc        func(ctx context.Context, auctionID string, amount float64) (bool, string, error)
	UpdateCurrentPriceFunc func(ctx context.Context, auctionID string, amount float64) error
}

func (m *MockAuctionService) CreateAuction(ctx context.Context, sellerID, title, description string, startPrice float64, startTime, endTime time.Time, category, imageURL string) (*domain.Auction, error) {
	if m.CreateAuctionFunc != nil {
		return m.CreateAuctionFunc(ctx, sellerID, title, description, startPrice, startTime, endTime, category, imageURL)
	}
	return nil, nil
}

func (m *MockAuctionService) GetAuction(ctx context.Context, id string) (*domain.Auction, error) {
	if m.GetAuctionFunc != nil {
		return m.GetAuctionFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockAuctionService) ListAuctions(ctx context.Context, page, limit int, status string, category string) ([]domain.Auction, int64, error) {
	if m.ListAuctionsFunc != nil {
		return m.ListAuctionsFunc(ctx, page, limit, status, category)
	}
	return nil, 0, nil
}

func (m *MockAuctionService) UpdateAuction(ctx context.Context, id string, title, description, imageURL string) (*domain.Auction, error) {
	if m.UpdateAuctionFunc != nil {
		return m.UpdateAuctionFunc(ctx, id, title, description, imageURL)
	}
	return nil, nil
}

func (m *MockAuctionService) CloseAuction(ctx context.Context, id string) error {
	if m.CloseAuctionFunc != nil {
		return m.CloseAuctionFunc(ctx, id)
	}
	return nil
}

func (m *MockAuctionService) ValidateBid(ctx context.Context, auctionID string, amount float64) (bool, string, error) {
	if m.ValidateBidFunc != nil {
		return m.ValidateBidFunc(ctx, auctionID, amount)
	}
	return false, "", nil
}

func (m *MockAuctionService) UpdateCurrentPrice(ctx context.Context, auctionID string, amount float64) error {
	if m.UpdateCurrentPriceFunc != nil {
		return m.UpdateCurrentPriceFunc(ctx, auctionID, amount)
	}
	return nil
}

func TestCreateAuction_Grpc(t *testing.T) {
	mockSvc := &MockAuctionService{
		CreateAuctionFunc: func(ctx context.Context, sellerID, title, description string, startPrice float64, startTime, endTime time.Time, category, imageURL string) (*domain.Auction, error) {
			return &domain.Auction{ID: "123"}, nil
		},
	}
	h := NewGrpcHandler(mockSvc)

	req := &pb.CreateAuctionRequest{
		SellerId:    "seller-1",
		Title:       "Test",
		Description: "Desc",
		StartPrice:  10.0,
		StartTime:   time.Now().Unix(),
		EndTime:     time.Now().Add(time.Hour).Unix(),
		Category:    "Cat",
		ImageUrl:    "url",
	}

	resp, err := h.CreateAuction(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Auction.Id != "123" {
		t.Errorf("expected id 123, got %s", resp.Auction.Id)
	}
}

func TestGetAuction_Grpc(t *testing.T) {
	mockSvc := &MockAuctionService{
		GetAuctionFunc: func(ctx context.Context, id string) (*domain.Auction, error) {
			if id == "found" {
				return &domain.Auction{ID: "found"}, nil
			}
			return nil, errors.New("not found")
		},
	}
	h := NewGrpcHandler(mockSvc)

	t.Run("Found", func(t *testing.T) {
		resp, err := h.GetAuction(context.Background(), &pb.GetAuctionRequest{Id: "found"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.Auction.Id != "found" {
			t.Errorf("expected id found, got %s", resp.Auction.Id)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		_, err := h.GetAuction(context.Background(), &pb.GetAuctionRequest{Id: "missing"})
		if status.Code(err) != codes.NotFound {
			t.Errorf("expected NotFound error, got %v", err)
		}
	})
}

func TestListAuctions_Grpc(t *testing.T) {
	mockSvc := &MockAuctionService{
		ListAuctionsFunc: func(ctx context.Context, page, limit int, status string, category string) ([]domain.Auction, int64, error) {
			return []domain.Auction{{ID: "1"}}, 1, nil
		},
	}
	h := NewGrpcHandler(mockSvc)

	resp, err := h.ListAuctions(context.Background(), &pb.ListAuctionsRequest{Page: 1, Limit: 10})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(resp.Auctions) != 1 {
		t.Errorf("expected 1 auction, got %d", len(resp.Auctions))
	}
}

func TestValidateBid_Grpc(t *testing.T) {
	mockSvc := &MockAuctionService{
		ValidateBidFunc: func(ctx context.Context, auctionID string, amount float64) (bool, string, error) {
			if amount > 100 {
				return true, "valid", nil
			}
			return false, "low bid", nil
		},
		GetAuctionFunc: func(ctx context.Context, id string) (*domain.Auction, error) {
			return &domain.Auction{CurrentPrice: 100}, nil
		},
	}
	h := NewGrpcHandler(mockSvc)

	t.Run("Valid", func(t *testing.T) {
		resp, err := h.ValidateBid(context.Background(), &pb.BidRequest{AuctionId: "1", Amount: 150})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !resp.IsValid {
			t.Error("expected valid")
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		resp, err := h.ValidateBid(context.Background(), &pb.BidRequest{AuctionId: "1", Amount: 50})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.IsValid {
			t.Error("expected invalid")
		}
	})
}
