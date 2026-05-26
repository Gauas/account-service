package controller

//func (ctrl *Controller) SendVerificationEmail(c echo.Context) error {
//	userIDStr := c.Param("user_id")
//	userID, err := uuid.Parse(userIDStr)
//	if err != nil {
//		return response.NewError(http.StatusBadRequest, "invalid user_id")
//	}
//
//	if err := ctrl.service.SendVerificationEmail(c.Request().Context(), userID); err != nil {
//		return response.Wrap(err)
//	}
//
//	return response.NoContent(c, "verification email sent")
//}
//
//func (ctrl *Controller) VerifyEmail(c echo.Context) error {
//	token := c.Param("token")
//	if token == "" {
//		return response.NewError(http.StatusBadRequest, "token is required")
//	}
//
//	if err := ctrl.service.VerifyEmail(c.Request().Context(), token); err != nil {
//		return response.Wrap(err)
//	}
//
//	return response.NoContent(c, "email verified successfully")
//}
