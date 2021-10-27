package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stockfolioofficial/back-editfolio/util/pointer"
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
	return ctx.JSON(http.StatusOK, CustomerSimpleInfoResponse{
		UserId:         userId,
		Name:           "더미 이름",
		SubscribeStart: pointer.Time(time.Now()),
		SubscribeEnd:   pointer.Time(time.Now().AddDate(0, 1, 0)),
		OrderableCount: 3,
		SimpleNotify:   CustomerSimpleNotifyNone,
	})
}