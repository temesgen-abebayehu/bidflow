package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/domain"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/service"
)

// MockCompanyRepository
type MockCompanyRepository struct {
	mock.Mock
}

func (m *MockCompanyRepository) CreateCompany(ctx context.Context, company *domain.Company) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func (m *MockCompanyRepository) GetCompanyByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Company), args.Error(1)
}

func (m *MockCompanyRepository) UpdateCompany(ctx context.Context, company *domain.Company) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func (m *MockCompanyRepository) VerifyCompany(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCompanyRepo := new(MockCompanyRepository)
	svc := service.NewUserService(mockRepo, mockCompanyRepo)

	userID := uuid.New()
	user := &domain.User{
		ID:    userID,
		Email: "test@example.com",
	}

	mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	dto, err := svc.GetProfile(context.Background(), userID.String())
	assert.NoError(t, err)
	assert.Equal(t, user.Email, dto.Email)
}

func TestCreateCompany(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCompanyRepo := new(MockCompanyRepository)
	svc := service.NewUserService(mockRepo, mockCompanyRepo)

	userID := uuid.New()
	user := &domain.User{
		ID:    userID,
		Email: "test@example.com",
	}

	req := auth.CreateCompanyRequest{Name: "New Company"}

	mockCompanyRepo.On("CreateCompany", mock.Anything, mock.MatchedBy(func(c *domain.Company) bool {
		return c.Name == req.Name
	})).Return(nil)

	mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
	mockRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)

	dto, err := svc.CreateCompany(context.Background(), userID.String(), req)
	assert.NoError(t, err)
	assert.Equal(t, req.Name, dto.Name)
}

func TestUpdateProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCompanyRepo := new(MockCompanyRepository)
	svc := service.NewUserService(mockRepo, mockCompanyRepo)

	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Email:    "test@example.com",
		FullName: "Old Name",
	}

	req := auth.UpdateProfileRequest{FullName: "New Name"}

	mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
	mockRepo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.FullName == "New Name"
	})).Return(nil)

	err := svc.UpdateProfile(context.Background(), userID.String(), req)
	assert.NoError(t, err)
}

func TestVerifyUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCompanyRepo := new(MockCompanyRepository)
	svc := service.NewUserService(mockRepo, mockCompanyRepo)

	userID := uuid.New()

	mockRepo.On("VerifyUser", mock.Anything, userID).Return(nil)

	err := svc.VerifyUser(context.Background(), userID.String())
	assert.NoError(t, err)
}

func TestUpdateCompany(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCompanyRepo := new(MockCompanyRepository)
	svc := service.NewUserService(mockRepo, mockCompanyRepo)

	companyID := uuid.New()
	company := &domain.Company{
		ID:   companyID,
		Name: "Old Company",
	}

	req := auth.UpdateCompanyRequest{Name: "Updated Company"}

	mockCompanyRepo.On("GetCompanyByID", mock.Anything, companyID).Return(company, nil)
	mockCompanyRepo.On("UpdateCompany", mock.Anything, mock.MatchedBy(func(c *domain.Company) bool {
		return c.Name == "Updated Company"
	})).Return(nil)

	err := svc.UpdateCompany(context.Background(), companyID.String(), req)
	assert.NoError(t, err)
}

func TestVerifyCompany(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCompanyRepo := new(MockCompanyRepository)
	svc := service.NewUserService(mockRepo, mockCompanyRepo)

	companyID := uuid.New()

	mockCompanyRepo.On("VerifyCompany", mock.Anything, companyID).Return(nil)

	err := svc.VerifyCompany(context.Background(), companyID.String())
	assert.NoError(t, err)
}

func TestGetCompany(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCompanyRepo := new(MockCompanyRepository)
	svc := service.NewUserService(mockRepo, mockCompanyRepo)

	companyID := uuid.New()
	company := &domain.Company{
		ID:        companyID,
		Name:      "Test Company",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockCompanyRepo.On("GetCompanyByID", mock.Anything, companyID).Return(company, nil)

	dto, err := svc.GetCompany(context.Background(), companyID.String())
	assert.NoError(t, err)
	assert.Equal(t, company.Name, dto.Name)
}
