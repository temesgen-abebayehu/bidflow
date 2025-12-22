package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
)

// AuthMiddleware creates a gin-middleware that validates the JWT token
func AuthMiddleware(tm *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the Authorization header (Format: Bearer <token>)
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		splitToken := strings.Split(header, "Bearer ")
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		tokenString := splitToken[1]

		// 2. Use our common/auth tool to verify the token
		claims, err := tm.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 3. Save the claims into the Context (using the tool we wrote in common/auth/context.go)
		// This makes UserID/Role available to all subsequent logic
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("company_id", claims.CompanyID)

		c.Next()
	}
}