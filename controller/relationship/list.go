package relationship

//func (h *Handler) List(c echo.Context) error {
//	var req request.ListRelationships
//	if err := req.Validate(c.QueryParam("status")); err != nil {
//		return httpresp.NewError(http.StatusBadRequest, err.Error())
//	}
//
//	user, err := h.Service.Repository.User.Take(c.Request().Context(), "key = ?", middlewares.UserID(c.Request().Context()))
//	if err != nil {
//		return err
//	}
//
//	data, err := h.Service.ListRelationships(c.Request().Context(), user, req.Status)
//	if err != nil {
//		return err
//	}
//
//	return httpresp.OK(c, data)
//}
