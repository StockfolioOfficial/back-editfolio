package di

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type echoBindWithValidate struct {
	echo.DefaultBinder
}

func (e *echoBindWithValidate) Bind(i interface{}, c echo.Context) (err error) {
	err = e.DefaultBinder.Bind(i, c)
	if err != nil {
		return
	}

	return c.Validate(i)
}

type echoValidator struct {
	v *validator.Validate
}

func (e *echoValidator) Validate(i interface{}) error {
	var wrapper struct {
		Value interface{} `validator:"dive"`
	}
	wrapper.Value = i
	return e.v.Struct(&wrapper)
}

func NewEcho() (e *echo.Echo) {
	e = echo.New()
	e.Binder = &echoBindWithValidate{}
	e.Validator = &echoValidator{v: newValidator()}
	return
}

type middlewares []echo.MiddlewareFunc

func NewMiddleware() middlewares {
	return middlewares{
		middleware.CORSWithConfig(middleware.CORSConfig{
			// todo debug 추후 production 모드일때 스크립트 형태로 외부에서 주입 받는 기능 추가 필요
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"*"},
			AllowMethods: []string{"*"},
		}),
	}
}
