package totp

import (
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GenerateTOTP(c echo.Context) error {
	data, err := h.Service.GenerateTOTP(c)
	if err != nil {
		return err
	}

	return httpresp.OK(c, data)
}
