package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/boykush/foresee/services/collector/internal/domain"
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

// Save persists a feed to PostgreSQL
func (r *FeedRepository) Save(ctx context.Context, feed domain.Feed) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO feeds (id, url, title, created_at) VALUES ($1, $2, $3, $4)",
		feed.ID(), feed.URL(), feed.Title(), feed.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert feed: %w", err)
	}
	return nil
}

// ListAll returns all feeds from PostgreSQL
func (r *FeedRepository) ListAll(ctx context.Context) ([]domain.Feed, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, url, title, created_at FROM feeds")
	if err != nil {
		return nil, fmt.Errorf("failed to query feeds: %w", err)
	}
	defer rows.Close()

	var feeds []domain.Feed
	for rows.Next() {
		var id, url, title string
		var createdAt time.Time
		if err := rows.Scan(&id, &url, &title, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan feed row: %w", err)
		}
		feed, err := domain.NewFeed(id, url, title, createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to create feed domain object: %w", err)
		}
		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating feed rows: %w", err)
	}

	return feeds, nil
}
