package service

import (
	"context"
	"errors"

	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/domain"
)

type UserService struct {
	repo        domain.UserRepository
	companyRepo domain.CompanyRepository
}

func NewUserService(r domain.UserRepository, cr domain.CompanyRepository) domain.UserService {
	return &UserService{repo: r, companyRepo: cr}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*auth.UserDTO, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &auth.UserDTO{
		ID:        u.ID.String(),
		Email:     u.Email,
		Username:  u.Username,
		FullName:  u.FullName,
		Role:      u.Role,
		CompanyID: u.CompanyID.String,
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req auth.UpdateProfileRequest) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if req.FullName != "" {
		u.FullName = req.FullName
	}
	if req.Username != "" {
		u.Username = req.Username
	}

	return s.repo.UpdateUser(ctx, u)
}

func (s *UserService) VerifyUser(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}
	return s.repo.VerifyUser(ctx, id)
}

func (s *UserService) CreateCompany(ctx context.Context, userID string, req auth.CreateCompanyRequest) (*auth.CompanyDTO, error) {
	// 1. Create Company
	companyID := uuid.New()
	company := &domain.Company{
		ID:        companyID,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.companyRepo.CreateCompany(ctx, company); err != nil {
		return nil, err
	}

	// 2. Link User to Company (Assuming user becomes admin/owner)
	// This part depends on business logic. Should we update the user's company_id?
	// For now, let's assume yes.
	uid, _ := uuid.Parse(userID)
	u, err := s.repo.GetByID(ctx, uid)
	if err == nil {
		u.CompanyID = sql.NullString{String: companyID.String(), Valid: true}
		if err := s.repo.UpdateUser(ctx, u); err != nil {
			// Log error but don't fail the request? Or fail?
			// Ideally transactional.
			return nil, err
		}
	}

	return &auth.CompanyDTO{
		ID:        company.ID.String(),
		Name:      company.Name,
		CreatedAt: company.CreatedAt.Format(time.RFC3339),
		UpdatedAt: company.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) UpdateCompany(ctx context.Context, companyID string, req auth.UpdateCompanyRequest) error {
	id, err := uuid.Parse(companyID)
	if err != nil {
		return errors.New("invalid company id")
	}

	c, err := s.companyRepo.GetCompanyByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		c.Name = req.Name
	}

	return s.companyRepo.UpdateCompany(ctx, c)
}

func (s *UserService) VerifyCompany(ctx context.Context, companyID string) error {
	id, err := uuid.Parse(companyID)
	if err != nil {
		return errors.New("invalid company id")
	}
	return s.companyRepo.VerifyCompany(ctx, id)
}

func (s *UserService) GetCompany(ctx context.Context, companyID string) (*auth.CompanyDTO, error) {
	id, err := uuid.Parse(companyID)
	if err != nil {
		return nil, errors.New("invalid company id")
	}

	c, err := s.companyRepo.GetCompanyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &auth.CompanyDTO{
		ID:         c.ID.String(),
		Name:       c.Name,
		IsVerified: c.IsVerified,
		CreatedAt:  c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  c.UpdatedAt.Format(time.RFC3339),
	}, nil
}
