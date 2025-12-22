package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/domain"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(s domain.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id") // Assuming middleware sets this
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id") // Assuming middleware sets this
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req auth.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProfile(c.Request.Context(), userID, req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Profile updated successfully"})
}

func (h *UserHandler) CreateCompany(c *gin.Context) {
	userID := c.GetString("user_id")
	var req auth.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	company, err := h.service.CreateCompany(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, company)
}

func (h *UserHandler) UpdateCompany(c *gin.Context) {
	companyID := c.Param("id")
	var req auth.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateCompany(c.Request.Context(), companyID, req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Company updated successfully"})
}

func (h *UserHandler) VerifyCompany(c *gin.Context) {
	companyID := c.Param("id")
	if err := h.service.VerifyCompany(c.Request.Context(), companyID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Company verified successfully"})
}

func (h *UserHandler) VerifyUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.VerifyUser(c.Request.Context(), userID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User verified successfully"})
}
