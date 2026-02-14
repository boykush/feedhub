package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	feedv1 "github.com/boykush/foresee/services/feed/gen/go"
	"github.com/boykush/foresee/services/feed/internal/infra/db"
	"github.com/boykush/foresee/services/feed/internal/infra/repository"
	"github.com/boykush/foresee/services/feed/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const defaultPort = "50052"

func main() {
	ctx := context.Background()

	port := os.Getenv("FEED_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Initialize database connection pool
	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	// Initialize repositories
	feedRepo := repository.NewFeedRepository(pool)
	articleRepo := repository.NewArticleRepository(pool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	feedServer := server.NewServer(feedRepo, articleRepo)

	feedv1.RegisterFeedServiceServer(grpcServer, feedServer)

	// Enable reflection for grpcurl and other tools
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on port %s", port)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down gRPC server...")
		pool.Close()
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
