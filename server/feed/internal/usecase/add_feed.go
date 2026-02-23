package usecase

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"

	"github.com/boykush/feedhub/server/feed/internal/domain/model"
	"github.com/boykush/feedhub/server/feed/internal/domain/repository"
)

type AddFeedUsecase struct {
	feedRepo repository.FeedRepository
}

func NewAddFeedUsecase(feedRepo repository.FeedRepository) *AddFeedUsecase {
	return &AddFeedUsecase{feedRepo: feedRepo}
}

func (u *AddFeedUsecase) Execute(ctx context.Context, feedURL string) (*model.Feed, error) {
	exists, err := u.feedRepo.ExistsByURL(ctx, feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to check feed existence: %w", err)
	}
	if exists {
		return nil, repository.ErrFeedAlreadyExists
	}

	parser := gofeed.NewParser()
	parsed, err := parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	return u.feedRepo.Save(ctx, feedURL, parsed.Title)
}
