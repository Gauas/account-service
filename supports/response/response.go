package response

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{Status: http.StatusOK, Data: data})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{Status: http.StatusCreated, Data: data})
}

func NoContent(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, Response{Status: http.StatusOK, Data: echo.Map{"message": message}})
}

func Wrap(err error) error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrorNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return NewError(http.StatusConflict, "conflict")
	}
	return NewError(http.StatusInternalServerError, "internal server error")
}
