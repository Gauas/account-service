package service

import (
	"github.com/gauas/account-service/dto"
	"github.com/labstack/echo/v4"
)

func (s *Service) Login(c echo.Context, req dto.LoginRequest) (echo.Map, error) {
	ctx := c.Request().Context()
	identity, err := s.Repository.Identity.Take(ctx, "email = ?", req.Email)
	if err != nil {
		return nil, err
	}

	hashed, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if hashed != *identity.Hash {
		return nil, echo.ErrUnauthorized
	}

	user, err := s.Repository.User.Take(ctx, "id = ?", identity.UserID)
	if err != nil {
		return nil, err
	}

	return s.TryAuthorize(c, user)
}
