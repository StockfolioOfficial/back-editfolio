package repository

import (
	"context"

	"github.com/google/uuid"
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

func (r *repo) FetchByIds(ctx context.Context, ids []uuid.UUID) (list []domain.Manager, err error) {
	err = r.db.WithContext(ctx).Find(&list, ids).Error
	return
}

func (r *repo) GetById(ctx context.Context, userId uuid.UUID) (manager *domain.Manager, err error) {
	var entity domain.Manager

	err = r.db.WithContext(ctx).First(&entity, userId).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	manager = &entity
	return
}

func (r *repo) GetByNickname(ctx context.Context, nickName string) (nickname *domain.Manager, err error) {
	var entity domain.Manager

	err = r.db.WithContext(ctx).First(&entity, nickName).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	nickname = &entity
	return
}

func (r *repo) Get() *gorm.DB {
	return r.db
}

func (r *repo) Save(ctx context.Context, manager *domain.Manager) error {
	return gormx.Upsert(ctx, r.db, manager)
}

func (r *repo) With(tx gormx.Tx) domain.ManagerTxRepository {
	return &repo{db: tx.Get()}
}
