package domain

import "context"

type OrderStateCode string

const (
	OrderStateCodeNone OrderStateCode = "NONE"
	OrderStateCodeDefault OrderStateCode = "DEFAULT"
	OrderStateCodeTake OrderStateCode = "TAKE"
	OrderStateCodeRequestEdit OrderStateCode = "REQUEST_EDIT"
	OrderStateCodeEditDone OrderStateCode = "EDIT_DONE"
	OrderStateCodeDone OrderStateCode = "DONE"
)

type OrderState struct {
	Id       uint8          `gorm:"primaryKey"`
	Code     OrderStateCode `gorm:"size:20;index;not null"`
	Content  string         `gorm:"size:150;index;not null"`
	ParentId *uint8         `gorm:"index"`
	Parent   *OrderState    `gorm:"foreignKey:ParentId"`
	Orders   []Order        `gorm:"foreignKey:State"`
}

func (OrderState) TableName() string {
	return "order_state"
}

type OrderStateRepository interface {
	GetById(ctx context.Context, id uint8) (*OrderState, error)

	FetchFull(ctx context.Context) ([]OrderState, error)
	FetchByIds(ctx context.Context, ids []uint8) ([]OrderState, error)

	GetByCode(ctx context.Context, code OrderStateCode) (*OrderState, error)
	FetchByParentId(ctx context.Context, parentId uint8) ([]OrderState, error)
}

type OrderStateInfo struct {
	Id      uint8
	Content string
}

type OrderStateUseCase interface {
	FetchFull(ctx context.Context) ([]OrderStateInfo, error)
	FetchByParentId(ctx context.Context, parentId uint8) ([]OrderStateInfo, error)
}