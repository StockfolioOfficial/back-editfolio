package domain

import "github.com/google/uuid"

type SubscribeTermUnit string

const (
	// SubscribeTermUnitNone 무한 혹은 없음
	SubscribeTermUnitNone SubscribeTermUnit = "nop"

	// SubscribeTermUnitYear 년 단위
	SubscribeTermUnitYear SubscribeTermUnit = "y"

	// SubscribeTermUnitMonth 월 단위
	SubscribeTermUnitMonth SubscribeTermUnit = "m"

	// SubscribeTermUnitDay 일 단위
	SubscribeTermUnitDay SubscribeTermUnit = "d"
)

type Product struct {
	Id                      uuid.UUID         `gorm:"type:char(36);primaryKey"`
	SubscribeTermNumberPart uint8             `gorm:"index;not null"`
	SubscribeTermUnitPart   SubscribeTermUnit `gorm:"type:char(3);index;not null"`
	OrderCount              uint8             `gorm:"not null"`
	EditCount               uint8             `gorm:"not null"`
	ExternalOriginId        uint64            `gorm:"type:not null;index"`
}
