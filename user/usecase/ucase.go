package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"time"
)

func NewUserUseCase(
	userRepo domain.UserRepository,
	timeout time.Duration,
) domain.UserUseCase {
	return &ucase{
		userRepo: userRepo,
		timeout:  timeout,
	}
}

type ucase struct {
	userRepo domain.UserRepository
	timeout  time.Duration
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


func (u *ucase) SignInUser(ctx context.Context, si domain.SignInUser) (jwt string, err error) {
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
		jwt = "generate jwt"
	} else {
		err = domain.UserWrongPassword
	}

	return
}


