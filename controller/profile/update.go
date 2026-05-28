package profile

import (
	"net/http"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateProfile(c echo.Context) error {
	var req dto.UpdateProfileRequest

	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request")
	}

	err := h.Service.UpdateProfile(c, req)
	if err != nil {
		return err
	}

	return response.OK(c, "successfully updated")
}
