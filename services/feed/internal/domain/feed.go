package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Feed represents an RSS/Atom feed source
type Feed struct {
	id        uuid.UUID
	url       string
	title     string
	createdAt time.Time
}

// NewFeed creates a new Feed with validation
func NewFeed(id uuid.UUID, url string, title string, createdAt time.Time) (Feed, error) {
	if id == uuid.Nil {
		return Feed{}, errors.New("id is required")
	}
	if url == "" {
		return Feed{}, errors.New("url is required")
	}
	if title == "" {
		return Feed{}, errors.New("title is required")
	}
	return Feed{
		id:        id,
		url:       url,
		title:     title,
		createdAt: createdAt,
	}, nil
}

// ID returns the feed ID
func (f Feed) ID() uuid.UUID {
	return f.id
}

// URL returns the feed URL
func (f Feed) URL() string {
	return f.url
}

// Title returns the feed title
func (f Feed) Title() string {
	return f.title
}

// CreatedAt returns when the feed was created
func (f Feed) CreatedAt() time.Time {
	return f.createdAt
}

// Article represents a single article from a feed
type Article struct {
	id          uuid.UUID
	feedID      uuid.UUID
	title       string
	url         string
	publishedAt time.Time
}

// NewArticle creates a new Article with validation
func NewArticle(id uuid.UUID, feedID uuid.UUID, title string, url string, publishedAt time.Time) (Article, error) {
	if id == uuid.Nil {
		return Article{}, errors.New("id is required")
	}
	if feedID == uuid.Nil {
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
func (a Article) ID() uuid.UUID {
	return a.id
}

// FeedID returns the feed ID this article belongs to
func (a Article) FeedID() uuid.UUID {
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

// PublishedAt returns when the article was published
func (a Article) PublishedAt() time.Time {
	return a.publishedAt
}
