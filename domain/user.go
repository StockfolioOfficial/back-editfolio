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

type UserCreateOption struct {
	Role     UserRole
	Username string
}

func CheckUserAlive(u *User, scope ...func(user User) bool) bool {
	if u == nil {
		return false
	}

	if u.IsDeleted() {
		return false
	}

	if len(scope) == 0 {
		return true
	}

	for i := range scope {
		if scope[i](*u) {
			return true
		}
	}

	return false
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

type User struct {
	Id        uuid.UUID  `gorm:"type:char(36);primaryKey"`
	Role      UserRole   `gorm:"size:30;index;not null"`
	Username  string     `gorm:"size:320;unique;not null"`
	Password  string     `gorm:"size:60;not null"`
	CreatedAt time.Time  `gorm:"type:datetime(6);not null"`
	UpdatedAt time.Time  `gorm:"type:datetime(6);not null"`
	DeletedAt *time.Time `gorm:"type:datetime(6);index"`
	Customer  *Customer  `gorm:"foreignKey:Id"`
	Manager   *Manager   `gorm:"foreignKey:Id"`
	MyJob     []Order    `gorm:"foreignKey:Orderer"`
	Ticket    []Order    `gorm:"foreignKey:Assignee"`
}

func (User) TableName() string {
	return "user"
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
		err = ErrItemNotFound
	}
	return
}

func (u *User) ComparePassword(plainPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPass)) == nil
}

func (u User) IsCustomer() bool {
	return u.HasRole(CustomerUserRole)
}

func (u User) IsAdmin() bool {
	return u.HasRole(AdminUserRole)
}

func (u User) IsSuperAdmin() bool {
	return u.HasRole(SuperAdminUserRole)
}

func (u User) HasRole(role UserRole) bool {
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
		err = ErrItemNotFound
	}
	return
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

type FetchAdminOption struct {
	Query string
}

type FetchCustomerOption struct {
	Query string
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	Transaction(ctx context.Context, fn func(userRepo UserTxRepository) error, options ...*sql.TxOptions) error

	GetByUsername(ctx context.Context, username string) (*User, error)
	GetById(ctx context.Context, userId uuid.UUID) (*User, error)

	FetchAllAdmin(ctx context.Context, option FetchAdminOption) ([]User, error)
	FetchAllCustomer(ctx context.Context, option FetchCustomerOption) ([]User, error)

	GetByIdWithCustomer(ctx context.Context, id uuid.UUID) (*User, error)
}

type UserTxRepository interface {
	UserRepository
	gormx.Tx
}

type SignInUser struct {
	Username string
	Password string
}

type CreateCustomerUser struct {
	Name   string
	Email  string
	Mobile string
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

type UpdateAdminInfo struct {
	UserId   uuid.UUID
	Name     string
	Username string
	Nickname string
}

type UpdateAdminPassword struct {
	UserId      uuid.UUID
	OldPassword string
	NewPassword string
}

type ForceUpdateAdminInfo struct {
	UserId   uuid.UUID
	Name     string
	Password string
	Username string
	Nickname string
}

type DeleteCustomerUser struct {
	UserId uuid.UUID
}

type DeleteAdminUser struct {
	UserId uuid.UUID
}

type CustomerInfoDetail struct {
	UserId         uuid.UUID
	Name           string
	ChannelName    string
	ChannelLink    string
	Email          string
	Mobile         string
	OrderableCount uint32
	PersonaLink    string
	OnedriveLink   string
	Memo           string
	SubscribeStart *time.Time
	SubscribeEnd   *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type AdminInfoData struct {
	UserId    uuid.UUID
	Name      string
	Nickname  string
	Email     string
	CreatedAt time.Time
}

type CustomerInfoData struct {
	UserId      uuid.UUID
	Name        string
	ChannelName string
	ChannelLink string
	Email       string
	Mobile      string
	CreatedAt   time.Time
}

type UserUseCase interface {
	SignInUser(ctx context.Context, in SignInUser) (string, error)

	CreateCustomerUser(ctx context.Context, in CreateCustomerUser) (uuid.UUID, error)
	CreateAdminUser(ctx context.Context, in CreateAdminUser) (uuid.UUID, error)

	UpdateCustomerUser(ctx context.Context, in UpdateCustomerUser) error
	UpdateAdminPassword(ctx context.Context, in UpdateAdminPassword) error
	UpdateAdminInfo(ctx context.Context, in UpdateAdminInfo) error
	ForceUpdateAdminInfo(ctx context.Context, in ForceUpdateAdminInfo) error

	DeleteCustomerUser(ctx context.Context, in DeleteCustomerUser) error
	DeleteAdminUser(ctx context.Context, in DeleteAdminUser) error

	GetCustomerInfoDetailByUserId(ctx context.Context, userId uuid.UUID) (CustomerInfoDetail, error)
	FetchAllAdmin(ctx context.Context, option FetchAdminOption) ([]AdminInfoData, error)
	FetchAllCustomer(ctx context.Context, option FetchCustomerOption) ([]CustomerInfoData, error)
}

type TokenGenerateAdapter interface {
	Generate(User) (string, error)
}