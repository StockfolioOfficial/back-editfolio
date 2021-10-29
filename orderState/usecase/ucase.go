package usecase

import (
	"context"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"time"
)

func NewOrderStateUseCase(
	orderStateRepo domain.OrderStateRepository,
	timeout time.Duration,
) domain.OrderStateUseCase {
	return &ucase{
		orderStateRepo: orderStateRepo,
		timeout:        timeout,
	}
}

type ucase struct {
	orderStateRepo domain.OrderStateRepository
	timeout time.Duration
}

func (u *ucase) FetchByParentId(ctx context.Context, parentId uint8) (res []domain.OrderStateInfo, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	list, err := u.orderStateRepo.FetchByParentId(c, parentId)
	if err != nil {
		return
	}

	res = domainToOrderStateInfoList(list)
	return
}

func (u *ucase) FetchFull(ctx context.Context) (res []domain.OrderStateInfo, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	list, err := u.orderStateRepo.FetchFull(c)
	if err != nil {
		return
	}

	res = domainToOrderStateInfoList(list)
	return
}

func domainToOrderStateInfo(src domain.OrderState) domain.OrderStateInfo {
	return domain.OrderStateInfo{
		Id:      src.Id,
		Content: src.Content,
	}
}

func domainToOrderStateInfoList(list []domain.OrderState) (res []domain.OrderStateInfo) {
	res = make([]domain.OrderStateInfo, len(list))
	for i := range list {
		res[i] = domainToOrderStateInfo(list[i])
	}

	return
}