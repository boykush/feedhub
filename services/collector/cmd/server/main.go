package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	collectorv1 "github.com/boykush/foresee/services/collector/gen/go"
	"github.com/boykush/foresee/services/collector/internal/infra/db"
	"github.com/boykush/foresee/services/collector/internal/infra/repository"
	"github.com/boykush/foresee/services/collector/internal/infra/rss"
	"github.com/boykush/foresee/services/collector/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const defaultPort = "50053"

func main() {
	ctx := context.Background()

	port := os.Getenv("COLLECTOR_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize PostgreSQL connection pool
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	// Initialize repositories
	feedRepo := repository.NewFeedRepository(pool)
	articleRepo := repository.NewArticleRepository(pool)

	// Initialize RSS fetcher
	fetcher := rss.NewFetcher()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	collectorServer := server.NewServer(feedRepo, articleRepo, fetcher)

	collectorv1.RegisterCollectorServiceServer(grpcServer, collectorServer)

	// Enable reflection for grpcurl and other tools
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on port %s", port)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
