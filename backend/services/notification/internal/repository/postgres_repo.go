package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/domain"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) domain.NotificationRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, n *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, message, resource_id, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}

	_, err := r.db.ExecContext(ctx, query,
		n.ID, n.UserID, n.Type, n.Title, n.Message, n.ResourceID, n.IsRead, n.CreatedAt,
	)
	return err
}

func (r *postgresRepo) ListByUserID(ctx context.Context, userID string, limit int) ([]domain.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, resource_id, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.ResourceID, &n.IsRead, &n.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *postgresRepo) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
