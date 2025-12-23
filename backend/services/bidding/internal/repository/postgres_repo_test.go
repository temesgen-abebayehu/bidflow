package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/domain"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	bid := &domain.Bid{
		ID:        "bid-1",
		AuctionID: "auction-1",
		BidderID:  "user-1",
		Amount:    100.0,
		Timestamp: time.Now(),
	}

	mock.ExpectExec("INSERT INTO bids").
		WithArgs(bid.ID, bid.AuctionID, bid.BidderID, bid.Amount, bid.Timestamp).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), bid)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	rows := sqlmock.NewRows([]string{"id", "auction_id", "bidder_id", "amount", "timestamp"}).
		AddRow("bid-1", "auction-1", "user-1", 100.0, time.Now())

	mock.ExpectQuery("SELECT id, auction_id, bidder_id, amount, timestamp FROM bids WHERE id = \\$1").
		WithArgs("bid-1").
		WillReturnRows(rows)

	bid, err := repo.GetByID(context.Background(), "bid-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if bid.ID != "bid-1" {
		t.Errorf("expected bid id 'bid-1', got '%s'", bid.ID)
	}
}

func TestListByAuctionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	rows := sqlmock.NewRows([]string{"id", "auction_id", "bidder_id", "amount", "timestamp"}).
		AddRow("bid-1", "auction-1", "user-1", 100.0, time.Now()).
		AddRow("bid-2", "auction-1", "user-2", 90.0, time.Now())

	mock.ExpectQuery("SELECT id, auction_id, bidder_id, amount, timestamp FROM bids WHERE auction_id = \\$1 ORDER BY amount DESC").
		WithArgs("auction-1").
		WillReturnRows(rows)

	bids, err := repo.ListByAuctionID(context.Background(), "auction-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(bids) != 2 {
		t.Errorf("expected 2 bids, got %d", len(bids))
	}
}
