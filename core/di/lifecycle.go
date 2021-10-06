package di

import (
	"github.com/labstack/echo/v4"
	"github.com/stockfolioofficial/back-editfolio/core/app"
	"github.com/stockfolioofficial/back-editfolio/core/di/scope"
	"github.com/stockfolioofficial/back-editfolio/helloworld/handler"
	handler2 "github.com/stockfolioofficial/back-editfolio/user/handler"
)

func OnStart(
	e *echo.Echo,
	mw middlewares,
	helloWorld *handler.HttpHandler,
	user *handler2.HttpHandler,
) app.OnStart {
	return func() error {
		// global middleware set
		e.Use(mw...)

		// routing
		bindEcho(
			e,
			helloWorld,
			user,
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

