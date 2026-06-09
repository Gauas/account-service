package auth

import (
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/gauas/account-service/packages/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Logout(c echo.Context) error {
	err := h.Service.Logout(c.Request().Context(), middlewares.RefreshToken(c), middlewares.DeviceID(c.Request().Context()))
	if err != nil {
		return err
	}

	session.Clear(c, h.Config)
	return httpresp.OK(c, "logged out successfully")
}
