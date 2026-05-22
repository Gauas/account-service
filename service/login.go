package service

import (
	"context"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/middlewares"
	"github.com/labstack/echo/v4"
)

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (echo.Map, error) {
	user, err := s.Repository.User.Take(ctx, "email = ?", req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || hashPassword(req.Password) != *user.Password {
		return nil, echo.NewHTTPError(401, "invalid credentials")
	}

	deviceID, _ := ctx.Value(middlewares.DeviceIDKey).(string)

	tokenPair, err := s.Infra.AuthSDK.CreateToken(ctx, user.UserID, user.Permission, deviceID)
	if err != nil {
		return nil, err
	}

	return echo.Map{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
	}, nil
}
