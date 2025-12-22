package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAuctionNotFound = errors.New("auction not found")
	ErrInvalidAuction  = errors.New("invalid auction data")
)

type AuctionStatus string

const (
	AuctionStatusActive    AuctionStatus = "ACTIVE"
	AuctionStatusClosed    AuctionStatus = "CLOSED"
	AuctionStatusPending   AuctionStatus = "PENDING"
	AuctionStatusCancelled AuctionStatus = "CANCELLED"
)

type Auction struct {
	ID           string        `json:"id" gorm:"primaryKey"`
	SellerID     string        `json:"seller_id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	StartPrice   float64       `json:"start_price"`
	CurrentPrice float64       `json:"current_price"` // Can be used as end_price if closed
	Status       AuctionStatus `json:"status"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Category     string        `json:"category"`
	ImageURL     string        `json:"image_url"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type AuctionRepository interface {
	Create(ctx context.Context, auction *Auction) error
	GetByID(ctx context.Context, id string) (*Auction, error)
	Update(ctx context.Context, auction *Auction) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int, status AuctionStatus, category string) ([]Auction, int64, error)
}
