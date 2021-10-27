package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	db.AutoMigrate(&domain.User{})
	return &repo{
		db: db,
	}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) FetchAllAdmin(ctx context.Context, option domain.FetchAdminOption) (list []domain.User, err error) {
	err = r.db.WithContext(ctx).
		Joins("Manager").
		Where("`role` = ?", domain.AdminUserRole).
		Or("`role` = ?", domain.SuperAdminUserRole).
		Find(&list).Error
	return
}

func (r *repo) FetchAllCustomer(ctx context.Context, option domain.FetchCustomerOption) (list []domain.User, err error) {
	err = r.db.WithContext(ctx).
		Joins("Customer").
		Where("`role` = ?", domain.CustomerUserRole).
		Find(&list).Error
	return
}

func (r *repo) GetByUsername(ctx context.Context, username string) (user *domain.User, err error) {
	var entity domain.User
	err = r.db.WithContext(ctx).
		Where("`username` = ?", username).
		First(&entity).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	user = &entity
	return
}

func (r *repo) GetById(ctx context.Context, userId uuid.UUID) (user *domain.User, err error) {
	var entity domain.User
	err = r.db.WithContext(ctx).First(&entity, userId).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	user = &entity
	return
}

func (r *repo) Save(ctx context.Context, user *domain.User) error {
	return gormx.Upsert(ctx, r.db, user)
}

func (r *repo) Get() *gorm.DB {
	return r.db
}

func (r *repo) Transaction(ctx context.Context, fn func(userRepo domain.UserTxRepository) error, options ...*sql.TxOptions) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&repo{db: tx})
	}, options...)
}

func (r *repo) With(tx gormx.Tx) domain.UserTxRepository {
	return &repo{db: tx.Get()}
}

func (r *repo) Delete(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}
