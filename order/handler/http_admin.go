package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
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
	res, alreadyResp, err := c.internalFetchOrder(ctx, domain.OrderGeneralStateReady, nil)
	if alreadyResp {
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
	if alreadyResp {
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
	if alreadyResp {
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

type orderDetailAssigneeInfoResponse struct {
	Id       uuid.UUID `json:"assignee" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name     string    `json:"assigneeName" example:"담당 편집자 이름"`
	Nickname string    `json:"assigneeNickname" example:"담당 편집자 닉네임"`
} // @name OrderDetailAssigneeInfoResponse

type OrderDetailInfoResponse struct {
	OrderId            uuid.UUID                        `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderedAt          time.Time                        `json:"orderedAt" validate:"required" example:"2021-10-27T04:44:18+00:00"`
	Orderer            uuid.UUID                        `json:"orderer" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	DueDate            *time.Time                       `json:"dueDate" example:"2021-10-30T00:00:00+00:00"`
	Assignee           *orderDetailAssigneeInfoResponse `json:"assignee"`
	OrderState         uint8                            `json:"orderState" validate:"required" example:"3"`
	OrderStateContent  string                           `json:"orderStateContent" validate:"required" example:"이펙트 추가 중"`
	RemainingEditCount uint8                            `json:"remainingEditCount" validate:"required" example:"2"`
	Requirement        string                           `json:"requirement"`
} // @name OrderDetailInfoResponse

// @Tags (Order) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 의뢰 상세 정보
// @Description 의뢰 상세 정보 가져오는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param order_id path string true "의뢰 식별 아이디(UUID)"
// @Success 200 {object} OrderDetailInfoResponse true "정보 가져오기 완료"
// @Router /order/{order_id} [get]
func (c *OrderController) getOrderDetailInfo(ctx echo.Context) error {

	var req struct {
		OrderId uuid.UUID `json:"-" param:"orderId"`
	}
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "get order detail info, request data bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	res, err := c.useCase.GetOrderDetailInfo(ctx.Request().Context(), req.OrderId)

	var assignee *orderDetailAssigneeInfoResponse
	if res.AssigneeInfo != nil {
		assignee = &orderDetailAssigneeInfoResponse{
			Id:       res.AssigneeInfo.Id,
			Name:     res.AssigneeInfo.Name,
			Nickname: res.AssigneeInfo.Nickname,
		}
	}

	return ctx.JSON(http.StatusOK, OrderDetailInfoResponse{
		OrderId:            res.OrderId,
		OrderedAt:          res.OrderedAt,
		Orderer:            res.Orderer,
		DueDate:            res.DueDate,
		Assignee:           assignee,
		OrderState:         res.OrderState,
		OrderStateContent:  res.OrderStateContent,
		RemainingEditCount: res.RemainingEditCount,
		Requirement:        res.Requirement,
	})
}

type UpdateOrderInfoRequest struct {
	OrderId    uuid.UUID `json:"-" param:"orderId" validate:"required" example:"150e8400-p11y-41d4-a716-446655440000"`
	DueDate    time.Time `json:"dueDate" validate:"required" example:"2021-10-30T00:00:00+00:00"`
	Assignee   uuid.UUID `json:"assignee" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderState uint8     `json:"orderState" validate:"required" example:"3"`
} // @name UpdateOrderInfoRequest

// @Tags (Order) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 의뢰 정보 수정
// @Description 의뢰 정보 수정하는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param order_id path string true "의뢰 식별 아이디(UUID)"
// @Param requestBody body UpdateOrderInfoRequest true "편집 의뢰 요청 데이터 구조"
// @Success 204 "정보 수정 완료"
// @Router /order/{order_id} [put]
func (c *OrderController) updateOrderInfo(ctx echo.Context) error {
	var req UpdateOrderInfoRequest
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "update order, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = c.useCase.UpdateOrderInfo(ctx.Request().Context(), domain.UpdateOrderInfo{
		OrderId:    req.OrderId,
		DueDate:    req.DueDate,
		Assignee:   req.Assignee,
		OrderState: req.OrderState,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusNoContent)
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
	case domain.ErrWeirdData:
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	default:
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

type OrderAssignSelfResponse struct {
	OrderId uuid.UUID `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// @Tags (Order) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 의뢰 나에게 업무 할당
// @Description 업무 나에게 할당 하는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param order_id path string true "의뢰 식별 아이디(UUID)"
// @Success 200 {object} OrderAssignSelfResponse true "수주 완료"
// @Router /order/{order_id}/assign-self [post]
func (c *OrderController) orderAssignSelf(ctx echo.Context, userId uuid.UUID) error {
	var req struct {
		OrderId uuid.UUID `json:"-" param:"orderId"`
	}
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "orderAssignSelf data binding error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	var in = domain.OrderAssignSelf{
		OrderId:  req.OrderId,
		Assignee: userId,
	}
	err = c.useCase.OrderAssignSelf(ctx.Request().Context(), in)

	switch err {
	case nil:
		return ctx.JSON(http.StatusOK, OrderAssignSelfResponse{
			OrderId: req.OrderId,
		})
	case domain.ErrItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ErrorResponse{Message: "assign conflict"})
	case domain.ErrNoPermission:
		return ctx.JSON(http.StatusUnauthorized, domain.NoPermissionResponse)
	default:
		log.WithError(err).
			WithField("in", in).
			Error(tag, "orderAssignSelf / unhandled error useCase.OrderAssignSelf")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}