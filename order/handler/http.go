package handler

import (
	"net/http"
	"time"

	"github.com/stockfolioofficial/back-editfolio/core/debug"
	"github.com/stockfolioofficial/back-editfolio/util/echox"

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

type OrderStateRequest struct {
	State uint8 `json:"-" query:"state" example:1`
}

type OrderStateResponse struct {
	OrderId     uuid.UUID  `json:"orderId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderedAt   time.Time  `json:"orderedAt" example:"2021-10-27 12:00"`
	Orderer     uuid.UUID  `json:"orderer" validate:"required" example:"437ae8fe-2349-4125-b4a3-b3154b63e8dc"`
	Assignee    *uuid.UUID `json:"orderer" validate:"required" example:"13aa33a3-1832-8819-k41d-jlkl490dfkjl"`
	Requirement *string    `json:"requirement" example:"예쁘게 잘 편집해 주세요~"`
}

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

// @Security Auth-Jwt-Bearer
// @Summary [어드민] 제작 의뢰 요청 목록
// @Description 제작 의뢰 요청 목록 API 기능
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

// @Security Auth-Jwt-Bearer
// @Summary [어드민] 제작 의뢰 진행중 목록
// @Description 제작 의뢰 진행중 목록 API 기능
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

// @Security Auth-Jwt-Bearer
// @Summary [어드민] 제작 의뢰 완료된 목록
// @Description 제작 의뢰 완료된 목록 API 기능
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


func (c *OrderController) Bind(e *echo.Echo) {
	//edit order request
	e.POST("/order", echox.UserID(c.createOrder), debug.JwtBypassOnDebug())

	// v1
	e.GET("/order/ready", c.fetchOrderToReady, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))
	e.GET("/order/processing", echox.UserID(c.fetchOrderToProcessing), debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))
	e.GET("/order/done", c.fetchOrderToDone, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))

}
