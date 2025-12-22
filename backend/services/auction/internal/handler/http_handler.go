package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/service"
)

type HttpHandler struct {
	service *service.AuctionService
}

func NewHttpHandler(service *service.AuctionService) *HttpHandler {
	return &HttpHandler{service: service}
}

type createAuctionRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	StartPrice  float64 `json:"start_price" binding:"required"`
	StartTime   int64   `json:"start_time" binding:"required"`
	EndTime     int64   `json:"end_time" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	ImageURL    string  `json:"image_url"`
}

func (h *HttpHandler) CreateAuction(c *gin.Context) {
	var req createAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get seller ID from context (set by auth middleware)
	sellerID := c.GetString("user_id")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startTime := time.Unix(req.StartTime, 0)
	endTime := time.Unix(req.EndTime, 0)

	auction, err := h.service.CreateAuction(
		c.Request.Context(),
		sellerID,
		req.Title,
		req.Description,
		req.StartPrice,
		startTime,
		endTime,
		req.Category,
		req.ImageURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, auction)
}

func (h *HttpHandler) GetAuction(c *gin.Context) {
	id := c.Param("id")
	auction, err := h.service.GetAuction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		return
	}
	c.JSON(http.StatusOK, auction)
}

func (h *HttpHandler) ListAuctions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	category := c.Query("category")

	auctions, total, err := h.service.ListAuctions(c.Request.Context(), page, limit, status, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  auctions,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *HttpHandler) UpdateAuction(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auction, err := h.service.UpdateAuction(c.Request.Context(), id, req.Title, req.Description, req.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, auction)
}

func (h *HttpHandler) CloseAuction(c *gin.Context) {
	id := c.Param("id")
	err := h.service.CloseAuction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "auction closed"})
}
