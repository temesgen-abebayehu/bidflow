package event

import (
	"context"
	"time"

	"github.com/temesgen-abebayehu/bidflow/backend/common/kafka"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/domain"
)

const (
	TopicBidPlaced = "bid.placed"
)

type BidPlacedEvent struct {
	BidID     string    `json:"bid_id"`
	AuctionID string    `json:"auction_id"`
	BidderID  string    `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type KafkaEventProducer struct {
	producer *kafka.Producer
}

func NewKafkaEventProducer(producer *kafka.Producer) domain.EventProducer {
	return &KafkaEventProducer{producer: producer}
}

func (p *KafkaEventProducer) PublishBidPlaced(ctx context.Context, bid *domain.Bid) error {
	event := BidPlacedEvent{
		BidID:     bid.ID,
		AuctionID: bid.AuctionID,
		BidderID:  bid.BidderID,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp,
	}
	// Keying by AuctionID ensures ordering for bids on the same auction
	return p.producer.Publish(ctx, TopicBidPlaced, bid.AuctionID, event)
}
