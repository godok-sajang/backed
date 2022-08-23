package middleware

import (
	"echo_sample/config"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
	eMiddleware "github.com/labstack/echo/v4/middleware"
)

var (
	TokenValidationMinutes   = 60 * 60 * 2
	RefreshValidationMinutes = 60 * 60 * 24 * 365

	JWTExpiredErrorStatus        = 9601
	DuplicateUsernameErrorStatus = 9305
)

func CreateToken(userid int64, duration int) (string, error) {
	var err error
	//Creating Access Token
	secret, _ := base64.StdEncoding.DecodeString(config.Config("SECRET"))
	atClaims := jwt.MapClaims{}
	atClaims, err = createJwt(strconv.FormatInt(userid, 10), time.Duration(duration), atClaims)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func JWTAuth() echo.MiddlewareFunc {
	jwtSecret, _ := base64.StdEncoding.DecodeString(os.Getenv("SECRET"))
	return eMiddleware.JWTWithConfig(eMiddleware.JWTConfig{
		ContinueOnIgnoredError: true,
		SigningKey:             jwtSecret,
		ErrorHandlerWithContext: func(err error, ctx echo.Context) error {
			ctx.Set("jwtError", err)
			return nil
		},
	})
}

func Protected() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u := c.Get("user")
			if u == nil {
				return JWTError(c)
			}
			token, ok := u.(*jwt.Token)
			if !ok {
				return JWTError(c)
			}
			_, ok = token.Claims.(jwt.MapClaims)
			if !ok {
				return JWTError(c)
			}
			return next(c)
		}
	}
}

type HTTPStatusResponse struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}

type JavaStatusResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorName    string `json:"error_name"`
	ErrorMessage string `json:"error_message"`
}

func JWTError(c echo.Context) error {
	err := c.Get("jwtError")
	if err != nil {
		if e, ok := err.(error); ok && e.Error() == "Missing or malformed JWT" {
			fmt.Fprintln(os.Stderr, "["+time.Now().String()+"] "+"jwtError: ", e.Error())
			return c.JSON(
				http.StatusUnauthorized,
				HTTPStatusResponse{
					Status:  http.StatusUnauthorized,
					Message: "Missing or malformed JWT",
				})

		}
	}
	return c.JSON(
		http.StatusUnauthorized,
		JavaStatusResponse{
			ErrorCode:    JWTExpiredErrorStatus,
			ErrorName:    "e-x-p-i-r-e-d_-j-w-t",
			ErrorMessage: "Expired JWT Exception",
		})
}

/* Private Functions */

func createJwt(subject string, expiration time.Duration, atclaims jwt.MapClaims) (jwt.MapClaims, error) {
	if subject != "" {
		atclaims["user"] = subject
	}
	atclaims["iat"] = time.Now().Unix()
	atclaims["exp"] = time.Now().Add(time.Second * expiration).Unix()

	return atclaims, nil
}
