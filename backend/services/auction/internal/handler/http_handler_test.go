package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/domain"
)

func TestCreateAuction_Http(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuctionService{
		CreateAuctionFunc: func(ctx context.Context, sellerID, title, description string, startPrice float64, startTime, endTime time.Time, category, imageURL string) (*domain.Auction, error) {
			return &domain.Auction{ID: "123"}, nil
		},
	}
	h := NewHttpHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := createAuctionRequest{
		Title:       "Test",
		Description: "Desc",
		StartPrice:  10.0,
		StartTime:   time.Now().Unix(),
		EndTime:     time.Now().Add(time.Hour).Unix(),
		Category:    "Cat",
		ImageURL:    "url",
	}
	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest("POST", "/auctions", bytes.NewBuffer(jsonBody))
	// Mock user ID in context (usually set by auth middleware)
	c.Set("user_id", "seller-1")

	h.CreateAuction(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestGetAuction_Http(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuctionService{
		GetAuctionFunc: func(ctx context.Context, id string) (*domain.Auction, error) {
			if id == "found" {
				return &domain.Auction{ID: "found"}, nil
			}
			return nil, errors.New("not found")
		},
	}
	h := NewHttpHandler(mockSvc)

	t.Run("Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "found"}}
		c.Request, _ = http.NewRequest("GET", "/auctions/found", nil)

		h.GetAuction(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "missing"}}
		c.Request, _ = http.NewRequest("GET", "/auctions/missing", nil)

		h.GetAuction(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestListAuctions_Http(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuctionService{
		ListAuctionsFunc: func(ctx context.Context, page, limit int, status string, category string) ([]domain.Auction, int64, error) {
			return []domain.Auction{{ID: "1"}}, 1, nil
		},
	}
	h := NewHttpHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/auctions?page=1&limit=10", nil)

	h.ListAuctions(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUpdateAuction_Http(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuctionService{
		UpdateAuctionFunc: func(ctx context.Context, id string, title, description, imageURL string) (*domain.Auction, error) {
			return &domain.Auction{ID: id}, nil
		},
	}
	h := NewHttpHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	reqBody := map[string]string{
		"title": "New Title",
	}
	jsonBody, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest("PUT", "/auctions/1", bytes.NewBuffer(jsonBody))

	h.UpdateAuction(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCloseAuction_Http(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuctionService{
		CloseAuctionFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}
	h := NewHttpHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("POST", "/auctions/1/close", nil)

	h.CloseAuction(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
