package usecase

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/google/uuid"

	"github.com/stockfolioofficial/back-editfolio/domain"
)

func NewOrderUseCase(
	orderRepo domain.OrderRepository,
	userRepo domain.UserRepository,
	managerRepo domain.ManagerRepository,
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

	if !domain.ExistsCustomer(user) {
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

func (u *ucase) Fetch(ctx context.Context, option domain.FetchOrderOption) (res []domain.OrderInfo, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	list, err := u.orderRepo.Fetch(c, option)
	res = make([]domain.OrderInfo, len(list))

	bufSize := int(float64(len(list)) * 0.7)

	statesIds := make([]uint8, 0, bufSize)
	managerIds := make([]uuid.UUID, 0, bufSize)
	customerIds := make([]uuid.UUID, 0, bufSize)

	stateDst := make(map[uint8][]*domain.OrderInfo)
	managerDst := make(map[uuid.UUID][]*domain.OrderInfo)
	customerDst := make(map[uuid.UUID][]*domain.OrderInfo)
	for i := range list {
		src := list[i]
		res[i] = domain.OrderInfo{
			OrderId:   src.Id,
			OrderedAt: src.OrderedAt,
			DoneAt:    src.DoneAt,
		}

		dst := &res[i]

		stateDst[src.State] = append(stateDst[src.State], dst)
		statesIds = append(statesIds, src.State)

		customerDst[src.Orderer] = append(customerDst[src.Orderer], dst)
		customerIds = append(customerIds, src.Orderer)

		if src.Assignee != nil {
			managerDst[*src.Assignee] = append(managerDst[*src.Assignee], dst)
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
			dsts := managerDst[src.Id]
			if len(dsts) == 0 {
				continue
			}

			for i := range dsts {
				dst := dsts[i]
				dst.AssigneeName = &src.Name
				dst.AssigneeNickname = &src.Nickname
			}
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
			dsts := customerDst[src.Id]
			if len(dsts) == 0 {
				continue
			}

			for i := range dsts {
				dst := dsts[i]
				dst.OrdererName = src.Name
				dst.OrdererChannelName = src.ChannelName
				dst.OrdererChannelLink = src.ChannelLink
			}
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
			dsts := stateDst[src.Id]
			if len(dsts) == 0 {
				continue
			}

			for i := range dsts {
				dst := dsts[i]
				dst.OrderState = src.Id
				dst.OrderStateContent = src.Content
			}
		}

		return nil
	})

	err = g.Wait()
	if err != nil { // 에러 일때 empty list 리턴
		res = []domain.OrderInfo{}
	}

	return
}

func (u *ucase) MyOrderDone(ctx context.Context, ud domain.OrderDone) (orderId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetById(c, ud.UserId)

	if err != nil {
		return
	}

	if domain.ExistsCustomer(user) {
		err = domain.ErrNoPermission
		return
	}

	order, err := u.orderRepo.GetRecentByOrdererId(c, ud.UserId)
	order.Done()

	err = u.orderRepo.Save(c, order)
  return
}

func (u *ucase) UpdateOrderDetailInfo(ctx context.Context, uo *domain.UpdateOrderInfo) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var (
		oExists *domain.Order
		aExists *domain.Manager
		sExists *domain.OrderState
	)

	g, gc := errgroup.WithContext(c)
	g.Go(func() (err error) {
		oExists, err = u.orderRepo.GetById(gc, uo.OrderId)

		if err != nil {
			return err
		}

		if oExists == nil {
			err = domain.ErrItemNotFound
		}
		return
	})

	g.Go(func() (err error) {
		aExists, err = u.managerRepo.GetById(gc, uo.Assignee)

		if err != nil {
			return err
		}

		if aExists == nil {
			err = domain.ErrWeirdData
		}
		return
	})

	g.Go(func() (err error) {
		sExists, err = u.orderStateRepo.GetById(gc, uo.OrderState)

		if err != nil {
			return err
		}

		if sExists == nil {
			err = domain.ErrWeirdData
		}
		return
	})

	err = g.Wait()

	oExists.DueDate = &uo.DueDate
	oExists.Assignee = &uo.Assignee
	oExists.State = uo.OrderState

	return u.orderRepo.Save(c, oExists)
}

func (u *ucase) GetRecentProcessingOrder(ctx context.Context, userId uuid.UUID) (ro domain.RecentOrderInfo, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var order *domain.Order
	var assignee *domain.Manager
	var state *domain.OrderState

	g, gc := errgroup.WithContext(c)
	g.Go(func() (err error) {
		order, err = u.orderRepo.GetRecentByOrdererId(gc, userId)
		return
	})
	g.Go(func() (err error) {
		assignee, err = u.managerRepo.GetById(gc, *order.Assignee)
		return
	})
	g.Go(func() (err error) {
		state, err = u.orderStateRepo.GetById(gc, order.State)
		return
	})
	err = g.Wait()
	if err != nil {
		return
	}

	ro.AssigneeNickname = &assignee.Nickname
	ro.DueDate = order.DueDate
	ro.OrderId = order.Id
	ro.OrderState = order.State
	ro.OrderStateContent = state.Content
	ro.OrderedAt = order.OrderedAt
	ro.RemainingEditCount = uint8(order.EditTotal - order.EditCount)

	return
}
