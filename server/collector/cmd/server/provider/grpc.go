package provider

import (
	"fmt"
	"log"
	"net"

	"github.com/samber/do/v2"

	collectorv1 "github.com/boykush/feedhub/server/collector/gen/go"
	"github.com/boykush/feedhub/server/collector/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer wraps *grpc.Server to implement do.Shutdownable.
type GRPCServer struct {
	*grpc.Server
}

func (s *GRPCServer) Shutdown() error {
	s.GracefulStop()
	return nil
}

// ProvideServer creates a new collector gRPC service server.
func ProvideServer(i do.Injector) (*server.Server, error) {
	return server.NewServer(), nil
}

// ProvideGRPCServer creates and starts a gRPC server with the collector service registered.
func ProvideGRPCServer(i do.Injector) (*GRPCServer, error) {
	cfg := do.MustInvoke[Config](i)
	srv := do.MustInvoke[*server.Server](i)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	collectorv1.RegisterCollectorServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	log.Printf("Starting Collector server on port %s", cfg.Port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Collector server failed to serve: %v", err)
		}
	}()

	return &GRPCServer{Server: grpcServer}, nil
}
