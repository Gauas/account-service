package auth

import (
	"net/http"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/gauas/account-service/packages/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) OAuth2(c echo.Context) error {
	var req request.OAuth2

	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request")
	}

	if err := req.Validate(); err != nil {
		return httpresp.NewError(http.StatusBadRequest, err.Error())
	}

	sessionData, err := h.Service.TryOAuth2(c.Request().Context(), req.Provider, req.Token, middlewares.DeviceID(c.Request().Context()))
	if err != nil {
		return err
	}

	session.Write(c, h.Config, sessionData)
	return httpresp.OK(c, session.Response(sessionData))
}
