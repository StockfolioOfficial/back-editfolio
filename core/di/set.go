package di

import (
	"github.com/google/wire"
	"github.com/stockfolioofficial/back-editfolio/core/app"
	"github.com/stockfolioofficial/back-editfolio/core/config"
	repository3 "github.com/stockfolioofficial/back-editfolio/customer/repository"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/helloworld/handler"
	repository2 "github.com/stockfolioofficial/back-editfolio/manager/repository"
	handler3 "github.com/stockfolioofficial/back-editfolio/order/handler"
	repository4 "github.com/stockfolioofficial/back-editfolio/order/repository"
	usecase2 "github.com/stockfolioofficial/back-editfolio/order/usecase"
	handler4 "github.com/stockfolioofficial/back-editfolio/orderState/handler"
	repository5 "github.com/stockfolioofficial/back-editfolio/orderState/repository"
	usecase3 "github.com/stockfolioofficial/back-editfolio/orderState/usecase"
	handler5 "github.com/stockfolioofficial/back-editfolio/orderTicket/handler"
	repository6 "github.com/stockfolioofficial/back-editfolio/orderTicket/repository"
	usecase4 "github.com/stockfolioofficial/back-editfolio/orderTicket/usecase"
	"github.com/stockfolioofficial/back-editfolio/user/adapter"
	handler2 "github.com/stockfolioofficial/back-editfolio/user/handler"
	"github.com/stockfolioofficial/back-editfolio/user/repository"
	"github.com/stockfolioofficial/back-editfolio/user/usecase"
	"time"
)

var DI = wire.NewSet(
	app.NewApp,
	infraSet,
	adapterSet,
	repositorySet,
	useCaseSet,
	controllerSet,
	lifecycleSet,
)

var infraSet = wire.NewSet(
	NewEcho,
	NewMiddleware,
	NewDatabase,

	// todo, 추후 별도로 config로 빼는게 좋을 듯
	// useCase timeout 3min
	wire.Value(time.Minute*3),
)

var adapterSet = wire.NewSet(
	wire.InterfaceValue(new(domain.TokenGenerateAdapter), adapter.NewTokenGenerateAdapter([]byte(config.JWTSecret))),
)

var repositorySet = wire.NewSet(
	repository.NewUserRepository,
	repository2.NewManagerRepository,
	repository3.NewCustomerRepository,
	repository4.NewOrderRepository,
	repository5.NewOrderStateRepository,
	repository6.NewOrderTicketRepository,
)

var useCaseSet = wire.NewSet(
	usecase.NewUserUseCase,
	usecase2.NewOrderUseCase,
	usecase3.NewOrderStateUseCase,
	usecase4.NewOrderTicketUseCase,
)

var controllerSet = wire.NewSet(
	handler.NewHelloWorldController,
	handler2.NewUserController,
	handler3.NewOrderController,
	handler4.NewOrderStateController,
	handler5.NewOrderTicketController,
)

var lifecycleSet = wire.NewSet(
	OnStart,
	OnClose,
)
