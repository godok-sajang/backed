package router

import (
	"net/http"

	"godok/domain/user/dto"
	userdto "godok/domain/user/dto"
	user "godok/domain/user/service"

	eu "godok/util/errorutil"

	"github.com/labstack/echo/v4"
)

var (
	userService = user.UserService{}
)

func MappingUrl(app *echo.Echo) {
	app.POST("/user/sign-up", UserSignUp)
	app.POST("/user/sign-in", UserSignIn)
}

func UserSignUp(c echo.Context) error {
	var req dto.UserInfoRequest
	if err := c.Bind(&req); err != nil {
		return eu.New().WithHttpCode(http.StatusBadRequest).WithMessage(err)
	}

	if !req.ValidateNickname() {
		return eu.New().WithHttpCode(http.StatusBadRequest).WithMessage("invalid nickname")
	}

	if !req.ValidateEmail() {
		return eu.New().WithHttpCode(http.StatusBadRequest).WithMessage("invalid email")
	}

	if !req.ValidatePassword() {
		return eu.New().WithHttpCode(http.StatusBadRequest).WithMessage("invalid password")
	}

	if req.ValidateBirth() != nil {
		return eu.New().WithHttpCode(http.StatusBadRequest).WithMessage("invalid birth")
	}

	token, err := userService.CreateUserInfo(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

func UserSignIn(c echo.Context) error {
	var u userdto.UserSignInRequest
	if err := c.Bind(u); err != nil {
		return eu.New().WithHttpCode(http.StatusBadRequest).WithMessage(err)
	}
	token, err := userService.SignIn(c.Request().Context(), userdto.UserSignInRequest{
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, token)
}
