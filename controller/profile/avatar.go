package profile

import (
	"net/http"

	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateAvatar(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return httpresp.NewError(http.StatusBadRequest, "file is required")
	}

	src, err := file.Open()
	if err != nil {
		return httpresp.NewError(http.StatusBadRequest, "failed to open uploaded file")
	}
	defer src.Close()

	url, err := h.Service.UpdateAvatar(
		c.Request().Context(),
		middlewares.UserID(c.Request().Context()),
		src,
		file.Header.Get("Content-Type"),
		file.Filename,
	)
	if err != nil {
		return err
	}

	return httpresp.OK(c, echo.Map{"avatar_url": url})
}
