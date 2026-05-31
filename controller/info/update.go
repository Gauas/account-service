package info

import (
	"net/http"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateInfo(c echo.Context) error {
	var req request.UpdateProfileRequest

	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request")
	}

	err := h.Service.UpdateInfo(c, req)
	if err != nil {
		return err
	}

	return httpresp.OK(c, "successfully updated")
}
