package app

import (
	"github.com/labstack/echo/v4"
	_ "github.com/stockfolioofficial/back-editfolio/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

type OnStart func() error
type OnClose func()

type App interface {
	Start() error
}

func NewApp(
	e *echo.Echo,
	db *gorm.DB,
	onStart OnStart,
	onClose OnClose,
) (res App) {
	res = &app{
		e: e,
		db: db,
		onStart: onStart,
		onClose: onClose,
	}
	return
}

type app struct {
	e *echo.Echo
	db *gorm.DB
	onStart OnStart
	onClose OnClose
}

func (a *app) Start() (err error) {
	err = a.onStart()
	if err != nil {
		return
	}
	defer a.onClose()

	var e = a.e

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	err = e.Start(":8000")
	return
}

