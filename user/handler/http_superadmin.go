package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"net/http"
)

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

// @Tags (User) 슈퍼어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [슈퍼어드민] 어드민 생성
// @Description 어드민을 생성하는 기능, 역할(role)이 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param requestBody body CreateAdminRequest true "어드민 생성 정보 데이터 구조"
// @Success 201 {object} CreatedUserResponse "어드민 생성 완료"
// @Router /admin [post]
func (h *UserController) createAdmin(ctx echo.Context) error {
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
		return ctx.JSON(http.StatusCreated, CreatedUserResponse{Id: newId})
	case domain.ErrItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ItemExist)
	default:
		log.WithError(err).Error(tag, "create admin, unhandled error useCase.CreateAdminUser")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}


type UpdateAdminInfoRequest struct {
	UserId   uuid.UUID `param:"userId" json:"-" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Password string    `json:"password" validate:"required,sf_password" example:"pass1234!@"`
	Email    string    `json:"email" validate:"required,email" example:"example@example.com"`
	Name     string    `json:"name" validate:"required,min=2,max=60" example:"sch"`
	Nickname string    `json:"nickname" validate:"required,min=2,max=60" example:"nickname"`
} // @name UpdateAdminInfoRequest

// @Tags (User) 슈퍼어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [슈퍼어드민] 어드민 정보 수정
// @Description 슈퍼 어드민이 어드민의 정보를 강제로 수정하는 기능, 역할(role)이 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param requestBody body UpdateAdminInfoRequest true "어드민 정보 수정 데이터 구조"
// @Param user_id path string true "어드민 식별 아이디(UUID)"
// @Success 204 "어드민 정보 수정 성공"
// @Router /admin/{user_id} [put]
func (h *UserController) updateAdminBySuperAdmin(ctx echo.Context) error {
	var req UpdateAdminInfoRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "force update admin, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.UpdateAdminInfoBySuperAdmin(ctx.Request().Context(), domain.UpdateAdminInfoBySuperAdmin{
		UserId:   req.UserId,
		Password: req.Password,
		Name:     req.Name,
		Username: req.Email,
		Nickname: req.Nickname,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusNoContent)
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
	case domain.ErrItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.EmailExistsResponse)
	default:
		log.WithError(err).Error(tag, "force-update admin, unhandled error useCase.ForceUpdateAdminInfoBySuperAdmin")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}


type DeleteAdminRequest struct {
	// Id, 어드민 Id
	Id uuid.UUID `param:"userId" json:"-" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// @Tags (User) 슈퍼어드민 기능
// @Security Auth-Jwt-Bearer
// @Summary [슈퍼어드민] 어드민 삭제
// @Description 어드민 유저를 삭제 하는 기능, 역할(role)이 'SUPER_ADMIN' 이여야함
// @Accept json
// @Produce json
// @Param user_id path string true "어드민 식별 아이디(UUID)"
// @Success 204 "삭제 완료"
// @Router /admin/{user_id} [delete]
func (h *UserController) deleteAdminBySuperAdmin(ctx echo.Context) error {
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
	case domain.ErrItemNotFound:
		return ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "delete customer failed")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}