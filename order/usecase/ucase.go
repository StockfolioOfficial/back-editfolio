package usecase

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/stockfolioofficial/back-editfolio/domain"
)

type ocase struct {
	orderRepo domain.OrderRepository
	timeout   time.Duration
}

func (o *ocase) VideoEditRequirement(ctx context.Context, vr domain.VideoEditRequirement) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	//exists, err :=o.userRepo.getbyId(vr.Id)
	//err handling
	//exists is customer,

	var orderOption domain.CreateOrderOption
	//orderOption.Orderer = exists
	if len(vr.Requirement) > 0 {
		orderOption.Requirement = &vr.Requirement
	}
	order := domain.CreateOrder(orderOption)
	newId = order.Id
	err =  o.orderRepo.Save(c, &order)
	return 
	//resp status 201, response {
	//	"orderId": "string,uuid"
	//}
	//order, err := o.orderRepo.GetById(c, vr.Id)
	//if order == nil {
	//	err = domain.ItemNotFound
	//	return
	//}
	//
	//err = order.LoadOrderInfo(c, o.orderRepo)
	//if err != nil {
	//	return
	//}
	//
	//order.UpdateVideoEditRequirement(vr.Requirement)
	//return o.orderRepo.Save(c, order)
}
