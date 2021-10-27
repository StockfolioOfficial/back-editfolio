package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"net/http"
	"time"
)

type OrderFetchRequest struct {
	Query        string `json:"-" query:"q"`
	ShowMyTicket bool   `json:"-" query:"smt" example:"false"`
} // @name OrderFetchRequest

type OrderReadyInfoResponse struct {
	OrderId            uuid.UUID `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderedAt          time.Time `json:"orderedAt" validate:"required"`
	OrdererName        string    `json:"ordererName" validate:"required"`
	OrdererChannelName string    `json:"ordererChannelName" validate:"required"`
	OrdererChannelLink string    `json:"ordererChannelLink" validate:"required"`
	OrderState         uint8     `json:"orderState" validate:"required"`
	OrderStateContent  string    `json:"orderStateContent" validate:"required"`
} // @name OrderReadyInfoResponse

type OrderReadyInfoListResponse []OrderReadyInfoResponse

// @Tags (Order) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 제작 의뢰 요청 목록
// @Description 제작 의뢰 요청 목록 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param q query string false "검색어"
// @Success 200 {object} OrderReadyInfoListResponse true "의뢰 요청 목록"
// @Router /order/ready [get]
func (c *OrderController) fetchOrderToReady(ctx echo.Context) error {
	res, alreadyResp, err := c.internalFetchOrder(ctx, domain.OrderGeneralStateDone, nil)
	if !alreadyResp {
		return err
	}

	resp := make(OrderReadyInfoListResponse, len(res))
	for i := range res {
		src := res[i]
		dst := &resp[i]
		*dst = OrderReadyInfoResponse{
			OrderId:            src.OrderId,
			OrderedAt:          src.OrderedAt,
			OrdererName:        src.OrdererName,
			OrdererChannelName: src.OrdererChannelName,
			OrdererChannelLink: src.OrdererChannelLink,
			OrderState:         src.OrderState,
			OrderStateContent:  src.OrderStateContent,
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

type OrderProcessingInfoResponse struct {
	OrderId            uuid.UUID `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderedAt          time.Time `json:"orderedAt" validate:"required"`
	OrdererName        string    `json:"ordererName" validate:"required"`
	OrdererChannelName string    `json:"ordererChannelName" validate:"required"`
	OrdererChannelLink string    `json:"ordererChannelLink" validate:"required"`
	OrderState         uint8     `json:"orderState" validate:"required"`
	OrderStateContent  string    `json:"orderStateContent" validate:"required"`
	AssigneeName       string    `json:"assigneeName" validate:"required"`
	AssigneeNickname   string    `json:"assigneeNickname" validate:"required"`
} // @name OrderProcessingInfoResponse

type OrderProcessingInfoListResponse []OrderProcessingInfoResponse

// @Tags (Order) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 제작 의뢰 진행중 목록
// @Description 제작 의뢰 진행중 목록 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param q query string false "검색어"
// @Param smt query boolean false "자기 업무만 보기"
// @Success 200 {object} OrderProcessingInfoListResponse true "진행중인 의뢰 목록"
// @Router /order/processing [get]
func (c *OrderController) fetchOrderToProcessing(ctx echo.Context, userId uuid.UUID) error {
	res, alreadyResp, err := c.internalFetchOrder(ctx, domain.OrderGeneralStateProcessing, &userId)
	if !alreadyResp {
		return err
	}

	resp := make(OrderProcessingInfoListResponse, len(res))
	for i := range res {
		src := res[i]
		dst := &resp[i]
		*dst = OrderProcessingInfoResponse{
			OrderId:            src.OrderId,
			OrderedAt:          src.OrderedAt,
			OrdererName:        src.OrdererName,
			OrdererChannelName: src.OrdererChannelName,
			OrdererChannelLink: src.OrdererChannelLink,
			OrderState:         src.OrderState,
			OrderStateContent:  src.OrderStateContent,
		}

		if src.AssigneeName == nil {
			dst.AssigneeName = "알 수 없음"
		} else {
			dst.AssigneeName = *src.AssigneeName
		}

		if src.AssigneeNickname == nil {
			dst.AssigneeNickname = "알 수 없음"
		} else {
			dst.AssigneeNickname = *src.AssigneeNickname
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

type OrderDoneInfoResponse struct {
	OrderId            uuid.UUID `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderedAt          time.Time `json:"orderedAt" validate:"required"`
	OrdererName        string    `json:"ordererName" validate:"required"`
	OrdererChannelName string    `json:"ordererChannelName" validate:"required"`
	OrdererChannelLink string    `json:"ordererChannelLink" validate:"required"`
	OrderState         uint8     `json:"orderState" validate:"required"`
	OrderStateContent  string    `json:"orderStateContent" validate:"required"`
} // @name OrderDoneInfoResponse

type OrderDoneInfoListResponse []OrderDoneInfoResponse

// @Tags (Order) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 제작 의뢰 완료된 목록
// @Description 제작 의뢰 완료된 목록 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param q query string false "검색어"
// @Success 200 {object} OrderDoneInfoListResponse true "완료 의뢰 목록"
// @Router /order/done [get]
func (c *OrderController) fetchOrderToDone(ctx echo.Context) error {
	res, alreadyResp, err := c.internalFetchOrder(ctx, domain.OrderGeneralStateDone, nil)
	if !alreadyResp {
		return err
	}

	resp := make(OrderDoneInfoListResponse, len(res))

	for i := range res {
		src := res[i]
		dst := &resp[i]
		*dst = OrderDoneInfoResponse{
			OrderId:            src.OrderId,
			OrderedAt:          src.OrderedAt,
			OrdererName:        src.OrdererName,
			OrdererChannelName: src.OrdererChannelName,
			OrdererChannelLink: src.OrdererChannelLink,
			OrderState:         src.OrderState,
			OrderStateContent:  src.OrderStateContent,
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (c *OrderController) internalFetchOrder(ctx echo.Context, state domain.OrderGeneralState, userId *uuid.UUID) (res []domain.OrderInfo, alreadyResp bool, err error) {
	var req OrderFetchRequest
	err = ctx.Bind(&req)
	if err != nil {
		alreadyResp = true
		log.WithError(err).Trace(tag, "fetch order request, request body bind error")
		err = ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if !req.ShowMyTicket {
		userId = nil
	}
	res, err = c.useCase.Fetch(ctx.Request().Context(), domain.FetchOrderOption{
		OrderState: state,
		Query:      req.Query,
		Assignee:   userId,
	})

	if err != nil {
		alreadyResp = true
		log.WithError(err).Error(tag, "fetch order, unhandled error useCase.Fetch")
		err = ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
		return
	}

	if len(res) == 0 {
		alreadyResp = true
		err = ctx.NoContent(http.StatusNoContent)
	}

	return
}
