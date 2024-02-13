package errs

import (
	"net/http"
)

type Code int

const (
	CodeAPIFailed      Code = 1000
	CodeBadRequest     Code = 1001
	CodeValidateFailed Code = 1002
)

var (
	ErrAPIFailed      Err = New(http.StatusInternalServerError, CodeAPIFailed, "api failed", nil)
	ErrBadRequest     Err = New(http.StatusBadRequest, CodeBadRequest, "bad request", nil)
	ErrValidateFailed Err = New(http.StatusUnprocessableEntity, CodeValidateFailed, "validate failed", nil)
)

type Err struct {
	Status  int    `json:"-"`
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func New(status int, code Code, message string, data any) Err {
	return Err{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (e Err) Error() string {
	return e.Message
}

func (e Err) WithStatus(httpStatusCode int) Err {
	e.Status = httpStatusCode
	return e
}

func (e Err) WithData(data any) Err {
	e.Data = data
	return e
}
