package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
)

type Order struct {
	Id          uuid.UUID  `gorm:"type:char(36);primaryKey"`
	OrderedAt   time.Time  `gorm:"size:6;index;not null"`
	Orderer     uuid.UUID  `gorm:"type:char(36);index;not null"`
	EditCount   uint16     `gorm:"not null"`
	EditTotal   uint16     `gorm:"not null"`
	DueDate     *time.Time `gorm:"type:date"`
	Assignee    *uuid.UUID `gorm:"type:char(36);index"`
	Requirement *string    `gorm:"size:2000"`
}

type CreateOrderOption struct {
	Orderer     User
	Requirement *string
}

func CreateOrder(option CreateOrderOption) Order {
	return Order{
		Id:          uuid.New(),
		OrderedAt:   time.Now(),
		Orderer:     option.Orderer.Id,
		Requirement: option.Requirement,
	}
}

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	Transaction(ctx context.Context, fn func(orderRepo OrderTxRepository) error, options ...*sql.TxOptions) error
	With(tx gormx.Tx) OrderTxRepository
	GetById(ctx context.Context, orderId uuid.UUID) (*Order, error)
	GetRecentByOrdererId(ctx context.Context, ordererId uuid.UUID) (*Order, error)
}

type OrderTxRepository interface {
	OrderRepository
	gormx.Tx
}

type RequestOrder struct {
	UserId      uuid.UUID
	Requirement string
}

type OrderUseCase interface {
	RequestOrder(ctx context.Context, or RequestOrder) (uuid.UUID, error)
}