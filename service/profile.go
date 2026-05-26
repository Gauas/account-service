package service

import (
	"github.com/gauas/account-service/dto"
	middleware "github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/supports"
	"github.com/labstack/echo/v4"
)

func (s *Service) UpdateProfile(c echo.Context, req dto.UpdateProfileRequest) error {
	ctx := c.Request().Context()

	user, err := s.Repository.User.Take(ctx, "id = ?", middleware.UserID(ctx))
	if err != nil {
		return err
	}

	if err := supports.Fill(user, req); err != nil {
		return err
	}

	err = s.Repository.User.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProfile(c echo.Context, id string) (*model.User, error) {
	ctx := c.Request().Context()
	user, err := s.Repository.User.Take(ctx, "id = ?", id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
