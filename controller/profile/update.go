package profile

import (
	"net/http"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Update(c echo.Context) error {
	var req request.UpdateProfile

	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request")
	}

	if err := req.Validate(); err != nil {
		return httpresp.NewError(http.StatusBadRequest, err.Error())
	}

	err := h.Service.UpdateProfile(c.Request().Context(), middlewares.UserID(c.Request().Context()), req.FullName, req.Dob, req.Gender)
	if err != nil {
		return err
	}

	return httpresp.OK(c, "successfully updated")
}
