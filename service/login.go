package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gauas/account-service/model/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *Service) Login(ctx context.Context, email types.Email, password string, deviceID string) (*Session, error) {
	email = email.Normalize()

	identity, err := s.Repository.Identity.Take(ctx, "provider = ? AND email = ?", types.EmailIdentityProvider, string(email))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError(http.StatusUnauthorized, "unauthorized")
	}

	if err != nil {
		return nil, err
	}

	if identity.Hash == nil || bcrypt.CompareHashAndPassword([]byte(*identity.Hash), []byte(password)) != nil {
		return nil, appError(http.StatusUnauthorized, "unauthorized")
	}

	user, err := s.Repository.User.Take(ctx, "id = ?", identity.UserID)
	if err != nil {
		return nil, err
	}

	return s.OpenSession(ctx, user, deviceID)
}
