package repository

import (
	"context"
	"fmt"

	"github.com/boykush/foresee/services/collector/internal/domain"
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

// BulkUpsert inserts articles, ignoring duplicates by URL
func (r *ArticleRepository) BulkUpsert(ctx context.Context, articles []domain.Article) error {
	if len(articles) == 0 {
		return nil
	}

	for _, article := range articles {
		_, err := r.pool.Exec(ctx,
			"INSERT INTO articles (id, feed_id, title, url, published_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (url) DO NOTHING",
			article.ID(), article.FeedID(), article.Title(), article.URL(), article.PublishedAt(),
		)
		if err != nil {
			return fmt.Errorf("failed to upsert article: %w", err)
		}
	}

	return nil
}
