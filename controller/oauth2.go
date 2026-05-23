package controller

import (
	"net/http"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (ctrl *Controller) LoginWithGoogle(c echo.Context) error {
	var req dto.Oauth2Request

	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request")
	}

	data, err := ctrl.service.TryWithGoogle(c, req)
	if err != nil {
		return response.Wrap(err)
	}

	return response.OK(c, data)
}
