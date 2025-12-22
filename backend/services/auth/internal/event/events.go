package event

import (
	"time"

	"github.com/google/uuid"
)

const (
	TopicUserRegistered = "user.registered"
	TopicUserVerified   = "user.verified"
)

type UserRegisteredEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullname"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}

type UserVerifiedEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}
