package main

import (
	"log"
	"net"
	"os"
	"pft/main/internal/app/config"
	"pft/main/internal/app/services"
	"pft/main/internal/app/validations"
	"pft/main/tools"

	pb "pft/main/internal/app/generated_proto"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

func main() {
	runArgs := tools.ParseCommand()

	if runArgs {
		return
	}

	// Load config and run Server
	configFile := os.Getenv("CONFIG_FILE")
	configs := config.LoadConfig(configFile)

	// Register validation
	validate := validator.New()
	validations.RegisterCustomValidation(validate)

	// Create service
	transactionService, err := services.NewTransactionService(configs.DatabaseConnection, validate)
	if err != nil {
		log.Fatalf("Failed to create transaction service: %v", err)
	}

	categoryService, err := services.NewCategoryService(configs.DatabaseConnection, validate)
	if err != nil {
		log.Fatalf("Failed to create category service: %v", err)
	}

	attachmentService, err := services.NewAttachmentService(configs.DatabaseConnection, validate)
	if err != nil {
		log.Fatalf("Failed to create attachment service: %v", err)
	}

	accountService, err := services.NewAccountService(configs.DatabaseConnection, validate)
	if err != nil {
		log.Fatalf("Failed to create account service: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTransactionServiceServer(s, transactionService)
	pb.RegisterCategoryServiceServer(s, categoryService)
	pb.RegisterAttachmentServiceServer(s, attachmentService)
	pb.RegisterAccountServiceServer(s, accountService)

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
