package echoutil

import (
	"fmt"
	"net/http"

	"godok/util/errorutil"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func HTTPErrorHandler(err error, ctx echo.Context) {
	switch v := errors.Cause(err).(type) {
	case errorutil.GodokError:
		ctx.JSON(v.GetHttpStatusCode(), errorutil.ErrorResponse{
			CustomCode: v.GetCustomCode(),
			Message:    v.GetMessage(),
		})
	default:
		var message string

		if ctx.Response().Committed {
			return
		}

		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
		}

		code := he.Code

		if m, ok := he.Message.(string); ok {
			message = m
		}

		cause := errors.Cause(err)
		fmt.Printf("%+v\n", cause)

		if ctx.Request().Method == http.MethodHead {
			err = ctx.NoContent(he.Code)
		} else {
			err = ctx.JSON(code, errorutil.ErrorResponse{
				Message: message,
			})
		}
	}
}
