package domain

import (
	"github.com/google/uuid"
	"time"
)

type OrderTicketType string

const (
	OrderTicketTypeSubscribe = "SUBSCRIBE"
	OrderTicketTypeOnceOrder = "VOLUME_BASED"
)

type OrderTicket struct {
	Id              uuid.UUID       `gorm:"type:char(36);primaryKey"`
	OwnerId         uuid.UUID       `gorm:"type:char(36);index;not null"`
	Type            OrderTicketType `gorm:"size:30;index;not null"`
	OrderCount      uint16          `gorm:"not null"`
	TotalOrderCount uint16          `gorm:"not null"`
	EditCount       uint8           `gorm:"not null"`
	BuyAt           time.Time       `gorm:"size:6;index;not null"`
	StartAt         *time.Time      `gorm:"size:6;index"`
	EndAt           *time.Time      `gorm:"size:6;index"`
}

type OrderTicketRepository interface {

}