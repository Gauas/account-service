package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/gauas/account-service/packages/response"
)

func (ctrl *Controller) GenerateTOTPQR(c echo.Context) error {
	userID, err := userIDFromContext(c)
	if err != nil {
		return response.NewError(http.StatusUnauthorized, err.Error())
	}

	setup, err := ctrl.service.GenerateTOTPSetup(c.Request().Context(), userID)
	if err != nil {
		return response.Wrap(err)
	}

	return response.OK(c, setup)
}

type enableTOTPRequest struct {
	OTPCode string `json:"otp_code"`
}

func (ctrl *Controller) EnableTOTP(c echo.Context) error {
	userID, err := userIDFromContext(c)
	if err != nil {
		return response.NewError(http.StatusUnauthorized, err.Error())
	}

	var req enableTOTPRequest
	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request body")
	}

	if err := ctrl.service.EnableTOTP(c.Request().Context(), userID, req.OTPCode); err != nil {
		return response.Wrap(err)
	}

	return response.NoContent(c, "TOTP enabled successfully")
}

type verifyTOTPRequest struct {
	OTPCode  string `json:"otp_code"`
	DeviceID string `json:"device_id"`
}

func (ctrl *Controller) VerifyTOTP(c echo.Context) error {
	userID, err := userIDFromContext(c)
	if err != nil {
		return response.NewError(http.StatusUnauthorized, err.Error())
	}

	var req verifyTOTPRequest
	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request body")
	}

	deviceID := deviceIDFromRequest(c, req.DeviceID)
	if deviceID == "" {
		return response.NewError(http.StatusBadRequest, "device_id is required")
	}

	accessToken, refreshToken, expiresAt, err := ctrl.service.VerifyTOTP(c.Request().Context(), userID, req.OTPCode, deviceID)
	if err != nil {
		return response.Wrap(err)
	}

	expiresIn := int(time.Until(expiresAt).Seconds())
	ctrl.setAccessCookie(c, accessToken, expiresIn)
	ctrl.setRefreshCookie(c, refreshToken, 30*24*60*60)

	return response.OK(c, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}
