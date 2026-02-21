package repository

import (
	"context"
	"errors"

	"github.com/boykush/feedhub/server/collector/internal/domain/model"
)

var ErrFeedAlreadyExists = errors.New("feed already exists")

type FeedRepository interface {
	Save(ctx context.Context, url, title string) (*model.Feed, error)
}
