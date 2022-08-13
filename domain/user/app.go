package user

import (
	"echo_sample/domain/user/router"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	router.MappingUrl(e)
}
