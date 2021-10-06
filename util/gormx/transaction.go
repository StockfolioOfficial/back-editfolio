package gormx

import "gorm.io/gorm"

type Tx interface {
	Get() *gorm.DB
}