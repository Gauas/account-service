package relationship

import (
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) List(c echo.Context) error {
	user, err := h.Service.Repository.User.Take(c.Request().Context(), "key = ?", middlewares.UserID(c.Request().Context()))
	if err != nil {
		return err
	}

	data, err := h.Service.ListRelationships(c.Request().Context(), user, c.QueryParam("status"))
	if err != nil {
		return err
	}

	return httpresp.OK(c, data)
}
