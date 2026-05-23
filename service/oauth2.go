package service

import (
	"context"
	"errors"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/supports/oauth2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Service) TryWithGoogle(c echo.Context, req dto.Oauth2Request) (echo.Map, error) {
	ctx := c.Request().Context()
	err := error(nil)

	data, err := oauth2.TryGoogle(req.Token)
	if err != nil {
		return nil, err
	}

	identity, err := s.Repository.Identity.Take(ctx, "provider = ? AND provider_user_id = ?", "google", data.Sub)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if identity == nil {
		return s.NewOAuthAccount(c, data)
	}

	user, err := s.Repository.User.Take(ctx, "id = ?", identity.UserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return s.TryAuthorize(c, user)
}

func (s *Service) NewOAuthAccount(c echo.Context, data *oauth2.GoogleUserInfo) (echo.Map, error) {
	err := error(nil)
	ctx := c.Request().Context()

	user := &model.User{
		Permission: "member",
		FullName:   &data.Name,
		AvatarURL:  &data.Picture,
	}

	if user.ID, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	identity := &model.Identity{
		UserID:         user.ID,
		Provider:       "google",
		ProviderUserID: data.Sub,
		Email:          &data.Email,
	}

	if identity.ID, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	err = s.Repository.Transaction(ctx, func(ctx context.Context) error {
		if _, err := s.Repository.User.Create(ctx, user); err != nil {
			return err
		}

		if _, err := s.Repository.Identity.Create(ctx, identity); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.TryAuthorize(c, user)
}
