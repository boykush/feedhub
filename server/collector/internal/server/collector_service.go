package server

import (
	"context"

	collectorv1 "github.com/boykush/feedhub/server/collector/gen/go"
)

type Server struct {
	collectorv1.UnimplementedCollectorServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) HealthCheck(ctx context.Context, req *collectorv1.HealthCheckRequest) (*collectorv1.HealthCheckResponse, error) {
	return &collectorv1.HealthCheckResponse{
		Status: collectorv1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Server) SyncFeeds(ctx context.Context, req *collectorv1.SyncFeedsRequest) (*collectorv1.SyncFeedsResponse, error) {
	return &collectorv1.SyncFeedsResponse{}, nil
}
