package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/stockfolioofficial/back-editfolio/core/debug"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/echox"
)

const (
	tag = "[ORDER] "
)

func NewOrderHttpHandler(useCase domain.OrderUseCase) *OrderController {
	return &OrderController{useCase: useCase}
}

type OrderController struct {
	useCase domain.OrderUseCase
}

func (c *OrderController) Bind(e *echo.Echo) {
	//edit order request
	e.POST("/order", echox.UserID(c.createOrder), debug.JwtBypassOnDebug())

	// v1 - fetch
	e.GET("/order/ready", c.fetchOrderToReady, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))
	e.GET("/order/processing", echox.UserID(c.fetchOrderToProcessing), debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))
	e.GET("/order/done", c.fetchOrderToDone, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))

}
