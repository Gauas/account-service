package auth

import (
	"net/http"

	"github.com/gauas/account-service/dto"
	response2 "github.com/gauas/account-service/supports/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) OAuth2(c echo.Context) error {
	var req dto.Oauth2Request

	if err := c.Bind(&req); err != nil {
		return response2.NewError(http.StatusBadRequest, "invalid request")
	}

	data, err := h.Service.TryOAuth2(c, req)
	if err != nil {
		return err
	}

	return response2.OK(c, data)
}
