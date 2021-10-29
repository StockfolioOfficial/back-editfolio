package domain

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
	"time"
)

type CreateOrderTicketOption struct {
	ExOrderId       string
	OwnerId         uuid.UUID
	TotalOrderCount uint8
	EditCount       uint8
	StartAt         *time.Time
	EndAt           *time.Time
}

func CreateOrderTicket(option CreateOrderTicketOption) OrderTicket {
	return OrderTicket{
		Id:              uuid.New(),
		ExOrderId:       option.ExOrderId,
		OwnerId:         option.OwnerId,
		TotalOrderCount: option.TotalOrderCount,
		EditCount:       option.EditCount,
		CreatedAt:       time.Now(),
		StartAt:         option.StartAt,
		EndAt:           option.EndAt,
	}
}

type OrderTicket struct {
	Id              uuid.UUID  `gorm:"type:char(36);primaryKey"`
	ExOrderId       string     `gorm:"size:90;unique;not null"`
	OwnerId         uuid.UUID  `gorm:"type:char(36);index;not null"`
	OrderCount      uint8      `gorm:"not null"`
	TotalOrderCount uint8      `gorm:"not null"`
	EditCount       uint8      `gorm:"not null"`
	CreatedAt       time.Time  `gorm:"size:datetime(6);index;not null"`
	StartAt         *time.Time `gorm:"size:datetime(6);index"`
	EndAt           *time.Time `gorm:"type:datetime(6);index"`
}

func (o *OrderTicket) UseOrder() {
	o.OrderCount++
}

func (o OrderTicket) RemainingOrderCount() uint8 {
	return o.TotalOrderCount - o.OrderCount
}

func (o OrderTicket) IsEmptyOrderCount() bool {
	return o.RemainingOrderCount() == 0
}

type OrderTicketRepository interface {
	Save(ctx context.Context, orderTicket *OrderTicket) error
	Transaction(ctx context.Context, fn func(orderTicketRepo OrderTicketTxRepository) error, options ...*sql.TxOptions) error

	GetById(ctx context.Context, id uuid.UUID) (*OrderTicket, error)
	GetByExOrderId(ctx context.Context, exId string) (*OrderTicket, error)
	GetEndByOwnerId(ctx context.Context, id uuid.UUID) (*OrderTicket, error)
	GetByOwnerIdBetweenStartAndEnd(ctx context.Context, id uuid.UUID, at time.Time) (*OrderTicket, error)
}

type OrderTicketTxRepository interface {
	OrderTicketRepository
	gormx.Tx
}

type SubscribeUnit string

const (
	SubscribeUnitMonth SubscribeUnit = "M"
	SubscribeUnitDay SubscribeUnit = "D"
)

type CreateSubscribeTicket struct {
	ExOrderId  string
	Username   string
	Value      uint16
	Unit       SubscribeUnit
	OrderCount uint8
	EditCount  uint8
}

type OrderTicketUseCase interface {
	CreateSubscribeTicket(ctx context.Context, in CreateSubscribeTicket) (uuid.UUID, error)
}