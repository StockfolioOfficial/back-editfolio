package usecase

import (
	"context"
	"errors"
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
	orderTicketRepo domain.OrderTicketRepository,
	timeout time.Duration,
) domain.OrderUseCase {
	return &ucase{
		orderRepo:       orderRepo,
		userRepo:        userRepo,
		managerRepo:     managerRepo,
		customerRepo:    customerRepo,
		orderStateRepo:  orderStateRepo,
		orderTicketRepo: orderTicketRepo,
		timeout:         timeout,
	}
}

type ucase struct {
	orderRepo       domain.OrderRepository
	userRepo        domain.UserRepository
	managerRepo     domain.ManagerRepository
	customerRepo    domain.CustomerRepository
	orderStateRepo  domain.OrderStateRepository
	orderTicketRepo domain.OrderTicketRepository
	timeout         time.Duration
}

func (u *ucase) RequestOrder(ctx context.Context, in domain.RequestOrder) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var defaultState uint8 = 1
	g, gc := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		exists, err := u.userRepo.GetById(gc, in.UserId)
		if err != nil {
			return
		}
		if !domain.CheckUserAlive(exists, domain.User.IsCustomer) {
			err = domain.ErrNoPermission
		}

		return
	})
	g.Go(func() (err error) {
		exists, _ := u.orderRepo.GetRecentByOrdererId(gc, in.UserId)
		if exists != nil && !exists.IsDone() {
			err = domain.ErrItemAlreadyExist
		}

		return
	})
	g.Go(func() error {
		exists, _ := u.orderStateRepo.GetByCode(gc, domain.OrderStateCodeDefault)
		if exists != nil {
			defaultState = exists.Id
		}

		return nil
	})
	err = g.Wait()
	if err != nil {
		return
	}

	err = u.orderTicketRepo.Transaction(c, func(otr domain.OrderTicketTxRepository) (err error) {
		orderOption := domain.CreateOrderOption{
			Orderer:     in.UserId,
			State:       defaultState,
		}
		if len(in.Requirement) > 0 {
			orderOption.Requirement = &in.Requirement
		}

		or := u.orderRepo.With(otr)

		ticket, err := otr.GetByOwnerIdBetweenStartAndEnd(c, in.UserId, time.Now())
		if err != nil {
			return
		}

		if ticket == nil || ticket.IsEmptyOrderCount() {
			return errors.New("no ticket") // todo error handling
		}

		ticket.UseOrder()
		orderOption.EditCount = ticket.EditCount
		order := domain.CreateOrder(orderOption)

		g, gc = errgroup.WithContext(c)
		g.Go(func() error {
			return otr.Save(gc, ticket)
		})
		g.Go(func() error {
			return or.Save(gc, &order)
		})
		err = g.Wait()
		if err != nil {
			return
		}

		newId = order.Id
		return
	})

	return
}

func (u *ucase) RequestEditOrder(ctx context.Context, in domain.RequestEditOrder) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var (
		order *domain.Order
		state *domain.OrderState
	)
	g, gc := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		exists, err := u.userRepo.GetById(gc, in.UserId)
		if err != nil {
			return
		}

		if !domain.CheckUserAlive(exists, domain.User.IsCustomer) {
			err = domain.ErrNoPermission
		}

		return
	})
	g.Go(func() (err error) {
		order, err = u.orderRepo.GetRecentByOrdererId(c, in.UserId)
		if err != nil {
			return
		}

		if order == nil {
			err = domain.ErrItemNotFound
			return
		}

		if order.IsEmptyEditCount() {
			err = domain.ErrWeirdData
			return
		}

		return
	})
	g.Go(func() (err error) {
		state, _ = u.orderStateRepo.GetByCode(gc, domain.OrderStateCodeRequestEdit)
		if state == nil {
			err = errors.New("orderStateRepo.GetByCode domain.OrderStateCodeRequestEdit not exists state")
		}
		return
	})
	err = g.Wait()
	if err != nil {
		return
	}

	if order.State == state.Id {
		err = domain.ErrItemAlreadyExist
		return
	}
	order.UseEdit()
	order.State = state.Id
	err = u.orderRepo.Save(c, order)
	return
}

func (u *ucase) OrderDone(ctx context.Context, in domain.OrderDone) (orderId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var (
		order *domain.Order
		state *domain.OrderState
	)
	g, gc := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		exists, err := u.userRepo.GetById(gc, in.UserId)
		if err != nil {
			return
		}

		if !domain.CheckUserAlive(exists, domain.User.IsCustomer) {
			err = domain.ErrNoPermission
		}

		return
	})
	g.Go(func() (err error) {
		order, err = u.orderRepo.GetRecentByOrdererId(c, in.UserId)
		if err != nil {
			return
		}

		if order == nil || order.IsDone() {
			err = domain.ErrItemNotFound
		}

		order.Done()
		return
	})
	g.Go(func() (err error) {
		state, _ = u.orderStateRepo.GetByCode(gc, domain.OrderStateCodeDone)
		if state == nil {
			err = errors.New("orderStateRepo.GetByCode domain.OrderStateCodeDone not exists state")
		}
		return
	})
	err = g.Wait()
	if err != nil {
		return
	}

	order.State = state.Id
	err = u.orderRepo.Save(c, order)
	if err != nil {
		return
	}

	orderId = order.Id
	return
}

func (u *ucase) UpdateOrderInfo(ctx context.Context, in domain.UpdateOrderInfo) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	order, err := u.orderRepo.GetById(c, in.OrderId)
	if err != nil {
		return err
	}

	if order == nil {
		err = domain.ErrItemNotFound
		return
	}


	var (
		aExists *domain.Manager
		sExists *domain.OrderState
	)

	g, gc := errgroup.WithContext(c)
	g.Go(func() (err error) {
		aExists, err = u.managerRepo.GetById(gc, in.Assignee)
		if err != nil {
			return err
		}

		if aExists == nil {
			err = domain.ErrWeirdData
		}
		return
	})

	g.Go(func() (err error) {
		if order.Assignee == nil {
			sExists, err = u.orderStateRepo.GetByCode(gc, domain.OrderStateCodeTake)
		} else {
			sExists, err = u.orderStateRepo.GetById(gc, in.OrderState)
		}
		if err != nil {
			return err
		}

		if sExists == nil {
			err = domain.ErrWeirdData
		}
		return
	})

	err = g.Wait()
	if err != nil {
		return
	}

	order.DueDate = &in.DueDate
	order.Assignee = &in.Assignee
	if sExists == nil {
		order.State = in.OrderState
	} else {
		order.State = sExists.Id
	}

	return u.orderRepo.Save(c, order)
}


func (u *ucase) OrderAssignSelf(ctx context.Context, in domain.OrderAssignSelf) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var (
		order *domain.Order
		state *domain.OrderState
	)
	g, gc := errgroup.WithContext(c)
	g.Go(func() (err error) {
		order, err = u.orderRepo.GetById(gc, in.OrderId)
		if err != nil {
			return
		}

		if order.Assignee != nil {
			err = domain.ErrItemAlreadyExist
			return
		}

		order.Assignee = &in.Assignee
		return
	})
	g.Go(func() (err error) {
		user, err := u.userRepo.GetById(gc, in.Assignee)
		if err != nil {
			return
		}

		if !domain.CheckUserAlive(user,
			domain.User.IsAdmin,
			domain.User.IsSuperAdmin) {
			err = domain.ErrNoPermission
		}

		return
	})
	g.Go(func() (err error) {
		state, _ = u.orderStateRepo.GetByCode(gc, domain.OrderStateCodeTake)
		if state == nil {
			err = errors.New("orderStateRepo.GetByCode domain.OrderStateCodeTake not exists state")
		}
		return
	})
	err = g.Wait()
	if err != nil {
		return
	}

	order.State = state.Id
	err = u.orderRepo.Save(c, order)
	return
}