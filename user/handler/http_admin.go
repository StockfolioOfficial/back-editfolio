package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"net/http"
)

type UpdateAdminMyInfoRequest struct {
	Email    string `json:"email" validate:"required,email" example:"example@example.com"`
	Name     string `json:"name" validate:"required,min=2,max=60" example:"sch"`
	Nickname string `json:"nickname" validate:"required,min=2,max=60" example:"nickname"`
} // @name UpdateAdminMyInfo

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 자기 정보 수정
// @Description 어드민이 자기자신의 정보를 수정하는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param requestBody body UpdateAdminMyInfoRequest true "어드민 정보 수정 데이터 구조"
// @Success 204 "정보 수정 성공"
// @Router /admin/me [put]
func (h *HttpHandler) updateAdminMyInfo(ctx echo.Context, userId uuid.UUID) error {
	var req UpdateAdminMyInfoRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "update admin, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	if err != nil {
		log.WithError(err).Trace(tag, "UUID error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.UpdateAdminInfo(ctx.Request().Context(), domain.UpdateAdminInfo{
		UserId:   userId,
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

type UpdateAdminMyPasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,sf_password" example:"abcd1234!@"`
	NewPassword string `json:"newPassword" validate:"required,sf_password" example:"pass1234!@"`
} // @name UpdateAdminMyPasswordRequest

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 자기 비밀번호 수정
// @Description 어드민이 자기자신의 비밀번호를 수정하는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param requestBody body UpdateAdminMyPasswordRequest true "비밀번호 수정 데이터 구조"
// @Success 204 "비밀번호 변경 성공"
// @Router /admin/me/pw [patch]
func (h *HttpHandler) updateAdminMyPassword(ctx echo.Context, userId uuid.UUID) error {
	var req UpdateAdminMyPasswordRequest
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "update password, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.UpdateAdminPassword(ctx.Request().Context(), domain.UpdateAdminPassword{
		UserId:      userId,
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

// =================== CUSTOMER CONTROL ===================

type CreateCustomerRequest struct {
	// Name, 길이 2~60 제한
	Name string `json:"name" validate:"required,min=2,max=60" example:"ljs"`

	// Email, 이메일 주소
	Email string `json:"email" validate:"required,email" example:"example@example.com"`

	// Mobile, 형식 : 01012345678
	Mobile string `json:"mobile" validate:"required,sf_mobile" example:"01012345678"`
} // @name CreateCustomerRequest

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 고객 생성
// @Description 고객을 생성하는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param requestBody body CreateCustomerRequest true "고객 생성 정보 데이터 구조"
// @Success 201 {object} CreatedUserResponse "고객 생성 완료"
// @Router /customer [post]
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
		return ctx.JSON(http.StatusCreated, CreatedUserResponse{Id: newId})
	default:
		log.WithError(err).Error(tag, "create customer, unhandled error useCase.CreateCustomerUser")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

type UpdateCustomerInfoRequest struct {
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
} //@name UpdateCustomerInfoRequest

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 고객 정보 수정
// @Description 고객 정보 수정하는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param user_id path string true "고객 식별 아이디(UUID)"
// @Param requestBody body UpdateCustomerInfoRequest true "고객 정보 수정 데이터 구조"
// @Success 204 "수정 완료"
// @Router /customer/{user_id} [put]
func (h *HttpHandler) updateCustomer(ctx echo.Context) error {
	var req UpdateCustomerInfoRequest

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

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 고객 삭제
// @Description 고객 삭제하는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param user_id path string true "고객 식별 아이디(UUID)"
// @Success 204 "삭제 완료"
// @Router /customer/{user_id} [delete]
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