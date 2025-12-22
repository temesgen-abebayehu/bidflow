package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"go.uber.org/zap"
)

func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is done
		end := time.Now()
		latency := end.Sub(start)

		log.Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		)
	}
}