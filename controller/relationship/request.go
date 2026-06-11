package relationship

//func (h *Handler) Request(c echo.Context) error {
//	var req request.Relationship
//	if err := c.Bind(&req); err != nil {
//		return httpresp.NewError(http.StatusBadRequest, "invalid request body")
//	}
//	if err := req.Validate(); err != nil {
//		return httpresp.NewError(http.StatusBadRequest, err.Error())
//	}
//
//	data, err := h.Service.RequestRelationship(c.Request().Context(), middlewares.UserID(c.Request().Context()), req.Target)
//	if err != nil {
//		return err
//	}
//
//	return httpresp.Created(c, data)
//}
