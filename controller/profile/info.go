package profile

import (
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUserInfo(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		id = middlewares.UserID(c.Request().Context())
	}

	user, err := h.Service.GetProfile(c, id)
	if err != nil {
		return err
	}

	return response.OK(c, user)
}
