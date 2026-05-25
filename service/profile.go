package service

import (
	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/supports"
	"github.com/labstack/echo/v4"
)

func (s *Service) UpdateProfile(c echo.Context, req dto.UpdateProfileRequest) error {
	ctx := c.Request().Context()

	user, err := s.Repository.User.Take(ctx, "id = ?", c.Get("user_id"))
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
