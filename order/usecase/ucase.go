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
		err = domain.ItemAlreadyExist
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

func (u *ucase) FetchOrderList(ctx context.Context, stmt uint8) (orderList []domain.Order, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	orderList, err = u.orderRepo.GetByOrderList(c, stmt)
	if err != nil {
		return
	}

	// orders에 아무 것도 조회가 안되면 No content ? or Not Found ? or OK ?
	return
}
