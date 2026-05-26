package auth

import (
	"net/http"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request body")
	}

	data, err := h.Service.NewAccount(c, req)
	if err != nil {
		return err
	}

	return response.OK(c, data)
}
