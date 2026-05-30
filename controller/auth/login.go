package auth

import (
	"net/http"

	"github.com/gauas/account-service/dto"
	response2 "github.com/gauas/account-service/supports/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return response2.NewError(http.StatusBadRequest, "invalid request")
	}

	data, err := h.Service.Login(c, req)
	if err != nil {
		return err
	}

	return response2.OK(c, data)
}
