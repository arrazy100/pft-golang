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

	err := s.validate.Struct(dbCategory)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}

		return nil, err
	}

	result := s.db.Create(dbCategory)
	if result.Error != nil {
		if pgError := result.Error.(*pgconn.PgError); errors.Is(result.Error, pgError) {
			if pgError.Code == "23505" && pgError.ConstraintName == "idx_category_name" {
				return nil, status.Errorf(codes.InvalidArgument, "Category must be unique")
			}
		}

		return nil, result.Error
	}

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

func (s *CategoryService) ListCategories(ctx context.Context, req *pb.ListCategoryRequest) (*pb.ListCategoryResponse, error) {
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

	var categories []*models.Category
	result := s.db.
		Where("user_id = ?", userId).
		Offset(skip).Limit(take).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	s.db.Model(&models.Category{}).Where("user_id = ?", userId).Count(&total)

	// Convert back to protobuf Accounts
	var pbCategories []*pb.Category
	for _, category := range categories {
		pbCategory := &pb.Category{
			Id:        category.Id.String(),
			Name:      category.Name,
			UserId:    category.UserId,
			CreatedBy: category.CreatedBy,
			CreatedAt: category.CreatedAt.String(),
		}

		pbCategories = append(pbCategories, pbCategory)
	}

	return &pb.ListCategoryResponse{Data: pbCategories, Total: total, Take: int32(take), Skip: int32(skip)}, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	id, err := utils.ValidateUUIDFromString(req.GetId())
	if err != nil {
		return nil, err
	}

	userId, err := utils.ValidateUUIDFromString(req.GetUserId())
	if err != nil {
		return nil, err
	}

	dbCategory := models.Category{}
	result := s.db.Where("id = ? AND user_id = ?", id, userId).First(&dbCategory)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.Category{
		Id:        dbCategory.Id.String(),
		Name:      dbCategory.Name,
		UserId:    dbCategory.UserId,
		CreatedBy: dbCategory.CreatedBy,
		CreatedAt: dbCategory.CreatedAt.String(),
	}, nil
}

func (s *CategoryService) EditCategory(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	id, err := utils.ValidateUUIDFromString(req.GetId())
	if err != nil {
		return nil, err
	}

	userId, err := utils.ValidateUUIDFromString(req.GetUserId())
	if err != nil {
		return nil, err
	}

	dbCategory := models.Category{}
	result := s.db.
		Where("id = ?", id).
		Where("user_id = ?", userId).
		First(&dbCategory)
	if result.Error != nil {
		return nil, result.Error
	}

	dbCategory.Name = req.GetName()

	dbCategory.SetAuditUpdate(*userId)

	result = s.db.Save(&dbCategory)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.Category{
		Id:        dbCategory.Id.String(),
		Name:      dbCategory.Name,
		UserId:    dbCategory.UserId,
		CreatedAt: dbCategory.CreatedAt.String(),
		CreatedBy: dbCategory.CreatedBy,
	}, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.DeleteCategoryResponse, error) {
	id, err := utils.ValidateUUIDFromString(req.GetId())
	if err != nil {
		return nil, err
	}

	userId, err := utils.ValidateUUIDFromString(req.GetUserId())
	if err != nil {
		return nil, err
	}

	dbCategory := models.Category{}
	result := s.db.Where("id = ? AND user_id = ?", id, userId).First(&dbCategory)
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.db.Delete(&dbCategory)
	if result.Error != nil {
		if pgError := result.Error.(*pgconn.PgError); errors.Is(result.Error, pgError) {
			if pgError.Code == "23503" && pgError.ConstraintName == "fk_transactions_category" {
				return nil, status.Errorf(codes.InvalidArgument, "Cannot delete category that has existing transactions")
			}
		}
		return nil, result.Error
	}

	return &pb.DeleteCategoryResponse{Message: "Category deleted successfully"}, nil
}
