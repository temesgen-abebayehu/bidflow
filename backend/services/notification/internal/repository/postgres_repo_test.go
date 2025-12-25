package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/domain"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepo(db)

	notification := &domain.Notification{
		ID:         "1",
		UserID:     "user-1",
		Type:       "INFO",
		Title:      "Test",
		Message:    "Message",
		ResourceID: "res-1",
		IsRead:     false,
		CreatedAt:  time.Now(),
	}

	mock.ExpectExec("INSERT INTO notifications").
		WithArgs(notification.ID, notification.UserID, notification.Type, notification.Title, notification.Message, notification.ResourceID, notification.IsRead, notification.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), notification)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepo(db)

	userID := "user-1"
	limit := 10
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "type", "title", "message", "resource_id", "is_read", "created_at"}).
		AddRow("1", userID, "INFO", "Test", "Message", "res-1", false, createdAt)

	mock.ExpectQuery("SELECT id, user_id, type, title, message, resource_id, is_read, created_at FROM notifications").
		WithArgs(userID, limit).
		WillReturnRows(rows)

	notifications, err := repo.ListByUserID(context.Background(), userID, limit)
	assert.NoError(t, err)
	assert.Len(t, notifications, 1)
	assert.Equal(t, "1", notifications[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkAsRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepo(db)

	id := "1"

	mock.ExpectExec("UPDATE notifications SET is_read = TRUE").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.MarkAsRead(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
