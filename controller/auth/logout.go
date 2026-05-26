package auth

import (
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Logout(c echo.Context) error {
	err := h.Service.Logout(c)
	if err != nil {
		return err
	}

	return response.OK(c, "logged out successfully")
}
