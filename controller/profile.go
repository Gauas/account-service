package controller

//func (ctrl *Controller) UpdateAvatar(c echo.Context) error {
//	userID, err := userIDFromContext(c)
//	if err != nil {
//		return response.NewError(http.StatusUnauthorized, err.Error())
//	}
//
//	user, err := ctrl.service.GetUserByID(c.Request().Context(), userID)
//	if err != nil {
//		return response.Wrap(err)
//	}
//
//	username := ""
//	if user.Username != nil {
//		username = *user.Username
//	}
//
//	file, err := c.FormFile("file")
//	if err == nil {
//		src, err := file.Open()
//		if err != nil {
//			return response.NewError(http.StatusBadRequest, "failed to open uploaded file")
//		}
//		defer src.Close()
//
//		buf := make([]byte, file.Size)
//		if _, err := src.Read(buf); err != nil {
//			return response.NewError(http.StatusBadRequest, "failed to read uploaded file")
//		}
//
//		cdnURL, err := ctrl.service.UpdateAvatarFromBytes(c.Request().Context(), userID, username, buf, file.Header.Get("Content-Type"))
//		if err != nil {
//			return response.Wrap(err)
//		}
//
//		return response.OK(c, echo.Map{"avatar_url": cdnURL})
//	}
//
//	var body struct {
//		ImageURL string `json:"image_url"`
//	}
//	if err := c.Bind(&body); err != nil || body.ImageURL == "" {
//		return response.NewError(http.StatusBadRequest, "provide file or image_url")
//	}
//
//	cdnURL, err := ctrl.service.UpdateAvatarFromURL(c.Request().Context(), userID, username, body.ImageURL)
//	if err != nil {
//		return response.Wrap(err)
//	}
//
//	return response.OK(c, echo.Map{"avatar_url": cdnURL})
//}
