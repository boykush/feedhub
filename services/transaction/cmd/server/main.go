package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	transactionv1 "github.com/boykush/foresee/services/transaction/gen/go"
	"github.com/boykush/foresee/services/transaction/internal/infra/repository"
	"github.com/boykush/foresee/services/transaction/internal/infra/storage"
	"github.com/boykush/foresee/services/transaction/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultPort   = "50051"
	defaultCSVKey = "transactions.csv"
)

func main() {
	ctx := context.Background()

	port := os.Getenv("TRANSACTION_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	csvKey := os.Getenv("STORAGE_CSV_KEY")
	if csvKey == "" {
		csvKey = defaultCSVKey
	}

	// Initialize storage client
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create storage client: %v", err)
	}

	// Initialize repository
	repo := repository.NewTransactionRepository(storageClient, csvKey)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	transactionServer := server.NewServer(repo)

	transactionv1.RegisterTransactionServiceServer(grpcServer, transactionServer)

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
