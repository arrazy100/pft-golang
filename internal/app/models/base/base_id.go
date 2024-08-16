package base_models

import "github.com/google/uuid"

type BaseId struct {
	Id uuid.UUID `gorm:"type:uuid;primary_key" validate:"required"`
}

func (b *BaseId) SetId() {
	b.Id = uuid.New()
}
