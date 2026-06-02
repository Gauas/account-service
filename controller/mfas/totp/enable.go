package totp

import (
	"net/http"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) EnableTOTP(c echo.Context) error {
	var req request.EnableTOTPRequest
	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.Service.EnableTOTP(c, req); err != nil {
		return err
	}

	return httpresp.NoContent(c, "TOTP enabled successfully")
}
