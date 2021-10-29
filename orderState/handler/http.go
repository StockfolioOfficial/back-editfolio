package handler

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"net/http"
)

const (
	tag = "[ORDER-STATE] "
)

func NewOrderStateController(useCase domain.OrderStateUseCase) *OrderStateController {
	return &OrderStateController{useCase: useCase}
}

type OrderStateController struct {
	useCase domain.OrderStateUseCase
}

type OrderStateInfoResponse struct {
	Id      uint8  `json:"id" validate:"required" example:"1"`
	Content string `json:"content" validate:"required" example:"편집자 배정 중"`
} // @name OrderStateInfoResponse

type OrderStateInfoListResponse []OrderStateInfoResponse

func useCaseToOrderStateInfoResponse(src domain.OrderStateInfo) OrderStateInfoResponse {
	return OrderStateInfoResponse{
		Id:      src.Id,
		Content: src.Content,
	}
}

func useCaseToOrderStateInfoListResponse(list []domain.OrderStateInfo) (res OrderStateInfoListResponse) {
	res = make(OrderStateInfoListResponse, len(list))
	for i := range list {
		res[i] = useCaseToOrderStateInfoResponse(list[i])
	}

	return
}

// @Tags 기타
// @Summary 제작 상태 목록 전부
// @Description 제작 상태 목록 전부 가져오는 기능
// @Accept json
// @Produce json
// @Success 200 {object} OrderStateInfoListResponse true "성공"
// @Router /order/state/full [get]
func (c *OrderStateController) fetchFull(ctx echo.Context) error {
	list, err := c.useCase.FetchFull(ctx.Request().Context())
	if err != nil{
		log.WithError(err).Error(tag, "fetch full, unhandled error useCase.FetchFull")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

	return ctx.JSON(http.StatusOK, useCaseToOrderStateInfoListResponse(list))
}

// @Tags 기타
// @Summary 제작 상태 서브 옵션 목록
// @Description 제작 상태 서브 옵션 목록 전부 가져오는 기능
// @Accept json
// @Produce json
// @Param order_state_id path int true "제작 상태 식별 아이디"
// @Success 200 {object} OrderStateInfoListResponse true "성공"
// @Success 204 "값이 없음"
// @Router /order/state/{order_state_id}/sub [get]
func (c *OrderStateController) fetchSub(ctx echo.Context) error {
	var req struct {
		OrderStateId uint8 `param:"orderStateId"`
	}

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "internalCreateTicket data binding error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	}

	list, err := c.useCase.FetchByParentId(ctx.Request().Context(), req.OrderStateId)
	if err != nil{
		log.WithError(err).Error(tag, "fetchSub, unhandled error useCase.FetchByParentId")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

	if len(list) == 0 {
		return ctx.NoContent(http.StatusNoContent)
	}

	return ctx.JSON(http.StatusOK, useCaseToOrderStateInfoListResponse(list))
}

func (c *OrderStateController) Bind(e *echo.Echo) {
	e.GET("/order/state/full", c.fetchFull)
	e.GET("/order/state/:orderStateId/sub", c.fetchSub)
}

