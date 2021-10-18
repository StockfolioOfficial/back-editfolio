package debug

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/core/config"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"net/http"
	"strings"
)

func JwtBypassOnDebug() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {

		if config.IsDebug {
			return handleJwtBypass(handlerFunc)
		}

		return func(ctx echo.Context) error {
			return handlerFunc(ctx)
		}
	}
}

func handleJwtBypass(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var jwtDummy struct {
			Sub string `json:"sub"`
		}

		fullValue := ctx.Request().Header.Get(echo.HeaderAuthorization)
		parts := strings.Split(fullValue, ".")
		if len(parts) != 3 {
			return ctx.JSON(http.StatusUnauthorized, domain.InvalidateTokenResponse)
		}

		err := json.Unmarshal([]byte(parts[1]), &jwtDummy)
		if err != nil {
			log.WithError(err).Error("bypass, jwt payload unmarshal failed")
			return ctx.JSON(http.StatusUnauthorized, domain.InvalidateTokenResponse)
		}

		ctx.Request().Header.Set("User-Id", jwtDummy.Sub)
		return handlerFunc(ctx)
	}
}