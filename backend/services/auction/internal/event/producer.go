package event

import (
	"context"
	"time"

	"github.com/temesgen-abebayehu/bidflow/backend/common/kafka"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/domain"
)

type KafkaEventProducer struct {
	producer *kafka.Producer
}

func NewKafkaEventProducer(producer *kafka.Producer) domain.EventProducer {
	return &KafkaEventProducer{producer: producer}
}

func (p *KafkaEventProducer) PublishAuctionCreated(ctx context.Context, auction *domain.Auction) error {
	event := AuctionCreatedEvent{
		AuctionID:  auction.ID,
		SellerID:   auction.SellerID,
		Title:      auction.Title,
		StartPrice: auction.StartPrice,
		StartTime:  auction.StartTime,
		EndTime:    auction.EndTime,
		Category:   auction.Category,
		Timestamp:  time.Now(),
	}
	return p.producer.Publish(ctx, TopicAuctionCreated, auction.ID, event)
}

func (p *KafkaEventProducer) PublishAuctionUpdated(ctx context.Context, auction *domain.Auction) error {
	event := AuctionUpdatedEvent{
		AuctionID:   auction.ID,
		Title:       auction.Title,
		Description: auction.Description,
		ImageURL:    auction.ImageURL,
		Timestamp:   time.Now(),
	}
	return p.producer.Publish(ctx, TopicAuctionUpdated, auction.ID, event)
}

func (p *KafkaEventProducer) PublishAuctionClosed(ctx context.Context, auction *domain.Auction, winnerID string) error {
	event := AuctionClosedEvent{
		AuctionID:  auction.ID,
		FinalPrice: auction.CurrentPrice,
		WinnerID:   winnerID,
		Timestamp:  time.Now(),
	}
	return p.producer.Publish(ctx, TopicAuctionClosed, auction.ID, event)
}
