package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/domain"
	"go.uber.org/zap"
)

// --- Mocks ---

type MockNotificationRepo struct {
	mock.Mock
}

func (m *MockNotificationRepo) Create(ctx context.Context, notification *domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockNotificationRepo) ListByUserID(ctx context.Context, userID string, limit int) ([]domain.Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}

func (m *MockNotificationRepo) MarkAsRead(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockHub struct {
	mock.Mock
}

func (m *MockHub) BroadcastToUser(userID string, message interface{}) {
	m.Called(userID, message)
}

func (m *MockHub) Run() {
	m.Called()
}

func (m *MockHub) Register(client domain.Client) {
	m.Called(client)
}

func (m *MockHub) Unregister(client domain.Client) {
	m.Called(client)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
func (m *MockLogger) Warn(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
func (m *MockLogger) Fatal(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
func (m *MockLogger) With(fields ...zap.Field) logger.Logger {
	args := m.Called(fields)
	return args.Get(0).(logger.Logger)
}
func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

// --- Test Suite ---

type NotificationServiceTestSuite struct {
	suite.Suite
	repo    *MockNotificationRepo
	hub     *MockHub
	logger  *MockLogger
	service domain.NotificationService
}

func (s *NotificationServiceTestSuite) SetupTest() {
	s.repo = new(MockNotificationRepo)
	s.hub = new(MockHub)
	s.logger = new(MockLogger)
	s.service = NewNotificationService(s.repo, s.hub, s.logger)
}

func (s *NotificationServiceTestSuite) TestSendNotification_Success() {
	notification := &domain.Notification{
		UserID:  "user-1",
		Type:    domain.NotificationTypeAuctionCreated,
		Title:   "New Auction",
		Message: "An auction has started",
	}

	// Expect Create to be called
	s.repo.On("Create", mock.Anything, mock.MatchedBy(func(n *domain.Notification) bool {
		return n.UserID == notification.UserID && n.Title == notification.Title
	})).Return(nil)

	// Expect BroadcastToUser to be called
	s.hub.On("BroadcastToUser", notification.UserID, mock.Anything).Return()

	err := s.service.SendNotification(context.Background(), notification)

	s.NoError(err)
	s.NotEmpty(notification.ID)        // ID should be generated
	s.False(notification.CreatedAt.IsZero()) // CreatedAt should be set
	s.repo.AssertExpectations(s.T())
	s.hub.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestSendNotification_RepoError() {
	notification := &domain.Notification{
		UserID: "user-1",
	}

	expectedErr := errors.New("db error")

	// Expect Create to fail
	s.repo.On("Create", mock.Anything, mock.Anything).Return(expectedErr)

	// Expect Logger Error
	s.logger.On("Error", "Failed to save notification", mock.Anything).Return()

	err := s.service.SendNotification(context.Background(), notification)

	s.Error(err)
	s.Equal(expectedErr, err)
	s.repo.AssertExpectations(s.T())
	s.hub.AssertNotCalled(s.T(), "BroadcastToUser") // Should not broadcast if DB fails
}

func (s *NotificationServiceTestSuite) TestGetUserNotifications() {
	userID := "user-1"
	notifications := []domain.Notification{
		{ID: "1", UserID: userID, Title: "Notif 1"},
		{ID: "2", UserID: userID, Title: "Notif 2"},
	}

	s.repo.On("ListByUserID", mock.Anything, userID, 50).Return(notifications, nil)

	result, err := s.service.GetUserNotifications(context.Background(), userID)

	s.NoError(err)
	s.Equal(notifications, result)
	s.repo.AssertExpectations(s.T())
}

func TestNotificationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceTestSuite))
}
