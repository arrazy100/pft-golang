package base_models

import (
	"time"

	"gorm.io/gorm"
)

type BaseAudit struct {
	CreatedAt time.Time `validate:"required"`
	CreatedBy string    `gorm:"type:uuid"`
	UpdatedAt *time.Time
	UpdatedBy *string `gorm:"type:uuid"`
	IsDeleted bool
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Set audit when creating data
func (b *BaseAudit) SetAuditCreate(createdBy string) {
	b.CreatedBy = createdBy
	b.CreatedAt = time.Now().UTC()
}

func (b *BaseAudit) SetAuditUpdate(updatedBy string) {
	now := time.Now().UTC()

	b.UpdatedBy = &updatedBy
	b.UpdatedAt = &now
}

func (b *BaseAudit) SetAuditDelete() {
	b.IsDeleted = true
	b.DeletedAt = gorm.DeletedAt{Time: time.Now().UTC()}
}
