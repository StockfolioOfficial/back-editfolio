package repository

import (
	"context"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/pointer"
	"gorm.io/gorm"
)

func NewOrderStateRepository(db *gorm.DB) domain.OrderStateRepository {
	db.AutoMigrate(&domain.OrderState{})
	bookedOrderState := []domain.OrderState{
		{
			Id:      1,
			Code:    domain.OrderStateCodeDefault,
			Content: "편집자 배정 중",
		},
		{
			Id:      2,
			Code:    domain.OrderStateCodeNone,
			Content: "영상 검토 중",
		},
		{
			Id:       3,
			Code:     domain.OrderStateCodeNone,
			Content:  "편집 중",
			ParentId: pointer.Uint8(2),
		},
		{
			Id:       4,
			Code:     domain.OrderStateCodeNone,
			Content:  "이펙트 추가 중",
			ParentId: pointer.Uint8(2),
		},
		{
			Id:       5,
			Code:     domain.OrderStateCodeNone,
			Content:  "완료",
			ParentId: pointer.Uint8(2),
		},
		{
			Id:      6,
			Code:    domain.OrderStateCodeDone,
			Content: "최종 완료",
		},
		{
			Id:      7,
			Code:    domain.OrderStateCodeRequestEdit,
			Content: "수정 중",
		},
		{
			Id:       8,
			Code:     domain.OrderStateCodeEditDone,
			Content:  "수정 완료",
			ParentId: pointer.Uint8(7),
		},
	}
	db.Create(bookedOrderState)
	return &repo{db: db}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) GetByCode(ctx context.Context, code string) (res *domain.OrderState, err error) {
	var entity domain.OrderState
	err = r.db.WithContext(ctx).First(&entity, "`code` = ?", code).Error
	if err == nil {
		res = &entity
	} else if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (r *repo) FetchByParentId(ctx context.Context, parentId uint8) (list []domain.OrderState, err error) {
	err = r.db.WithContext(ctx).
		Order("`id` asc").
		Where("`parent_id` = ?", parentId).
		Find(&list).
		Error
	return
}

func (r *repo) GetById(ctx context.Context, id uint8) (res *domain.OrderState, err error) {
	var entity domain.OrderState
	err = r.db.WithContext(ctx).First(&entity, id).Error
	if err == nil {
		res = &entity
	} else if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
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


