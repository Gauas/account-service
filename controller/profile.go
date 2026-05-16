package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/gauas/account-service/service"
	"github.com/gauas/account-service/packages/response"
)

func (ctrl *Controller) GetProfile(c echo.Context) error {
	userID, err := userIDFromContext(c)
	if err != nil {
		return response.NewError(http.StatusUnauthorized, err.Error())
	}

	user, err := ctrl.service.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return response.Wrap(err)
	}

	return response.OK(c, user)
}

type updateProfileRequest struct {
	Username    *string    `json:"username"`
	FullName    *string    `json:"fullname"`
	Email       *string    `json:"email"`
	Phone       *string    `json:"phone"`
	Gender      *string    `json:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	FacebookURL *string    `json:"facebook_url"`
	GithubURL   *string    `json:"github_url"`
}

func (ctrl *Controller) UpdateProfile(c echo.Context) error {
	userID, err := userIDFromContext(c)
	if err != nil {
		return response.NewError(http.StatusUnauthorized, err.Error())
	}

	var req updateProfileRequest
	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request body")
	}

	user, err := ctrl.service.UpdateProfile(c.Request().Context(), userID, service.UpdateProfileRequest{
		Username:    req.Username,
		FullName:    req.FullName,
		Email:       req.Email,
		Phone:       req.Phone,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
		FacebookURL: req.FacebookURL,
		GithubURL:   req.GithubURL,
	})
	if err != nil {
		return response.Wrap(err)
	}

	return response.OK(c, user)
}

func (ctrl *Controller) UpdateAvatar(c echo.Context) error {
	userID, err := userIDFromContext(c)
	if err != nil {
		return response.NewError(http.StatusUnauthorized, err.Error())
	}

	user, err := ctrl.service.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return response.Wrap(err)
	}

	username := ""
	if user.Username != nil {
		username = *user.Username
	}

	file, err := c.FormFile("file")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return response.NewError(http.StatusBadRequest, "failed to open uploaded file")
		}
		defer src.Close()

		buf := make([]byte, file.Size)
		if _, err := src.Read(buf); err != nil {
			return response.NewError(http.StatusBadRequest, "failed to read uploaded file")
		}

		cdnURL, err := ctrl.service.UpdateAvatarFromBytes(c.Request().Context(), userID, username, buf, file.Header.Get("Content-Type"))
		if err != nil {
			return response.Wrap(err)
		}

		return response.OK(c, echo.Map{"avatar_url": cdnURL})
	}

	var body struct {
		ImageURL string `json:"image_url"`
	}
	if err := c.Bind(&body); err != nil || body.ImageURL == "" {
		return response.NewError(http.StatusBadRequest, "provide file or image_url")
	}

	cdnURL, err := ctrl.service.UpdateAvatarFromURL(c.Request().Context(), userID, username, body.ImageURL)
	if err != nil {
		return response.Wrap(err)
	}

	return response.OK(c, echo.Map{"avatar_url": cdnURL})
}
