package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/samber/do/v2"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	collectorpb "github.com/boykush/feedhub/server/bff/gen/go/collector"
	feedpb "github.com/boykush/feedhub/server/bff/gen/go/feed"
)

// ProvideHTTPServer creates and starts an HTTP server with gRPC-Gateway handlers.
func ProvideHTTPServer(i do.Injector) (*http.Server, error) {
	cfg := do.MustInvoke[Config](i)

	ctx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := feedpb.RegisterFeedServiceHandlerFromEndpoint(ctx, mux, cfg.FeedServiceAddr, opts); err != nil {
		return nil, fmt.Errorf("failed to register feed service handler: %w", err)
	}

	if err := collectorpb.RegisterCollectorServiceHandlerFromEndpoint(ctx, mux, cfg.CollectorServiceAddr, opts); err != nil {
		return nil, fmt.Errorf("failed to register collector service handler: %w", err)
	}

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

	return httpServer, nil
}
