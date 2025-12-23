package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/domain"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) domain.BidRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, bid *domain.Bid) error {
	query := `
		INSERT INTO bids (id, auction_id, bidder_id, amount, timestamp)
		VALUES ($1, $2, $3, $4, $5)
	`
	if bid.Timestamp.IsZero() {
		bid.Timestamp = time.Now()
	}

	_, err := r.db.ExecContext(ctx, query,
		bid.ID, bid.AuctionID, bid.BidderID, bid.Amount, bid.Timestamp,
	)
	return err
}

func (r *postgresRepo) GetByID(ctx context.Context, id string) (*domain.Bid, error) {
	query := `SELECT id, auction_id, bidder_id, amount, timestamp FROM bids WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var b domain.Bid
	err := row.Scan(&b.ID, &b.AuctionID, &b.BidderID, &b.Amount, &b.Timestamp)
	if err == sql.ErrNoRows {
		return nil, domain.ErrBidNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *postgresRepo) ListByAuctionID(ctx context.Context, auctionID string) ([]domain.Bid, error) {
	query := `SELECT id, auction_id, bidder_id, amount, timestamp FROM bids WHERE auction_id = $1 ORDER BY amount DESC`
	rows, err := r.db.QueryContext(ctx, query, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []domain.Bid
	for rows.Next() {
		var b domain.Bid
		if err := rows.Scan(&b.ID, &b.AuctionID, &b.BidderID, &b.Amount, &b.Timestamp); err != nil {
			return nil, err
		}
		bids = append(bids, b)
	}
	return bids, nil
}

func (r *postgresRepo) GetHighestBid(ctx context.Context, auctionID string) (*domain.Bid, error) {
	query := `SELECT id, auction_id, bidder_id, amount, timestamp FROM bids WHERE auction_id = $1 ORDER BY amount DESC LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, auctionID)

	var b domain.Bid
	err := row.Scan(&b.ID, &b.AuctionID, &b.BidderID, &b.Amount, &b.Timestamp)
	if err == sql.ErrNoRows {
		return nil, nil // No bids yet
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}
