package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/domain"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/service"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update2FA(ctx context.Context, userID uuid.UUID, enabled bool, secret string) error {
	args := m.Called(ctx, userID, enabled, secret)
	return args.Error(0)
}

func (m *MockUserRepository) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestRegister(t *testing.T) {
	mockRepo := new(MockUserRepository)
	tm := auth.NewTokenManager("secret")
	svc := service.NewAuthService(mockRepo, tm)

	req := auth.RegisterRequest{
		Email:    "test@example.com",
		Password: "password",
		Role:     "BIDDER",
	}

	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == req.Email && u.Role == req.Role
	})).Return(nil)

	err := svc.Register(context.Background(), req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	tm := auth.NewTokenManager("secret")
	svc := service.NewAuthService(mockRepo, tm)

	hashedPassword, _ := auth.HashPassword("password")
	user := &domain.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     "BIDDER",
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	userDTO, token, mfa, err := svc.Login(context.Background(), "test@example.com", "password")
	assert.NoError(t, err)
	assert.NotNil(t, userDTO)
	assert.NotEmpty(t, token)
	assert.False(t, mfa)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	tm := auth.NewTokenManager("secret")
	svc := service.NewAuthService(mockRepo, tm)

	hashedPassword, _ := auth.HashPassword("password")
	user := &domain.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: hashedPassword,
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	_, _, _, err := svc.Login(context.Background(), "test@example.com", "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestVerify2FA(t *testing.T) {
	mockRepo := new(MockUserRepository)
	tm := auth.NewTokenManager("secret")
	svc := service.NewAuthService(mockRepo, tm)

	// Generate a real secret for testing
	key, _ := totp.Generate(totp.GenerateOpts{Issuer: "Test", AccountName: "test@example.com"})
	secret := key.Secret()
	code, _ := totp.GenerateCode(secret, time.Now())

	user := &domain.User{
		ID:              uuid.New(),
		Email:           "test@example.com",
		TwoFactorSecret: sql.NullString{String: secret, Valid: true},
		Role:            "BIDDER",
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	token, err := svc.Verify2FA(context.Background(), "test@example.com", code)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestToggle2FA(t *testing.T) {
	mockRepo := new(MockUserRepository)
	tm := auth.NewTokenManager("secret")
	svc := service.NewAuthService(mockRepo, tm)

	userID := uuid.New()

	mockRepo.On("Update2FA", mock.Anything, userID, true, mock.AnythingOfType("string")).Return(nil)

	secret, err := svc.Toggle2FA(context.Background(), userID.String(), true)
	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
}
