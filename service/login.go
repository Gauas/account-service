package service

import (
	"time"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/middlewares"
	"github.com/labstack/echo/v4"
)

func (s *Service) Login(c echo.Context, req dto.LoginRequest) (echo.Map, error) {
	ctx := c.Request().Context()
	user, err := s.Repository.User.Take(ctx, "email = ?", req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || hashPassword(req.Password) != *user.Password {
		return nil, echo.NewHTTPError(401, "invalid credentials")
	}

	tokens, err := s.Infra.AuthSDK.CreateToken(ctx, user.UserID, user.Permission, middlewares.DeviceID(ctx))
	if err != nil {
		return nil, err
	}

	s.SetCookie(c, "access_token", tokens.AccessToken, time.Until(tokens.ExpiresAt))
	s.SetCookie(c, "refresh_token", tokens.RefreshToken, time.Until(tokens.RefreshExpiresAt))

	return echo.Map{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    tokens.ExpiresIn,
	}, nil
}
