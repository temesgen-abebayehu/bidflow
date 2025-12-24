package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/domain"
	"go.uber.org/zap"
)

type notificationService struct {
	repo domain.NotificationRepository
	hub  domain.Hub
	log  logger.Logger
}

func NewNotificationService(repo domain.NotificationRepository, hub domain.Hub, log logger.Logger) domain.NotificationService {
	return &notificationService{
		repo: repo,
		hub:  hub,
		log:  log,
	}
}

func (s *notificationService) SendNotification(ctx context.Context, notification *domain.Notification) error {
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now()
	}

	// 1. Save to database
	if err := s.repo.Create(ctx, notification); err != nil {
		s.log.Error("Failed to save notification", zap.Error(err))
		return err
	}

	// 2. Push to WebSocket
	s.hub.BroadcastToUser(notification.UserID, notification)

	return nil
}

func (s *notificationService) GetUserNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	return s.repo.ListByUserID(ctx, userID, 50) // Default limit 50
}
