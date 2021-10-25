package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
)

type Customer struct {
	Id             uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name           string    `gorm:"size:320;index;not null"`
	ChannelName    string    `gorm:"size:100;index;not null"`
	ChannelLink    string    `gorm:"size:2048;not null"`
	Email          string    `gorm:"size:320;index;not null"`
	Mobile         string    `gorm:"size:24;index;not null"`
	OrderableCount int       `gorm:"column:orderable_cnt;default:0"`
	PersonaLink    string    `gorm:"size:2048;not null"`
	OnedriveLink   string    `gorm:"size:2048;not null"`
	Memo           string    `gorm:"type:text"`
}

func (Customer) TableName() string {
	return "customer"
}

type CustomerCreateOption struct {
	User           *User
	Name           string
	ChannelName    string
	ChannelLink    string
	Email          string
	Mobile         string
	OrderableCount int
	PersonaLink    string
	OnedriveLink   string
	Memo           string
}

func CreateCustomer(option CustomerCreateOption) Customer {
	return Customer{
		Id:             option.User.Id,
		Name:           option.Name,
		ChannelName:    option.ChannelName,
		ChannelLink:    option.ChannelLink,
		Email:          option.Email,
		Mobile:         option.Mobile,
		OrderableCount: option.OrderableCount,
		PersonaLink:    option.PersonaLink,
		OnedriveLink:   option.OnedriveLink,
		Memo:           option.Memo,
	}
}

type CustomerRepository interface {
	Save(ctx context.Context, customer *Customer) error
	GetById(ctx context.Context, userId uuid.UUID) (*Customer, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	With(tx gormx.Tx) CustomerTxRepository
}

type CustomerTxRepository interface {
	CustomerRepository
	gormx.Tx
}
