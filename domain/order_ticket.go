package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type CreateOrderTicketOption struct {
	OwnerId         uuid.UUID
	TotalOrderCount uint8
	EditCount       uint8
	StartAt         *time.Time
	EndAt           *time.Time
}

func CreateOrderTicket(option CreateOrderTicketOption) OrderTicket {
	return OrderTicket{
		Id:              uuid.New(),
		OwnerId:         option.OwnerId,
		OrderCount:      0,
		TotalOrderCount: option.TotalOrderCount,
		EditCount:       option.EditCount,
		BuyAt:           time.Now(),
		StartAt:         option.StartAt,
		EndAt:           option.EndAt,
	}
}

type OrderTicket struct {
	Id              uuid.UUID  `gorm:"type:char(36);primaryKey"`
	OwnerId         uuid.UUID  `gorm:"type:char(36);index;not null"`
	OrderCount      uint8      `gorm:"not null"`
	TotalOrderCount uint8      `gorm:"not null"`
	EditCount       uint8      `gorm:"not null"`
	BuyAt           time.Time  `gorm:"size:6;index;not null"`
	StartAt         *time.Time `gorm:"size:6;index"`
	EndAt           *time.Time `gorm:"size:6;index"`
}

type OrderTicketRepository interface {
	Save(ctx context.Context, orderTicket *OrderTicket) error
	GetById(ctx context.Context, id uuid.UUID) (*OrderTicket, error)
}

type OrderTicketUseCase interface {

}