package domain

import "context"

// FeedRepository defines the interface for feed data access
type FeedRepository interface {
	// Save persists a feed
	Save(ctx context.Context, feed Feed) error
	// ListAll returns all feeds
	ListAll(ctx context.Context) ([]Feed, error)
}

// ArticleRepository defines the interface for article data access
type ArticleRepository interface {
	// BulkUpsert inserts articles, ignoring duplicates by URL
	BulkUpsert(ctx context.Context, articles []Article) error
}
