package models

import (
	base_models "pft/main/internal/app/models/base"

	"gorm.io/gorm"
)

type TransactionType int

const (
	Income TransactionType = iota
	Expense
	Transfer
)

type Transaction struct {
	gorm.Model
	base_models.BaseAudit
	Id           string          `gorm:"type:uuid;primary_key"`
	UserId       string          `gorm:"type:uuid;not null"`
	CategoryId   string          `gorm:"type:uuid;not null"`
	Description  string          `gorm:"type:varchar(255);not null"`
	AccountId    string          `gorm:"type:uuid;not null"`
	AttachmentId string          `gorm:"type:uuid;not null"`
	Amount       string          `gorm:"type:decimal(10,4);not null"`
	Type         TransactionType `gorm:"type:int;not null"`
}
