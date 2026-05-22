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

	return s.Infra.AuthSDK.RevokeToken(ctx, refreshToken, middlewares.DeviceID(ctx))
}
