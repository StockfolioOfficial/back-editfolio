package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/stockfolioofficial/back-editfolio/domain"
)

type ocase struct {
	orderRepo domain.OrderRepository
	userRepo  domain.UserRepository
	timeout   time.Duration
}

func (o *ocase) CreateOrder(ctx context.Context, or domain.OrderRequirement) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	user, err := o.userRepo.GetById(c, or.UserId)
	if err != nil {
		return
	}

	if user == nil || !user.IsCustomer() {
		err = domain.ItemNotFound
		return
	}

	var orderOption domain.CreateOrderOption

	if orderOption.Orderer.Id != user.Id {
		err = domain.ItemNotFound
		return
	}

	if len(or.Requirement) > 0 {
		orderOption.Requirement = &or.Requirement
	}

	order := domain.CreateOrder(orderOption)
	newId = order.Id
	err = o.orderRepo.Save(c, &order)
	return
	// resp status 201, response {
	// 	"orderId": "string,uuid"
	// }
	// order, err := o.orderRepo.GetById(c, or.Id)
	// if order == nil {
	// 	err = domain.ItemNotFound
	// 	return
	// }

	// err = order.LoadOrderInfo(c, o.orderRepo)
	// if err != nil {
	// 	return
	// }

	// order.UpdateVideoEditRequirement(or.Requirement)
	// return o.orderRepo.Save(c, order)
}
