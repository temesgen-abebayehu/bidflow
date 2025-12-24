package domain

import (
	"context"
	"time"
)

type NotificationType string

const (
	NotificationTypeAuctionCreated NotificationType = "AUCTION_CREATED"
	NotificationTypeBidPlaced      NotificationType = "BID_PLACED"
	NotificationTypeAuctionClosed  NotificationType = "AUCTION_CLOSED"
	NotificationTypeOutbid         NotificationType = "OUTBID"
)

type Notification struct {
	ID         string           `json:"id"`
	UserID     string           `json:"user_id"`
	Type       NotificationType `json:"type"`
	Title      string           `json:"title"`
	Message    string           `json:"message"`
	ResourceID string           `json:"resource_id"` // e.g., AuctionID
	IsRead     bool             `json:"is_read"`
	CreatedAt  time.Time        `json:"created_at"`
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	ListByUserID(ctx context.Context, userID string, limit int) ([]Notification, error)
	MarkAsRead(ctx context.Context, id string) error
}

type NotificationService interface {
	SendNotification(ctx context.Context, notification *Notification) error
	GetUserNotifications(ctx context.Context, userID string) ([]Notification, error)
}

type Hub interface {
	BroadcastToUser(userID string, message interface{})
	Run()
	Register(client Client)
	Unregister(client Client)
}

type Client interface {
	ReadPump()
	WritePump()
}
