package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	Update2FA(ctx context.Context, userID uuid.UUID, enabled bool, secret string) error
	VerifyUser(ctx context.Context, userID uuid.UUID) error
}

type CompanyRepository interface {
	CreateCompany(ctx context.Context, company *Company) error
	GetCompanyByID(ctx context.Context, id uuid.UUID) (*Company, error)
	UpdateCompany(ctx context.Context, company *Company) error
	VerifyCompany(ctx context.Context, id uuid.UUID) error
}

type AuthService interface {
	Register(ctx context.Context, req auth.RegisterRequest) error
	Login(ctx context.Context, email, password string) (*auth.UserDTO, string, bool, error)
	Verify2FA(ctx context.Context, email, code string) (string, error)
	Toggle2FA(ctx context.Context, userID string, enable bool) (string, error) // Returns secret if enabling
}

type UserService interface {
	GetProfile(ctx context.Context, userID string) (*auth.UserDTO, error)
	UpdateProfile(ctx context.Context, userID string, req auth.UpdateProfileRequest) error
	VerifyUser(ctx context.Context, userID string) error

	CreateCompany(ctx context.Context, userID string, req auth.CreateCompanyRequest) (*auth.CompanyDTO, error)
	UpdateCompany(ctx context.Context, companyID string, req auth.UpdateCompanyRequest) error
	VerifyCompany(ctx context.Context, companyID string) error
	GetCompany(ctx context.Context, companyID string) (*auth.CompanyDTO, error)
}

type EventProducer interface {
	PublishUserRegistered(ctx context.Context, user *User) error
	PublishUserVerified(ctx context.Context, userID uuid.UUID) error
}
