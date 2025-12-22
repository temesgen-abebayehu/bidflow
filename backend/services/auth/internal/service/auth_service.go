package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/domain"
)

type AuthService struct {
	repo         domain.UserRepository
	tokenManager *auth.TokenManager
	producer     domain.EventProducer
}

func NewAuthService(r domain.UserRepository, tm *auth.TokenManager, p domain.EventProducer) domain.AuthService {
	return &AuthService{repo: r, tokenManager: tm, producer: p}
}

func (s *AuthService) Register(ctx context.Context, req auth.RegisterRequest) error {
	hashedPassword, _ := auth.HashPassword(req.Password)
	user := &domain.User{
		Email:    req.Email,
		Username: req.Email, // Default username to email
		FullName: req.Email, // Default fullname to email
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: true,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}
	return s.producer.PublishUserRegistered(ctx, user)
}

// Login returns (userDTO, token, mfaRequired, error)
func (s *AuthService) Login(ctx context.Context, email, password string) (*auth.UserDTO, string, bool, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", false, errors.New("user not found")
	}

	if !auth.CheckPasswordHash(password, u.Password) {
		return nil, "", false, errors.New("invalid credentials")
	}

	userDTO := &auth.UserDTO{
		ID:        u.ID.String(),
		Email:     u.Email,
		Username:  u.Username,
		FullName:  u.FullName,
		Role:      u.Role,
		CompanyID: u.CompanyID.String,
	}

	// If 2FA is on, don't give the JWT yet
	if u.TwoFactorEnabled {
		return userDTO, "", true, nil
	}

	token, _ := s.tokenManager.GenerateToken(u.ID.String(), "", u.Role)
	return userDTO, token, false, nil
}

func (s *AuthService) Verify2FA(ctx context.Context, email, code string) (string, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil || !u.TwoFactorSecret.Valid {
		return "", errors.New("2FA not configured")
	}

	valid := totp.Validate(code, u.TwoFactorSecret.String)
	if !valid {
		return "", errors.New("invalid OTP code")
	}

	return s.tokenManager.GenerateToken(u.ID.String(), "", u.Role)
}

func (s *AuthService) Toggle2FA(ctx context.Context, userID string, enable bool) (string, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return "", errors.New("invalid user id")
	}

	var secret string
	if enable {
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "BidFlow",
			AccountName: userID, // Should ideally be email
		})
		if err != nil {
			return "", err
		}
		secret = key.Secret()
	}

	// If disabling, secret is empty (or we keep it but disable flag)
	// Let's assume we clear it or just disable flag.
	// Repo Update2FA takes secret.

	err = s.repo.Update2FA(ctx, id, enable, secret)
	return secret, err
}
