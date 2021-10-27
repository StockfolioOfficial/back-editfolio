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
	State       uint8      `gorm:"not null"`
	DueDate     *time.Time `gorm:"type:date"`
	Assignee    *uuid.UUID `gorm:"type:char(36);index"`
	Requirement *string    `gorm:"size:2000"`
	DoneAt      *time.Time `gorm:"size:6;index"`
}

func (Order) TableName() string{
	return "order"
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

type OrderGeneralState uint8

const (
	// OrderGeneralStateReady 의뢰 요청만 된 리스트
	OrderGeneralStateReady OrderGeneralState = iota

	// OrderGeneralStateProcessing 의뢰 진행 중인 리스트
	OrderGeneralStateProcessing

	// OrderGeneralStateDone 의뢰가 끝난 리스트
	OrderGeneralStateDone
)

type FetchOrderOption struct {
	OrderState OrderGeneralState
	Query      string
	Assignee   *uuid.UUID
	//TODO Sort OrderedAt, Name, Assignee, State
	//TODO Pagination
	//Limit    int64
}

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	Transaction(ctx context.Context, fn func(orderRepo OrderTxRepository) error, options ...*sql.TxOptions) error
	With(tx gormx.Tx) OrderTxRepository
	GetById(ctx context.Context, orderId uuid.UUID) (*Order, error)
	GetRecentByOrdererId(ctx context.Context, ordererId uuid.UUID) (*Order, error)
	Fetch(ctx context.Context, option FetchOrderOption) ([]Order, error)
}

type OrderTxRepository interface {
	OrderRepository
	gormx.Tx
}

type RequestOrder struct {
	UserId      uuid.UUID
	Requirement string
}

type OrderInfo struct {
	OrderId            uuid.UUID
	OrderedAt          time.Time
	OrdererName        string
	OrdererChannelName string
	OrdererChannelLink string
	AssigneeName       *string
	AssigneeNickname   *string
	OrderState         uint8
	OrderStateContent  string
	DoneAt             *time.Time
}

type UpdateOrderInfo struct {
	OrderId    uuid.UUID
	DueDate    time.Time
	Assignee   uuid.UUID
	OrderState uint8
}

type OrderUseCase interface {
	RequestOrder(ctx context.Context, or RequestOrder) (uuid.UUID, error)
	Fetch(ctx context.Context, option FetchOrderOption) (res []OrderInfo, err error)
	UpdateOrderDetailInfo(ctx context.Context, uo *UpdateOrderInfo) (err error)
}
