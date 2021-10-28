package repository

import (
	"context"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"gorm.io/gorm"
)

func NewOrderStateRepository(db *gorm.DB) domain.OrderStateRepository {
	db.AutoMigrate(&domain.OrderState{})
	bookedOrderState := []domain.OrderState{
		{
			Id:      1,
			Content: "편집자 배정 중",
		},
		{
			Id:      2,
			Content: "편집 중",
		},
		{
			Id:      3,
			Content: "이펙트 추가 중",
		},
		{
			Id:      4,
			Content: "수정 중",
		},
		{
			Id:      5,
			Content: "완료",
		},
	}
	db.Create(bookedOrderState)
	return &repo{db: db}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) FetchFull(ctx context.Context) (list []domain.OrderState, err error) {
	err = r.db.WithContext(ctx).
		Order("`id` asc").
		Find(&list).Error
	return
}

func (r *repo) FetchByIds(ctx context.Context, ids []uint8) (list []domain.OrderState, err error) {
	err = r.db.WithContext(ctx).
		Order("`id` asc").
		Find(&list, ids).Error
	return
}

func (r *repo) GetById(ctx context.Context, id uint8) (res *domain.OrderState, err error) {
	var entity domain.OrderState
	err = r.db.WithContext(ctx).Find(&entity, id).Error
	if err == nil {
		res = &entity
	} else if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

