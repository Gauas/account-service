package totp

import (
	"net/http"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GenerateTOTP(c echo.Context) error {
	data, err := h.Service.GenerateTOTP(c)
	if err != nil {
		return err
	}

	return response.OK(c, data)
}

func (h *Handler) EnableTOTP(c echo.Context) error {
	var req dto.EnableTOTPRequest
	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request body")
	}
	if req.OTPCode == "" {
		return response.NewError(http.StatusBadRequest, "otp_code is required")
	}

	if err := h.Service.EnableTOTP(c, req.OTPCode); err != nil {
		return err
	}

	return response.NoContent(c, "TOTP enabled successfully")
}

func (h *Handler) VerifyTOTP(c echo.Context) error {
	var req dto.VerifyTOTPRequest
	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request body")
	}
	if req.OTPCode == "" {
		return response.NewError(http.StatusBadRequest, "otp_code is required")
	}

	data, err := h.Service.VerifyTOTP(c, req.OTPCode)
	if err != nil {
		return err
	}

	return response.OK(c, data)
}
