package domain

import (
	"github.com/google/uuid"
)

type Customer struct {
	Id             uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name           string    `gorm:"size:320;index;not null"`
	ChannelName    string    `gorm:"size:100;index;not null"`
	ChannelLink    string    `gorm:"size:2048;not null"`
	Email          string    `gorm:"size:320;index;not null"`
	Mobile         string    `gorm:"size:24;index;not null"`
	OrderableCount int       `gorm:"column:orderable_cnt;not null"`
	PersonaLink    string    `gorm:"size:2048;not null"`
	OnedriveLink   string    `gorm:"size:2048;not null"`
	Memo           string    `gorm:"type:text;not null"`
}

func (Customer) TableName() string {
	return "customer"
}
