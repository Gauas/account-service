package relationship

import (
	"net/http"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Cancel(c echo.Context) error {
	var req request.RelationshipRequest
	if err := c.Bind(&req); err != nil {
		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.Service.CancelRelationship(c.Request().Context(), middlewares.UserID(c.Request().Context()), req.Target); err != nil {
		return err
	}

	return httpresp.NoContent(c, "relationship request cancelled")
}
