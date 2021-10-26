package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

const (
	tag = "[ORDER] "
)

func NewOrderHttpHandler(useCase domain.OrderUseCase) *HttpHandler {
	return &HttpHandler{useCase: useCase}
}

type HttpHandler struct {
	useCase domain.OrderUseCase
}

type CreateOrderRequest struct {
	// Id, 오더 Id
	Id uuid.UUID `param:"userId" json:"-" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Requirement, 요청사항
	Requirement string `json:"requirement" validate:"required,max=2000" example:"알잘딱깔센"`
} // @name CreateOrderRequest

// @Summary 편집 요청 사항 입력
// @Description 고객이 영상 편집 요청을 입력하는 기능
// @Accept json
// @Produce json
// @Param createOrderBody body CreateOrderRequest true "Create Order Requirement Body"
// @Success 201
// @Router /order [post]
func (h *HttpHandler) CreateOrder(ctx echo.Context) error {
	var req CreateOrderRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "create order request, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.VideoEditRequirement(ctx.Request().Context(), domain.OrderRequirement{
		OrderId:     req.Id,
		Requirement: req.Requirement,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusNoContent, domain.ErrorResponse{Message: err.Error()})
	case domain.ItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "video requirement failed")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

func (h *HttpHandler) Bind(e *echo.Echo) {
	//Video Edit Requirement
	e.POST("/order", h.CreateOrder)
}
