package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/middleware"
)

func SetupRouter(h *HttpHandler, tm *auth.TokenManager) *gin.Engine {
	r := gin.Default()

	// Global Middleware
	// You might want to add logger and recovery middleware from common here as well if needed
	// r.Use(middleware.LoggerMiddleware())
	// r.Use(middleware.RecoveryMiddleware())

	api := r.Group("/api/v1/auctions")
	{
		// Public routes
		api.GET("", h.ListAuctions)
		api.GET("/:id", h.GetAuction)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(tm))
		{
			protected.POST("", h.CreateAuction)
			protected.PUT("/:id", h.UpdateAuction)
			protected.POST("/:id/close", h.CloseAuction)
		}
	}

	return r
}
