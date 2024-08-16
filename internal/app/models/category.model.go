package models

import (
	base_models "pft/main/internal/app/models/base"
)

type Category struct {
	Name string `gorm:"type:varchar(255);uniqueIndex:idx_category_name;not null" validate:"required,min=3,max=255"`

	// Embedded
	base_models.BaseId
	base_models.BaseAudit
	base_models.BaseUser
}
