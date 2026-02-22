package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	collectorv1 "github.com/boykush/feedhub/server/collector/gen/go"
	"github.com/boykush/feedhub/server/collector/internal/infra/ent"
	infrarepo "github.com/boykush/feedhub/server/collector/internal/infra/repository"
	"github.com/boykush/feedhub/server/collector/internal/server"
	"github.com/boykush/feedhub/server/collector/internal/usecase"
	"github.com/caarlos0/env/v11"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := env.ParseAs[config]()
	if err != nil {
		log.Fatalf("failed to parse environment variables: %v", err)
	}

	client, err := ent.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	defer client.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	feedRepo := infrarepo.NewFeedRepository(client)
	addFeedUsecase := usecase.NewAddFeedUsecase(feedRepo)

	grpcServer := grpc.NewServer()
	collectorv1.RegisterCollectorServiceServer(grpcServer, server.NewServer(addFeedUsecase))
	reflection.Register(grpcServer)

	log.Printf("Starting Collector server on port %s", cfg.Port)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down Collector server...")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Collector server failed to serve: %v", err)
	}
}
