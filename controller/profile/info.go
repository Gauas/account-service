package profile

import (
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUserInfo(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		id = c.Get("user_id").(string)
	}
	user, err := h.Service.GetProfile(c, id)
	if err != nil {
		return response.Wrap(err)
	}

	return response.OK(c, user)
}
