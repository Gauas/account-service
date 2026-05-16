package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/gauas/account-service/service"
	"github.com/gauas/account-service/packages/response"
)

type registerRequest struct {
	Username    *string   `json:"username"`
	Password    string    `json:"password"`
	Email       *string   `json:"email"`
	Phone       *string   `json:"phone"`
	FullName    string    `json:"fullname"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

func (ctrl *Controller) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request body")
	}

	user, err := ctrl.service.Register(c.Request().Context(), service.RegisterRequest{
		Username:    req.Username,
		Password:    req.Password,
		Email:       req.Email,
		Phone:       req.Phone,
		FullName:    req.FullName,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
	})
	if err != nil {
		return response.Wrap(err)
	}

	return response.Created(c, echo.Map{"message": "registration successful", "user_id": user.UserID})
}

type loginRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password *string `json:"password"`
}

func (ctrl *Controller) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return response.NewError(http.StatusBadRequest, "invalid request body")
	}
	if req.Password == nil {
		return response.NewError(http.StatusBadRequest, "password is required")
	}

	deviceID := deviceIDFromRequest(c, "")
	if deviceID == "" {
		return response.NewError(http.StatusBadRequest, "X-Device-ID header is required")
	}

	identifierType, identifier := resolveIdentifier(req)
	if identifierType == "" {
		return response.NewError(http.StatusBadRequest, "email/username/phone is required")
	}

	accessToken, refreshToken, expiresAt, err := ctrl.service.Login(c.Request().Context(), identifierType, identifier, *req.Password, deviceID)
	if err != nil {
		return response.Wrap(err)
	}

	expiresIn := int(time.Until(expiresAt).Seconds())
	ctrl.setAccessCookie(c, accessToken, expiresIn)
	ctrl.setRefreshCookie(c, refreshToken, 30*24*60*60)

	return response.OK(c, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}

func (ctrl *Controller) Logout(c echo.Context) error {
	refreshToken := c.Request().Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		if cookie, err := c.Cookie("refresh_token"); err == nil {
			refreshToken = cookie.Value
		}
	}
	if refreshToken == "" {
		return response.NewError(http.StatusBadRequest, "no refresh token provided")
	}

	deviceID := deviceIDFromRequest(c, "")
	if deviceID == "" {
		return response.NewError(http.StatusBadRequest, "X-Device-ID header is required")
	}

	if err := ctrl.service.RevokeToken(c.Request().Context(), refreshToken, deviceID); err != nil {
		return response.Wrap(err)
	}

	c.SetCookie(&http.Cookie{Name: "access_token", Value: "", MaxAge: -1, Path: "/"})
	c.SetCookie(&http.Cookie{Name: "refresh_token", Value: "", MaxAge: -1, Path: "/"})

	return response.NoContent(c, "logout successful")
}
