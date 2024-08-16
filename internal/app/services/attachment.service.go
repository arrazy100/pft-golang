package services

import (
	"context"
	"fmt"
	pb "pft/main/internal/app/generated_proto"
	"pft/main/internal/app/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AttachmentService struct {
	db       *gorm.DB
	validate *validator.Validate
	pb.UnimplementedAttachmentServiceServer
}

func NewAttachmentService(db *gorm.DB, validate *validator.Validate) (*AttachmentService, error) {
	return &AttachmentService{db: db, validate: validate}, nil
}

func (s *AttachmentService) CreateAttachment(ctx context.Context, req *pb.CreateAttachmentRequest) (*pb.CreateAttachmentResponse, error) {
	attachment := req.GetData()

	if attachment == nil {
		return nil, fmt.Errorf("no attachment data provided")
	}

	userId := "597c5e7d-63dc-4df1-9954-fefe8b415634"

	dbAttachment := &models.Attachment{
		Type:       models.AttachmentType(attachment.Type),
		ContentUrl: attachment.ContentUrl,
	}

	dbAttachment.SetId()
	dbAttachment.SetAuditCreate(userId)
	dbAttachment.SetUser(userId)

	// validation
	err := s.validate.Struct(dbAttachment)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}

		return nil, err
	}

	// Save to database
	result := s.db.Create(dbAttachment)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert back to protobuf Transaction
	createdAttachment := &pb.Attachment{
		Id:         dbAttachment.Id.String(),
		Type:       pb.AttachmentType(dbAttachment.Type),
		ContentUrl: dbAttachment.ContentUrl,
		UserId:     dbAttachment.UserId,
		CreatedAt:  dbAttachment.CreatedAt.String(),
		CreatedBy:  dbAttachment.CreatedBy,
	}

	return &pb.CreateAttachmentResponse{
		Message: "Attachment created successfully",
		Data:    createdAttachment,
	}, nil
}
