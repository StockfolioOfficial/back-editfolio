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

	// Email, 이메일 주소
	Email string `json:"email" validate:"required,email" example:"example@example.com"`

	// Mobile, 형식 : 01012345678
	Mobile string `json:"mobile" validate:"required,sf_mobile" example:"01012345678"`
} // @name CreateCustomerUserRequest

type CreatedCustomerResp struct {
	Id uuid.UUID `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name CreatedCustomerResponse

type UpdatePasswordRequest struct {
	UserId string `json:"-" header:"User-Id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Mobile, 형식 : 01012345678
	OldPassword string `json:"oldPassword" validate:"required,sf_password" example:"abcd1234!@"`
	NewPassword string `json:"newPassword" validate:"required,sf_password" example:"pass1234!@"`
} // @name UpdatePasswordRequest

type UdpateAdminInfomationRequest struct {
	UserId   string `json:"-" header:"User-Id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email    string `json:"email" validate:"required,email" example:"example@example.com"`
	Name     string `json:"name" validate:"required,min=2,max=60" example:"sch"`
	Nickname string `json:"nickname" validate:"required,min=2,max=60" example:"nickname"`
} // @name UpdateAdminInfomationRequest

type ForceUpdateAdminInfomationRequest struct {
	UserId   string `param:"userId" json:"-" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Password string `json:"password" validate:"required,sf_password" example:"pass1234!@"`
	Email    string `json:"email" validate:"required,email" example:"example@example.com"`
	Name     string `json:"name" validate:"required,min=2,max=60" example:"sch"`
	Nickname string `json:"nickname" validate:"required,min=2,max=60" example:"nickname"`
} // @name ForceUpdateAdminInfomationRequest

// @Security Auth-Jwt-Bearer
// @Summary 고객 유저 생성
// @Description 고객 유저를 생성하는 기능
// @Accept json
// @Produce json
// @Param customerUserBody body CreateCustomerRequest true "Customer User Body"
// @Success 201 {object} CreatedCustomerResp "유저 생성"
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

type UpdateCustomerRequest struct {
	// UserId,
	UserId uuid.UUID `json:"-" param:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Name, 길이 2~60 제한
	Name string `json:"name" validate:"required,min=2,max=60" example:"ljs"`

	// ChannelName, 길이 2~100 제한
	ChannelName string `json:"channelName" validate:"required,min=2,max=100" example:"밥굽남"`

	// ChannelLink, 길이 2048 제한
	ChannelLink string `json:"channelLink" validate:"required,max=2048" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`

	// Email, 이메일 주소
	Email string `json:"email" validate:"required,email" example:"example@example.com"`

	// Mobile, 형식 : 01012345678
	Mobile string `json:"mobile" validate:"required,sf_mobile" example:"01012345678"`

	//PersonaLink, 길이 2048 제한
	PersonaLink string `json:"personaLink" validate:"required,max=2048" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`

	//OnedriveLink, 길이 2048 제한
	OnedriveLink string `json:"onedriveLink" validate:"required,max=2048" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`

	//Memo, 형식 : text
	Memo string `json:"memo" example:"이사람 까다로움"`
} //@name UpdateCustomerUserRequest

// @Security Auth-Jwt-Bearer
// @Summary 고객 유저 정보 수정하는 기능
// @Description 고객 유저 정보 수정하는 기능
// @Accept json
// @Produce json
// @Param user_id path string true "Customer User Id"
// @Param customerUserBody body UpdateCustomerRequest true "Customer User Body"
// @Success 204 "수정 완료"
// @Router /user/customer/{user_id} [put]
func (h *HttpHandler) updateCustomer(ctx echo.Context) error {
	var req UpdateCustomerRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "update customer, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.UpdateCustomerUser(ctx.Request().Context(), domain.UpdateCustomerUser{
		UserId:       req.UserId,
		Name:         req.Name,
		ChannelName:  req.ChannelName,
		ChannelLink:  req.ChannelLink,
		Email:        req.Email,
		Mobile:       req.Mobile,
		PersonaLink:  req.PersonaLink,
		OnedriveLink: req.OnedriveLink,
		Memo:         req.Memo,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusNoContent)
	case domain.ItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ItemNotFound)
	case domain.ItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ItemAlreadyExist)
	default:
		log.WithError(err).Error(tag, "update customer, unhandled error useCase.UpdateCustomerUser")
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
// @Param user_id path string true "Customer User Id"
// @Success 204 "삭제 완료"
// @Router /user/customer/{user_id} [delete]
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

type DeleteAdminRequest struct {
	// Id, 어드민 Id
	Id uuid.UUID `param:"adminId" json:"-" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// @Security Auth-Jwt-Bearer
// @Summary 어드민 유저 삭제
// @Description 어드민 유저를 삭제하는 기능
// @Accept json
// @Produce json
// @Param user_id path string true "Admin User Id"
// @Success 204
// @Router /user/admin/{user_id} [delete]
func (h *HttpHandler) deleteAdminUser(ctx echo.Context) error {
	var req DeleteAdminRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "delete admin, request body error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}
	err = h.useCase.DeleteAdminUser(ctx.Request().Context(), domain.DeleteAdminUser{
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
// @Summary 어드민 정보 수정
// @Description 어드민 유저의 정보를 수정하는 API
// @Accept json
// @Produce json
// @Param updateAdmin body UdpateAdminInfomationRequest true "Update Admin Info"
// @Success 204 "정보 수정 성공"
// @Router /user/admin [put]
func (h *HttpHandler) updateAdmin(ctx echo.Context) error {
	var req UdpateAdminInfomationRequest
	req.UserId = ctx.Request().Header.Get("User-Id")

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "update admin, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	userUUID, err := uuid.Parse(req.UserId)

	if err != nil {
		log.WithError(err).Trace(tag, "UUID error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.UpdateAdminInfo(ctx.Request().Context(), domain.UpdateAdminInfo{
		UserId:   userUUID,
		Name:     req.Name,
		Username: req.Email,
		Nickname: req.Nickname,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusNoContent)
	case domain.ItemNotFound:
		return ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
	case domain.ItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ItemExist)
	default:
		log.WithError(err).Error(tag, "create admin, unhandled error useCase.UpdateAdminInfo")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

// @Security Auth-Jwt-Bearer
// @Summary 어드민 정보 강제 수정
// @Description 슈퍼 어드민이 어드민 유저의 정보를 강제로 수정하는 API
// @Accept json
// @Produce json
// @Param updateAdminBySuperAdmin body ForceUpdateAdminInfomationRequest true "Force-Update Admin Info"
// @Param user_id path string true "Admin User Id"
// @Success 204 "정보 수정 성공"
// @Router /user/super-admin/{user_id} [put]
func (h *HttpHandler) updateAdminBySuperAdmin(ctx echo.Context) error {
	var req ForceUpdateAdminInfomationRequest
	req.UserId = ctx.Request().Header.Get("User-Id")

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "force update admin, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	userUUID, err := uuid.Parse(req.UserId)

	if err != nil {
		log.WithError(err).Trace(tag, "UUID error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.ForceUpdateAdminInfoBySuperAdmin(ctx.Request().Context(), domain.ForceUpdateAdminInfo{
		UserId:   userUUID,
		Password: req.Password,
		Name:     req.Name,
		Username: req.Email,
		Nickname: req.Nickname,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusNoContent)
	case domain.ItemNotFound:
		return ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
	case domain.ItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ItemExist)
	default:
		log.WithError(err).Error(tag, "force-update admin, unhandled error useCase.ForceUpdateAdminInfoBySuperAdmin")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

func (h *HttpHandler) Bind(e *echo.Echo) {
	//CRUD, customer or admin
	e.POST("/user/customer", h.createCustomer, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))
	e.PUT("/user/customer/:userId", h.updateCustomer, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))

	//sign, auth
	e.POST("/user/sign", h.signInUser)

	e.DELETE("/user/customer/:userId", h.deleteCustomerUser, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))

	//Update Admin Password
	e.PATCH("/user/admin/pw", h.updateAdminPassword, debug.JwtBypassOnDebug())

	//Update Admin Infomation
	e.PUT("/user/admin", h.updateAdmin, debug.JwtBypassOnDebug())

	//Update Admin Infomation By SuperAdmin
	e.PUT("/user/super-admin/:userId", h.updateAdminBySuperAdmin, debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))

	//create admin
	e.POST("/user/admin", h.createAdmin, debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))

	//Delete admin
	e.DELETE("/user/admin/:adminId", h.deleteAdminUser, debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))
}
