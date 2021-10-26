package echox

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UserID(wrapper func(ctx echo.Context, userID uuid.UUID) error) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		strID := ctx.Request().Header.Get("User-ID")
		id, err := uuid.Parse(strID)
		if err != nil { // TODO error handling
			return err
		}
		return wrapper(ctx, id)
	}
}

func OptionalUserID(wrapper func(ctx echo.Context, userID *uuid.UUID) error) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		strID := ctx.Request().Header.Get("User-ID")
		id, err := uuid.Parse(strID)
		if err != nil {
			return wrapper(ctx, nil)
		}
		return wrapper(ctx, &id)
	}
}