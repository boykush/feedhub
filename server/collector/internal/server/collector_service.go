package server

import (
	"context"

	collectorv1 "github.com/boykush/feedhub/server/collector/gen/go"
	"github.com/boykush/feedhub/server/collector/internal/infra/ent"
)

type Server struct {
	collectorv1.UnimplementedCollectorServiceServer
	client *ent.Client
}

func NewServer(client *ent.Client) *Server {
	return &Server{client: client}
}

func (s *Server) HealthCheck(ctx context.Context, req *collectorv1.HealthCheckRequest) (*collectorv1.HealthCheckResponse, error) {
	return &collectorv1.HealthCheckResponse{
		Status: collectorv1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Server) AddFeed(ctx context.Context, req *collectorv1.AddFeedRequest) (*collectorv1.AddFeedResponse, error) {
	return &collectorv1.AddFeedResponse{}, nil
}

func (s *Server) SyncFeeds(ctx context.Context, req *collectorv1.SyncFeedsRequest) (*collectorv1.SyncFeedsResponse, error) {
	return &collectorv1.SyncFeedsResponse{}, nil
}
