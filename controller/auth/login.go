package auth

import (
	"net/http"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/gauas/account-service/packages/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	var req request.LoginRequest

	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request")
	}

	sessionData, err := h.Service.Login(c.Request().Context(), req, middlewares.DeviceID(c.Request().Context()))
	if err != nil {
		return err
	}

	session.Write(c, h.Config, sessionData)
	return httpresp.OK(c, session.Response(sessionData))
}
