package models

import (
	base_models "pft/main/internal/app/models/base"
)

type BalanceTotal struct {
	IncomeTotal  string `gorm:"type:decimal(19,4);not null" validate:"required"`
	ExpenseTotal string `gorm:"type:decimal(19,4);not null" validate:"required"`
	Month        string `gorm:"type:varchar(2);not null;uniqueIndex:idx_balance_total_month_year" validate:"required"`
	Year         string `gorm:"type:varchar(4);not null;uniqueIndex:idx_balance_total_month_year" validate:"required"`
	UserId       string `gorm:"type:uuid;not null;uniqueIndex:idx_balance_total_month_year" validate:"required"`

	// Embedded
	base_models.BaseId
	base_models.BaseAudit
}
