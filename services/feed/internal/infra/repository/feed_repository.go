package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/boykush/foresee/services/feed/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FeedRepository implements domain.FeedRepository using PostgreSQL
type FeedRepository struct {
	pool *pgxpool.Pool
}

// NewFeedRepository creates a new FeedRepository
func NewFeedRepository(pool *pgxpool.Pool) *FeedRepository {
	return &FeedRepository{pool: pool}
}

// List returns all feeds from PostgreSQL
func (r *FeedRepository) List(ctx context.Context) ([]domain.Feed, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, url, title, created_at FROM feeds ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query feeds: %w", err)
	}
	defer rows.Close()

	var feeds []domain.Feed
	for rows.Next() {
		var (
			id        uuid.UUID
			url       string
			title     string
			createdAt time.Time
		)
		if err := rows.Scan(&id, &url, &title, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan feed row: %w", err)
		}

		feed, err := domain.NewFeed(id, url, title, createdAt)
		if err != nil {
			continue
		}
		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating feed rows: %w", err)
	}

	return feeds, nil
}
