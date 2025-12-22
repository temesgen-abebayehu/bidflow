package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/handler"
)

// MockUserService is a mock implementation of domain.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetProfile(ctx context.Context, userID string) (*auth.UserDTO, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.UserDTO), args.Error(1)
}

func (m *MockUserService) UpdateProfile(ctx context.Context, userID string, req auth.UpdateProfileRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserService) VerifyUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserService) CreateCompany(ctx context.Context, userID string, req auth.CreateCompanyRequest) (*auth.CompanyDTO, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.CompanyDTO), args.Error(1)
}

func (m *MockUserService) UpdateCompany(ctx context.Context, companyID string, req auth.UpdateCompanyRequest) error {
	args := m.Called(ctx, companyID, req)
	return args.Error(0)
}

func (m *MockUserService) VerifyCompany(ctx context.Context, companyID string) error {
	args := m.Called(ctx, companyID)
	return args.Error(0)
}

func (m *MockUserService) GetCompany(ctx context.Context, companyID string) (*auth.CompanyDTO, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.CompanyDTO), args.Error(1)
}

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)
		r := gin.Default()
		r.GET("/profile", func(c *gin.Context) {
			c.Set("user_id", "user123")
			h.GetProfile(c)
		})

		userDTO := &auth.UserDTO{ID: "user123", Email: "test@example.com"}
		mockSvc.On("GetProfile", mock.Anything, "user123").Return(userDTO, nil)

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp auth.UserDTO
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "user123", resp.ID)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)
		r := gin.Default()
		r.GET("/profile", h.GetProfile) // No user_id set

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)
		r := gin.Default()
		r.PUT("/profile", func(c *gin.Context) {
			c.Set("user_id", "user123")
			h.UpdateProfile(c)
		})

		reqBody := auth.UpdateProfileRequest{FullName: "New Name"}
		mockSvc.On("UpdateProfile", mock.Anything, "user123", reqBody).Return(nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCreateCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)
		r := gin.Default()
		r.POST("/company", func(c *gin.Context) {
			c.Set("user_id", "user123")
			h.CreateCompany(c)
		})

		reqBody := auth.CreateCompanyRequest{Name: "New Company"}
		companyDTO := &auth.CompanyDTO{ID: "comp123", Name: "New Company"}
		mockSvc.On("CreateCompany", mock.Anything, "user123", reqBody).Return(companyDTO, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/company", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestUpdateCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)
		r := gin.Default()
		r.PUT("/company/:id", h.UpdateCompany)

		reqBody := auth.UpdateCompanyRequest{Name: "Updated Company"}
		mockSvc.On("UpdateCompany", mock.Anything, "comp123", reqBody).Return(nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/company/comp123", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestVerifyCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)
		r := gin.Default()
		r.POST("/company/:id/verify", h.VerifyCompany)

		mockSvc.On("VerifyCompany", mock.Anything, "comp123").Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/company/comp123/verify", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestVerifyUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)
		r := gin.Default()
		r.POST("/users/:id/verify", h.VerifyUser)

		mockSvc.On("VerifyUser", mock.Anything, "user123").Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/users/user123/verify", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
