package services

import (
	"context"
	"time"

	pb "pft/main/internal/app/generated_proto"
	"pft/main/internal/app/models"
	base_models "pft/main/internal/app/models/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	db *gorm.DB
	pb.UnimplementedTransactionServiceServer
}

func NewTransactionService(db *gorm.DB) (*TransactionService, error) {
	return &TransactionService{db: db}, nil
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	transaction := req.GetTransaction()

	id := uuid.New().String()
	createdAt := time.Now()

	dbTransaction := &models.Transaction{
		Id:           id,
		UserId:       transaction.UserId,
		CategoryId:   transaction.CategoryId,
		Description:  transaction.Description,
		AccountId:    transaction.AccountId,
		AttachmentId: transaction.AttachmentId,
		Amount:       transaction.Amount,
		Type:         models.TransactionType(transaction.Type),
		BaseAudit:    base_models.BaseAudit{CreatedAt: createdAt},
	}

	// Save to database
	result := s.db.Create(dbTransaction)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert back to protobuf Transaction
	createdTransaction := &pb.Transaction{
		Id:           dbTransaction.Id,
		UserId:       dbTransaction.UserId,
		CategoryId:   dbTransaction.CategoryId,
		Description:  dbTransaction.Description,
		AccountId:    dbTransaction.AccountId,
		AttachmentId: dbTransaction.AttachmentId,
		Amount:       dbTransaction.Amount,
		Type:         pb.TransactionType(dbTransaction.Type),
	}

	return &pb.CreateTransactionResponse{
		Message:     "Transaction created successfully",
		Transaction: createdTransaction,
	}, nil
}
