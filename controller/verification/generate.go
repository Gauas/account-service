package verification

import (
	"net/http"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Generate(c echo.Context) error {
	var req request.GenerateVerification
	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
	}

	if err := req.Validate(); err != nil {
		return httpresp.NewError(http.StatusBadRequest, err.Error())
	}

	if err := h.Service.GenerateVerification(c.Request().Context(), middlewares.UserID(c.Request().Context()), req.Type, req.Target); err != nil {
		return err
	}

	return httpresp.NoContent(c, "verification code sent")
}
