package services

import (
	"context"
	"errors"
	"fmt"
	pb "pft/main/internal/app/generated_proto"
	"pft/main/internal/app/models"
	"pft/main/internal/app/utils"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	dbAccount := &models.Account{
		Type:    models.AccountType(account.Type),
		Name:    account.Name,
		Balance: account.Balance,
	}

	dbAccount.SetId()
	dbAccount.SetAuditCreate(account.UserId)
	dbAccount.SetUser(account.UserId)

	err := s.validate.Struct(dbAccount)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}

		return nil, err
	}

	result := s.db.Create(dbAccount)
	if result.Error != nil {
		return nil, result.Error
	}

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

func (s *AccountService) ListAccounts(ctx context.Context, req *pb.ListAccountRequest) (*pb.ListAccountResponse, error) {
	parsePagination, err := utils.ParsePaginationFilter(int(req.GetTake()), int(req.GetSkip()))
	if err != nil {
		return nil, err
	}
	take := parsePagination.Take
	skip := parsePagination.Skip
	userId := req.GetUserId()

	if userId == "" {
		return nil, fmt.Errorf("no user id provided")
	}

	var accounts []*models.Account
	result := s.db.
		Where("user_id = ?", userId).
		Offset(skip).Limit(take).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	s.db.Model(&models.Account{}).Where("user_id = ?", userId).Count(&total)

	// Convert back to protobuf Accounts
	var pbAccounts []*pb.Account
	for _, account := range accounts {
		pbAccount := &pb.Account{
			Id:        account.Id.String(),
			Type:      pb.AccountType(account.Type),
			Balance:   account.Balance,
			Name:      account.Name,
			UserId:    account.UserId,
			CreatedAt: account.CreatedAt.String(),
			CreatedBy: account.CreatedBy,
		}

		pbAccounts = append(pbAccounts, pbAccount)
	}

	return &pb.ListAccountResponse{Data: pbAccounts, Total: total, Take: int32(take), Skip: int32(skip)}, nil
}

func (s *AccountService) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	id, err := utils.ValidateUUIDFromString(req.GetId())
	if err != nil {
		return nil, err
	}

	userId, err := utils.ValidateUUIDFromString(req.GetUserId())
	if err != nil {
		return nil, err
	}

	dbAccount := models.Account{}
	result := s.db.Where("id = ? AND user_id = ?", id, userId).First(&dbAccount)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.Account{
		Id:        dbAccount.Id.String(),
		Name:      dbAccount.Name,
		Balance:   dbAccount.Balance,
		UserId:    dbAccount.UserId,
		CreatedAt: dbAccount.CreatedAt.String(),
		CreatedBy: dbAccount.CreatedBy,
	}, nil
}

func (s *AccountService) EditAccount(ctx context.Context, req *pb.Account) (*pb.Account, error) {
	id, err := utils.ValidateUUIDFromString(req.GetId())
	if err != nil {
		return nil, err
	}

	userId, err := utils.ValidateUUIDFromString(req.GetUserId())
	if err != nil {
		return nil, err
	}

	dbAccount := models.Account{}
	result := s.db.
		Where("id = ?", id).
		Where("user_id = ?", userId).
		First(&dbAccount)
	if result.Error != nil {
		return nil, result.Error
	}

	dbAccount.Name = req.GetName()
	dbAccount.Type = models.AccountType(req.GetType())

	dbAccount.SetAuditUpdate(*userId)

	result = s.db.Save(&dbAccount)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.Account{
		Id:        dbAccount.Id.String(),
		Name:      dbAccount.Name,
		Balance:   dbAccount.Balance,
		UserId:    dbAccount.UserId,
		CreatedAt: dbAccount.CreatedAt.String(),
		CreatedBy: dbAccount.CreatedBy,
	}, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	id, err := utils.ValidateUUIDFromString(req.GetId())
	if err != nil {
		return nil, err
	}

	userId, err := utils.ValidateUUIDFromString(req.GetUserId())
	if err != nil {
		return nil, err
	}

	dbAccount := models.Account{}
	result := s.db.Where("id = ? AND user_id = ?", id, userId).First(&dbAccount)
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.db.Delete(&dbAccount)
	if result.Error != nil {
		if pgError := result.Error.(*pgconn.PgError); errors.Is(result.Error, pgError) {
			if pgError.Code == "23503" && pgError.ConstraintName == "fk_transactions_account" {
				return nil, status.Errorf(codes.InvalidArgument, "Cannot delete account that has existing transactions")
			}
		}
		return nil, result.Error
	}

	return &pb.DeleteAccountResponse{Message: "Account deleted successfully"}, nil
}
