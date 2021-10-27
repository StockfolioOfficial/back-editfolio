package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

type SignInRequest struct {
	// Username, 아이디
	Username string `json:"username" validate:"required,min=8" example:"example@example.com"`

	// Password, 패스워드
	Password string `json:"password" validate:"required,min=8" example:"abcd12!@"`
} // @name SignInRequest

type TokenResponse struct {
	Token string `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
} // @name TokenResponse

// @Tags (Auth) 공용 기능
// @Summary 로그인 기능
// @Description 로그인하여 jwt 토큰을 받아오는 기능
// @Accept json
// @Produce json
// @Param signInUserBody body SignInRequest true "로그인 데이터 정보"
// @Success 200 {object} TokenResponse "로그인 완료"
// @Router /sign-in [post]
func (h *UserController) signInUser(ctx echo.Context) error {
	var req SignInRequest
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "sign in user, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	token, err := h.useCase.SignInUser(ctx.Request().Context(), domain.SignInUser{
		Username: req.Username,
		Password: req.Password,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusOK, TokenResponse{Token: token})
	case domain.ErrItemNotFound, domain.ErrUserWrongPassword:
		return ctx.JSON(http.StatusUnauthorized, domain.UserSignInFailedResponse)
	default:
		log.WithError(err).Error(tag, "sign in user, unhandled error useCase.SignInUser")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}
