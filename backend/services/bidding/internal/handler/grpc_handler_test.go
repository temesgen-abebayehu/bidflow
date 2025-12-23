package handler

import (
	"context"
	"testing"
	"time"

	pb "github.com/temesgen-abebayehu/bidflow/backend/proto/pb"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/domain"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/service"
)

func TestPlaceBidGrpc(t *testing.T) {
	repo := &MockBidRepo{
		CreateFunc: func(ctx context.Context, bid *domain.Bid) error {
			bid.ID = "bid-123"
			bid.Timestamp = time.Now()
			return nil
		},
	}
	// MockAuctionClient and MockEventProducer are defined in http_handler_test.go
	// and are available here since they are in the same package (handler)
	svc := service.NewBiddingService(repo, &MockEventProducer{}, &MockAuctionClient{})
	h := NewGrpcHandler(svc)

	req := &pb.PlaceBidRequest{
		AuctionId: "auction-1",
		BidderId:  "user-1",
		Amount:    100.0,
	}

	resp, err := h.PlaceBid(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Bid.AuctionId != req.AuctionId {
		t.Errorf("expected auction id %s, got %s", req.AuctionId, resp.Bid.AuctionId)
	}
	if resp.Bid.Amount != req.Amount {
		t.Errorf("expected amount %f, got %f", req.Amount, resp.Bid.Amount)
	}
}

func TestGetBidsByAuctionGrpc(t *testing.T) {
	repo := &MockBidRepo{
		ListByAuctionIDFunc: func(ctx context.Context, auctionID string) ([]domain.Bid, error) {
			return []domain.Bid{
				{ID: "1", AuctionID: auctionID, Amount: 100, Timestamp: time.Now()},
				{ID: "2", AuctionID: auctionID, Amount: 90, Timestamp: time.Now()},
			}, nil
		},
	}
	svc := service.NewBiddingService(repo, &MockEventProducer{}, &MockAuctionClient{})
	h := NewGrpcHandler(svc)

	req := &pb.GetBidsByAuctionRequest{
		AuctionId: "auction-1",
	}

	resp, err := h.GetBidsByAuction(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Bids) != 2 {
		t.Errorf("expected 2 bids, got %d", len(resp.Bids))
	}
}
