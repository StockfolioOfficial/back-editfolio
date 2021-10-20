package repository

import (
	"context"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
	"gorm.io/gorm"
)

func NewManagerRepository(db *gorm.DB) domain.ManagerRepository {
	db.AutoMigrate(&domain.Manager{})
	return &repo{db: db}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) Get() *gorm.DB {
	return r.db
}

func (r *repo) Save(ctx context.Context, manager *domain.Manager) error {
	return r.db.WithContext(ctx).Save(manager).Error
}

func (r *repo) With(tx gormx.Tx) domain.ManagerTxRepository {
	return &repo{db: tx.Get()}
}

