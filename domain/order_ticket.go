package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type OrderTicket struct {
	Id              uuid.UUID       `gorm:"type:char(36);primaryKey"`
	OwnerId         uuid.UUID       `gorm:"type:char(36);index;not null"`
	OrderCount      uint8          `gorm:"not null"`
	TotalOrderCount uint8          `gorm:"not null"`
	EditCount       uint8           `gorm:"not null"`
	BuyAt           time.Time       `gorm:"size:6;index;not null"`
	StartAt         *time.Time      `gorm:"size:6;index"`
	EndAt           *time.Time      `gorm:"size:6;index"`
}

type CreateOrderTicketOption struct {
	OwnerId uuid.UUID
	Product *Product
	StartAt *time.Time
}

func CreateOrderTicket(option CreateOrderTicketOption) (res OrderTicket, err error) {
	if option.Product == nil {
		err = errors.New("isn't nil")
		return
	}

	if option.Product.SubscribeTermUnitPart != SubscribeTermUnitNone && option.StartAt == nil {
		err = errors.New("weird request data")
		return
	}



	res = OrderTicket{
		Id:              uuid.New(),
		OwnerId:         option.OwnerId,
		OrderCount:      0,
		TotalOrderCount: option.Product.OrderCount,
		EditCount:       option.Product.EditCount,
		BuyAt:           time.Now(),
	}
	return
}

type OrderTicketRepository interface {
	Save(ctx context.Context, orderTicket *OrderTicket) error
	GetById(ctx context.Context, id uuid.UUID) (orderTicket *OrderTicket, err error)
}