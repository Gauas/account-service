package info

import (
	"github.com/gauas/account-service/dto/response"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetInfo(c echo.Context) error {
	key := c.QueryParam("key")
	if key == "" {
		key = middlewares.UserID(c.Request().Context())
	}

	user, err := h.Service.GetInfo(c, key)
	if err != nil {
		return err
	}

	return httpresp.OK(c, httpresp.Refine[*model.User, response.ProfileResponse](user))
}
