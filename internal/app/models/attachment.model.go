package models

import (
	base_models "pft/main/internal/app/models/base"
)

type AttachmentType int

const (
	GOOGLE_DRIVE AttachmentType = iota
)

type Attachment struct {
	Type       AttachmentType `gorm:"type:int;not null;"`
	ContentUrl string         `gorm:"type:varchar(255);not null;" validate:"required,url"`

	// Embedded
	base_models.BaseId
	base_models.BaseAudit
	base_models.BaseUser
}
