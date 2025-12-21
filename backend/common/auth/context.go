package auth

import "context"

type contextKey string

const claimsKey contextKey = "user_claims"

// ToContext attaches UserClaims to a context (Used by Middleware)
func ToContext(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// FromContext retrieves UserClaims from a context (Used by Logic/Handlers)
func FromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(claimsKey).(*UserClaims)
	return claims, ok
}