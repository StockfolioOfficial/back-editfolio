package usecase

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

func NewUserUseCase(
	userRepo domain.UserRepository,
	tokenAdapter domain.TokenGenerateAdapter,
	managerRepo domain.ManagerRepository,
	customerRepo domain.CustomerRepository,
	timeout time.Duration,
) domain.UserUseCase {
	return &ucase{
		userRepo:     userRepo,
		tokenAdapter: tokenAdapter,
		managerRepo:  managerRepo,
		customerRepo: customerRepo,
		timeout:      timeout,
	}
}

type ucase struct {
	userRepo     domain.UserRepository
	tokenAdapter domain.TokenGenerateAdapter
	managerRepo  domain.ManagerRepository
	customerRepo domain.CustomerRepository
	timeout      time.Duration
}

func (u *ucase) CreateCustomerUser(ctx context.Context, cu domain.CreateCustomerUser) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	email, err := u.userRepo.GetByUsername(c, cu.Email)

	if email != nil {
		err = domain.ItemAlreadyExist
		return
	}

	var user = createUser(domain.CustomerUserRole, cu.Email, cu.Mobile)
	var customer = domain.CreateCustomer(domain.CustomerCreateOption{
		User:   &user,
		Name:   cu.Name,
		Email:  cu.Email,
		Mobile: cu.Mobile,
	})
	err = u.userRepo.Transaction(c, func(ur domain.UserTxRepository) error {
		mr := u.customerRepo.With(ur)
		g, gc := errgroup.WithContext(c)
		g.Go(func() error {
			return ur.Save(gc, &user)
		})
		g.Go(func() error {
			return mr.Save(gc, &customer)
		})
		return g.Wait()
	})
	newId = user.Id
	return
}

func (u *ucase) UpdateCustomerUser(ctx context.Context, cu domain.UpdateCustomerUser) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	exists, err := u.userRepo.GetByUsername(c, cu.Email)
	if err != nil {
		return
	}

	var user *domain.User
	if exists != nil {
		if exists.Id == cu.UserId {
			user = exists
		} else {
			err = domain.ItemAlreadyExist
			return
		}
	}

	if user == nil {
		user, err = u.userRepo.GetById(c, cu.UserId)
		if err != nil {
			return
		}
	}

	if !domain.ExistsCustomer(user) {
		err = domain.ItemNotFound
		return
	}

	err = user.LoadCustomerInfo(c, u.customerRepo)
	if err != nil {
		return
	}

	user.UpdateCustomerInfo(
		cu.Name,
		cu.ChannelName,
		cu.ChannelLink,
		cu.Email,
		cu.Mobile,
		cu.PersonaLink,
		cu.OnedriveLink,
		cu.Memo,
	)

	g, gc := errgroup.WithContext(c)
	g.Go(func() error {
		return u.userRepo.Save(gc, user)
	})
	g.Go(func() error {
		return u.customerRepo.Save(c, user.Customer)
	})
	return g.Wait()
}

func (u *ucase) UpdateAdminPassword(ctx context.Context, up domain.UpdateAdminPassword) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetById(c, up.UserId)
	if user == nil || user.IsDeleted() || !user.IsAdmin() {
		err = domain.ItemNotFound
		return
	}

	if !user.ComparePassword(up.OldPassword) {
		err = domain.UserWrongPassword
		return
	}

	user.UpdatePassword(up.NewPassword)
	return u.userRepo.Save(c, user)
}

func (u *ucase) SignInUser(ctx context.Context, si domain.SignInUser) (token string, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetByUsername(c, si.Username)
	if err != nil {
		return
	}

	if user == nil {
		err = domain.ItemNotFound
		return
	}

	if user.ComparePassword(si.Password) {
		// token generate
		token, err = u.tokenAdapter.Generate(*user)
	} else {
		err = domain.UserWrongPassword
	}

	return
}

func (u *ucase) DeleteCustomerUser(ctx context.Context, du domain.DeleteCustomerUser) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetById(c, du.Id)

	if user == nil || !user.IsCustomer() {
		err = domain.ItemNotFound
		return
	}

	user.Delete()
	return u.userRepo.Save(ctx, user)
}

func (u *ucase) CreateAdminUser(ctx context.Context, au domain.CreateAdminUser) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	email, err := u.userRepo.GetByUsername(c, au.Email)

	if email != nil {
		err = domain.ItemAlreadyExist
		return
	}

	var user = createUser(domain.AdminUserRole, au.Email, au.Password)
	var manager = domain.CreateManager(domain.ManagerCreateOption{
		User:     &user,
		Name:     au.Name,
		Nickname: au.Nickname,
	})

	err = u.userRepo.Transaction(c, func(ur domain.UserTxRepository) error {
		mr := u.managerRepo.With(ur)
		g, gc := errgroup.WithContext(c)
		g.Go(func() error {
			return ur.Save(gc, &user)
		})
		g.Go(func() error {
			return mr.Save(gc, &manager)
		})
		return g.Wait()
	})
	newId = user.Id
	return
}

func (u *ucase) loadManager(ctx context.Context, userId uuid.UUID) (*domain.User, error) {
	var user *domain.User
	var manager *domain.Manager

	g, c := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		user, err = u.userRepo.GetById(c, userId)
		return
	})
	g.Go(func() (err error) {
		manager, err = u.managerRepo.GetById(c, userId)
		return
	})
	err := g.Wait()
	if err != nil {
		return nil, err
	}

	user.Manager = manager
	return user, nil
}

func createUser(role domain.UserRole, username, password string) (user domain.User) {
	user = domain.CreateUser(domain.UserCreateOption{
		Role:     role,
		Username: username,
	})

	user.UpdatePassword(password)
	return
}
