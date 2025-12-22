package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/domain"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	auction := &domain.Auction{
		ID:          "1",
		SellerID:    "seller-1",
		Title:       "Test",
		Description: "Desc",
		StartPrice:  10.0,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Category:    "Cat",
		ImageURL:    "url",
	}

	mock.ExpectExec("INSERT INTO auctions").
		WithArgs(auction.ID, auction.SellerID, auction.Title, auction.Description, auction.StartPrice, auction.CurrentPrice, auction.Status, auction.StartTime, auction.EndTime, auction.Category, auction.ImageURL, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), auction)
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

	rows := sqlmock.NewRows([]string{"id", "seller_id", "title", "description", "start_price", "current_price", "status", "start_time", "end_time", "category", "image_url", "created_at", "updated_at"}).
		AddRow("1", "seller-1", "Test", "Desc", 10.0, 10.0, "ACTIVE", time.Now(), time.Now().Add(time.Hour), "Cat", "url", time.Now(), time.Now())

	mock.ExpectQuery("SELECT .* FROM auctions WHERE id = \\$1").
		WithArgs("1").
		WillReturnRows(rows)

	auction, err := repo.GetByID(context.Background(), "1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if auction.ID != "1" {
		t.Errorf("expected id 1, got %s", auction.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	auction := &domain.Auction{
		ID:          "1",
		Title:       "Updated",
		Description: "Desc",
		Status:      domain.AuctionStatusActive,
	}

	mock.ExpectExec("UPDATE auctions SET").
		WithArgs(auction.Title, auction.Description, auction.CurrentPrice, auction.Status, auction.ImageURL, sqlmock.AnyArg(), auction.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), auction)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	mock.ExpectExec("DELETE FROM auctions WHERE id = \\$1").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), "1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	rows := sqlmock.NewRows([]string{"id", "seller_id", "title", "description", "start_price", "current_price", "status", "start_time", "end_time", "category", "image_url", "created_at", "updated_at"}).
		AddRow("1", "seller-1", "Test", "Desc", 10.0, 10.0, "ACTIVE", time.Now(), time.Now().Add(time.Hour), "Cat", "url", time.Now(), time.Now())

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM auctions").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery("SELECT id, seller_id, title, description").
		WillReturnRows(rows)

	auctions, count, err := repo.List(context.Background(), 1, 10, "", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
	if len(auctions) != 1 {
		t.Errorf("expected 1 auction, got %d", len(auctions))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
