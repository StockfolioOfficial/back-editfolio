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
	Orderer     string     `gorm:"size:36;not null"`
	EditCount   uint16     `gorm:"not null"`
	EditTotal   uint16     `gorm:"not null"`
	DueDate     *time.Time `gorm:"size:6"`
	Assignee    *string    `gorm:"size:36"`
	Requirement string     `gorm:size:2000`
}

func MakeOrder(orderer string) Order {
	return Order{
		Id:        uuid.New(),
		OrderedAt: time.Now(),
		Orderer:   orderer,
		DueDate:   nil,
		Assignee:  nil,
	}
}

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	Transaction(ctx context.Context, fn func(orderRepo OrderTxRepository) error, options ...*sql.TxOptions) error
	With(tx gormx.Tx) OrderTxRepository
	GetById(ctx context.Context, orderId uuid.UUID) (*Order, error)
}

type OrderTxRepository interface {
	OrderRepository
	gormx.Tx
}

type VideoEditRequirement struct {
	Id          uuid.UUID
	Requirement string
}

type OrderUseCase interface {
	VideoEditRequirement(ctx context.Context, vr VideoEditRequirement) error
}

func (o *Order) LoadOrderInfo(ctx context.Context, repo OrderRepository) (err error) {
	o, err = repo.GetById(ctx, o.Id)
	if err != nil {
		return
	}

	if o == nil {
		err = ItemNotFound
	}
	return
}

func (o *Order) UpdateVideoEditRequirement(requirement string) {
	o.Requirement = requirement
}
