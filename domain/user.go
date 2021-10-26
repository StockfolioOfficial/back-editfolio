package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
	"github.com/stockfolioofficial/back-editfolio/util/pointer"
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
	Customer  *Customer  `gorm:"foreignKey:Id"`
	Manager   *Manager   `gorm:"foreignKey:Id"`
	MyOrder   []Order    `gorm:"foreignKey:orderer"`
	Ticket    []Order    `gorm:"foreignKey:assignee"`
}

func (User) TableName() string {
	return "user"
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

func (u *User) UpdateUsername(username string) {
	u.Username = username
	u.stampUpdate()
}

func (u *User) LoadManagerInfo(ctx context.Context, repo ManagerRepository) (err error) {
	u.Manager, err = repo.GetById(ctx, u.Id)
	if err != nil {
		return
	}

	if u.Manager == nil {
		err = ItemNotFound
	}
	return
}

func (u *User) ComparePassword(plainPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPass)) == nil
}

func (u *User) IsCustomer() bool {
	return u.HasRole(CustomerUserRole)
}

func (u *User) IsAdmin() bool {
	return u.HasRole(AdminUserRole)
}

func (u *User) IsSuperAdmin() bool {
	return u.HasRole(SuperAdminUserRole)
}

func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

func (u *User) LoadCustomerInfo(ctx context.Context, repo CustomerRepository) (err error) {
	u.Customer, err = repo.GetById(ctx, u.Id)
	if err != nil {
		return
	}

	if u.Customer == nil {
		err = ItemNotFound
	}
	return
}

func ExistsAdmin(u *User) bool {
	return u != nil && !u.IsDeleted() && u.IsAdmin()
}

func (u *User) UpdatePassword(plainPass string) {
	generated, _ := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost+2)
	u.Password = string(generated)
	u.stampUpdate()
}

func (u *User) StampUpdate() {
	u.stampUpdate()
}

func (u *User) UpdateManagerInfo(username, name, nickname string) {
	defer u.stampUpdate()
	u.UpdateUsername(username)
	if u.Manager == nil {
		return
	}
	u.Manager.Name = name
	u.Manager.Nickname = nickname
}

func (u *User) stampUpdate() {
	u.UpdatedAt = time.Now()
}

func (u *User) Delete() {
	u.DeletedAt = pointer.Time(time.Now())
}

func (u *User) UpdateCustomerInfo(name, channelName, channelLink, email, mobile, personaLink, onedriveLink, memo string) {
	defer u.stampUpdate()
	u.UpdateUsername(email)
	u.UpdatePassword(mobile)

	var customer = u.Customer
	if customer == nil {
		return
	}

	customer.Name = name
	customer.ChannelName = channelName
	customer.ChannelLink = channelLink
	customer.Email = email
	customer.Mobile = mobile
	customer.PersonaLink = personaLink
	customer.OnedriveLink = onedriveLink
	customer.Memo = memo
}

func ExistsCustomer(u *User) bool {
	return u != nil && !u.IsDeleted() && u.IsCustomer()
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetById(ctx context.Context, userId uuid.UUID) (*User, error)
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

type UpdateAdminPassword struct {
	OldPassword string
	NewPassword string
	UserId      uuid.UUID
}

type UpdateAdminInfo struct {
	UserId   uuid.UUID
	Name     string
	Username string
	Nickname string
}

type UpdateAdminInfoBySuperAdmin struct {
	UserId   uuid.UUID
	Name     string
	Password string
	Username string
	Nickname string
}

type CreateAdminUser struct {
	Name     string
	Email    string
	Password string
	Nickname string
}

type UpdateCustomerUser struct {
	UserId       uuid.UUID
	Name         string
	ChannelName  string
	ChannelLink  string
	Email        string
	Mobile       string
	PersonaLink  string
	OnedriveLink string
	Memo         string
}

type UserUseCase interface {
	CreateCustomerUser(ctx context.Context, cu CreateCustomerUser) (uuid.UUID, error)
	UpdateCustomerUser(ctx context.Context, cu UpdateCustomerUser) error
	UpdateAdminPassword(ctx context.Context, up UpdateAdminPassword) error
	UpdateAdminInfo(ctx context.Context, ui UpdateAdminInfo) error
	UpdateAdminInfoBySuperAdmin(ctx context.Context, fu UpdateAdminInfoBySuperAdmin) error
	SignInUser(ctx context.Context, si SignInUser) (string, error)
	DeleteCustomerUser(ctx context.Context, du DeleteCustomerUser) error
	CreateAdminUser(ctx context.Context, au CreateAdminUser) (uuid.UUID, error)
	DeleteAdminUser(ctx context.Context, da DeleteAdminUser) error
}

type TokenGenerateAdapter interface {
	Generate(User) (string, error)
}

type DeleteCustomerUser struct {
	Id uuid.UUID
}

type DeleteAdminUser struct {
	Id uuid.UUID
}
