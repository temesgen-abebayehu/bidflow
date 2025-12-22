package domain

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID         uuid.UUID
	Name       string
	LogoURL    string
	FoundedDate string
	Area       string
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
