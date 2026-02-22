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

	"github.com/caarlos0/env/v11"
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
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return fmt.Errorf("failed to parse environment variables: %w", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux()

	// Connect to backend services
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := feedpb.RegisterFeedServiceHandlerFromEndpoint(ctx, mux, cfg.FeedServiceAddr, opts); err != nil {
		return fmt.Errorf("failed to register feed service handler: %w", err)
	}

	if err := collectorpb.RegisterCollectorServiceHandlerFromEndpoint(ctx, mux, cfg.CollectorServiceAddr, opts); err != nil {
		return fmt.Errorf("failed to register collector service handler: %w", err)
	}

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting BFF server on port %s", cfg.HTTPPort)
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

	log.Println("BFF server stopped")
	return nil
}
