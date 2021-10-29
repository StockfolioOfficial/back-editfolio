package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

const (
	tag = "[ORDER-TICKET] "
)

func NewOrderTicketController(useCase domain.OrderTicketUseCase) *OrderTicketController {
	return &OrderTicketController{useCase: useCase}
}

type OrderTicketController struct {
	useCase domain.OrderTicketUseCase
}

func (c *OrderTicketController) Bind(e *echo.Echo) {
	e.POST("/internal/order/ticket", c.internalCreateTicket)
}