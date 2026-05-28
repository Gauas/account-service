package profile

import (
	"net/http"

	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateAvatar(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.NewError(http.StatusBadRequest, "file is required")
	}

	url, err := h.Service.UpdateAvatar(c, file)
	if err != nil {
		return err
	}

	return response.OK(c, echo.Map{"avatar_url": url})
}
