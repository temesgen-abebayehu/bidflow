package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/domain"
)

type BiddingService struct {
	repo          domain.BidRepository
	eventProducer domain.EventProducer
	auctionClient domain.AuctionClient
}

func NewBiddingService(repo domain.BidRepository, eventProducer domain.EventProducer, auctionClient domain.AuctionClient) *BiddingService {
	return &BiddingService{
		repo:          repo,
		eventProducer: eventProducer,
		auctionClient: auctionClient,
	}
}

func (s *BiddingService) PlaceBid(ctx context.Context, auctionID, bidderID string, amount float64) (*domain.Bid, error) {
	// 1. Validate with Auction Service
	isValid, msg, err := s.auctionClient.ValidateBid(ctx, auctionID, amount, bidderID)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, errors.New(msg)
	}

	// 2. Create Bid
	bid := &domain.Bid{
		ID:        uuid.New().String(),
		AuctionID: auctionID,
		BidderID:  bidderID,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	// 3. Save to DB
	if err := s.repo.Create(ctx, bid); err != nil {
		return nil, err
	}

	// 4. Update Auction Price (Synchronous for consistency)
	if err := s.auctionClient.UpdateAuctionPrice(ctx, auctionID, amount); err != nil {
		// Log error but don't fail the bid? Or fail?
		// If we fail here, we have an inconsistency (Bid saved, Price not updated).
		// Ideally we should rollback the bid.
		// For now, we return error.
		return nil, err
	}

	// 5. Publish Event
	if err := s.eventProducer.PublishBidPlaced(ctx, bid); err != nil {
		// In a real system, we might want to use the outbox pattern here
		return nil, err
	}

	return bid, nil
}

func (s *BiddingService) GetBidsByAuction(ctx context.Context, auctionID string) ([]domain.Bid, error) {
	return s.repo.ListByAuctionID(ctx, auctionID)
}
