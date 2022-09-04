package errorutil

import "github.com/pkg/errors"

type ErrorResponse struct {
	CustomCode string `json:"code,omitempty"`
	Message    string `json:"message"`
}

type GodokError interface {
	Error() string
	GetHttpStatusCode() int
	GetCustomCode() string
	GetMessage() string
	WithHttpCode(code int) GodokError
	WithCustomCode(code string) GodokError
	WithMessage(message interface{}) GodokError
}

type godokError struct {
	httpStatusCode int
	customCode     string
	message        string
}

func New() GodokError {
	return &godokError{
		httpStatusCode: 0,
		customCode:     "",
		message:        "",
	}
}

func (g *godokError) Error() string {
	return ""
}

func (g *godokError) GetHttpStatusCode() int {
	return g.httpStatusCode
}

func (g *godokError) GetCustomCode() string {
	return g.customCode
}

func (g *godokError) GetMessage() string {
	return g.message
}

func (g *godokError) WithHttpCode(code int) GodokError {
	g.httpStatusCode = code
	return g
}

func (g *godokError) WithCustomCode(code string) GodokError {
	g.customCode = code
	return g
}

func (g *godokError) WithMessage(message interface{}) GodokError {
	switch v := message.(type) {
	case string:
		g.message = v
	case error:
		g.message = v.Error()
	default:
		g.message = "unkown message"
	}
	return g
}

func HasCode(err error, code string) bool {
	var gerr GodokError

	gerr, ok := errors.Cause(err).(GodokError)
	if !ok {
		return false
	}

	if gerr.GetCustomCode() != code {
		return false
	}

	return false
}

func InternalError(v interface{}) error {
	switch v := v.(type) {
	case string:
		return errors.New(v)
	case error:
		return errors.New(v.Error())
	default:
		return errors.New("unkown type")
	}
}
