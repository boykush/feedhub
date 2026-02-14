package server

import (
	"context"

	feedv1 "github.com/boykush/foresee/services/feed/gen/go"
	"github.com/boykush/foresee/services/feed/internal/domain"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the FeedServiceServer interface
type Server struct {
	feedv1.UnimplementedFeedServiceServer
	feedRepo    domain.FeedRepository
	articleRepo domain.ArticleRepository
}

// NewServer creates a new instance of the feed service server
func NewServer(feedRepo domain.FeedRepository, articleRepo domain.ArticleRepository) *Server {
	return &Server{
		feedRepo:    feedRepo,
		articleRepo: articleRepo,
	}
}

// HealthCheck implements the health check endpoint
func (s *Server) HealthCheck(ctx context.Context, req *feedv1.HealthCheckRequest) (*feedv1.HealthCheckResponse, error) {
	return &feedv1.HealthCheckResponse{
		Status: feedv1.HealthCheckResponse_SERVING,
	}, nil
}

// ListFeeds returns all feeds
func (s *Server) ListFeeds(ctx context.Context, req *feedv1.ListFeedsRequest) (*feedv1.ListFeedsResponse, error) {
	feeds, err := s.feedRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	pbFeeds := make([]*feedv1.Feed, len(feeds))
	for i, f := range feeds {
		pbFeeds[i] = toProtoFeed(f)
	}

	return &feedv1.ListFeedsResponse{
		Feeds: pbFeeds,
	}, nil
}

// ListArticles returns articles, optionally filtered by feed ID
func (s *Server) ListArticles(ctx context.Context, req *feedv1.ListArticlesRequest) (*feedv1.ListArticlesResponse, error) {
	var feedID *uuid.UUID
	if req.FeedId != nil {
		parsed, err := uuid.Parse(*req.FeedId)
		if err != nil {
			return nil, err
		}
		feedID = &parsed
	}

	articles, err := s.articleRepo.List(ctx, feedID)
	if err != nil {
		return nil, err
	}

	pbArticles := make([]*feedv1.Article, len(articles))
	for i, a := range articles {
		pbArticles[i] = toProtoArticle(a)
	}

	return &feedv1.ListArticlesResponse{
		Articles: pbArticles,
	}, nil
}

func toProtoFeed(f domain.Feed) *feedv1.Feed {
	return &feedv1.Feed{
		Id:        f.ID().String(),
		Url:       f.URL(),
		Title:     f.Title(),
		CreatedAt: timestamppb.New(f.CreatedAt()),
	}
}

func toProtoArticle(a domain.Article) *feedv1.Article {
	return &feedv1.Article{
		Id:          a.ID().String(),
		FeedId:      a.FeedID().String(),
		Title:       a.Title(),
		Url:         a.URL(),
		PublishedAt: timestamppb.New(a.PublishedAt()),
	}
}
