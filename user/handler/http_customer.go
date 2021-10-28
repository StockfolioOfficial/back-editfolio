package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"net/http"
	"time"
)

type CustomerSimpleNotify string

const (
	CustomerSimpleNotifyNone CustomerSimpleNotify = "NONE"
	CustomerSimpleNotifyNeedBuySubscribe CustomerSimpleNotify = "NEED_BUY_SUBSCRIBE"
	CustomerSimpleNotifyNeedBuyOneEdit CustomerSimpleNotify = "NEED_BUY_ONE_EDIT"
)

type CustomerSimpleInfoResponse struct {
	UserId         uuid.UUID            `json:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name           string               `json:"name" validate:"required" example:"나 고객"`
	SubscribeStart *time.Time           `json:"subscribeStart" example:"2021-10-27T04:44:18+00:00"`
	SubscribeEnd   *time.Time           `json:"subscribeEnd" example:"2021-10-27T04:44:18+00:00"`
	OrderableCount uint32               `json:"orderableCount" example:"4"`

	// Simple notification :
	// * NONE - 없음
	// * NEED_BUY_SUBSCRIBE - 구독권 구매 필요
	// * NEED_BUY_ONE_EDIT - 1회 편집권 구매필요
	SimpleNotify   CustomerSimpleNotify `json:"simpleNotify" example:"NONE" enums:"NONE,NEED_BUY_SUBSCRIBE,NEED_BUY_ONE_EDIT"`
} // @name CustomerSimpleInfoResponse

// @Tags (User) 고객 기능
// @Security Auth-Jwt-Bearer
// @Summary [고객] 내 간단 정보 가져오기
// @Description 내 간단 정보 가져오는 기능, 역할(role)이 'CUSTOMER' 이여야함
// @Accept json
// @Produce json
// @Success 200 {object} CustomerSimpleInfoResponse "성공"
// @Router /customer/me.simply [get]
func (c *UserController) getCustomerMyInfoSimply(ctx echo.Context, userId uuid.UUID) error {
	detail, err := c.useCase.GetCustomerInfoDetailByUserId(ctx.Request().Context(), userId)
	if err != nil {
		log.WithError(err).Error(tag, "get customer detail, unhandled error useCase.GetCustomerInfoDetailByUserId")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

	res := CustomerSimpleInfoResponse{
		UserId:         detail.UserId,
		Name:           detail.Name,
		SubscribeStart: detail.SubscribeStart,
		SubscribeEnd:   detail.SubscribeEnd,
		OrderableCount: detail.OrderableCount,
		SimpleNotify:   CustomerSimpleNotifyNone,
	}

	now := time.Now()
	if detail.SubscribeStart == nil && detail.SubscribeEnd == nil {
		res.SimpleNotify = CustomerSimpleNotifyNeedBuySubscribe
	} else if detail.SubscribeStart != nil &&
		detail.SubscribeEnd != nil &&
		now.After(*detail.SubscribeEnd) {
		res.SimpleNotify = CustomerSimpleNotifyNeedBuySubscribe
	} else if detail.SubscribeStart == nil || detail.SubscribeEnd == nil {
		log.WithField("detail", detail).Error("의도 하지 않는 구독 일자")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

	if detail.OrderableCount == 0 {
		res.SimpleNotify = CustomerSimpleNotifyNeedBuyOneEdit
	}

	return ctx.JSON(http.StatusOK, res)
}