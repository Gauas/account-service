package service

import (
	"errors"

	"github.com/gauas/account-service/dto"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *Service) Login(c echo.Context, req dto.LoginRequest) (echo.Map, error) {
	ctx := c.Request().Context()
	identity, err := s.Repository.Identity.Take(ctx, "email = ?", req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, echo.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}

	if identity.Hash == nil || bcrypt.CompareHashAndPassword([]byte(*identity.Hash), []byte(req.Password)) != nil {
		return nil, echo.ErrUnauthorized
	}

	user, err := s.Repository.User.Take(ctx, "id = ?", identity.UserID)
	if err != nil {
		return nil, err
	}

	return s.TryAuthorize(c, user)
}
