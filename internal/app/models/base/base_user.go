package base_models

type BaseUser struct {
	UserId string `gorm:"type:uuid;not null" validate:"required"`
}

func (b *BaseUser) SetUser(userId string) {
	b.UserId = userId
}
