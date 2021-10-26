package di

import (
	"time"

	repository2 "github.com/stockfolioofficial/back-editfolio/manager/repository"

	"github.com/google/wire"
	"github.com/stockfolioofficial/back-editfolio/core/app"
	"github.com/stockfolioofficial/back-editfolio/core/config"
	repository3 "github.com/stockfolioofficial/back-editfolio/customer/repository"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/helloworld/handler"
	"github.com/stockfolioofficial/back-editfolio/user/adapter"
	handler2 "github.com/stockfolioofficial/back-editfolio/user/handler"
	"github.com/stockfolioofficial/back-editfolio/user/repository"
	"github.com/stockfolioofficial/back-editfolio/user/usecase"
)

var DI = wire.NewSet(
	app.NewApp,
	infraSet,
	adapterSet,
	repositorySet,
	useCaseSet,
	handlerSet,
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
)

var useCaseSet = wire.NewSet(
	usecase.NewUserUseCase,
)

var handlerSet = wire.NewSet(
	handler.NewHelloWorldHttpHandler,
	handler2.NewUserHttpHandler,
)

var lifecycleSet = wire.NewSet(
	OnStart,
	OnClose,
)
