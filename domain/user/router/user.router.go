package router

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"echo_sample/domain/user/dto"
	userdto "echo_sample/domain/user/dto"
	user "echo_sample/domain/user/service"
	"echo_sample/util"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	userService = user.UserService{}
)

func MappingUrl(app *echo.Echo) {
	app.GET("/user/info/:id", GetUserInfo)
	app.POST("/user/info", UserSignUp)
	app.POST("/user/signIn", UserSignIn)
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

	if !(len(req.Nickname) >= 1 && len(req.Nickname) <= 15) {
		fmt.Fprintln(os.Stderr, "["+time.Now().String()+"]"+"Nickname: ", req.Nickname)
		return c.JSON(http.StatusBadRequest, "invalid nickname")
	}
	if req.Email == nil || !util.ValidateEmail(*req.Email) {
		return c.JSON(http.StatusBadRequest, "invalid email")
	}
	if req.Password == nil || !util.ValidatePassword(*req.Password) {
		return c.JSON(http.StatusBadRequest, "invalid password")
	}
	var birth time.Time
	birth, err = time.Parse(time.RFC3339, req.Birth)
	if err != nil {
		birth, err = time.Parse("2006-01-02", req.Birth)
	}

	token, err := userService.CreateUserInfo(c.Request().Context(), dto.UserInfo{
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: req.Password,
		Birth:    birth,
		Gender:   req.Gender,
	})
	if err != nil {
		fmt.Printf("%+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, token)
}

func UserSignIn(c echo.Context) (err error) {
	u := new(userdto.UserVerified)
	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	token, err := userService.SignIn(c.Request().Context(), userdto.UserVerified{
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Cause(err))
	}
	return c.JSON(http.StatusOK, token)
}

// 테스트용 API
func MockData() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mock Data를 생성한다.
		list := map[string]string{
			"1": "고양이",
			"2": "사자",
			"3": "호랑이",
		}
		return c.JSON(http.StatusOK, list)
	}
}
