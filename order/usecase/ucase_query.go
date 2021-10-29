package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/pointer"
	"github.com/stockfolioofficial/back-editfolio/util/safe"
	"golang.org/x/sync/errgroup"
)

func (u *ucase) GetRecentProcessingOrder(ctx context.Context, userId uuid.UUID) (res domain.RecentOrderInfo, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	order, err := u.orderRepo.GetRecentByOrdererId(c, userId)
	if err != nil {
		return
	}

	if order == nil || order.IsDone() {
		err = domain.ErrItemNotFound
		return
	}

	res = domain.RecentOrderInfo{
		OrderId:            order.Id,
		OrderedAt:          order.OrderedAt,
		DueDate:            order.DueDate,
		OrderState:         order.State,
		OrderStateContent:  "알 수 없는 상태", // todo string resource
		RemainingEditCount: order.RemainingEditCount(),
	}

	g, gc := errgroup.WithContext(c)
	g.Go(func() (err error) {
		if order.Assignee == nil {
			return
		}

		assignee, err := u.managerRepo.GetById(gc, *order.Assignee)
		if err != nil {
			return
		}

		if assignee != nil {
			res.AssigneeNickname = &assignee.Nickname
		} else {
			res.AssigneeNickname = pointer.String("알 수 없는 편집자") // todo string resource
		}
		return
	})

	g.Go(func() (err error) {
		state, err := u.orderStateRepo.GetById(gc, order.State)
		if err != nil {
			return
		}

		if state != nil {
			res.OrderStateContent = state.Content
		}
		return
	})
	err = g.Wait()
	if err != nil {
		res = domain.RecentOrderInfo{}
		return
	}

	return
}

func (u *ucase) GetOrderDetailInfo(ctx context.Context, orderId uuid.UUID) (res domain.OrderDetailInfo, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	order, err := u.orderRepo.GetById(c, orderId)
	if err != nil {
		return
	}

	if order == nil {
		err = domain.ErrItemNotFound
		return
	}

	res = domain.OrderDetailInfo{
		OrderId:            order.Id,
		OrderedAt:          order.OrderedAt,
		Orderer:            order.Orderer,
		DueDate:            order.DueDate,
		AssigneeInfo:       nil,
		OrderState:         order.State,
		OrderStateContent:  "알 수 없는 상태", // todo string resource
		RemainingEditCount: order.RemainingEditCount(),
		Requirement:        safe.StringOrZero(order.Requirement),
	}

	g, gc := errgroup.WithContext(c)
	g.Go(func() (err error) {
		if order.Assignee == nil {
			return
		}

		assignee, err := u.managerRepo.GetById(gc, *order.Assignee)
		if err != nil {
			return
		}

		if assignee != nil {
			res.AssigneeInfo = &domain.OrderAssigneeInfo{
				Id:       assignee.Id,
				Name:     assignee.Name,
				Nickname: assignee.Nickname,
			}
		} else {
			res.AssigneeInfo = &domain.OrderAssigneeInfo{
				Id:       assignee.Id,
				Name:     "알 수 없는 편집자", // todo string resource
				Nickname: "알 수 없는 편집자", // todo string resource
			}
		}
		return
	})
	g.Go(func() (err error) {
		state, err := u.orderStateRepo.GetById(gc, order.State)
		if err != nil {
			return
		}

		if state != nil {
			res.OrderStateContent = state.Content
		}
		return
	})
	err = g.Wait()
	if err != nil {
		return
	}

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
	if err != nil {
		res = []domain.OrderInfo{}
	}

	return
}