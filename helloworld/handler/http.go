package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewHelloWorldController() *HelloWorldController {
	return &HelloWorldController{}
}

type HelloWorldController struct{}

func (h *HelloWorldController) HelloWorld(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"hello": "world",
	})
}

func (h *HelloWorldController) Bind(e *echo.Echo) {
	e.GET("/", h.HelloWorld)
}
