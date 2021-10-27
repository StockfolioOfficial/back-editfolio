package domain

import "context"

type OrderState struct {
	Id      uint8   `gorm:"primaryKey"` // AutoIncreament 설정 ?
	Content string  `gorm:"size:150;not null;index"`
	Orders  []Order `gorm:"foreignKey:State"`
}

type OrderStateRepository interface {
	FetchByIds(ctx context.Context, ids []uint8) ([]OrderState, error)
	GetById(ctx context.Context, id uint8) (*OrderState, error)
}