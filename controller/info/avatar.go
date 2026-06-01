package info

import (
	"net/http"

	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateAvatar(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return httpresp.NewError(http.StatusBadRequest, "file is required")
	}

	url, err := h.Service.UpdateAvatar(c, file)
	if err != nil {
		return err
	}

	return httpresp.OK(c, echo.Map{"avatar_url": url})
}
