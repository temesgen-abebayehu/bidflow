package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID
	Email            string
	Username         string
	FullName         string
	Password         string
	Role             string
	CompanyID        sql.NullString
	IsVerified       bool
	IsActive         bool
	TwoFactorEnabled bool
	TwoFactorSecret  sql.NullString
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
