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
		if exists != nil {
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

func (u *ucase) OrderDone(ctx context.Context, in domain.OrderDone) (orderId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var order *domain.Order
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
		}

		return
	})
	err = g.Wait()
	if err != nil {
		return
	}

	order.Done()
	err = u.orderRepo.Save(c, order)
	if err != nil {
		return
	}

	orderId = order.Id
	return
}

func (u *ucase) UpdateOrderInfo(ctx context.Context, in *domain.UpdateOrderInfo) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var (
		oExists *domain.Order
		aExists *domain.Manager
		sExists *domain.OrderState
	)

	g, gc := errgroup.WithContext(c)
	g.Go(func() (err error) {
		oExists, err = u.orderRepo.GetById(gc, in.OrderId)
		if err != nil {
			return err
		}

		if oExists == nil {
			err = domain.ErrItemNotFound
		}
		return
	})

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
		sExists, err = u.orderStateRepo.GetById(gc, in.OrderState)
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

	oExists.DueDate = &in.DueDate
	oExists.Assignee = &in.Assignee
	oExists.State = in.OrderState

	return u.orderRepo.Save(c, oExists)
}
