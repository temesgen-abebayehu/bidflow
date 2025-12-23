package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/domain"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/service"
)

// Mocks for Service Dependencies (reused from service test logic, but simplified here)
type MockBidRepo struct {
	CreateFunc          func(ctx context.Context, bid *domain.Bid) error
	ListByAuctionIDFunc func(ctx context.Context, auctionID string) ([]domain.Bid, error)
}

func (m *MockBidRepo) Create(ctx context.Context, bid *domain.Bid) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, bid)
	}
	return nil
}
func (m *MockBidRepo) GetByID(ctx context.Context, id string) (*domain.Bid, error) { return nil, nil }
func (m *MockBidRepo) ListByAuctionID(ctx context.Context, auctionID string) ([]domain.Bid, error) {
	if m.ListByAuctionIDFunc != nil {
		return m.ListByAuctionIDFunc(ctx, auctionID)
	}
	return nil, nil
}
func (m *MockBidRepo) GetHighestBid(ctx context.Context, auctionID string) (*domain.Bid, error) {
	return nil, nil
}

type MockEventProducer struct{}

func (m *MockEventProducer) PublishBidPlaced(ctx context.Context, bid *domain.Bid) error { return nil }

type MockAuctionClient struct{}

func (m *MockAuctionClient) ValidateBid(ctx context.Context, auctionID string, amount float64, bidderID string) (bool, string, error) {
	return true, "valid", nil
}
func (m *MockAuctionClient) UpdateAuctionPrice(ctx context.Context, auctionID string, amount float64) error {
	return nil
}

func TestPlaceBidHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &MockBidRepo{}
	svc := service.NewBiddingService(repo, &MockEventProducer{}, &MockAuctionClient{})
	h := NewHttpHandler(svc)

	r := gin.Default()
	// Mock Auth Middleware by setting user_id manually
	r.POST("/bids", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		h.PlaceBid(c)
	})

	reqBody := map[string]interface{}{
		"auction_id": "auction-1",
		"amount":     150.0,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/bids", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestGetBidsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &MockBidRepo{
		ListByAuctionIDFunc: func(ctx context.Context, auctionID string) ([]domain.Bid, error) {
			return []domain.Bid{{ID: "1", Amount: 100}}, nil
		},
	}
	svc := service.NewBiddingService(repo, &MockEventProducer{}, &MockAuctionClient{})
	h := NewHttpHandler(svc)

	r := gin.Default()
	r.GET("/bids/:auction_id", h.GetBids)

	req, _ := http.NewRequest("GET", "/bids/auction-1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
