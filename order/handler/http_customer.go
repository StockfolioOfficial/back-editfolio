package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

type CreateOrderRequest struct {
	// Requirement, 요청사항
	Requirement string `json:"requirement" validate:"required,max=2000" example:"알잘딱깔센"`
} // @name CreateOrderRequest

type CreateOrderResponse struct {
	Id uuid.UUID `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name CreateOrderResponse

// @Tags (Order) 고객 기능
// @Security Auth-Jwt-Bearer
// @Summary [고객] 편집 의뢰 요청
// @Description 고객이 편집 의뢰를 하는 기능, 역할(role)이 'CUSTOMER' 이여야함
// @Accept json
// @Produce json
// @Param requestBody body CreateOrderRequest true "편집 의뢰 요청 데이터 구조"
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
		return ctx.JSON(http.StatusCreated, CreateOrderResponse{Id: orderId})
	case domain.ErrNoPermission:
		return ctx.JSON(http.StatusUnauthorized, domain.NoPermissionResponse)
	case domain.ErrItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "video requirement failed")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

type RecentOrderInfoResponse struct {
	OrderId            uuid.UUID  `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderedAt          time.Time  `json:"orderedAt" validate:"required" example:"2021-10-27T04:44:18+00:00"`
	DueDate            *time.Time `json:"dueDate" example:"2021-10-30T00:00:00+00:00"`
	AssigneeNickname   *string    `json:"assigneeNickname" example:"담당 편집자 닉네임"`
	OrderState         uint8      `json:"orderState" validate:"required" example:"3"`
	OrderStateContent  string     `json:"orderStateContent" validate:"required" example:"이펙트 추가 중"`
	RemainingEditCount uint16     `json:"remainingEditCount" validate:"required" example:"2"`
} //@name RecentOrderInfoResponse

// @Tags (Order) 고객 기능
// @Security Auth-Jwt-Bearer
// @Summary [고객] 진행중인 최근 편집 의뢰 정보
// @Description 고객이 진행중인 최근 편집 의뢰 정보를 가져오는 기능, 역할(role)이 'CUSTOMER' 이여야함
// @Accept json
// @Success 200 {object} RecentOrderInfoResponse true "의뢰 정보 가져오기 완료"
// @Router /order/recent-processing [get]
func (c *OrderController) getRecentProcessingOrder(ctx echo.Context, userId uuid.UUID) error {

	res, err := c.useCase.GetRecentProcessingOrder(ctx.Request().Context(), userId)

	if err != nil {
		return domain.ErrItemNotFound
	}

	return ctx.JSON(http.StatusOK, RecentOrderInfoResponse{
		OrderedAt:          res.OrderedAt,
		DueDate:            res.DueDate,
		AssigneeNickname:   res.AssigneeNickname,
		OrderState:         res.OrderState,
		OrderStateContent:  res.OrderStateContent,
		RemainingEditCount: uint16(res.RemainingEditCount),
	})
}

type DoneOrderResponse struct {
	Id uuid.UUID `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name DoneOrderResponse

// @Tags (Order) 고객 기능
// @Security Auth-Jwt-Bearer
// @Summary [고객] 진행중인 편집 의뢰 완료
// @Description 고객이 진행중인 편집 의뢰 완료 기능, 역할(role)이 'CUSTOMER' 이여야함
// @Accept json
// @Success 200 {object} DoneOrderResponse true "의뢰 완료 요청 성공"
// @Router /order/recent-processing/done [post]
func (c *OrderController) myOrderDone(ctx echo.Context, userId uuid.UUID) error {
	//TODO 채우세요
	return ctx.JSON(http.StatusOK, DoneOrderResponse{Id: uuid.New()})
}
