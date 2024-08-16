package models

import (
	base_models "pft/main/internal/app/models/base"
)

type AccountType int

const (
	BANK AccountType = iota
	EWALLET
	CARD
)

type Account struct {
	Type    AccountType `gorm:"type:int;not null"`
	Balance string      `gorm:"type:decimal(19,4);not null" validate:"required"`
	Name    string      `gorm:"type:varchar(255);not null" validate:"required"`

	// Embedded
	base_models.BaseId
	base_models.BaseAudit
	base_models.BaseUser
}
