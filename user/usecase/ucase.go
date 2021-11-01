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
	orderTicketRepo domain.OrderTicketRepository,
	timeout time.Duration,
) domain.UserUseCase {
	return &ucase{
		userRepo:        userRepo,
		tokenAdapter:    tokenAdapter,
		managerRepo:     managerRepo,
		customerRepo:    customerRepo,
		orderTicketRepo: orderTicketRepo,
		timeout:         timeout,
	}
}

type ucase struct {
	userRepo        domain.UserRepository
	tokenAdapter    domain.TokenGenerateAdapter
	managerRepo     domain.ManagerRepository
	customerRepo    domain.CustomerRepository
	orderTicketRepo domain.OrderTicketRepository
	timeout         time.Duration
}

func (u *ucase) SignInUser(ctx context.Context, si domain.SignInUser) (token string, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetByUsername(c, si.Username)
	if err != nil {
		return
	}

	if user == nil {
		err = domain.ErrItemNotFound
		return
	}

	if user.ComparePassword(si.Password) {
		// token generate
		token, err = u.tokenAdapter.Generate(*user)
	} else {
		err = domain.ErrUserWrongPassword
	}

	return
}

func (u *ucase) CreateSuperAdminUser(ctx context.Context, in domain.CreateSuperAdminUser) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	//TODO 나중에 유저네임 이미 있는거 체크도 필요할듯
	exists, err := u.userRepo.ExistsSuperUser(c)
	if err != nil {
		return
	}
	if exists {
		err = domain.ErrItemAlreadyExist
		return
	}

	var user = createUser(domain.SuperAdminUserRole, in.Email, in.Password)
	var manager = domain.CreateManager(domain.ManagerCreateOption{
		User:     &user,
		Name:     in.Name,
		Nickname: in.Nickname,
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


func (u *ucase) CreateCustomerUser(ctx context.Context, in domain.CreateCustomerUser) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	exists, err := u.userRepo.GetByUsername(c, in.Email)
	if err != nil {
		return
	}
	if exists != nil {
		err = domain.ErrItemAlreadyExist
		return
	}

	var user = createUser(domain.CustomerUserRole, in.Email, in.Mobile)
	var customer = domain.CreateCustomer(domain.CustomerCreateOption{
		User:   &user,
		Name:   in.Name,
		Email:  in.Email,
		Mobile: in.Mobile,
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


func (u *ucase) CreateAdminUser(ctx context.Context, in domain.CreateAdminUser) (newId uuid.UUID, err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	email, err := u.userRepo.GetByUsername(c, in.Email)

	if email != nil {
		err = domain.ErrItemAlreadyExist
		return
	}

	var user = createUser(domain.AdminUserRole, in.Email, in.Password)
	var manager = domain.CreateManager(domain.ManagerCreateOption{
		User:     &user,
		Name:     in.Name,
		Nickname: in.Nickname,
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

func (u *ucase) UpdateCustomerUser(ctx context.Context, in domain.UpdateCustomerUser) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	exists, err := u.userRepo.GetByUsername(c, in.Email)
	if err != nil {
		return
	}

	var user *domain.User
	if exists != nil {
		if exists.Id == in.UserId {
			user = exists
		} else {
			err = domain.ErrItemAlreadyExist
			return
		}
	}

	if user == nil {
		user, err = u.userRepo.GetById(c, in.UserId)
		if err != nil {
			return
		}
	}

	if !domain.CheckUserAlive(user,
		domain.User.IsCustomer) {
		err = domain.ErrItemNotFound
		return
	}

	err = user.LoadCustomerInfo(c, u.customerRepo)
	if err != nil {
		return
	}

	user.UpdateCustomerInfo(
		in.Name,
		in.ChannelName,
		in.ChannelLink,
		in.Email,
		in.Mobile,
		in.PersonaLink,
		in.OnedriveLink,
		in.Memo,
	)

	return u.userRepo.Transaction(c, func(ur domain.UserTxRepository) error {
		g, gc := errgroup.WithContext(c)
		g.Go(func() error {
			return u.userRepo.Save(gc, user)
		})
		g.Go(func() error {
			return u.customerRepo.Save(gc, user.Customer)
		})
		return g.Wait()
	})
}

func (u *ucase) UpdateAdminPassword(ctx context.Context, in domain.UpdateAdminPassword) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetById(c, in.UserId)
	if !domain.CheckUserAlive(user,
		domain.User.IsAdmin,
		domain.User.IsSuperAdmin) {
		err = domain.ErrItemNotFound
		return
	}

	if !user.ComparePassword(in.OldPassword) {
		err = domain.ErrUserWrongPassword
		return
	}

	user.UpdatePassword(in.NewPassword)
	return u.userRepo.Save(c, user)
}

func (u *ucase) UpdateAdminInfo(ctx context.Context, in domain.UpdateAdminInfo) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	exists, err := u.userRepo.GetByUsername(c, in.Username)
	if err != nil {
		return
	}

	var user *domain.User
	if exists != nil {
		if exists.Id == in.UserId {
			user = exists
		} else {
			err = domain.ErrItemAlreadyExist
			return
		}
	}

	if user == nil {
		user, err = u.userRepo.GetById(c, in.UserId)
		if err != nil {
			return
		}
	}

	if !domain.CheckUserAlive(user,
		domain.User.IsAdmin,
		domain.User.IsSuperAdmin) {
		err = domain.ErrItemNotFound
		return
	}

	err = user.LoadManagerInfo(c, u.managerRepo)
	if err != nil {
		return
	}

	user.UpdateManagerInfo(in.Username, in.Name, in.Nickname)
	return u.userRepo.Transaction(c, func(ur domain.UserTxRepository) error {
		g, gc := errgroup.WithContext(c)
		g.Go(func() error {
			return u.userRepo.Save(gc, user)
		})
		g.Go(func() error {
			return u.managerRepo.Save(gc, user.Manager)
		})
		return g.Wait()
	})
}

func (u *ucase) ForceUpdateAdminInfo(ctx context.Context, in domain.ForceUpdateAdminInfo) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	exists, err := u.userRepo.GetByUsername(c, in.Username)
	if err != nil {
		return
	}

	var user *domain.User
	if exists != nil {
		if exists.Id == in.UserId {
			user = exists
		} else {
			err = domain.ErrItemAlreadyExist
			return
		}
	}

	if user == nil {
		user, err = u.userRepo.GetById(c, in.UserId)
		if err != nil {
			return
		}
	}

	if !domain.CheckUserAlive(user,
		domain.User.IsAdmin,
		domain.User.IsSuperAdmin) {
		err = domain.ErrItemNotFound
		return
	}

	err = user.LoadManagerInfo(c, u.managerRepo)
	if err != nil {
		return
	}

	user.UpdatePassword(in.Password)
	user.UpdateManagerInfo(in.Username, in.Name, in.Nickname)
	return u.userRepo.Transaction(c, func(ur domain.UserTxRepository) error {
		g, gc := errgroup.WithContext(c)
		g.Go(func() error {
			return u.userRepo.Save(gc, user)
		})
		g.Go(func() error {
			return u.managerRepo.Save(gc, user.Manager)
		})
		return g.Wait()
	})
}

func (u *ucase) DeleteCustomerUser(ctx context.Context, in domain.DeleteCustomerUser) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetById(c, in.UserId)
	if err != nil {
		return
	}

	if !domain.CheckUserAlive(user, domain.User.IsCustomer) {
		err = domain.ErrItemNotFound
		return
	}

	user.Delete()
	return u.userRepo.Save(c, user)
}

func (u *ucase) DeleteAdminUser(ctx context.Context, in domain.DeleteAdminUser) (err error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.GetById(c, in.UserId)

	if !domain.CheckUserAlive(user, domain.User.IsAdmin) {
		err = domain.ErrItemNotFound
		return
	}

	user.Delete()
	return u.userRepo.Save(c, user)
}

func createUser(role domain.UserRole, username, password string) (user domain.User) {
	user = domain.CreateUser(domain.UserCreateOption{
		Role:     role,
		Username: username,
	})

	user.UpdatePassword(password)
	return
}