package controller

import (
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (ctrl *Controller) Logout(c echo.Context) error {
	err := ctrl.service.Logout(c)
	if err != nil {
		return response.Wrap(err)
	}

	return response.OK(c, "logged out successfully")
}
