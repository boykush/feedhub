package provider

import (
	"github.com/samber/do/v2"

	"github.com/boykush/feedhub/server/feed/internal/domain/repository"
	infrarepo "github.com/boykush/feedhub/server/feed/internal/infra/repository"
)

// ProvideFeedRepository creates a new FeedRepository implementation.
func ProvideFeedRepository(i do.Injector) (repository.FeedRepository, error) {
	ec := do.MustInvoke[*EntClient](i)
	return infrarepo.NewFeedRepository(ec.Client), nil
}
