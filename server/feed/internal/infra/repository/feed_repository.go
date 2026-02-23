package repository

import (
	"context"

	"github.com/boykush/feedhub/server/feed/internal/domain/model"
	domainrepo "github.com/boykush/feedhub/server/feed/internal/domain/repository"
	"github.com/boykush/feedhub/server/feed/internal/infra/ent"
	"github.com/boykush/feedhub/server/feed/internal/infra/ent/feed"
)

type feedRepository struct {
	client *ent.Client
}

func NewFeedRepository(client *ent.Client) domainrepo.FeedRepository {
	return &feedRepository{client: client}
}

func (r *feedRepository) ExistsByURL(ctx context.Context, url string) (bool, error) {
	return r.client.Feed.Query().Where(feed.URLEQ(url)).Exist(ctx)
}

func (r *feedRepository) Save(ctx context.Context, url, title string) (*model.Feed, error) {
	created, err := r.client.Feed.Create().
		SetURL(url).
		SetTitle(title).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, domainrepo.ErrFeedAlreadyExists
		}
		return nil, err
	}
	return &model.Feed{
		ID:    created.ID,
		URL:   created.URL,
		Title: created.Title,
	}, nil
}
