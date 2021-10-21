package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
)

type Manager struct {
	Id        uuid.UUID  `gorm:"type:char(36);primaryKey"`
	Name      string     `gorm:"size:60;index;not null"`
	Nickname  string     `gorm:"size:60;index;not null"`
	CreatedAt time.Time  `gorm:"size:6;not null"`
	UpdatedAt time.Time  `gorm:"size:6;not null"`
	DeletedAt *time.Time `gorm:"size:6;index"`
}

func (Manager) TableName() string {
	return "manager"
}

type ManagerCreateOption struct {
	User     *User
	Name     string
	Nickname string
}

func CreateManager(option ManagerCreateOption) Manager {
	return Manager{
		Id:       option.User.Id,
		Name:     option.Name,
		Nickname: option.Nickname,
	}
}

type ManagerRepository interface {
	Save(ctx context.Context, manager *Manager) error
	GetById(ctx context.Context, userId uuid.UUID) (*Manager, error)
	With(tx gormx.Tx) ManagerTxRepository
}

type ManagerTxRepository interface {
	ManagerRepository
	gormx.Tx
}
