package domain

import (
	"errors"
	"time"
)

// Article represents an article from an RSS feed
type Article struct {
	id          string
	feedID      string
	title       string
	url         string
	publishedAt time.Time
}

// NewArticle creates a new Article with validation
func NewArticle(id, feedID, title, url string, publishedAt time.Time) (Article, error) {
	if id == "" {
		return Article{}, errors.New("id is required")
	}
	if feedID == "" {
		return Article{}, errors.New("feed_id is required")
	}
	if title == "" {
		return Article{}, errors.New("title is required")
	}
	if url == "" {
		return Article{}, errors.New("url is required")
	}
	return Article{
		id:          id,
		feedID:      feedID,
		title:       title,
		url:         url,
		publishedAt: publishedAt,
	}, nil
}

// ID returns the article ID
func (a Article) ID() string {
	return a.id
}

// FeedID returns the feed ID this article belongs to
func (a Article) FeedID() string {
	return a.feedID
}

// Title returns the article title
func (a Article) Title() string {
	return a.title
}

// URL returns the article URL
func (a Article) URL() string {
	return a.url
}

// PublishedAt returns the article publication time
func (a Article) PublishedAt() time.Time {
	return a.publishedAt
}
