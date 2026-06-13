package relationship

//
//func (h *Handler) Decline(c echo.Context) error {
//	var req request.Relationship
//	if err := c.Bind(&req); err != nil {
//		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
//	}
//	if err := req.Validate(); err != nil {
//		return httpresp.NewError(http.StatusBadRequest, err.Error())
//	}
//
//	if err := h.Service.DeclineRelationship(c.Request().Context(), middlewares.UserID(c.Request().Context()), req.Target); err != nil {
//		return err
//	}
//
//	return httpresp.NoContent(c, "relationship request declined")
//}
