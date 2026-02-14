package domain

import (
	"errors"
	"time"
)

// Feed represents an RSS feed
type Feed struct {
	id        string
	url       string
	title     string
	createdAt time.Time
}

// NewFeed creates a new Feed with validation
func NewFeed(id, url, title string, createdAt time.Time) (Feed, error) {
	if id == "" {
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
func (f Feed) ID() string {
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

// CreatedAt returns the feed creation time
func (f Feed) CreatedAt() time.Time {
	return f.createdAt
}
