package user

import (
	"godok/domain/user/router"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	router.MappingUrl(e)
}
