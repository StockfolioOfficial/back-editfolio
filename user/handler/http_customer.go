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
	UserId              uuid.UUID  `json:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name                string     `json:"name" validate:"required" example:"나 고객"`
	SubscribeStart      *time.Time `json:"subscribeStart" example:"2021-10-27T04:44:18+00:00"`
	SubscribeEnd        *time.Time `json:"subscribeEnd" example:"2021-10-27T04:44:18+00:00"`
	RemainingOrderCount uint8      `json:"remainingOrderCount" validate:"required" example:"4"`
	OnedriveLink        string     `json:"onedriveLink" validate:"required" example:"(대충 링크)"`

	// Simple notification :
	// * NONE - 없음
	// * NEED_BUY_SUBSCRIBE - 구독권 구매 필요
	// * NEED_BUY_ONE_EDIT - 1회 편집권 구매필요
	SimpleNotify CustomerSimpleNotify `json:"simpleNotify" example:"NONE" enums:"NONE,NEED_BUY_SUBSCRIBE,NEED_BUY_ONE_EDIT"`
} // @name CustomerSimpleInfoResponse

// @Tags (User) 고객 기능
// @Security Auth-Jwt-Bearer
// @Summary [고객] 내 정보 가져오기
// @Description 내 정보 가져오는 기능, 역할(role)이 'CUSTOMER' 이여야함
// @Accept json
// @Produce json
// @Success 200 {object} CustomerSimpleInfoResponse "성공"
// @Router /customer/me [get]
func (c *UserController) getMyCustomerInfo(ctx echo.Context, userId uuid.UUID) error {
	out, err := c.useCase.CustomerSubscribeInfoByUserId(ctx.Request().Context(), userId)
	if err != nil {
		log.WithError(err).
			WithField("in", userId).
			Error(tag, "getMyCustomerInfo, unhandled error useCase.GetCustomerInfoDetailByUserId")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

	res := CustomerSimpleInfoResponse{
		UserId:              out.UserId,
		Name:                out.Name,
		SubscribeStart:      out.SubscribeStart,
		SubscribeEnd:        out.SubscribeEnd,
		RemainingOrderCount: out.RemainingOrderCount,
		OnedriveLink:        out.OnedriveLink,
		SimpleNotify:        CustomerSimpleNotifyNone,
	}

	if out.RemainingOrderCount == 0 {
		res.SimpleNotify = CustomerSimpleNotifyNeedBuyOneEdit
	}

	now := time.Now()
	if out.SubscribeStart == nil && out.SubscribeEnd == nil {
		res.SimpleNotify = CustomerSimpleNotifyNeedBuySubscribe
	} else if out.SubscribeStart != nil &&
		out.SubscribeEnd != nil &&
		now.After(*out.SubscribeEnd) {
		res.SimpleNotify = CustomerSimpleNotifyNeedBuySubscribe
	} else if out.SubscribeStart == nil || out.SubscribeEnd == nil {
		log.WithField("out", out).Error("의도 하지 않는 구독 일자")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

	return ctx.JSON(http.StatusOK, res)
}