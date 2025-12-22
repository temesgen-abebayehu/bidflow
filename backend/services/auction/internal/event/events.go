package event

import (
	"time"
)

const (
	TopicAuctionCreated = "auction.created"
	TopicAuctionUpdated = "auction.updated"
	TopicAuctionClosed  = "auction.closed"
)

type AuctionCreatedEvent struct {
	AuctionID  string    `json:"auction_id"`
	SellerID   string    `json:"seller_id"`
	Title      string    `json:"title"`
	StartPrice float64   `json:"start_price"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Category   string    `json:"category"`
	Timestamp  time.Time `json:"timestamp"`
}

type AuctionUpdatedEvent struct {
	AuctionID   string    `json:"auction_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Timestamp   time.Time `json:"timestamp"`
}

type AuctionClosedEvent struct {
	AuctionID  string    `json:"auction_id"`
	FinalPrice float64   `json:"final_price"`
	WinnerID   string    `json:"winner_id,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}
