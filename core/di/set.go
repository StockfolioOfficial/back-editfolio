package di

import (
	"github.com/google/wire"
	"github.com/stockfolioofficial/back-editfolio/core/app"
	"github.com/stockfolioofficial/back-editfolio/helloworld/handler"
	handler2 "github.com/stockfolioofficial/back-editfolio/user/handler"
	"github.com/stockfolioofficial/back-editfolio/user/repository"
	"github.com/stockfolioofficial/back-editfolio/user/usecase"
	"time"
)

var DI = wire.NewSet(
	app.NewApp,
	infraSet,
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
	wire.Value(time.Minute * 3),
)

var repositorySet = wire.NewSet(
	repository.NewUserRepository,
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
