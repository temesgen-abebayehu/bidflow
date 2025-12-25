package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/domain"
	"go.uber.org/zap"
)

// --- Mocks ---

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendNotification(ctx context.Context, notification *domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockNotificationService) GetUserNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, fields ...zap.Field) {}
func (m *MockLogger) Info(msg string, fields ...zap.Field)  {}
func (m *MockLogger) Warn(msg string, fields ...zap.Field)  {}
func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
func (m *MockLogger) Fatal(msg string, fields ...zap.Field)  {}
func (m *MockLogger) With(fields ...zap.Field) logger.Logger { return m }
func (m *MockLogger) Sync() error                            { return nil }

// --- Tests ---

func TestGetNotifications_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNotificationService)
	mockLogger := new(MockLogger)

	handler := NewNotificationHandler(mockService, nil, nil, mockLogger)

	userID := "user-1"
	notifications := []domain.Notification{
		{ID: "1", Title: "Test"},
	}

	mockService.On("GetUserNotifications", mock.Anything, userID).Return(notifications, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/notifications", nil)
	c.Set("user_id", userID)

	handler.GetNotifications(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetNotifications_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewNotificationHandler(nil, nil, nil, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/notifications", nil)
	// No user_id set

	handler.GetNotifications(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetNotifications_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNotificationService)
	mockLogger := new(MockLogger)

	handler := NewNotificationHandler(mockService, nil, nil, mockLogger)

	userID := "user-1"
	expectedErr := errors.New("db error")

	mockService.On("GetUserNotifications", mock.Anything, userID).Return(nil, expectedErr)
	mockLogger.On("Error", "Failed to get notifications", mock.Anything).Return()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/notifications", nil)
	c.Set("user_id", userID)

	handler.GetNotifications(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
