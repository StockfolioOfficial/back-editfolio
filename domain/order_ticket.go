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
	Id   uuid.UUID
	Type OrderTicketType `gorm:"size:30;index;not null"`
	CreatedAt time.Time `gorm:"size:6;not null;"`
}
