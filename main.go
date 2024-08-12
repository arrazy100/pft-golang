package main

import (
	"log"
	"net"
	"os"
	"pft/main/internal/app/config"
	"pft/main/internal/app/services"

	pb "pft/main/internal/app/generated_proto"

	"google.golang.org/grpc"
)

func main() {
	// Load config and run Server
	configFile := os.Getenv("CONFIG_FILE")
	configs := config.LoadConfig(configFile)

	// Create service
	transactionService, err := services.NewTransactionService(configs.DatabaseConnection)
	if err != nil {
		log.Fatalf("Failed to create transaction service: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTransactionServiceServer(s, transactionService)
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
