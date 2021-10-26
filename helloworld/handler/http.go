package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewHelloWorldHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

type HttpHandler struct{}

func (h *HttpHandler) HelloWorld(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"hello": "world",
	})
}

func (h *HttpHandler) Bind(e *echo.Echo) {
	e.GET("/", h.HelloWorld)
}
