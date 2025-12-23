package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrBidNotFound = errors.New("bid not found")
	ErrInvalidBid  = errors.New("invalid bid")
)

type Bid struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	AuctionID string    `json:"auction_id"`
	BidderID  string    `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type BidRepository interface {
	Create(ctx context.Context, bid *Bid) error
	GetByID(ctx context.Context, id string) (*Bid, error)
	ListByAuctionID(ctx context.Context, auctionID string) ([]Bid, error)
	GetHighestBid(ctx context.Context, auctionID string) (*Bid, error)
}

type EventProducer interface {
	PublishBidPlaced(ctx context.Context, bid *Bid) error
}

type AuctionClient interface {
	ValidateBid(ctx context.Context, auctionID string, amount float64, bidderID string) (bool, string, error)
	UpdateAuctionPrice(ctx context.Context, auctionID string, amount float64) error
}
