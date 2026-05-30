package info

import (
	"net/http"

	response2 "github.com/gauas/account-service/supports/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateAvatar(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response2.NewError(http.StatusBadRequest, "file is required")
	}

	url, err := h.Service.UpdateAvatar(c, file)
	if err != nil {
		return err
	}

	return response2.OK(c, echo.Map{"avatar_url": url})
}
