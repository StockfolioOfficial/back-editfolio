package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"net/http"
	"time"
)

type AdminSimpleInfoResponse struct {
	UserId   uuid.UUID `json:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role     []string  `json:"roles" validate:"required" example:"ADMIN"`
	Name     string    `json:"name" validate:"required" example:"이름"`
	Username string    `json:"username" validate:"required" example:"example@example.com"`
	Nickname string    `json:"nickname" validate:"required" example:"(주)스톡폴리오"`
} // @name AdminSimpleInfoResponse

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 자기 정보 가져오기
// @Description 어드민이 자기 정보 가져오는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Success 200 {object} AdminSimpleInfoResponse "성공"
// @Router /admin/me [get]
func (c *UserController) getAdminMyInfo(ctx echo.Context, userId uuid.UUID) error {

	res, err := c.useCase.GetAdminInfoDetailByUserId(ctx.Request().Context(), userId)

	switch err {
	case nil:
		return ctx.JSON(http.StatusOK, AdminSimpleInfoResponse{
			UserId:   res.UserId,
			Role:     []string{string(res.Role)},
			Name:     res.Name,
			Username: res.Username,
			Nickname: res.Nickname,
		})
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "getAdminMyInfo, unhandled error useCase.GetAdminInfoDetailByUserId")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}


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
func (c *UserController) updateAdminMyInfo(ctx echo.Context, userId uuid.UUID) error {
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

	err = c.useCase.UpdateAdminInfo(ctx.Request().Context(), domain.UpdateAdminInfo{
		UserId:   userId,
		Name:     req.Name,
		Username: req.Email,
		Nickname: req.Nickname,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusNoContent)
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
	case domain.ErrItemAlreadyExist:
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
func (c *UserController) updateAdminMyPassword(ctx echo.Context, userId uuid.UUID) error {
	var req UpdateAdminMyPasswordRequest
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "update password, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = c.useCase.UpdateAdminPassword(ctx.Request().Context(), domain.UpdateAdminPassword{
		UserId:      userId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusNoContent)
	case domain.ErrUserWrongPassword:
		return ctx.JSON(http.StatusUnauthorized, domain.UserWrongPasswordToUpdatePassword)
	case domain.ErrItemNotFound:
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
func (c *UserController) createCustomer(ctx echo.Context) error {
	var req CreateCustomerRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "create customer, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	newId, err := c.useCase.CreateCustomerUser(ctx.Request().Context(), domain.CreateCustomerUser{
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
	ChannelName string `json:"channelName" validate:"max=100" example:"밥굽남"`

	// ChannelLink, 길이 2048 제한
	ChannelLink string `json:"channelLink" validate:"max=2048" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`

	// Email, 이메일 주소
	Email string `json:"email" validate:"required,email" example:"example@example.com"`

	// Mobile, 형식 : 01012345678
	Mobile string `json:"mobile" validate:"required,sf_mobile" example:"01012345678"`

	//PersonaLink, 길이 2048 제한
	PersonaLink string `json:"personaLink" validate:"max=2048" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`

	//OnedriveLink, 길이 2048 제한
	OnedriveLink string `json:"onedriveLink" validate:"max=2048" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`

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
func (c *UserController) updateCustomer(ctx echo.Context) error {
	var req UpdateCustomerInfoRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "update customer, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = c.useCase.UpdateCustomerUser(ctx.Request().Context(), domain.UpdateCustomerUser{
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
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ErrItemNotFound) // TODO refactor
	case domain.ErrItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ErrItemAlreadyExist) // TODO refactor
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
func (c *UserController) deleteCustomerUser(ctx echo.Context) error {
	var req DeleteCustomerRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "delete customer, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}
	err = c.useCase.DeleteCustomerUser(ctx.Request().Context(), domain.DeleteCustomerUser{
		UserId: req.Id,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusNoContent, domain.ErrorResponse{Message: err.Error()})
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "delete customer failed")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

type FetchCustomerRequest struct {
	Query string `json:"-" query:"q"`
}

type CustomerInfoResponse struct {
	UserId      uuid.UUID `json:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" validate:"required" example:"(대충 고객 이름)"`
	ChannelName string    `json:"channelName" validate:"required" example:"(대충 채널 이름)"`
	ChannelLink string    `json:"channelLink" validate:"required" example:"(대충 채널 url 링크)"`
	Email       string    `json:"email" validate:"required,email" example:"example@example.com"`
	Mobile      string    `json:"mobile" validate:"required" example:"01012345678"`
	CreatedAt   time.Time `json:"createdAt" validate:"required" example:"2021-10-27T04:44:18+00:00"`
} // @name CustomerInfoResponse

type CustomerInfoListResponse []CustomerInfoResponse

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 고객 목록
// @Description 고객 목록 가져오는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param q query string false "검색어"
// @Success 200 {object} CustomerInfoListResponse "성공"
// @Router /customer [get]
func (c *UserController) fetchCustomer(ctx echo.Context) error {
	var req FetchCustomerRequest
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "fetch full customer, request data bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	list, err := c.useCase.FetchAllCustomer(ctx.Request().Context(), domain.FetchCustomerOption{
		Query: req.Query,
	})

	if err != nil {
		log.WithError(err).Error(tag, "fetch full customer, unhandled error useCase.FetchAllCustomer")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

	if len(list) == 0 {
		return ctx.NoContent(http.StatusNoContent)
	}

	res := make(CustomerInfoListResponse, len(list))

	for i := range list {
		src := list[i]
		res[i] = CustomerInfoResponse{
			UserId:      src.UserId,
			Name:        src.Name,
			ChannelName: src.ChannelName,
			ChannelLink: src.ChannelLink,
			Email:       src.Email,
			Mobile:      src.Mobile,
			CreatedAt:   src.CreatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, res)
}


type CustomerDetailInfoResponse struct {
	UserId       uuid.UUID `json:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string    `json:"name" validate:"required" example:"(대충 고객 이름)"`
	ChannelName  string    `json:"channelName" validate:"required" example:"(대충 채널 이름)"`
	ChannelLink  string    `json:"channelLink" validate:"required" example:"(대충 채널 url 링크)"`
	Email        string    `json:"email" validate:"required" example:"example@example.com"`
	Mobile       string    `json:"mobile" validate:"required" example:"01012345678"`
	PersonaLink  string    `json:"personaLink" validate:"required" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`
	OnedriveLink string    `json:"onedriveLink" validate:"required" example:"https://www.youtube.com/channel/UCdfhK0yIMjmhcQ3gP-qpXRw"`
	Memo         string    `json:"memo" example:"이사람 까다로움"`
} // @name CustomerDetailInfoResponse

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 고객 상세 정보
// @Description 고객 상제 정보 가져오는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param user_id path string true "고객 식별 아이디(UUID)"
// @Success 200 {object} CustomerDetailInfoResponse "성공"
// @Router /customer/{user_id} [get]
func (c *UserController) getCustomerDetailInfo(ctx echo.Context) error {
	var req struct {
		UserId uuid.UUID `json:"-" param:"userId"`
	}
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "get customer detail info, request data bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	detail, err := c.useCase.GetCustomerInfoDetailByUserId(ctx.Request().Context(), req.UserId)

	switch err {
	case nil:
		return ctx.JSON(http.StatusOK, CustomerDetailInfoResponse{
			UserId:       detail.UserId,
			Name:         detail.Name,
			ChannelName:  detail.ChannelName,
			ChannelLink:  detail.ChannelLink,
			Email:        detail.Email,
			Mobile:       detail.Mobile,
			PersonaLink:  detail.PersonaLink,
			OnedriveLink: detail.OnedriveLink,
			Memo:         detail.Memo,
		})
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "fetch full customer, unhandled error useCase.FetchAllCustomer")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

type FetchAdminRequest struct {
	Query string `json:"-" query:"q"`
}

type AdminInfoResponse struct {
	UserId    uuid.UUID `json:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" validate:"required" example:"(대충 어드민 이름)"`
	Nickname  string    `json:"nickname" validate:"required" example:"(대충 어드민 닉네임)"`
	Email     string    `json:"email" validate:"required" example:"example@example.com"`
	CreatedAt time.Time `json:"createdAt" validate:"required" example:"2021-10-27T04:44:18+00:00"`
} // @name AdminInfoResponse

type AdminInfoListResponse []AdminInfoResponse

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 어드민 목록
// @Description 어드민 목록 가져오는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param q query string false "검색어"
// @Success 200 {object} AdminInfoListResponse "성공"
// @Router /admin [get]
func (c *UserController) fetchAdmin(ctx echo.Context) error {
	var req FetchAdminRequest
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "fetch full admin, request data bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	list, err := c.useCase.FetchAllAdmin(ctx.Request().Context(), domain.FetchAdminOption{
		Query: req.Query,
	})

	if err != nil {
		log.WithError(err).Error(tag, "fetch full customer, unhandled error useCase.FetchAllCustomer")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

	if len(list) == 0 {
		return ctx.NoContent(http.StatusNoContent)
	}

	res := make(AdminInfoListResponse, len(list))

	for i := range list {
		src := list[i]
		res[i] = AdminInfoResponse{
			UserId:    src.UserId,
			Name:      src.Name,
			Nickname:  src.Nickname,
			Email:     src.Email,
			CreatedAt: src.CreatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, res)
}

type AdminCreatorInfoResponse struct {
	UserId   uuid.UUID `json:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name     string    `json:"name" validate:"required" example:"(대충 편집자 이름)"`
	Nickname string    `json:"nickname" validate:"required" example:"(대충 편집자 닉네임)"`
} //@name AdminCreatorInfoResponse

type AdminCreatorInfoListResponse []AdminCreatorInfoResponse

// @Tags (User) 어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [어드민] 편집자 목록
// @Description 편집자 목록 가져오는 기능, 역할(role)이 'ADMIN', 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param q query string false "검색어"
// @Success 200 {object} AdminCreatorInfoListResponse "성공"
// @Router /admin/creator [get]
func (c *UserController) fetchAdminCreator(ctx echo.Context) error {
	var req FetchAdminRequest
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "fetch full admin, request data bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	list, err := c.useCase.FetchAllAdmin(ctx.Request().Context(), domain.FetchAdminOption{
		Query: req.Query,
	})

	if err != nil {
		log.WithError(err).Error(tag, "fetch full customer, unhandled error useCase.FetchAllCustomer")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

	if len(list) == 0 {
		return ctx.NoContent(http.StatusNoContent)
	}

	res := make(AdminCreatorInfoListResponse, len(list))

	for i := range list {
		src := list[i]
		res[i] = AdminCreatorInfoResponse{
			UserId:    src.UserId,
			Name:      src.Name,
			Nickname:  src.Nickname,
		}
	}

	return ctx.JSON(http.StatusOK, res)
}
