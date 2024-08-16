package services

import (
	"context"
	"fmt"
	pb "pft/main/internal/app/generated_proto"
	"pft/main/internal/app/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AccountService struct {
	db       *gorm.DB
	validate *validator.Validate
	pb.UnimplementedAccountServiceServer
}

func NewAccountService(db *gorm.DB, validate *validator.Validate) (*AccountService, error) {
	return &AccountService{db: db, validate: validate}, nil
}

func (s *AccountService) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	account := req.GetData()

	if account == nil {
		return nil, fmt.Errorf("no account data provided")
	}

	userId := "597c5e7d-63dc-4df1-9954-fefe8b415634"

	dbAccount := &models.Account{
		Type:    models.AccountType(account.Type),
		Name:    account.Name,
		Balance: account.Balance,
	}

	dbAccount.SetId()
	dbAccount.SetAuditCreate(userId)
	dbAccount.SetUser(userId)

	// validation
	err := s.validate.Struct(dbAccount)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}

		return nil, err
	}

	// Save to database
	result := s.db.Create(dbAccount)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert back to protobuf Transaction
	createdAcount := &pb.Account{
		Id:        dbAccount.Id.String(),
		Type:      pb.AccountType(dbAccount.Type),
		Balance:   dbAccount.Balance,
		Name:      dbAccount.Name,
		UserId:    dbAccount.UserId,
		CreatedAt: dbAccount.CreatedAt.String(),
		CreatedBy: dbAccount.CreatedBy,
	}

	return &pb.CreateAccountResponse{
		Message: "Account created successfully",
		Data:    createdAcount,
	}, nil
}
