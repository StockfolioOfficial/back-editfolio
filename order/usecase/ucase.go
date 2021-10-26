package usecase

import (
	"context"
	"time"

	"github.com/stockfolioofficial/back-editfolio/domain"
)

type ocase struct {
	orderRepo domain.OrderRepository
	timeout   time.Duration
}

func (o *ocase) VideoEditRequirement(ctx context.Context, vr domain.VideoEditRequirement) (err error) {
	c, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	order, err := o.orderRepo.GetById(c, vr.Id)
	if order == nil {
		err = domain.ItemNotFound
		return
	}

	err = order.LoadOrderInfo(c, o.orderRepo)
	if err != nil {
		return
	}

	order.UpdateVideoEditRequirement(vr.Requirement)
	return o.orderRepo.Save(c, order)
}
