package services

import (
	"context"
	"fmt"
	"time"

	pb "pft/main/internal/app/generated_proto"
	"pft/main/internal/app/models"
	"pft/main/internal/app/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	db       *gorm.DB
	validate *validator.Validate
	pb.UnimplementedTransactionServiceServer
}

func NewTransactionService(db *gorm.DB, validate *validator.Validate) (*TransactionService, error) {
	return &TransactionService{db: db, validate: validate}, nil
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	transaction := req.GetData()

	if transaction == nil {
		return nil, fmt.Errorf("no transaction data provided")
	}

	userId := "597c5e7d-63dc-4df1-9954-fefe8b415634"

	transactionDate, err := time.Parse(time.RFC3339, transaction.TransactionDate)
	if err != nil {
		return nil, err
	}

	dbTransaction := &models.Transaction{
		Description:     transaction.Description,
		Amount:          transaction.Amount,
		Type:            models.TransactionType(transaction.Type),
		TransactionDate: transactionDate,
		CategoryId:      uuid.MustParse(transaction.CategoryId),
		AccountId:       uuid.MustParse(transaction.AccountId),
		AttachmentId:    uuid.MustParse(transaction.AttachmentId),
	}

	dbTransaction.SetId()
	dbTransaction.SetAuditCreate(userId)
	dbTransaction.SetUser(userId)

	// validation
	err = s.validate.Struct(dbTransaction)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}

		return nil, err
	}

	// Save to database
	result := s.db.Create(dbTransaction)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert back to protobuf Transaction
	createdTransaction := &pb.Transaction{
		Id:              dbTransaction.Id.String(),
		Description:     dbTransaction.Description,
		Amount:          dbTransaction.Amount,
		Type:            pb.TransactionType(dbTransaction.Type),
		TransactionDate: dbTransaction.TransactionDate.String(),
		CategoryId:      dbTransaction.CategoryId.String(),
		AccountId:       dbTransaction.AccountId.String(),
		AttachmentId:    dbTransaction.AttachmentId.String(),
		UserId:          dbTransaction.UserId,
		CreatedAt:       dbTransaction.CreatedAt.String(),
		CreatedBy:       dbTransaction.CreatedBy,
	}

	return &pb.CreateTransactionResponse{
		Message: "Transaction created successfully",
		Data:    createdTransaction,
	}, nil
}

func (s *TransactionService) ListTransaction(ctx context.Context, req *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	take := int(req.GetTake())
	skip := int(req.GetSkip())
	startDate := req.GetStartDate()
	endDate := req.GetEndDate()
	userId := req.GetUserId()
	categoryId := req.GetCategoryId()
	timezone := 420

	parsedRequest, err := utils.ParseRequestDateTimeFilter(take, skip, startDate, endDate, int(timezone))

	if err != nil {
		return nil, err
	}

	dbTransaction := []models.Transaction{}
	query := s.db.
		Where("user_id = ?", userId).
		Where("transaction_date BETWEEN ? AND ?", parsedRequest.StartDate, parsedRequest.EndDate)

	category, err := uuid.Parse(categoryId)
	if err == nil {
		query = query.Where("category_id = ?", category)
	}

	result := query.
		Preload("Category").
		Preload("Account").
		Preload("Attachment").
		Limit(parsedRequest.Take).
		Offset(parsedRequest.Skip).
		Find(&dbTransaction)

	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	result = s.db.
		Model(&models.Transaction{}).
		Where("user_id = ?", userId).
		Count(&total)

	if result.Error != nil {
		return nil, result.Error
	}

	var transactions []*pb.Transaction
	for _, dbTransaction := range dbTransaction {
		transaction := &pb.Transaction{
			Id:              dbTransaction.Id.String(),
			Description:     dbTransaction.Description,
			Amount:          dbTransaction.Amount,
			Type:            pb.TransactionType(dbTransaction.Type),
			TransactionDate: dbTransaction.TransactionDate.String(),
			CategoryId:      dbTransaction.CategoryId.String(),
			AccountId:       dbTransaction.AccountId.String(),
			AttachmentId:    dbTransaction.AttachmentId.String(),
			UserId:          dbTransaction.UserId,
			CreatedAt:       dbTransaction.CreatedAt.String(),
			CreatedBy:       dbTransaction.CreatedBy,
			Category:        &pb.CategoryMini{Id: dbTransaction.Category.Id.String(), Name: dbTransaction.Category.Name},
			Account:         &pb.AccountMini{Id: dbTransaction.Account.Id.String(), Name: dbTransaction.Account.Name, Balance: dbTransaction.Account.Balance},
			Attachment:      &pb.AttachmentMini{Id: dbTransaction.Attachment.Id.String(), ContentUrl: dbTransaction.Attachment.ContentUrl, Type: pb.AttachmentType(dbTransaction.Attachment.Type)},
		}

		transactions = append(transactions, transaction)
	}

	return &pb.ListTransactionResponse{
		Data:      transactions,
		Total:     total,
		Take:      int32(parsedRequest.Take),
		Skip:      int32(parsedRequest.Skip),
		StartDate: utils.TimeFormatAsDate(parsedRequest.StartDate, timezone),
		EndDate:   utils.TimeFormatAsDate(parsedRequest.EndDate, timezone),
	}, nil
}
