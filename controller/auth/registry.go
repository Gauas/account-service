package auth

import (
	"net/http"

	dtoReq "github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(c echo.Context) error {
	var req dtoReq.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
	}

	data, err := h.Service.NewAccount(c, req)
	if err != nil {
		return err
	}

	return httpresp.OK(c, data)
}
