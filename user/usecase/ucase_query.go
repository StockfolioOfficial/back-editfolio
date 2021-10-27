package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
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

func (u *ucase) GetCustomerInfoDetailByUserId(ctx context.Context, userId uuid.UUID) (domain.CustomerInfoDetail, error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	detail, err := u.userRepo.GetByIdWithCustomer(c, userId)
	if err != nil {
		return domain.CustomerInfoDetail{}, err
	}

	if detail.Customer == nil {
		return domain.CustomerInfoDetail{}, errors.New("join failed customer info data")
	}

	return domain.CustomerInfoDetail{
		UserId:         detail.Id,
		Name:           detail.Customer.Name,
		ChannelName:    detail.Customer.ChannelName,
		ChannelLink:    detail.Customer.ChannelLink,
		Email:          detail.Customer.Email,
		Mobile:         detail.Customer.Mobile,
		OrderableCount: detail.Customer.OrderableCount,
		PersonaLink:    detail.Customer.PersonaLink,
		OnedriveLink:   detail.Customer.OnedriveLink,
		Memo:           detail.Customer.Memo,
		SubscribeStart: detail.Customer.SubscribeStart,
		SubscribeEnd:   detail.Customer.SubscribeEnd,
		CreatedAt:      detail.CreatedAt,
		UpdatedAt:      detail.UpdatedAt,
	}, nil
}