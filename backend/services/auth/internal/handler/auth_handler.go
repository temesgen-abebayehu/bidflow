package handler

// AuthHandler handles authentication requests
import (
	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/domain"
)

type AuthHandler struct {
	service domain.AuthService
}

func NewAuthHandler(s domain.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(c.Request.Context(), req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, token, mfaRequired, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	if mfaRequired {
		c.JSON(200, gin.H{"mfa_required": true, "message": "Please enter OTP"})
		return
	}

	c.JSON(200, auth.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req auth.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Verify2FA(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *AuthHandler) Toggle2FA(c *gin.Context) {
	userID := c.GetString("user_id")
	var req auth.Toggle2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	secret, err := h.service.Toggle2FA(c.Request.Context(), userID, req.Enable)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if req.Enable {
		c.JSON(200, gin.H{"message": "2FA enabled", "secret": secret})
	} else {
		c.JSON(200, gin.H{"message": "2FA disabled"})
	}
}
