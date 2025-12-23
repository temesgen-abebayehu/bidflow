package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/service"
)

type HttpHandler struct {
	service *service.BiddingService
}

func NewHttpHandler(service *service.BiddingService) *HttpHandler {
	return &HttpHandler{service: service}
}

type placeBidRequest struct {
	AuctionID string  `json:"auction_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}

func (h *HttpHandler) PlaceBid(c *gin.Context) {
	var req placeBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bid, err := h.service.PlaceBid(c.Request.Context(), req.AuctionID, userID.(string), req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bid)
}

func (h *HttpHandler) GetBids(c *gin.Context) {
	auctionID := c.Param("auction_id")
	if auctionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auction_id is required"})
		return
	}

	bids, err := h.service.GetBidsByAuction(c.Request.Context(), auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}
