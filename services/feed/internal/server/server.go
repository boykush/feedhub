package server

import (
	"context"

	feedv1 "github.com/boykush/foresee/services/feed/gen/go"
)

type Server struct {
	feedv1.UnimplementedFeedServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) HealthCheck(ctx context.Context, req *feedv1.HealthCheckRequest) (*feedv1.HealthCheckResponse, error) {
	return &feedv1.HealthCheckResponse{
		Status: feedv1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Server) ListFeeds(ctx context.Context, req *feedv1.ListFeedsRequest) (*feedv1.ListFeedsResponse, error) {
	return &feedv1.ListFeedsResponse{}, nil
}

func (s *Server) ListArticles(ctx context.Context, req *feedv1.ListArticlesRequest) (*feedv1.ListArticlesResponse, error) {
	return &feedv1.ListArticlesResponse{}, nil
}
