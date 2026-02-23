package provider

import (
	"github.com/samber/do/v2"

	"github.com/boykush/feedhub/server/feed/internal/domain/repository"
	"github.com/boykush/feedhub/server/feed/internal/usecase"
)

// ProvideAddFeedUsecase creates a new AddFeedUsecase.
func ProvideAddFeedUsecase(i do.Injector) (*usecase.AddFeedUsecase, error) {
	feedRepo := do.MustInvoke[repository.FeedRepository](i)
	return usecase.NewAddFeedUsecase(feedRepo), nil
}
