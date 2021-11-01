package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"golang.org/x/sync/errgroup"
	"time"
)

func (u *ucase) FetchAllAdmin(ctx context.Context, option domain.FetchAdminOption) (res []domain.AdminInfoData, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	list, err := u.userRepo.FetchAllAdmin(c, option)
	if err != nil {
		return
	}

	res = make([]domain.AdminInfoData, len(list))
	for i := range list {
		src := list[i]
		if src.Manager == nil {
			res = []domain.AdminInfoData{}
			err = errors.New("join failed manager info data")
			return
		}
		res[i] = domain.AdminInfoData{
			UserId:    src.Id,
			Name:      src.Manager.Name,
			Nickname:  src.Manager.Nickname,
			Email:     src.Username,
			CreatedAt: src.CreatedAt,
		}
	}

	return
}

func (u *ucase) FetchAllCustomer(ctx context.Context, option domain.FetchCustomerOption) (res []domain.CustomerInfoData, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	list, err := u.userRepo.FetchAllCustomer(c, option)
	if err != nil {
		return
	}

	res = make([]domain.CustomerInfoData, len(list))
	for i := range list {
		src := list[i]
		if src.Customer == nil {
			res = []domain.CustomerInfoData{}
			err = errors.New("join failed customer info data")
			return
		}
		res[i] = domain.CustomerInfoData{
			UserId:      src.Id,
			Name:        src.Customer.Name,
			ChannelName: src.Customer.ChannelName,
			ChannelLink: src.Customer.ChannelLink,
			Email:       src.Customer.Email,
			Mobile:      src.Customer.Mobile,
			CreatedAt:   src.CreatedAt,
		}
	}

	return
}

func (u *ucase) GetAdminInfoDetailByUserId(ctx context.Context, userId uuid.UUID) (res domain.AdminInfoDetailData, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetByIdWithManager(c, userId)
	if err != nil {
		return
	}

	if user == nil {
		err = domain.ErrItemNotFound
		return
	}

	if user.Manager == nil {
		err = errors.New("join failed manager info data")
		return
	}

	res = domain.AdminInfoDetailData{
		UserId:    uuid.UUID{},
		Username:  user.Username,
		Name:      user.Manager.Name,
		Nickname:  user.Manager.Nickname,
		CreatedAt: user.CreatedAt,
	}

	return
}


func (u *ucase) GetCustomerInfoDetailByUserId(ctx context.Context, userId uuid.UUID) (res domain.CustomerInfoDetailData, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	detail, err := u.userRepo.GetByIdWithCustomer(c, userId)
	if err != nil {
		return
	}

	if detail == nil {
		err = domain.ErrItemNotFound
		return
	}

	if detail.Customer == nil {
		err = errors.New("join failed customer info data")
		return
	}

	res = domain.CustomerInfoDetailData{
		UserId:         detail.Id,
		Name:           detail.Customer.Name,
		ChannelName:    detail.Customer.ChannelName,
		ChannelLink:    detail.Customer.ChannelLink,
		Email:          detail.Customer.Email,
		Mobile:         detail.Customer.Mobile,
		PersonaLink:    detail.Customer.PersonaLink,
		OnedriveLink:   detail.Customer.OnedriveLink,
		Memo:           detail.Customer.Memo,
		CreatedAt:      detail.CreatedAt,
		UpdatedAt:      detail.UpdatedAt,
	}
	return
}

func (u *ucase) CustomerSubscribeInfoByUserId(ctx context.Context, userId uuid.UUID) (res domain.CustomerSubscribeInfoData, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	g, gc := errgroup.WithContext(c)
	g.Go(func() (err error) {
		exists, err := u.userRepo.GetByIdWithCustomer(gc, userId)
		if err != nil {
			return
		}

		if exists == nil {
			err = domain.ErrItemNotFound
			return
		}

		if exists.Customer == nil {
			err = errors.New("join failed customer info data")
			return
		}


		res.UserId = userId
		res.Name = exists.Customer.Name
		res.OnedriveLink = exists.Customer.OnedriveLink
		return
	})
	g.Go(func() (err error) {
		ticket, err := u.orderTicketRepo.GetByOwnerIdBetweenStartAndEnd(gc, userId, time.Now())
		if err != nil {
			return
		}

		if ticket != nil {
			res.SubscribeStart = ticket.StartAt
			res.SubscribeEnd = ticket.EndAt
			res.RemainingOrderCount = ticket.RemainingOrderCount()
		}

		return
	})
	err = g.Wait()
	if err != nil {
		res = domain.CustomerSubscribeInfoData{}
		return
	}

	return
}
