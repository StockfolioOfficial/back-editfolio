package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	SuperAdminUserRole UserRole = "SUPER_ADMIN"
	AdminUserRole      UserRole = "ADMIN"
	CustomerUserRole   UserRole = "CUSTOMER"
)

type User struct {
	Id        uuid.UUID  `gorm:"type:char(36);primaryKey"`
	Role      UserRole   `gorm:"size:30;index;not null"`
	Username  string     `gorm:"size:320;unique;not null"`
	Password  string     `gorm:"size:60;not null"`
	CreatedAt time.Time  `gorm:"size:6;not null"`
	UpdatedAt time.Time  `gorm:"size:6;not null"`
	DeletedAt *time.Time `gorm:"size:6;index"`
}

type UserCreateOption struct {
	Role     UserRole
	Username string
}

func CreateUser(option UserCreateOption) User {
	return User{
		Id:        uuid.New(),
		Role:      option.Role,
		Username:  option.Username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
}

func (u *User) ComparePassword(plainPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPass)) == nil
}

func (u *User) UpdatePassword(plainPass string) {
	generated, _ := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost+2)
	u.Password = string(generated)
	u.stampUpdate()
}

func (u *User) stampUpdate() {
	u.UpdatedAt = time.Now()
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	Transaction(ctx context.Context, fn func(userRepo UserTxRepository) error, options ...*sql.TxOptions) error
	With(tx gormx.Tx) UserTxRepository
}

type UserTxRepository interface {
	UserRepository
	gormx.Tx
}

type CreateCustomerUser struct {
	Name   string
	Email  string
	Mobile string
}

type SignInUser struct {
	Username string
	Password string
}

type UserUseCase interface {
	CreateCustomerUser(ctx context.Context, cu CreateCustomerUser) (uuid.UUID, error)
	SignInUser(ctx context.Context, si SignInUser) (string, error)
}
