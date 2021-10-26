package domain

import (
	"time"

	"github.com/google/uuid"
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
