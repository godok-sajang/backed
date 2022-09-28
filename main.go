package main

import (
	"godok/config"
	"godok/db"
	"godok/domain/user"
	userService "godok/domain/user/service"
	"godok/middleware"
	"godok/util/echoutil"
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

	// Custom error handler
	e.HTTPErrorHandler = echoutil.HTTPErrorHandler

	// Service init
	user.Init(e)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello")
	})

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
