package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/stockfolioofficial/back-editfolio/domain"
)

func NewOrderUseCase(
	orderRepo domain.OrderRepository,
	userRepo domain.UserRepository,
	timeout time.Duration,
) domain.OrderUseCase {
	return &ucase{
		orderRepo: orderRepo,
		userRepo:  userRepo,
		timeout:   timeout,
	}
}

type ucase struct {
	orderRepo domain.OrderRepository
	userRepo  domain.UserRepository
	timeout   time.Duration
}

func (u *ucase) RequestOrder(ctx context.Context, or domain.RequestOrder) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetById(c, or.UserId)
	if err != nil {
		return
	}

	if domain.ExistsCustomer(user) {
		err = domain.ErrNoPermission
		return
	}

	exists, err := u.orderRepo.GetRecentByOrdererId(c, or.UserId)
	if err != nil {
		return
	}

	if exists != nil {
		err = domain.ErrItemAlreadyExist
		return
	}

	var orderOption domain.CreateOrderOption
	orderOption.Orderer = *user
	if len(or.Requirement) > 0 {
		orderOption.Requirement = &or.Requirement
	}

	order := domain.CreateOrder(orderOption)
	newId = order.Id
	err = u.orderRepo.Save(c, &order)
	return
}
