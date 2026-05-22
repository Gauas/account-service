package controller

import (
	"net/http"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (ctrl *Controller) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request body")
	}

	err := ctrl.service.Register(c.Request().Context(), req)
	if err != nil {
		return response.Wrap(err)
	}

	return response.Created(c, "account created successfully")
}
