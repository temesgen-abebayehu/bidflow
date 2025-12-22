package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"go.uber.org/zap"
)

func RecoveryMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Recovery from panic",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}