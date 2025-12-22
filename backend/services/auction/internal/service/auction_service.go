package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/domain"
	"go.uber.org/zap"
)

type AuctionService struct {
	repo     domain.AuctionRepository
	producer domain.EventProducer
	log      logger.Logger
}

func NewAuctionService(repo domain.AuctionRepository, producer domain.EventProducer, log logger.Logger) *AuctionService {
	return &AuctionService{
		repo:     repo,
		producer: producer,
		log:      log,
	}
}

func (s *AuctionService) CreateAuction(ctx context.Context, sellerID, title, description string, startPrice float64, startTime, endTime time.Time, category, imageURL string) (*domain.Auction, error) {
	if startTime.After(endTime) {
		return nil, errors.New("start time must be before end time")
	}

	if startPrice < 0 {
		return nil, errors.New("start price cannot be negative")
	}

	auction := &domain.Auction{
		ID:           uuid.New().String(),
		SellerID:     sellerID,
		Title:        title,
		Description:  description,
		StartPrice:   startPrice,
		CurrentPrice: startPrice,
		Status:       domain.AuctionStatusPending, // Or ACTIVE if start time is now
		StartTime:    startTime,
		EndTime:      endTime,
		Category:     category,
		ImageURL:     imageURL,
	}

	if startTime.Before(time.Now()) {
		auction.Status = domain.AuctionStatusActive
	}

	err := s.repo.Create(ctx, auction)
	if err != nil {
		return nil, err
	}

	if err := s.producer.PublishAuctionCreated(ctx, auction); err != nil {
		s.log.Error("failed to publish auction created event", zap.Error(err), zap.String("auction_id", auction.ID))
	}

	return auction, nil
}

func (s *AuctionService) GetAuction(ctx context.Context, id string) (*domain.Auction, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AuctionService) ListAuctions(ctx context.Context, page, limit int, status string, category string) ([]domain.Auction, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.List(ctx, page, limit, domain.AuctionStatus(status), category)
}

func (s *AuctionService) UpdateAuction(ctx context.Context, id string, title, description, imageURL string) (*domain.Auction, error) {
	auction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if auction.Status == domain.AuctionStatusClosed || auction.Status == domain.AuctionStatusCancelled {
		return nil, errors.New("cannot update closed or cancelled auction")
	}

	if title != "" {
		auction.Title = title
	}
	if description != "" {
		auction.Description = description
	}
	if imageURL != "" {
		auction.ImageURL = imageURL
	}

	err = s.repo.Update(ctx, auction)
	if err != nil {
		return nil, err
	}

	if err := s.producer.PublishAuctionUpdated(ctx, auction); err != nil {
		s.log.Error("failed to publish auction updated event", zap.Error(err), zap.String("auction_id", auction.ID))
	}

	return auction, nil
}

func (s *AuctionService) CloseAuction(ctx context.Context, id string) error {
	auction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if auction.Status == domain.AuctionStatusClosed {
		return nil
	}

	auction.Status = domain.AuctionStatusClosed
	err = s.repo.Update(ctx, auction)
	if err != nil {
		return err
	}

	if err := s.producer.PublishAuctionClosed(ctx, auction, ""); err != nil {
		s.log.Error("failed to publish auction closed event", zap.Error(err), zap.String("auction_id", auction.ID))
	}

	return nil
}

func (s *AuctionService) ValidateBid(ctx context.Context, auctionID string, amount float64) (bool, string, error) {
	auction, err := s.repo.GetByID(ctx, auctionID)
	if err != nil {
		return false, "Auction not found", err
	}

	if auction.Status != domain.AuctionStatusActive {
		return false, "Auction is not active", nil
	}

	if time.Now().After(auction.EndTime) {
		// Should probably close it here or have a background job
		return false, "Auction has ended", nil
	}

	if amount <= auction.CurrentPrice {
		return false, "Bid amount must be higher than current price", nil
	}

	return true, "Valid bid", nil
}

func (s *AuctionService) UpdateCurrentPrice(ctx context.Context, auctionID string, amount float64) error {
	auction, err := s.repo.GetByID(ctx, auctionID)
	if err != nil {
		return err
	}

	auction.CurrentPrice = amount
	return s.repo.Update(ctx, auction)
}
