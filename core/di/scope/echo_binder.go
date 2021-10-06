package scope

import "github.com/labstack/echo/v4"

type EchoBinder interface {
	Bind(*echo.Echo)
}