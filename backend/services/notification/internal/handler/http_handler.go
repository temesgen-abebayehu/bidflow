package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/domain"
	ws "github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/websocket"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	service      domain.NotificationService
	hub          *ws.Hub
	tokenManager *auth.TokenManager
	log          logger.Logger
	upgrader     websocket.Upgrader
}

func NewNotificationHandler(service domain.NotificationService, hub *ws.Hub, tm *auth.TokenManager, log logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		service:      service,
		hub:          hub,
		tokenManager: tm,
		log:          log,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
	}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notifications, err := h.service.GetUserNotifications(c.Request.Context(), userID.(string))
	if err != nil {
		h.log.Error("Failed to get notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	claims, err := h.tokenManager.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("Failed to upgrade websocket", zap.Error(err))
		return
	}

	client := ws.NewClient(h.hub, conn, claims.UserID, h.log)
	h.hub.Register(client)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}
