package auth

import (
	"net/http"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) OAuth2(c echo.Context) error {
	var req request.Oauth2Request

	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request")
	}

	data, err := h.Service.TryOAuth2(c, req)
	if err != nil {
		return err
	}

	return httpresp.OK(c, data)
}
