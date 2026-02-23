package server

import (
	"context"
	"errors"

	feedv1 "github.com/boykush/feedhub/server/feed/gen/go"
	"github.com/boykush/feedhub/server/feed/internal/domain/repository"
	"github.com/boykush/feedhub/server/feed/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	feedv1.UnimplementedFeedServiceServer
	addFeedUsecase *usecase.AddFeedUsecase
}

func NewServer(addFeedUsecase *usecase.AddFeedUsecase) *Server {
	return &Server{addFeedUsecase: addFeedUsecase}
}

func (s *Server) HealthCheck(ctx context.Context, req *feedv1.HealthCheckRequest) (*feedv1.HealthCheckResponse, error) {
	return &feedv1.HealthCheckResponse{
		Status: feedv1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Server) AddFeed(ctx context.Context, req *feedv1.AddFeedRequest) (*feedv1.AddFeedResponse, error) {
	feed, err := s.addFeedUsecase.Execute(ctx, req.Url)
	if err != nil {
		if errors.Is(err, repository.ErrFeedAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "feed with this URL already exists")
		}
		return nil, status.Errorf(codes.InvalidArgument, "failed to add feed: %v", err)
	}
	return &feedv1.AddFeedResponse{
		FeedId: feed.ID.String(),
		Title:  feed.Title,
	}, nil
}

func (s *Server) ListFeeds(ctx context.Context, req *feedv1.ListFeedsRequest) (*feedv1.ListFeedsResponse, error) {
	return &feedv1.ListFeedsResponse{}, nil
}

func (s *Server) ListArticles(ctx context.Context, req *feedv1.ListArticlesRequest) (*feedv1.ListArticlesResponse, error) {
	return &feedv1.ListArticlesResponse{}, nil
}
