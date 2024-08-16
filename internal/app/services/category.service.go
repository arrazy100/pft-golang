package services

import (
	"context"
	"errors"
	"fmt"
	pb "pft/main/internal/app/generated_proto"
	"pft/main/internal/app/models"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type CategoryService struct {
	db       *gorm.DB
	validate *validator.Validate
	pb.UnimplementedCategoryServiceServer
}

func NewCategoryService(db *gorm.DB, validate *validator.Validate) (*CategoryService, error) {
	return &CategoryService{db: db, validate: validate}, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	category := req.GetData()

	if category == nil {
		return nil, fmt.Errorf("no category data provided")
	}

	userId := "597c5e7d-63dc-4df1-9954-fefe8b415634"

	dbCategory := &models.Category{
		Name: category.Name,
	}

	dbCategory.SetId()
	dbCategory.SetAuditCreate(userId)
	dbCategory.SetUser(userId)

	// validation
	err := s.validate.Struct(dbCategory)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}

		return nil, err
	}

	// Save to database
	result := s.db.Create(dbCategory)
	if result.Error != nil {
		if pgError := result.Error.(*pgconn.PgError); errors.Is(result.Error, pgError) {
			if pgError.Code == "23505" && pgError.ConstraintName == "idx_category_name" {
				return nil, status.Errorf(codes.InvalidArgument, "Category must be unique")
			}
		}

		return nil, result.Error
	}

	// Convert back to protobuf Transaction
	createdCategory := &pb.Category{
		Id:        dbCategory.Id.String(),
		Name:      dbCategory.Name,
		UserId:    dbCategory.UserId,
		CreatedAt: dbCategory.CreatedAt.String(),
		CreatedBy: dbCategory.CreatedBy,
	}

	return &pb.CreateCategoryResponse{
		Message: "Category created successfully",
		Data:    createdCategory,
	}, nil
}
