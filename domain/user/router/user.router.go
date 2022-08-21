package router

import (
	"database/sql"
	"net/http"
	"strconv"

	"echo_sample/domain/user/dto"
	userdto "echo_sample/domain/user/dto"
	user "echo_sample/domain/user/service"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	userService = user.UserService{}
)

func MappingUrl(app *echo.Echo) {
	app.GET("/user/info/:id", GetUserInfo)
	app.POST("/user/sign-up", UserSignUp)
	app.POST("/user/sign-in", UserSignIn)
}

func GetUserInfo(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user, err := userService.GetUserInfo(c.Request().Context(), idInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func UserSignUp(c echo.Context) error {
	var req dto.UserInfoRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if req.ValidateNickname() {
		return c.JSON(http.StatusBadRequest, "invalid nickname")
	}

	if req.ValidateEmail() {
		return c.JSON(http.StatusBadRequest, "invalid email")
	}

	if req.ValidatePassword() {
		return c.JSON(http.StatusBadRequest, "invalid password")
	}

	if req.ValidateBirth() != nil {
		return c.JSON(http.StatusBadRequest, "invalid Birth")
	}

	token, err := userService.CreateUserInfo(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, token)
}

func UserSignIn(c echo.Context) (err error) {
	var u userdto.UserSignInRequest
	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	token, err := userService.SignIn(c.Request().Context(), userdto.UserSignInRequest{
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Cause(err))
	}
	return c.JSON(http.StatusOK, token)
}
