package totp

import (
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Generate(c echo.Context) error {
	data, err := h.Service.GenerateTOTP(c.Request().Context(), middlewares.UserID(c.Request().Context()))
	if err != nil {
		return err
	}

	return httpresp.OK(c, data)
}
