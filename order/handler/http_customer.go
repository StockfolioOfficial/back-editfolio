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
	OrderId uuid.UUID `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
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
		return ctx.JSON(http.StatusCreated, CreateOrderResponse{OrderId: orderId})
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
	// OrderId 주문 식별아이디 (UUID)
	OrderId            uuid.UUID  `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`

	// OrderedAt 주문 일자 (Datetime) RFC3339 datetime format
	OrderedAt          time.Time  `json:"orderedAt" validate:"required" example:"2021-10-27T04:44:18+00:00"`

	// DueDate 완료 예정일 (Date) RFC3339 datetime format
	DueDate            *time.Time `json:"dueDate" example:"2021-10-30T00:00:00+00:00"`

	// AssigneeNickname 담당 편집자 이름
	AssigneeNickname   *string    `json:"assigneeNickname" example:"담당 편집자 닉네임"`

	// OrderState 주문 상태 식별 번호
	OrderState         uint8      `json:"orderState" validate:"required" example:"3"`

	// OrderStateContent 주문 상태 명
	OrderStateContent  string     `json:"orderStateContent" validate:"required" example:"이펙트 추가 중"`

	// RemainingEditCount 남은 수정 횟수
	RemainingEditCount uint8      `json:"remainingEditCount" validate:"required" example:"2"`
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

	switch err {
	case nil:
		return ctx.JSON(http.StatusOK, RecentOrderInfoResponse{
			OrderId:            res.OrderId,
			OrderedAt:          res.OrderedAt,
			DueDate:            res.DueDate,
			AssigneeNickname:   res.AssigneeNickname,
			OrderState:         res.OrderState,
			OrderStateContent:  res.OrderStateContent,
			RemainingEditCount: res.RemainingEditCount,
		})
	case domain.ErrItemNotFound:
		return ctx.NoContent(http.StatusNoContent)
	default:
		log.WithError(err).Error(tag, "order done requirement failed")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

// @Tags (Order) 고객 기능
// @Security Auth-Jwt-Bearer
// @Summary [고객] 진행중인 편집 수정 의뢰
// @Description 고객이 진행중인 편집 수정 의뢰 기능, 역할(role)이 'CUSTOMER' 이여야함
// @Accept json
// @Success 202 "수정 요청 성공"
// @Router /order/recent-processing/edit [post]
func (c *OrderController) myOrderEdit(ctx echo.Context, userId uuid.UUID) error {
	err := c.useCase.RequestEditOrder(ctx.Request().Context(), domain.RequestEditOrder{
		UserId: userId,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusAccepted)
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
	case domain.ErrWeirdData:
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "empty remaining edit count"})
	case domain.ErrItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ErrorResponse{Message: "already requested edit"})
	case domain.ErrNoPermission:
		return ctx.JSON(http.StatusUnauthorized, domain.NoPermissionResponse)
	default:
		log.WithError(err).
			WithField("in", userId).
			Error(tag, "myOrderEdit, unhandled error useCase.RequestEditOrder")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}


type DoneOrderResponse struct {
	// OrderId 주문 식별아이디 (UUID)
	OrderId uuid.UUID `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name DoneOrderResponse

// @Tags (Order) 고객 기능
// @Security Auth-Jwt-Bearer
// @Summary [고객] 진행중인 편집 의뢰 완료
// @Description 고객이 진행중인 편집 의뢰 완료 기능, 역할(role)이 'CUSTOMER' 이여야함
// @Accept json
// @Success 200 {object} DoneOrderResponse true "의뢰 완료 요청 성공"
// @Router /order/recent-processing/done [post]
func (c *OrderController) myOrderDone(ctx echo.Context, userId uuid.UUID) error {

	orderId, err := c.useCase.OrderDone(ctx.Request().Context(), domain.OrderDone{
		UserId: userId,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusOK, DoneOrderResponse{OrderId: orderId})
	case domain.ErrNoPermission:
		return ctx.JSON(http.StatusUnauthorized, domain.NoPermissionResponse)
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "not exists order"})
	default:
		log.WithError(err).
			WithField("in", userId).
			Error(tag, "myOrderDone, unhandled error useCase.OrderDone")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}
