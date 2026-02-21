package server

import (
	"context"
	"errors"

	collectorv1 "github.com/boykush/feedhub/server/collector/gen/go"
	"github.com/boykush/feedhub/server/collector/internal/domain/repository"
	"github.com/boykush/feedhub/server/collector/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	collectorv1.UnimplementedCollectorServiceServer
	addFeedUsecase *usecase.AddFeedUsecase
}

func NewServer(addFeedUsecase *usecase.AddFeedUsecase) *Server {
	return &Server{addFeedUsecase: addFeedUsecase}
}

func (s *Server) HealthCheck(ctx context.Context, req *collectorv1.HealthCheckRequest) (*collectorv1.HealthCheckResponse, error) {
	return &collectorv1.HealthCheckResponse{
		Status: collectorv1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Server) AddFeed(ctx context.Context, req *collectorv1.AddFeedRequest) (*collectorv1.AddFeedResponse, error) {
	feed, err := s.addFeedUsecase.Execute(ctx, req.Url)
	if err != nil {
		if errors.Is(err, repository.ErrFeedAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "feed with this URL already exists")
		}
		return nil, status.Errorf(codes.InvalidArgument, "failed to add feed: %v", err)
	}
	return &collectorv1.AddFeedResponse{
		FeedId: feed.ID.String(),
		Title:  feed.Title,
	}, nil
}

func (s *Server) SyncFeeds(ctx context.Context, req *collectorv1.SyncFeedsRequest) (*collectorv1.SyncFeedsResponse, error) {
	return &collectorv1.SyncFeedsResponse{}, nil
}
