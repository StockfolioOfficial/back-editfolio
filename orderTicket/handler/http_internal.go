package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"net/http"
)

func (c *OrderTicketController) internalCreateTicket(ctx echo.Context) error {
	var req struct {
		ExOrderId  string `json:"exOrderId" validate:"required"`
		Username   string `json:"username" validate:"required,email"`
		Value      uint16 `json:"value" validate:"required,max=30000"`
		Unit       string `json:"unit" validate:"required,eq=M|eq=D"`
		OrderCount uint8  `json:"orderCount" validate:"required,max=30"`
		EditCount  uint8  `json:"editCount" validate:"required,max=60"`
	}

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "internalCreateTicket data binding error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	}

	newId, err := c.useCase.CreateSubscribeTicket(ctx.Request().Context(), domain.CreateSubscribeTicket{
		ExOrderId:  req.ExOrderId,
		Username:   req.Username,
		Value:      req.Value,
		Unit:       domain.SubscribeUnit(req.Unit),
		OrderCount: req.OrderCount,
		EditCount:  req.EditCount,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusOK, echo.Map{
			"ticketId": newId,
		})
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: fmt.Sprintf("user=%s, not found", req.Username),
		})
	case domain.ErrItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ErrorResponse{
			Message: fmt.Sprintf("ex_order_id=%s, exists", req.ExOrderId),
		})
	default:
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}