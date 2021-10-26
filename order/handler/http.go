package handler

import (
	"github.com/stockfolioofficial/back-editfolio/core/debug"
	"github.com/stockfolioofficial/back-editfolio/util/echox"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
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

type CreateOrderRequest struct {
	// Requirement, 요청사항
	Requirement string `json:"requirement" validate:"required,max=2000" example:"알잘딱깔센"`
} // @name CreateOrderRequest

type CreateOrderResponse struct {
	OrderId uuid.UUID `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name CreateOrderResponse

// @Security Auth-Jwt-Bearer
// @Summary [고객] 편집 의뢰 요청
// @Description 고객이 편집 의뢰를 요청하는 API 기능
// @Accept json
// @Produce json
// @Param createOrderBody body CreateOrderRequest true "편집 의뢰 요청 데이터 구조"
// @Success 201 {object} CreateOrderResponse true "의뢰 요청 성공"
// @Router /order [post]
func (c *OrderController) createOrder(ctx echo.Context, userId uuid.UUID) error {
	var req CreateOrderRequest
	err := ctx.Bind(&req)

	if err != nil {
		log.WithError(err).Trace(tag, "create order request, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	orderId, err := c.useCase.RequestOrder(ctx.Request().Context(), domain.RequestOrder{
		UserId:      userId,
		Requirement: req.Requirement,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusCreated, CreateOrderResponse{OrderId: orderId})
	case domain.ErrNoPermission:
		return ctx.JSON(http.StatusUnauthorized, domain.NoPermissionResponse)
	case domain.ItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "video requirement failed")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

func (c *OrderController) Bind(e *echo.Echo) {
	//edit order request
	e.POST("/order", echox.UserID(c.createOrder), debug.JwtBypassOnDebug())
}
