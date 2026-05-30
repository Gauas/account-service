package totp

import (
	"net/http"

	"github.com/gauas/account-service/dto"
	response2 "github.com/gauas/account-service/supports/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GenerateTOTP(c echo.Context) error {
	data, err := h.Service.GenerateTOTP(c)
	if err != nil {
		return err
	}

	return response2.OK(c, data)
}

func (h *Handler) EnableTOTP(c echo.Context) error {
	var req dto.EnableTOTPRequest
	if err := c.Bind(&req); err != nil {
		return response2.NewError(http.StatusBadRequest, "invalid request body")
	}
	if req.OTPCode == "" {
		return response2.NewError(http.StatusBadRequest, "otp_code is required")
	}

	if err := h.Service.EnableTOTP(c, req.OTPCode); err != nil {
		return err
	}

	return response2.NoContent(c, "TOTP enabled successfully")
}

func (h *Handler) VerifyTOTP(c echo.Context) error {
	var req dto.VerifyTOTPRequest
	if err := c.Bind(&req); err != nil {
		return response2.NewError(http.StatusBadRequest, "invalid request body")
	}
	if req.OTPCode == "" {
		return response2.NewError(http.StatusBadRequest, "otp_code is required")
	}

	data, err := h.Service.VerifyTOTP(c, req.OTPCode)
	if err != nil {
		return err
	}

	return response2.OK(c, data)
}
