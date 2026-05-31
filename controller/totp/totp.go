package totp

import (
	"net/http"

	"github.com/gauas/account-service/dto/request"
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

func (h *Handler) EnableTOTP(c echo.Context) error {
	var req request.EnableTOTPRequest
	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
	}
	if req.OTPCode == "" {
		return httpresp.NewError(http.StatusBadRequest, "otp_code is required")
	}

	if err := h.Service.EnableTOTP(c, req.OTPCode); err != nil {
		return err
	}

	return httpresp.NoContent(c, "TOTP enabled successfully")
}

func (h *Handler) VerifyTOTP(c echo.Context) error {
	var req request.VerifyTOTPRequest
	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
	}
	if req.OTPCode == "" {
		return httpresp.NewError(http.StatusBadRequest, "otp_code is required")
	}

	data, err := h.Service.VerifyTOTP(c, req.OTPCode)
	if err != nil {
		return err
	}

	return httpresp.OK(c, data)
}
