package usecase

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/google/uuid"

	"github.com/stockfolioofficial/back-editfolio/domain"
)

func NewOrderUseCase(
	orderRepo domain.OrderRepository,
	userRepo domain.UserRepository,
	managerRepo  domain.ManagerRepository,
	customerRepo domain.CustomerRepository,
	orderStateRepo domain.OrderStateRepository,
	timeout time.Duration,
) domain.OrderUseCase {
	return &ucase{
		orderRepo:      orderRepo,
		userRepo:       userRepo,
		managerRepo:    managerRepo,
		customerRepo:   customerRepo,
		orderStateRepo: orderStateRepo,
		timeout:        timeout,
	}
}

type ucase struct {
	orderRepo      domain.OrderRepository
	userRepo       domain.UserRepository
	managerRepo    domain.ManagerRepository
	customerRepo   domain.CustomerRepository
	orderStateRepo domain.OrderStateRepository
	timeout        time.Duration
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

func (u *ucase) Fetch(ctx context.Context, option domain.FetchOrderOption) (res []domain.OrderInfo, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	list, err := u.orderRepo.Fetch(c, option)
	res = make([]domain.OrderInfo, len(list))

	bufSize := int(float64(len(list)) * 0.7)

	statesIds := make([]uint8, 0, bufSize)
	managerIds := make([]uuid.UUID, 0, bufSize)
	customerIds := make([]uuid.UUID, 0, bufSize)

	stateDst := make(map[uint8]*domain.OrderInfo)
	managerDst := make(map[uuid.UUID]*domain.OrderInfo)
	customerDst := make(map[uuid.UUID]*domain.OrderInfo)
	for i := range list {
		src := list[i]
		res[i] = domain.OrderInfo{
			OrderId:            src.Id,
			OrderedAt:          src.OrderedAt,
			DoneAt:             src.DoneAt,
		}

		dst := &res[i]

		stateDst[src.State] = dst
		statesIds = append(statesIds, src.State)

		managerDst[src.Orderer] = dst
		customerIds = append(customerIds, src.Orderer)

		if src.Assignee != nil {
			customerDst[*src.Assignee] = dst
			managerIds = append(managerIds, *src.Assignee)
		}
	}

	g, gc := errgroup.WithContext(c)
	g.Go(func() error {
		mList, err := u.managerRepo.FetchByIds(gc, managerIds)
		if err != nil {
			return err
		}

		for i := range mList {
			src := mList[i]
			dst, ok := managerDst[src.Id]
			if !ok {
				continue
			}

			dst.AssigneeName = &src.Name
			dst.AssigneeNickname = &src.Nickname
		}

		return nil
	})

	g.Go(func() error {
		cList, err := u.customerRepo.FetchByIds(gc, customerIds)
		if err != nil {
			return err
		}

		for i := range cList {
			src := cList[i]
			dst, ok := customerDst[src.Id]
			if !ok {
				continue
			}

			dst.OrdererName = src.Name
			dst.OrdererChannelName = src.ChannelName
			dst.OrdererChannelLink = src.ChannelLink
		}

		return nil
	})
	g.Go(func() error {
		sList, err := u.orderStateRepo.FetchByIds(gc, statesIds)
		if err != nil {
			return err
		}

		for i := range sList {
			src := sList[i]
			dst, ok := stateDst[src.Id]
			if !ok {
				continue
			}

			dst.OrderState = src.Id
			dst.OrderStateContent = src.Content
		}

		return nil
	})

	err = g.Wait()
	if err != nil { // 에러 일때 empty list 리턴
		res = []domain.OrderInfo{}
	}

	return
}