package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/boykush/foresee/services/feed/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ArticleRepository implements domain.ArticleRepository using PostgreSQL
type ArticleRepository struct {
	pool *pgxpool.Pool
}

// NewArticleRepository creates a new ArticleRepository
func NewArticleRepository(pool *pgxpool.Pool) *ArticleRepository {
	return &ArticleRepository{pool: pool}
}

// List returns articles from PostgreSQL, optionally filtered by feed ID
func (r *ArticleRepository) List(ctx context.Context, feedID *uuid.UUID) ([]domain.Article, error) {
	var query string
	var args []any

	if feedID != nil {
		query = "SELECT id, feed_id, title, url, published_at FROM articles WHERE feed_id = $1 ORDER BY published_at DESC"
		args = append(args, *feedID)
	} else {
		query = "SELECT id, feed_id, title, url, published_at FROM articles ORDER BY published_at DESC"
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	var articles []domain.Article
	for rows.Next() {
		var (
			id          uuid.UUID
			artFeedID   uuid.UUID
			title       string
			url         string
			publishedAt time.Time
		)
		if err := rows.Scan(&id, &artFeedID, &title, &url, &publishedAt); err != nil {
			return nil, fmt.Errorf("failed to scan article row: %w", err)
		}

		article, err := domain.NewArticle(id, artFeedID, title, url, publishedAt)
		if err != nil {
			continue
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating article rows: %w", err)
	}

	return articles, nil
}
