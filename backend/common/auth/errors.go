package auth

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrNoToken      = errors.New("authorization token is missing")
	ErrUnauthorized = errors.New("user is not authorized for this action")
	ErrForbidden    = errors.New("user does not have the required permissions")
)