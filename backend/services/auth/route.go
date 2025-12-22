package main

import (
	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/middleware"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/handler"
)

func SetupRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, tm *auth.TokenManager) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/verify-otp", authHandler.VerifyOTP)
			authGroup.POST("/2fa/toggle", middleware.AuthMiddleware(tm), authHandler.Toggle2FA)
		}

		userGroup := api.Group("/users")
		userGroup.Use(middleware.AuthMiddleware(tm))
		{
			userGroup.GET("/profile", userHandler.GetProfile)
			userGroup.PUT("/profile", userHandler.UpdateProfile)

			// Admin routes (should check role, but for now just auth)
			userGroup.POST("/verify/:id", userHandler.VerifyUser)

			// Company routes
			userGroup.POST("/company", userHandler.CreateCompany)
			userGroup.PUT("/company/:id", userHandler.UpdateCompany)
			userGroup.POST("/company/:id/verify", userHandler.VerifyCompany)
		}
	}

	return r
}
