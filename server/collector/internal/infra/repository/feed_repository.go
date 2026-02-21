package repository

import (
	"context"

	"github.com/boykush/feedhub/server/collector/internal/domain/model"
	domainrepo "github.com/boykush/feedhub/server/collector/internal/domain/repository"
	"github.com/boykush/feedhub/server/collector/internal/infra/ent"
)

type feedRepository struct {
	client *ent.Client
}

func NewFeedRepository(client *ent.Client) domainrepo.FeedRepository {
	return &feedRepository{client: client}
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
