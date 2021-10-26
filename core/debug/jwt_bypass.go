package debug

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/core/config"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

func JwtBypassOnDebug() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {

		if config.IsDebug {
			return handleJwtBypass(handlerFunc, nil)
		}

		return func(ctx echo.Context) error {
			return handlerFunc(ctx)
		}
	}
}

func JwtBypassOnDebugWithRole(role domain.UserRole) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {

		if config.IsDebug {
			return handleJwtBypass(handlerFunc, (*string)(&role))
		}

		return func(ctx echo.Context) error {
			return handlerFunc(ctx)
		}
	}
}

func handleJwtBypass(handlerFunc echo.HandlerFunc, role *string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var jwtDummy struct {
			Sub   string   `json:"sub"`
			Roles []string `json:"roles"`
		}

		fullValue := ctx.Request().Header.Get(echo.HeaderAuthorization)
		parts := strings.Split(fullValue, ".")
		if len(parts) != 3 {
			return ctx.JSON(http.StatusUnauthorized, domain.InvalidateTokenResponse)
		}

		decodedPart, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			log.WithError(err).Error("bypass, jwt payload base64 decode failed")
			return ctx.JSON(http.StatusUnauthorized, domain.InvalidateTokenResponse)
		}
		err = json.Unmarshal(decodedPart, &jwtDummy)
		if err != nil {
			log.WithError(err).Error("bypass, jwt payload unmarshal failed")
			return ctx.JSON(http.StatusUnauthorized, domain.InvalidateTokenResponse)
		}

		if role != nil && !hasRole(jwtDummy.Roles, *role) {
			return ctx.JSON(http.StatusUnauthorized, domain.NoPermissionResponse)
		}

		ctx.Request().Header.Set("User-Id", jwtDummy.Sub)
		return handlerFunc(ctx)
	}
}

func hasRole(roles []string, role string) bool {
	for _, v := range roles {
		if v == role {
			return true
		}
	}

	return false
}
