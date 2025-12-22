package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/domain"
)

type companyRepo struct {
	db *sql.DB
}

func NewCompanyRepo(db *sql.DB) domain.CompanyRepository {
	return &companyRepo{db: db}
}

func (r *companyRepo) CreateCompany(ctx context.Context, c *domain.Company) error {
	query := `INSERT INTO companies (id, name, is_verified, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, c.ID, c.Name, c.IsVerified, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *companyRepo) GetCompanyByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	c := &domain.Company{}
	query := `SELECT id, name, is_verified, created_at, updated_at FROM companies WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.IsVerified, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *companyRepo) UpdateCompany(ctx context.Context, c *domain.Company) error {
	query := `UPDATE companies SET name = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, c.Name, time.Now(), c.ID)
	return err
}

func (r *companyRepo) VerifyCompany(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE companies SET is_verified = true, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
