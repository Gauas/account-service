package mfas

import (
	"net/http"

	dtoReq "github.com/gauas/account-service/dto/request"
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
	var req dtoReq.EnableTOTPRequest
	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.Service.EnableTOTP(c, req); err != nil {
		return err
	}

	return httpresp.NoContent(c, "TOTP enabled successfully")
}

func (h *Handler) VerifyTOTP(c echo.Context) error {
	var req dtoReq.VerifyTOTPRequest
	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
	}

	data, err := h.Service.VerifyTOTP(c, req)
	if err != nil {
		return err
	}

	return httpresp.OK(c, data)
}
