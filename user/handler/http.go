package handler

import (
	"net/http"

	"github.com/stockfolioofficial/back-editfolio/core/debug"

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

	// ChannelName, 길이 2~100 제한
	ChannelName string `json:"channelName" example:"밥굽남"`

	// ChannelLink, 길이 2048 제한
	ChannelLink string `json:"channelLink" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`

	// Email, 이메일 주소
	Email string `json:"email" validate:"required,email" example:"example@example.com"`

	// Mobile, 형식 : 01012345678
	Mobile string `json:"mobile" validate:"required,sf_mobile" example:"01012345678"`

	//OrderableCount, 주문 가능 횟수
	//OrderableCount int `json:"orderableCount" default: "4" example: "4"`

	//PersonaLink, 길이 2048 제한
	PersonaLink string `json:"personaLink" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`

	//OnedriveLink, 길이 2048 제한
	OnedriveLink string `json:"onedriveLink" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`

	//Memo, 형식 : text
	//Memo string `json:"Memo" example:"편집 잘 부탁 드립니다. 따로 요청 사항은 없어요~"`
} // @name CreateCustomerUserRequest

type CreatedCustomerResp struct {
	Id uuid.UUID `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name CreatedCustomerResponse

type UpdatePasswordRequest struct {
	UserId string `json:"-" header:"User-Id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Mobile, 형식 : 01012345678
	OldPassword string `json:"oldPassword" validate:"required" example:"01012345678"`
	NewPassword string `json:"newPassword" validate:"required" example:"01087654321"`
} // @name UpdatePasswordRequest

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

	newId, err := h.useCase.CreateCustomerUser(ctx.Request().Context(), domain.CreateCustomerInformation{
		Name:        req.Name,
		ChannelName: req.ChannelName,
		ChannelLink: req.ChannelLink,
		Email:       req.Email,
		Mobile:      req.Mobile,
		//OrderableCount: req.OrderableCount,
		PersonaLink:  req.PersonaLink,
		OnedriveLink: req.OnedriveLink,
		//Memo:           req.Memo,
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

// @Security Auth-Jwt-Bearer
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
	case domain.ItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "delete customer failed")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

// @Security Auth-Jwt-Bearer
// @Summary 어드민 비밀번호 수정
// @Description 어드민 유저의 비밀번호를 수정하는 API
// @Accept json
// @Produce json
// @Param updateAdminPassword body UpdatePasswordRequest true "Update Admin Password"
// @Success 204 "비밀번호 변경 성공"
// @Router /user/admin/pw [patch]
func (h *HttpHandler) updateAdminPassword(ctx echo.Context) error {
	var req UpdatePasswordRequest

	req.UserId = ctx.Request().Header.Get("User-Id")
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "update password, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.UpdateAdminPassword(ctx.Request().Context(), domain.UpdateAdminPassword{
		UserId:      uuid.MustParse(req.UserId),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusNoContent)
	case domain.UserWrongPassword:
		return ctx.JSON(http.StatusUnauthorized, domain.UserWrongPasswordToUpdatePassword)
	case domain.ItemNotFound:
		return ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "update password, unhandled error useCase.UpdateAdminPassword")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

type CreateAdminRequest struct {
	// Name, 길이 2~60 제한
	Name string `json:"name" validate:"required,min=2,max=60" example:"ljs"`

	// Email, 이메일 주소
	Email string `json:"email" validate:"required,email" example:"example@example.com"`

	// Password, 형식 : 1234qwer!@
	Password string `json:"password" validate:"required,sf_password" example:"1234qwer!@"`

	// Nickname, 길이 2~60 제한
	Nickname string `json:"nickname" validate:"required,min=2,max=60" example:"광대버기"`
} // @name CreateAdminRequest

// @Security Auth-Jwt-Bearer
// @Summary Admin 유저 생성
// @Description Admin 유저를 생성하는 기능
// @Accept json
// @Produce json
// @Param AdminUserBody body CreateAdminRequest true "Admin User Body"
// @Success 201 {object} CreatedCustomerResp
// @Router /user/admin [post]
func (h *HttpHandler) createAdmin(ctx echo.Context) error {
	var req CreateAdminRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "create admin, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	newId, err := h.useCase.CreateAdminUser(ctx.Request().Context(), domain.CreateAdminUser{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusCreated, CreatedCustomerResp{Id: newId})
	case domain.ItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ItemExist)
	default:
		log.WithError(err).Error(tag, "create admin, unhandled error useCase.CreateAdminUser")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

func (h *HttpHandler) Bind(e *echo.Echo) {
	//CRUD, customer or admin
	e.POST("/user/customer", h.createCustomer)
	//sign, auth
	e.POST("/user/sign", h.signInUser)

	// todo debug.JwtBypassOnDebugWithRole 추후 추가해주세요
	e.DELETE("/user/customer/:userId", h.deleteCustomerUser)

	//Update Admin Password
	e.PATCH("/user/admin/pw", h.updateAdminPassword, debug.JwtBypassOnDebug())

	//create admin
	e.POST("/user/admin", h.createAdmin, debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))
}
