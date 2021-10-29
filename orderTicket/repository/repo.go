package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
	"gorm.io/gorm"
	"time"
)

func NewOrderTicketRepository(db *gorm.DB) domain.OrderTicketRepository {
	db.AutoMigrate(&domain.OrderTicket{})
	return &repo{db: db}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) Transaction(ctx context.Context, fn func(orderTicketRepo domain.OrderTicketTxRepository) error, options ...*sql.TxOptions) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&repo{db: tx})
	}, options...)
}

func (r *repo) Save(ctx context.Context, orderTicket *domain.OrderTicket) error {
	return gormx.Upsert(ctx, r.db, orderTicket)
}

func (r *repo) GetById(ctx context.Context, id uuid.UUID) (res *domain.OrderTicket, err error) {
	var entity domain.OrderTicket
	err = r.db.WithContext(ctx).First(&entity, id).Error
	if err == nil {
		res = &entity
	} else if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (r *repo) GetByExOrderId(ctx context.Context, exId string) (res *domain.OrderTicket, err error) {
	var entity domain.OrderTicket
	err = r.db.WithContext(ctx).
		Where("`ex_order_id` = ?", exId).
		First(&entity).Error

	if err == nil {
		res = &entity
	} else if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (r *repo) GetEndByOwnerId(ctx context.Context, id uuid.UUID) (res *domain.OrderTicket, err error) {
	var entity domain.OrderTicket
	err = r.db.WithContext(ctx).
		Order("`end_at` desc").
		Where("`owner_id` = ? AND `end_at` IS NOT NULL", id).
		First(&entity).Error
	if err == nil {
		res = &entity
	} else if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (r *repo) GetByOwnerIdBetweenStartAndEnd(ctx context.Context, id uuid.UUID, at time.Time) (res *domain.OrderTicket, err error) {
	var entity domain.OrderTicket
	err = r.db.WithContext(ctx).
		Where("`owner_id` = ?", id).
		Where("? between `start_at` AND `end_at`", at).
		First(&entity).Error
	if err == nil {
		res = &entity
	} else if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (r *repo) Get() *gorm.DB {
	return r.db
}
