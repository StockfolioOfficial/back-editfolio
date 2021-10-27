package di

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/core/app"
	"github.com/stockfolioofficial/back-editfolio/core/config"
	"github.com/stockfolioofficial/back-editfolio/core/di/scope"
	"github.com/stockfolioofficial/back-editfolio/helloworld/handler"
	handler3 "github.com/stockfolioofficial/back-editfolio/order/handler"
	handler4 "github.com/stockfolioofficial/back-editfolio/orderState/handler"
	handler2 "github.com/stockfolioofficial/back-editfolio/user/handler"
)

func OnStart(
	e *echo.Echo,
	mw middlewares,
	helloWorld *handler.HelloWorldController,
	user *handler2.UserController,
	order *handler3.OrderController,
	orderState *handler4.OrderStateController,
) app.OnStart {
	return func() error {
		logLevel := log.ErrorLevel
		if config.IsDebug {
			logLevel = log.TraceLevel
		}
		log.SetLevel(logLevel)

		// global middleware set
		e.Use(mw...)

		// routing
		bindEcho(
			e,
			helloWorld,
			user,
			order,
			orderState,
		)
		return nil
	}
}

func bindEcho(e *echo.Echo, binders ...scope.EchoBinder) {
	for i := range binders {
		binders[i].Bind(e)
	}
}

func OnClose() app.OnClose {
	return func() {

	}
}
