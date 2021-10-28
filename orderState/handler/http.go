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
}

type OrderStateInfoListResponse []OrderStateInfoResponse

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

	res := make(OrderStateInfoListResponse, len(list))
	for i := range list {
		src := list[i]
		res[i] = OrderStateInfoResponse{
			Id:      src.Id,
			Content: src.Content,
		}
	}

	return ctx.JSON(http.StatusOK, res)
}

func (c *OrderStateController) Bind(e *echo.Echo) {
	e.GET("/order/state/full", c.fetchFull)
}


