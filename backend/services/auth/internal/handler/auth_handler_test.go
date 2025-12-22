package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/handler"
)

// MockAuthService is a mock implementation of domain.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req auth.RegisterRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (*auth.UserDTO, string, bool, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Bool(2), args.Error(3)
	}
	return args.Get(0).(*auth.UserDTO), args.String(1), args.Bool(2), args.Error(3)
}

func (m *MockAuthService) Verify2FA(ctx context.Context, email, code string) (string, error) {
	args := m.Called(ctx, email, code)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Toggle2FA(ctx context.Context, userID string, enable bool) (string, error) {
	args := m.Called(ctx, userID, enable)
	return args.String(0), args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)
		r := gin.Default()
		r.POST("/register", h.Register)

		reqBody := auth.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Username: "testuser",
			FullName: "Test User",
			Role:     "BIDDER",
		}
		mockSvc.On("Register", mock.Anything, reqBody).Return(nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)
		r := gin.Default()
		r.POST("/register", h.Register)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)
		r := gin.Default()
		r.POST("/register", h.Register)

		reqBody := auth.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		mockSvc.On("Register", mock.Anything, reqBody).Return(errors.New("email already exists"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)
		r := gin.Default()
		r.POST("/login", h.Login)

		reqBody := auth.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		userDTO := &auth.UserDTO{ID: "1", Email: "test@example.com"}
		mockSvc.On("Login", mock.Anything, reqBody.Email, reqBody.Password).Return(userDTO, "token123", false, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp auth.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "token123", resp.Token)
		assert.Equal(t, "test@example.com", resp.User.Email)
	})

	t.Run("MFARequired", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)
		r := gin.Default()
		r.POST("/login", h.Login)

		reqBody := auth.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		mockSvc.On("Login", mock.Anything, reqBody.Email, reqBody.Password).Return(nil, "", true, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, true, resp["mfa_required"])
	})
}
