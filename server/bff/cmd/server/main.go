package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	collectorpb "github.com/boykush/feedhub/server/bff/gen/go/collector"
	feedpb "github.com/boykush/feedhub/server/bff/gen/go/feed"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Get configuration from environment variables
	httpPort := getEnv("BFF_HTTP_PORT", "8080")
	feedServiceAddr := getEnv("FEED_SERVICE_ADDR", "feed-service:50052")
	collectorServiceAddr := getEnv("COLLECTOR_SERVICE_ADDR", "collector-service:50053")

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux()

	// Connect to backend services
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := feedpb.RegisterFeedServiceHandlerFromEndpoint(ctx, mux, feedServiceAddr, opts); err != nil {
		return fmt.Errorf("failed to register feed service handler: %w", err)
	}

	if err := collectorpb.RegisterCollectorServiceHandlerFromEndpoint(ctx, mux, collectorServiceAddr, opts); err != nil {
		return fmt.Errorf("failed to register collector service handler: %w", err)
	}

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting BFF server on port %s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigCh
	log.Println("Shutting down BFF server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Server stopped")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
