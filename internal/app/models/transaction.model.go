package models

import (
	base_models "pft/main/internal/app/models/base"
	"time"

	"github.com/google/uuid"
)

type TransactionType int

const (
	Income TransactionType = iota
	Expense
)

type Transaction struct {
	Description     string          `gorm:"type:varchar(255);not null" validate:"required,min=1,max=255"`
	Amount          string          `gorm:"type:decimal(19,4);not null" validate:"required"`
	Type            TransactionType `gorm:"type:int;not null"`
	TransactionDate time.Time       `gorm:"index" validate:"required"`
	CategoryId      uuid.UUID       `gorm:"type:uuid;not null" validate:"required"`
	AccountId       uuid.UUID       `gorm:"type:uuid;not null" validate:"required"`
	AttachmentId    uuid.UUID       `gorm:"type:uuid;not null"`

	// Embedded
	base_models.BaseAudit
	base_models.BaseId
	base_models.BaseUser
	Category   *Category   `gorm:"foreignKey:CategoryId"`
	Account    *Account    `gorm:"foreignKey:AccountId"`
	Attachment *Attachment `gorm:"foreignKey:AttachmentId"`
}
