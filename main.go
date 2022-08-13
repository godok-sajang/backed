package main

import (
	"echo_sample/config"
	"echo_sample/db"
	"echo_sample/domain/user"
	userService "echo_sample/domain/user/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	//JWT Authorization
	// e.Use(middleware.JWT([]byte("secret")))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Service init
	user.Init(e)
	e.GET("/ping", ping)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

func ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "ping")
}
