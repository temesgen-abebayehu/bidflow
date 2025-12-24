package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/temesgen-abebayehu/bidflow/backend/common/kafka"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/domain"
	"go.uber.org/zap"
)

type NotificationConsumer struct {
	consumer *kafka.Consumer
	service  domain.NotificationService
	log      logger.Logger
}

func NewNotificationConsumer(consumer *kafka.Consumer, service domain.NotificationService, log logger.Logger) *NotificationConsumer {
	return &NotificationConsumer{
		consumer: consumer,
		service:  service,
		log:      log,
	}
}

func (c *NotificationConsumer) Start(ctx context.Context) {
	c.log.Info("Starting notification consumer")

	// Start consuming using the common consumer's Start method
	c.consumer.Start(ctx, c.handleMessage)
}

func (c *NotificationConsumer) handleMessage(ctx context.Context, topic string, key, value []byte) error {
	c.log.Info("Received message", zap.String("topic", topic), zap.String("key", string(key)))

	switch topic {
	case TopicAuctionCreated:
		return c.handleAuctionCreated(ctx, value)
	case TopicBidPlaced:
		return c.handleBidPlaced(ctx, value)
	default:
		c.log.Warn("Unknown topic", zap.String("topic", topic))
		return nil
	}
}

func (c *NotificationConsumer) handleAuctionCreated(ctx context.Context, value []byte) error {
	var event AuctionCreatedEvent
	if err := json.Unmarshal(value, &event); err != nil {
		c.log.Error("Failed to unmarshal AuctionCreatedEvent", zap.Error(err))
		return nil // Don't retry on unmarshal error
	}

	notification := &domain.Notification{
		UserID:     event.SellerID,
		Type:       domain.NotificationTypeAuctionCreated,
		Title:      "Auction Created",
		Message:    fmt.Sprintf("Your auction '%s' has been successfully created.", event.Title),
		ResourceID: event.AuctionID,
	}

	if err := c.service.SendNotification(ctx, notification); err != nil {
		c.log.Error("Failed to send notification for AuctionCreated", zap.Error(err))
		return err
	}
	return nil
}

func (c *NotificationConsumer) handleBidPlaced(ctx context.Context, value []byte) error {
	var event BidPlacedEvent
	if err := json.Unmarshal(value, &event); err != nil {
		c.log.Error("Failed to unmarshal BidPlacedEvent", zap.Error(err))
		return nil // Don't retry on unmarshal error
	}

	// Notify the bidder
	notification := &domain.Notification{
		UserID:     event.BidderID,
		Type:       domain.NotificationTypeBidPlaced,
		Title:      "Bid Placed",
		Message:    fmt.Sprintf("You placed a bid of %.2f on auction %s.", event.Amount, event.AuctionID),
		ResourceID: event.AuctionID,
	}

	if err := c.service.SendNotification(ctx, notification); err != nil {
		c.log.Error("Failed to send notification for BidPlaced", zap.Error(err))
		return err
	}
	return nil
}
