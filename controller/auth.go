package controller

type loginRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password *string `json:"password"`
}

//func (ctrl *Controller) Logout(c echo.Context) error {
//	refreshToken := c.Request().Header.Get("X-Refresh-Token")
//	if refreshToken == "" {
//		if cookie, err := c.Cookie("refresh_token"); err == nil {
//			refreshToken = cookie.Value
//		}
//	}
//	if refreshToken == "" {
//		return response.NewError(http.StatusBadRequest, "no refresh token provided")
//	}
//
//	deviceID := deviceIDFromRequest(c, "")
//	if deviceID == "" {
//		return response.NewError(http.StatusBadRequest, "X-Device-ID header is required")
//	}
//
//	if err := ctrl.Config.RevokeToken(c.Request().Context(), refreshToken, deviceID); err != nil {
//		return response.Wrap(err)
//	}
//
//	c.SetCookie(&http.Cookie{Name: "access_token", Value: "", MaxAge: -1, Path: "/"})
//	c.SetCookie(&http.Cookie{Name: "refresh_token", Value: "", MaxAge: -1, Path: "/"})
//
//	return response.NoContent(c, "logout successful")
//}
