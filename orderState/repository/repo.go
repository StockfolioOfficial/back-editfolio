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
			Id:          1,
			Code:        domain.OrderStateCodeDefault,
			Content:     "편집자 배정 중",
			LongContent: "영상에 알맞는 편집자를\n 배정 중입니다.",
			Emoji:       "🤔",
		},
		{
			Id:          2,
			Code:        domain.OrderStateCodeTake,
			Content:     "영상 검토 중",
			LongContent: "배정된 편집자가 영상을\n 열심히 확인하고 있어요",
			Emoji:       "👀",
		},
		{
			Id:          3,
			Code:        domain.OrderStateCodeNone,
			Content:     "편집 중",
			LongContent: "영상을 이쁘게 자르고 붙이는 중...",
			Emoji:       "😍",
			ParentId:    pointer.Uint8(2),
		},
		{
			Id:          4,
			Code:        domain.OrderStateCodeNone,
			Content:     "이펙트 추가 중",
			LongContent: "아주 환상적인 이펙트를 입히는 중입니다.",
			Emoji:       "🎇",
			ParentId:    pointer.Uint8(2),
		},
		{
			Id:          5,
			Code:        domain.OrderStateCodeNone,
			Content:     "완료",
			LongContent: "영상편집이 완료되었습니다",
			Emoji:       "😘",
			ParentId:    pointer.Uint8(2),
		},
		{
			Id:          6,
			Code:        domain.OrderStateCodeDone,
			LongContent: "영상편집이 완료되었습니다",
			Emoji:       "😘",
			Content:     "최종 완료",
		},
		{
			Id:          7,
			Code:        domain.OrderStateCodeRequestEdit,
			Content:     "수정 중",
			LongContent: "요청하신 수정사항을 작업중입니다.",
			Emoji:       "🛠",
		},
		{
			Id:          8,
			Code:        domain.OrderStateCodeEditDone,
			Content:     "수정 완료",
			LongContent: "영상편집이 완료되었습니다",
			Emoji:       "😘",
			ParentId:    pointer.Uint8(7),
		},
	}
	db.Create(bookedOrderState)
	return &repo{db: db}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) GetByCode(ctx context.Context, code domain.OrderStateCode) (res *domain.OrderState, err error) {
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


