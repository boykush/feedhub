package usecase

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"

	"github.com/boykush/feedhub/server/collector/internal/domain/model"
	"github.com/boykush/feedhub/server/collector/internal/domain/repository"
)

type AddFeedUsecase struct {
	feedRepo repository.FeedRepository
}

func NewAddFeedUsecase(feedRepo repository.FeedRepository) *AddFeedUsecase {
	return &AddFeedUsecase{feedRepo: feedRepo}
}

func (u *AddFeedUsecase) Execute(ctx context.Context, feedURL string) (*model.Feed, error) {
	parser := gofeed.NewParser()
	parsed, err := parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	return u.feedRepo.Save(ctx, feedURL, parsed.Title)
}
