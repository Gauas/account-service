package auth

import (
	"net/http"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) LoginWithOAuth2(c echo.Context) error {
	var req dto.Oauth2Request

	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request")
	}

	data, err := h.Service.TryOAuth2(c, req)
	if err != nil {
		return response.Wrap(err)
	}

	return response.OK(c, data)
}
