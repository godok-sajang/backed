package main

import (
	"godok/config"
	"godok/db"
	"godok/domain/user"
	userService "godok/domain/user/service"
	"godok/middleware"
	"godok/util/echoutil"

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

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
