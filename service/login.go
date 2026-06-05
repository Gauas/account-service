package service

import (
	"errors"

	dtoReq "github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/model/types"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *Service) Login(c echo.Context, req dtoReq.LoginRequest) (echo.Map, error) {
	ctx := c.Request().Context()
	email := req.Email.Normalize()

	identity, err := s.Repository.Identity.Take(ctx, "provider = ? AND email = ?", types.EmailIdentityProvider, string(email))
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
