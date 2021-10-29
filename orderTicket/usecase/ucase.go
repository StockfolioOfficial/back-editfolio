package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"golang.org/x/sync/errgroup"
	"time"
)

func NewOrderTicketUseCase(
	orderTicketRepo domain.OrderTicketRepository,
	userRepo domain.UserRepository,
	timeout time.Duration,
) domain.OrderTicketUseCase {
	return &ucase{
		orderTicketRepo: orderTicketRepo,
		userRepo:        userRepo,
		timeout:         timeout,
	}
}

type ucase struct {
	orderTicketRepo domain.OrderTicketRepository
	userRepo        domain.UserRepository
	timeout         time.Duration
}

func (u *ucase) CreateSubscribeTicket(ctx context.Context, in domain.CreateSubscribeTicket) (ticketId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var userId uuid.UUID
	g, gc := errgroup.WithContext(c)
	g.Go(func() (err error) {
		exists, err := u.userRepo.GetByUsername(gc, in.Username)
		if err != nil {
			return
		}

		if exists == nil {
			err = domain.ErrItemNotFound
			return
		}

		userId = exists.Id
		return
	})
	g.Go(func() (err error) {
		exists, err := u.orderTicketRepo.GetByExOrderId(ctx, in.ExOrderId)
		if err != nil {
			return
		}

		if exists != nil {
			err = domain.ErrItemAlreadyExist
		}

		return
	})
	err = g.Wait()
	if err != nil {
		return
	}

	ticket, err :=  u.orderTicketRepo.GetEndByOwnerId(c, userId)
	if err != nil {
		return
	}

	var (
		startAt = time.Now()
		endAt time.Time
	)
	if ticket != nil && ticket.EndAt != nil && ticket.EndAt.After(startAt) {
		startAt = *ticket.EndAt
	}

	switch in.Unit {
	case domain.SubscribeUnitMonth:
		monthEndDay := time.Date(startAt.Year(), startAt.Month() + time.Month(in.Value) + 1, 0,
			startAt.Hour(), startAt.Minute(), startAt.Second(),
			startAt.Nanosecond(), startAt.Location(),
		)
		if startAt.Day() < monthEndDay.Day() {
			endAt = startAt.AddDate(0, int(in.Value), 0)
		} else {
			endAt = monthEndDay
		}
	case domain.SubscribeUnitDay:
		endAt = startAt.AddDate(0, 0, int(in.Value))
	}

	newTicket := domain.CreateOrderTicket(domain.CreateOrderTicketOption{
		ExOrderId:       in.ExOrderId,
		OwnerId:         userId,
		TotalOrderCount: in.OrderCount,
		EditCount:       in.EditCount,
		StartAt:         &startAt,
		EndAt:           &endAt,
	})

	err = u.orderTicketRepo.Save(c, &newTicket)
	if err != nil {
		return
	}

	ticketId = newTicket.Id
	return
}

