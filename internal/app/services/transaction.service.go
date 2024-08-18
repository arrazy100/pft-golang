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

func (s *TransactionService) SetAttachmentToTransaction(attachmentId string, transactionId string) error {
	_, err := utils.ValidateUUIDFromString(attachmentId)
	if err != nil {
		return err
	}

	_, err = utils.ValidateUUIDFromString(transactionId)
	if err != nil {
		return err
	}

	dbAttachment := &models.Attachment{}
	result := s.db.Where("id = ?", attachmentId).First(&dbAttachment)
	if result.Error != nil {
		return result.Error
	}

	if dbAttachment.TransactionId != uuid.Nil {
		return fmt.Errorf("attachment already has a transaction associated with it")
	}

	dbAttachment.TransactionId = uuid.MustParse(transactionId)

	result = s.db.Save(&dbAttachment)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	transaction := req.GetData()

	if transaction == nil {
		return nil, fmt.Errorf("no transaction data provided")
	}

	userId := transaction.UserId

	transactionDate, err := time.Parse(time.RFC3339, transaction.TransactionDate)
	if err != nil {
		return nil, err
	}

	transactionDate = transactionDate.Add(utils.GetTimezoneOffset(int(transaction.Timezone)))

	dbTransaction := &models.Transaction{
		Description:     transaction.Description,
		Amount:          transaction.Amount,
		Type:            models.TransactionType(transaction.Type),
		TransactionDate: transactionDate,
		CategoryId:      uuid.MustParse(transaction.CategoryId),
		AccountId:       uuid.MustParse(transaction.AccountId),
	}

	dbTransaction.SetAttachment(transaction.AttachmentId)
	dbTransaction.SetId()
	dbTransaction.SetAuditCreate(userId)
	dbTransaction.SetUser(userId)

	err = s.validate.Struct(dbTransaction)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}

		return nil, err
	}

	result := s.db.Create(dbTransaction)
	if result.Error != nil {
		return nil, result.Error
	}

	if transaction.AttachmentId != "" && dbTransaction.AttachmentId != uuid.Nil {
		err := s.SetAttachmentToTransaction(transaction.AttachmentId, dbTransaction.Id.String())
		if err != nil {
			return nil, err
		}

		dbTransaction.AttachmentId = uuid.MustParse(transaction.AttachmentId)
	}

	createdTransaction := &pb.Transaction{
		Id:              dbTransaction.Id.String(),
		Description:     dbTransaction.Description,
		Amount:          dbTransaction.Amount,
		Type:            pb.TransactionType(dbTransaction.Type),
		TransactionDate: dbTransaction.TransactionDate.String(),
		CategoryId:      dbTransaction.CategoryId.String(),
		AccountId:       dbTransaction.AccountId.String(),
		UserId:          dbTransaction.UserId,
		CreatedAt:       dbTransaction.CreatedAt.String(),
		CreatedBy:       dbTransaction.CreatedBy,
	}

	if dbTransaction.AttachmentId != uuid.Nil {
		createdTransaction.AttachmentId = dbTransaction.AttachmentId.String()
	}

	return &pb.CreateTransactionResponse{
		Message: "Transaction created successfully",
		Data:    createdTransaction,
	}, nil
}

func (s *TransactionService) ListTransactions(ctx context.Context, req *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	take := int(req.GetTake())
	skip := int(req.GetSkip())
	startDate := req.GetStartDate()
	endDate := req.GetEndDate()
	categoryId := req.GetCategoryId()
	timezone := req.GetTimezone()

	userId, err := utils.ValidateUUIDFromString(req.GetUserId())
	if err != nil {
		return nil, err
	}

	parsedRequest, err := utils.ParseRequestDateTimeFilter(take, skip, startDate, endDate, int(timezone))
	if err != nil {
		return nil, err
	}

	dbTransaction := []models.Transaction{}
	query := s.db.
		Where("user_id = ?", userId).
		Where("transaction_date BETWEEN ? AND ?", parsedRequest.StartDate, parsedRequest.EndDate)

	var total int64
	queryTotal := s.db.
		Model(&models.Transaction{}).
		Where("user_id = ?", userId).
		Where("transaction_date BETWEEN ? AND ?", parsedRequest.StartDate, parsedRequest.EndDate)

	category, err := uuid.Parse(categoryId)
	if err == nil {
		query = query.Where("category_id = ?", category)
		queryTotal = queryTotal.Where("category_id = ?", category)
	}

	result := query.
		Preload("Category").
		Preload("Account").
		Preload("Attachment").
		Limit(parsedRequest.Take).Offset(parsedRequest.Skip).Find(&dbTransaction)
	if result.Error != nil {
		return nil, result.Error
	}

	result = queryTotal.Count(&total)
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
			UserId:          dbTransaction.UserId,
			CreatedAt:       dbTransaction.CreatedAt.String(),
			CreatedBy:       dbTransaction.CreatedBy,
			Category:        &pb.CategoryMini{Id: dbTransaction.Category.Id.String(), Name: dbTransaction.Category.Name},
			Account:         &pb.AccountMini{Id: dbTransaction.Account.Id.String(), Name: dbTransaction.Account.Name, Balance: dbTransaction.Account.Balance},
		}

		if dbTransaction.AttachmentId != uuid.Nil {
			transaction.AttachmentId = dbTransaction.AttachmentId.String()
			transaction.Attachment = &pb.AttachmentMini{Id: dbTransaction.Attachment.Id.String(), ContentUrl: dbTransaction.Attachment.ContentUrl, Type: pb.AttachmentType(dbTransaction.Attachment.Type)}
		}

		transactions = append(transactions, transaction)
	}

	return &pb.ListTransactionResponse{
		Data:      transactions,
		Total:     total,
		Take:      int32(parsedRequest.Take),
		Skip:      int32(parsedRequest.Skip),
		StartDate: utils.TimeFormatAsDate(parsedRequest.StartDate, parsedRequest.Timezone),
		EndDate:   utils.TimeFormatAsDate(parsedRequest.EndDate, parsedRequest.Timezone),
	}, nil
}
