package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/domain"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) domain.UserRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateUser(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (email, username, full_name, password, role, is_active, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		u.Email, u.Username, u.FullName, u.Password, u.Role, u.IsActive, time.Now(), time.Now()).Scan(&u.ID)

	return err
}

func (r *postgresRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, password, role, two_factor_enabled, two_factor_secret, full_name, username, company_id, is_active, is_verified FROM users WHERE email = $1",
		email).Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.TwoFactorEnabled, &u.TwoFactorSecret, &u.FullName, &u.Username, &u.CompanyID, &u.IsActive, &u.IsVerified)
	return u, err
}

func (r *postgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, password, role, two_factor_enabled, two_factor_secret, full_name, username, company_id, is_active, is_verified FROM users WHERE id = $1",
		id).Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.TwoFactorEnabled, &u.TwoFactorSecret, &u.FullName, &u.Username, &u.CompanyID, &u.IsActive, &u.IsVerified)
	return u, err
}

func (r *postgresRepo) UpdateUser(ctx context.Context, u *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET full_name = $1, username = $2, company_id = $3 WHERE id = $4",
		u.FullName, u.Username, u.CompanyID, u.ID)
	return err
}

func (r *postgresRepo) Update2FA(ctx context.Context, userID uuid.UUID, enabled bool, secret string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET two_factor_enabled = $1, two_factor_secret = $2 WHERE id = $3",
		enabled, secret, userID)
	return err
}

func (r *postgresRepo) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET is_verified = true WHERE id = $1", userID)
	return err
}
