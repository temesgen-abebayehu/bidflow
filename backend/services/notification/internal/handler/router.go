package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/middleware"
)

func SetupRouter(h *NotificationHandler, tm *auth.TokenManager) *gin.Engine {
	r := gin.Default()

	// Global Middleware
	// r.Use(middleware.CORSMiddleware())

	api := r.Group("/api/v1/notifications")
	{
		// WebSocket endpoint (needs custom auth handling inside handler usually, or query param auth)
		api.GET("/ws", h.HandleWebSocket)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(tm))
		{
			protected.GET("", h.GetNotifications)
		}
	}

	return r
}
