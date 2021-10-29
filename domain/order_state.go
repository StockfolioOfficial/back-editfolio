package domain

import "context"

type OrderState struct {
	Id      uint8   `gorm:"primaryKey"`
	Content string  `gorm:"size:150;not null;index"`
	Orders  []Order `gorm:"foreignKey:State"`
}

func (OrderState) TableName() string {
	return "order_state"
}

type OrderStateRepository interface {
	GetById(ctx context.Context, id uint8) (*OrderState, error)

	FetchFull(ctx context.Context) ([]OrderState, error)
	FetchByIds(ctx context.Context, ids []uint8) ([]OrderState, error)
}

type OrderStateInfo struct {
	Id      uint8
	Content string
}

type OrderStateUseCase interface {
	FetchFull(ctx context.Context) ([]OrderStateInfo, error)
}