package routes

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/config"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "API Gateway is running"})
	})

	// Helper function to create a reverse proxy handler
	proxy := func(target string) gin.HandlerFunc {
		return func(c *gin.Context) {
			remote, err := url.Parse(target)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid service URL"})
				return
			}

			proxy := httputil.NewSingleHostReverseProxy(remote)

			// Update the director to ensure the host header is set correctly
			originalDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				originalDirector(req)
				req.Host = remote.Host // Set Host header to the target service
			}

			proxy.ServeHTTP(c.Writer, c.Request)
		}
	}

	api := r.Group("/api/v1")
	{
		// Auth Service (handles both auth and users)
		api.Any("/auth/*any", proxy(cfg.AuthServiceURL))
		api.Any("/users/*any", proxy(cfg.AuthServiceURL))

		// Auction Service
		api.Any("/auctions/*any", proxy(cfg.AuctionServiceURL))

		// Bidding Service
		api.Any("/bids/*any", proxy(cfg.BiddingServiceURL))

		// Notification Service
		api.Any("/notifications/*any", proxy(cfg.NotificationServiceURL))
	}

	return r
}
