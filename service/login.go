package service

import (
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

	tokenPair, err := s.Infra.AuthSDK.CreateToken(ctx, user.UserID, user.Permission, middlewares.DeviceID(ctx))
	if err != nil {
		return nil, err
	}

	return echo.Map{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
	}, nil
}
