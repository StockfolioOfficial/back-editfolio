package gormx

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Upsert(ctx context.Context, db *gorm.DB, model interface{}) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(model).Error
}
