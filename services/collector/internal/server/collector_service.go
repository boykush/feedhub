package server

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	collectorv1 "github.com/boykush/foresee/services/collector/gen/go"
	"github.com/boykush/foresee/services/collector/internal/domain"
	"github.com/boykush/foresee/services/collector/internal/infra/rss"
)

// Server implements the CollectorServiceServer interface
type Server struct {
	collectorv1.UnimplementedCollectorServiceServer
	feedRepo    domain.FeedRepository
	articleRepo domain.ArticleRepository
	fetcher     *rss.Fetcher
}

// NewServer creates a new instance of the collector service server
func NewServer(feedRepo domain.FeedRepository, articleRepo domain.ArticleRepository, fetcher *rss.Fetcher) *Server {
	return &Server{
		feedRepo:    feedRepo,
		articleRepo: articleRepo,
		fetcher:     fetcher,
	}
}

// HealthCheck implements the health check endpoint
func (s *Server) HealthCheck(ctx context.Context, req *collectorv1.HealthCheckRequest) (*collectorv1.HealthCheckResponse, error) {
	return &collectorv1.HealthCheckResponse{
		Status: collectorv1.HealthCheckResponse_SERVING,
	}, nil
}

// AddFeed adds a new RSS feed and fetches its articles
func (s *Server) AddFeed(ctx context.Context, req *collectorv1.AddFeedRequest) (*collectorv1.AddFeedResponse, error) {
	feedID := uuid.New().String()

	// Fetch RSS to get title and articles
	result, err := s.fetcher.Fetch(ctx, req.GetUrl(), feedID)
	if err != nil {
		return nil, err
	}

	// Save feed to DB
	feed, err := domain.NewFeed(feedID, req.GetUrl(), result.Title, time.Now())
	if err != nil {
		return nil, err
	}

	if err := s.feedRepo.Save(ctx, feed); err != nil {
		return nil, err
	}

	// Save articles to DB
	if err := s.articleRepo.BulkUpsert(ctx, result.Articles); err != nil {
		return nil, err
	}

	return &collectorv1.AddFeedResponse{
		FeedId: feedID,
		Title:  result.Title,
	}, nil
}

// SyncFeeds syncs all registered feeds and fetches new articles
func (s *Server) SyncFeeds(ctx context.Context, req *collectorv1.SyncFeedsRequest) (*collectorv1.SyncFeedsResponse, error) {
	feeds, err := s.feedRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var syncedCount int32
	for _, feed := range feeds {
		result, err := s.fetcher.Fetch(ctx, feed.URL(), feed.ID())
		if err != nil {
			log.Printf("failed to fetch feed %s (%s): %v", feed.ID(), feed.URL(), err)
			continue
		}

		if err := s.articleRepo.BulkUpsert(ctx, result.Articles); err != nil {
			log.Printf("failed to upsert articles for feed %s: %v", feed.ID(), err)
			continue
		}

		syncedCount++
	}

	return &collectorv1.SyncFeedsResponse{
		SyncedCount: syncedCount,
	}, nil
}
