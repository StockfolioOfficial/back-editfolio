package domain

import (
	"github.com/google/uuid"
)

type Manager struct {
	Id       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name     string    `gorm:"size:60;index;not null"`
	NickName string    `gorm:"size:60;index;not null"`
}
