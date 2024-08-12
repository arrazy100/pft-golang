package base_models

import (
	"time"

	"gorm.io/gorm"
)

type BaseAudit struct {
	IsDeleted bool           `json:"is_deleted"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
