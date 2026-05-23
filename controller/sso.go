package controller

type loginWithGoogleRequest struct {
	Token string `json:"token"`
}

//func (ctrl *Controller) LoginWithGoogle(c echo.Context) error {
//	var req loginWithGoogleRequest
//	if err := c.Bind(&req); err != nil {
//		return response.NewError(http.StatusBadRequest, "invalid request body")
//	}
//
//	deviceID := deviceIDFromRequest(c, "")
//	if deviceID == "" {
//		return response.NewError(http.StatusBadRequest, "X-Device-ID header is required")
//	}
//
//	user, err := ctrl.service.LoginWithGoogle(c.Request().Context(), req.Token)
//	if err != nil {
//		return response.Wrap(err)
//	}
//
//	accessToken, refreshToken, expiresAt, err := ctrl.service.CreateToken(c.Request().Context(), user.UserID, user.Permission, deviceID)
//	if err != nil {
//		return response.Wrap(err)
//	}
//
//	expiresIn := int(time.Until(expiresAt).Seconds())
//	ctrl.setAccessCookie(c, accessToken, expiresIn)
//	ctrl.setRefreshCookie(c, refreshToken, 30*24*60*60)
//
//	return response.OK(c, echo.Map{
//		"access_token":  accessToken,
//		"refresh_token": refreshToken,
//		"expires_in":    expiresIn,
//	})
//}
