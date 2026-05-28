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

func (s *Service) TryOAuth2(c echo.Context, req dto.Oauth2Request) (echo.Map, error) {
	ctx := c.Request().Context()

	provider, ok := oauth2.Providers[req.Provider]
	if !ok {
		return nil, echo.NewHTTPError(400, "unsupported oauth2 provider")
	}

	data, err := provider.GetUser(req.Token)
	if err != nil {
		return nil, err
	}

	identity, err := s.Repository.Identity.Take(ctx, "provider = ? AND provider_user_id = ?", data.Provider, data.ProviderUserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.NewOAuthAccount(c, data)
	}
	if err != nil {
		return nil, err
	}

	user, err := s.Repository.User.Take(ctx, "id = ?", identity.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, echo.NewHTTPError(404, "user not found")
	}

	return s.TryAuthorize(c, user)
}

func (s *Service) NewOAuthAccount(c echo.Context, data *oauth2.UserInfo) (echo.Map, error) {
	err := error(nil)
	ctx := c.Request().Context()

	user := &model.User{
		Permission: "member",
		FullName:   &data.Name,
	}

	if user.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	avatarURL, err := s.oauthAvatarURL(c, user.Key.String(), data.AvatarURL)
	if err != nil {
		return nil, err
	}
	if avatarURL != nil {
		user.AvatarURL = avatarURL
	}

	identity := &model.Identity{
		Provider:       data.Provider,
		ProviderUserID: data.ProviderUserID,
		Email:          data.Email,
	}

	if identity.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	err = s.Repository.Transaction(ctx,
		func(ctx context.Context) error {
			if _, err = s.Repository.User.Create(ctx, user); err != nil {
				return err
			}

			identity.UserID = user.ID
			if _, err = s.Repository.Identity.Create(ctx, identity); err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return s.TryAuthorize(c, user)
}

func (s *Service) oauthAvatarURL(c echo.Context, seed, sourceURL string) (*string, error) {
	if sourceURL == "" {
		return nil, nil
	}

	url, err := s.UploadAvatarFromURL(c.Request().Context(), seed, sourceURL)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
