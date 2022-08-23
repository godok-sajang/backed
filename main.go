package main

import (
	"echo_sample/config"
	"echo_sample/db"
	"echo_sample/domain/user"
	userService "echo_sample/domain/user/service"
	"echo_sample/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	eMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	config.Init()

	db.Connect()
	defer db.Close()

	InitServices()

	InitWebServices()
}

func InitServices() {
	userService.Init()
}

func InitWebServices() {
	// Echo instance
	e := echo.New()

	//JWT Authorization``
	e.Use(middleware.JWTAuth())

	// Middleware
	e.Use(eMiddleware.Logger())
	e.Use(eMiddleware.Recover())

	// Service init
	user.Init(e)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

func ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "ping")
}
