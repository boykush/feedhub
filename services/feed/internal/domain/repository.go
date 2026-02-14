package domain

import (
	"context"

	"github.com/google/uuid"
)

// FeedRepository defines the interface for feed data access
type FeedRepository interface {
	// List returns all feeds
	List(ctx context.Context) ([]Feed, error)
}

// ArticleRepository defines the interface for article data access
type ArticleRepository interface {
	// List returns all articles, optionally filtered by feed ID
	List(ctx context.Context, feedID *uuid.UUID) ([]Article, error)
}
