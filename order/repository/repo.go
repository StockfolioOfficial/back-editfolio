package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
	"gorm.io/gorm"
)

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	db.AutoMigrate(&domain.Order{})
	return &repo{
		db: db,
	}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) Save(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *repo) Get() *gorm.DB {
	return r.db
}

func (r *repo) Transaction(ctx context.Context, fn func(orderRepo domain.OrderTxRepository) error, options ...*sql.TxOptions) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&repo{db: tx})
	}, options...)
}

func (r *repo) With(tx gormx.Tx) domain.OrderTxRepository {
	return &repo{db: tx.Get()}
}

func (r *repo) GetById(ctx context.Context, orderId uuid.UUID) (order *domain.Order, err error) {
	var entity domain.Order
	err = r.db.WithContext(ctx).First(&entity, orderId).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	order = &entity
	return
}
