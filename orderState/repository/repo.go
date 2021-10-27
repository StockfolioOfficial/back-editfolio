package repository

import (
	"context"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"gorm.io/gorm"
)

func NewOrderStateRepository(db *gorm.DB) domain.OrderStateRepository {
	return &repo{db: db}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) FetchByIds(ctx context.Context, ids []uint8) (list []domain.OrderState, err error) {
	err = r.db.WithContext(ctx).Find(&list, ids).Error
	return
}

func (r *repo) GetById(ctx context.Context, id uint8) (res *domain.OrderState, err error) {
	var entity domain.OrderState
	err = r.db.WithContext(ctx).Find(&entity, id).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	} else if err == nil {
		res = &entity
	}
	return
}

