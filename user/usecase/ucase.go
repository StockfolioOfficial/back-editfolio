package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

func NewUserUseCase(
	userRepo domain.UserRepository,
	tokenAdapter domain.TokenGenerateAdapter,
	timeout time.Duration,
) domain.UserUseCase {
	return &ucase{
		userRepo:     userRepo,
		tokenAdapter: tokenAdapter,
		timeout:      timeout,
	}
}

type ucase struct {
	userRepo     domain.UserRepository
	tokenAdapter domain.TokenGenerateAdapter
	timeout      time.Duration
}

func (u *ucase) CreateCustomerUser(ctx context.Context, cu domain.CreateCustomerUser) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var user = domain.CreateUser(domain.UserCreateOption{
		Role:     domain.CustomerUserRole,
		Username: cu.Email,
	})
	user.UpdatePassword(cu.Mobile)
	err = u.userRepo.Transaction(c, func(userRepo domain.UserTxRepository) error {
		return userRepo.Save(c, &user)
		//TODO customer 테이블 만들어서 연결필요
	})

	newId = user.Id

	return
}

func (u *ucase) SignInUser(ctx context.Context, si domain.SignInUser) (token string, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetByUsername(c, si.Username)
	if err != nil {
		return
	}

	if user == nil {
		err = domain.ItemNotFound
		return
	}

	if user.ComparePassword(si.Password) {
		// token generate
		token, err = u.tokenAdapter.Generate(*user)
	} else {
		err = domain.UserWrongPassword
	}

	return
}

func (u *ucase) DeleteCustomerUser(ctx context.Context, du domain.DeleteCustomerUser) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetById(c, du.Id)

	if user == nil || !user.IsUserRole() {
		err = domain.ItemNotFound
		return
	}

	user.Delete()
	return u.userRepo.Save(ctx, user)
}
