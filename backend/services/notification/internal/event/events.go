package event

import "time"

const (
	TopicAuctionCreated = "auction.created"
	TopicBidPlaced      = "bid.placed"
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

type BidPlacedEvent struct {
	BidID     string    `json:"bid_id"`
	AuctionID string    `json:"auction_id"`
	BidderID  string    `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}
