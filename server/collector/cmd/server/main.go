package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	collectorv1 "github.com/boykush/foresee/server/collector/gen/go"
	"github.com/boykush/foresee/server/collector/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const defaultPort = "50053"

func main() {
	port := os.Getenv("COLLECTOR_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	collectorv1.RegisterCollectorServiceServer(grpcServer, server.NewServer())
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on port %s", port)

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
