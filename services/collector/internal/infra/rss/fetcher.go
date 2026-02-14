package rss

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"

	"github.com/boykush/foresee/services/collector/internal/domain"
)

// Fetcher fetches and parses RSS feeds
type Fetcher struct {
	parser *gofeed.Parser
}

// NewFetcher creates a new RSS Fetcher
func NewFetcher() *Fetcher {
	return &Fetcher{
		parser: gofeed.NewParser(),
	}
}

// FetchResult contains the parsed RSS feed data
type FetchResult struct {
	Title    string
	Articles []domain.Article
}

// Fetch fetches an RSS feed from URL and returns parsed articles
func (f *Fetcher) Fetch(ctx context.Context, feedURL string, feedID string) (*FetchResult, error) {
	feed, err := f.parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	articles := make([]domain.Article, 0, len(feed.Items))
	for _, item := range feed.Items {
		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}

		title := item.Title
		if title == "" {
			continue
		}

		link := item.Link
		if link == "" {
			continue
		}

		article, err := domain.NewArticle(uuid.New().String(), feedID, title, link, publishedAt)
		if err != nil {
			continue
		}
		articles = append(articles, article)
	}

	return &FetchResult{
		Title:    feed.Title,
		Articles: articles,
	}, nil
}
