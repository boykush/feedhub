package repository

import (
	"context"
	"errors"

	"github.com/boykush/feedhub/server/feed/internal/domain/model"
)

var ErrFeedAlreadyExists = errors.New("feed already exists")

type FeedRepository interface {
	ExistsByURL(ctx context.Context, url string) (bool, error)
	Save(ctx context.Context, url, title string) (*model.Feed, error)
}
