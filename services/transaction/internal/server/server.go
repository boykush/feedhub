package server

import (
	"context"

	transactionv1 "github.com/boykush/foresee/services/transaction/gen/go"
	v1 "github.com/boykush/foresee/services/transaction/gen/go/grpc/health/v1"
)

// Server implements the TransactionServiceServer interface
type Server struct {
	transactionv1.UnimplementedTransactionServiceServer
}

// NewServer creates a new instance of the transaction service server
func NewServer() *Server {
	return &Server{}
}

// HealthCheck implements the health check endpoint
func (s *Server) HealthCheck(ctx context.Context, req *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	return &v1.HealthCheckResponse{
		Status: v1.HealthCheckResponse_SERVING,
	}, nil
}
