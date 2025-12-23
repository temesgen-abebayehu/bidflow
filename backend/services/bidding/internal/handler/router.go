package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/middleware"
)

func SetupRouter(h *HttpHandler, tm *auth.TokenManager) *gin.Engine {
	r := gin.Default()

	// Global Middleware
	// r.Use(middleware.LoggerMiddleware())

	api := r.Group("/api/v1/bids")
	{
		// Public routes
		api.GET("/:auction_id", h.GetBids)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(tm))
		{
			protected.POST("", h.PlaceBid)
		}
	}

	return r
}
