package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

const (
	tag = "[USER] "
)

func NewUserHttpHandler(useCase domain.UserUseCase) *HttpHandler {
	return &HttpHandler{useCase: useCase}
}

type HttpHandler struct {
	useCase domain.UserUseCase
}

type CreateCustomerRequest struct {
	// Name, 길이 2~60 제한
	Name string `json:"name" validate:"required,min=2,max=60" example:"ljs"`

	// Email, 이메일 주소
	Email string `json:"email" validate:"required,email" example:"example@example.com"`

	// Mobile, 형식 : 01012345678
	Mobile string `json:"mobile" validate:"required,sf_mobile" example:"01012345678"`
} // @name CreateCustomerUserRequest

type CreatedCustomerResp struct {
	Id uuid.UUID `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name CreatedCustomerResponse

// @Summary 고객 유저 생성
// @Description 고객 유저를 생성하는 기능
// @Accept json
// @Produce json
// @Param customerUserBody body CreateCustomerRequest true "Customer User Body"
// @Success 201 {object} CreatedCustomerResp
// @Router /user/customer [post]
func (h *HttpHandler) createCustomer(ctx echo.Context) error {
	var req CreateCustomerRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "create customer, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	newId, err := h.useCase.CreateCustomerUser(ctx.Request().Context(), domain.CreateCustomerUser{
		Name:   req.Name,
		Email:  req.Email,
		Mobile: req.Mobile,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusCreated, CreatedCustomerResp{Id: newId})
	default:
		log.WithError(err).Error(tag, "create customer, unhandled error useCase.CreateCustomerUser")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

type DeleteCustomerRequest struct {
	// Id, 유저 Id
	Id uuid.UUID `param:"userId" json:"-" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} //@name DeleteCustomerRequest

// @Summary 고객 유저 삭제
// @Description 고객 유저를 삭제하는 기능
// @Accept json
// @Produce json
// @Param customerUserBody body DeleteCustomerRequest true "Customer User Body"
// @Success 204
// @Router /user/customer/:userId [delete]
func (h *HttpHandler) deleteCustomerUser(ctx echo.Context) error {
	var req DeleteCustomerRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "delete customer, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.DeleteCustomerUser(ctx.Request().Context(), domain.DeleteCustomerUser{
		Id: req.Id,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusNoContent, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "delete customer")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

func (h *HttpHandler) Bind(e *echo.Echo) {
	//CRUD, customer or admin
	e.POST("/user/customer", h.createCustomer)

	//sign, auth
	e.POST("/user/sign", h.signInUser)

	e.DELETE("/user/customer/:userId", h.deleteCustomerUser)
}
