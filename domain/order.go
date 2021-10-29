package domain

import (
	"context"
	"database/sql"
	"github.com/stockfolioofficial/back-editfolio/util/pointer"
	"time"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
)

type CreateOrderOption struct {
	Orderer     uuid.UUID
	EditCount   uint8
	State       uint8
	Requirement *string
}

func CreateOrder(option CreateOrderOption) Order {
	return Order{
		Id:             uuid.New(),
		OrderedAt:      time.Now(),
		Orderer:        option.Orderer,
		TotalEditCount: option.EditCount,
		State:          option.State,
		Requirement:    option.Requirement,
	}
}

type Order struct {
	Id             uuid.UUID  `gorm:"type:char(36);primaryKey"`
	OrderedAt      time.Time  `gorm:"type:datetime(6);index;not null"`
	Orderer        uuid.UUID  `gorm:"type:char(36);index;not null"`
	EditCount      uint8      `gorm:"not null"`
	TotalEditCount uint8      `gorm:"not null"`
	State          uint8      `gorm:"not null"`
	DueDate        *time.Time `gorm:"type:date"`
	Assignee       *uuid.UUID `gorm:"type:char(36);index"`
	Requirement    *string    `gorm:"size:2000"`
	DoneAt         *time.Time `gorm:"type:datetime(6);index"`
}

func (Order) TableName() string {
	return "order"
}

func (o *Order) RemainingEditCount() uint8 {
	return o.TotalEditCount - o.EditCount
}

func (o *Order) Done() {
	o.DoneAt = pointer.Time(time.Now())
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

type OrderDone struct {
	UserId uuid.UUID
}

type UpdateOrderInfo struct {
	OrderId    uuid.UUID
	DueDate    time.Time
	Assignee   uuid.UUID
	OrderState uint8
}

type OrderAssignSelf struct {
	OrderId  uuid.UUID
	Assignee uuid.UUID
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

type RecentOrderInfo struct {
	OrderId            uuid.UUID
	OrderedAt          time.Time
	DueDate            *time.Time
	AssigneeNickname   *string
	OrderState         uint8
	OrderStateContent  string
	RemainingEditCount uint8
}

type OrderAssigneeInfo struct {
	Id       uuid.UUID
	Name     string
	Nickname string
}

type OrderDetailInfo struct {
	OrderId            uuid.UUID
	OrderedAt          time.Time
	Orderer            uuid.UUID
	DueDate            *time.Time
	AssigneeInfo       *OrderAssigneeInfo
	OrderState         uint8
	OrderStateContent  string
	RemainingEditCount uint8
	Requirement        string
}

type OrderUseCase interface {
	RequestOrder(ctx context.Context, in RequestOrder) (uuid.UUID, error)
	OrderDone(ctx context.Context, in OrderDone) (uuid.UUID, error)

	UpdateOrderInfo(ctx context.Context, in UpdateOrderInfo) error
	OrderAssignSelf(ctx context.Context, in OrderAssignSelf) error

	GetRecentProcessingOrder(ctx context.Context, userId uuid.UUID) (RecentOrderInfo, error)
	GetOrderDetailInfo(ctx context.Context, orderId uuid.UUID) (OrderDetailInfo, error)

	Fetch(ctx context.Context, option FetchOrderOption) ([]OrderInfo, error)
}