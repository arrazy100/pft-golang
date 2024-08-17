package base_models

import (
	"time"
)

type BaseAudit struct {
	CreatedAt time.Time `validate:"required"`
	CreatedBy string    `gorm:"type:uuid"`
	UpdatedAt *time.Time
	UpdatedBy *string `gorm:"type:uuid"`
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
