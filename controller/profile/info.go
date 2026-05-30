package profile

import (
	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/supports/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUserInfo(c echo.Context) error {
	key := c.QueryParam("key")
	if key == "" {
		key = middlewares.UserID(c.Request().Context())
	}

	user, err := h.Service.GetProfile(c, key)
	if err != nil {
		return err
	}

	return response.OK(c, response.Refine[*model.User, dto.ProfileResponse](user))
}
