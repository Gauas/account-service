package service

import (
	"fmt"

	"github.com/gauas/account-service/middlewares"
	"github.com/labstack/echo/v4"
)

func (s *Service) Logout(c echo.Context) error {
	ctx := c.Request().Context()
	refreshToken := middlewares.RefreshToken(c)
	if refreshToken == "" {
		return fmt.Errorf("no refresh token provided")
	}

	s.SetCookie(c, "access_token", "", 0)
	s.SetCookie(c, "refresh_token", "", 0)

	return s.Infra.AuthSDK.RevokeToken(ctx, refreshToken, middlewares.DeviceID(ctx))
}
