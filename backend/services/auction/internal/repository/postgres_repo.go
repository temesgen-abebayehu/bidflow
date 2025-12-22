package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/domain"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) domain.AuctionRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, auction *domain.Auction) error {
	query := `
		INSERT INTO auctions (
			id, seller_id, title, description, start_price, current_price, 
			status, start_time, end_time, category, image_url, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	now := time.Now()
	auction.CreatedAt = now
	auction.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		auction.ID, auction.SellerID, auction.Title, auction.Description,
		auction.StartPrice, auction.CurrentPrice, auction.Status,
		auction.StartTime, auction.EndTime, auction.Category,
		auction.ImageURL, auction.CreatedAt, auction.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) GetByID(ctx context.Context, id string) (*domain.Auction, error) {
	query := `
		SELECT id, seller_id, title, description, start_price, current_price, 
		       status, start_time, end_time, category, image_url, created_at, updated_at
		FROM auctions WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var a domain.Auction
	err := row.Scan(
		&a.ID, &a.SellerID, &a.Title, &a.Description, &a.StartPrice, &a.CurrentPrice,
		&a.Status, &a.StartTime, &a.EndTime, &a.Category, &a.ImageURL, &a.CreatedAt, &a.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAuctionNotFound
	}
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *postgresRepo) Update(ctx context.Context, auction *domain.Auction) error {
	query := `
		UPDATE auctions SET 
			title = $1, description = $2, current_price = $3, status = $4, 
			image_url = $5, updated_at = $6
		WHERE id = $7
	`

	auction.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		auction.Title, auction.Description, auction.CurrentPrice, auction.Status,
		auction.ImageURL, auction.UpdatedAt, auction.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrAuctionNotFound
	}

	return nil
}

func (r *postgresRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM auctions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrAuctionNotFound
	}

	return nil
}

func (r *postgresRepo) List(ctx context.Context, page, limit int, status domain.AuctionStatus, category string) ([]domain.Auction, int64, error) {
	offset := (page - 1) * limit

	baseQuery := `SELECT id, seller_id, title, description, start_price, current_price, 
		                 status, start_time, end_time, category, image_url, created_at, updated_at
		          FROM auctions WHERE 1=1`

	countQuery := `SELECT COUNT(*) FROM auctions WHERE 1=1`

	var args []interface{}
	argID := 1

	if status != "" {
		filter := fmt.Sprintf(" AND status = $%d", argID)
		baseQuery += filter
		countQuery += filter
		args = append(args, status)
		argID++
	}

	if category != "" {
		filter := fmt.Sprintf(" AND category = $%d", argID)
		baseQuery += filter
		countQuery += filter
		args = append(args, category)
		argID++
	}

	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var auctions []domain.Auction
	for rows.Next() {
		var a domain.Auction
		err := rows.Scan(
			&a.ID, &a.SellerID, &a.Title, &a.Description, &a.StartPrice, &a.CurrentPrice,
			&a.Status, &a.StartTime, &a.EndTime, &a.Category, &a.ImageURL, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		auctions = append(auctions, a)
	}

	return auctions, total, nil
}
