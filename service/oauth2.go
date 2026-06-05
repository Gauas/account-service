package service

import (
	"context"
	"errors"
	"net/http"

	dtoReq "github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/supports/oauth2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Service) TryOAuth2(c echo.Context, req dtoReq.Oauth2Request) (echo.Map, error) {
	ctx := c.Request().Context()

	provider, ok := oauth2.Providers[req.Provider]
	if !ok {
		return nil, echo.NewHTTPError(400, "unsupported oauth2 provider")
	}

	data, err := provider.GetUser(req.Token)
	if err != nil {
		return nil, err
	}
	if data.Email != nil {
		email := data.Email.Normalize()
		data.Email = &email
	}

	identity, err := s.Repository.Identity.Take(ctx, "provider = ? AND provider_user_id = ?", data.Provider, data.ProviderUserID)
	if err == nil {
		return s.OpenSession(c, identity.UserID)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if data.Email == nil || *data.Email == "" || !data.EmailVerified {
		return s.NewOAuthAccount(c, data)
	}

	identity, err = s.Repository.Identity.Take(ctx, "email = ?", string(*data.Email))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.NewOAuthAccount(c, data)
	}
	if err != nil {
		return nil, err
	}

	if err = s.LinkIdentify(ctx, identity.UserID, data); err != nil {
		return nil, err
	}

	return s.OpenSession(c, identity.UserID)
}

func (s *Service) OpenSession(c echo.Context, userID int64) (echo.Map, error) {
	user, err := s.Repository.User.Take(c.Request().Context(), "id = ?", userID)
	if err != nil {
		return nil, err
	}

	return s.TryAuthorize(c, user)
}

func (s *Service) LinkIdentify(ctx context.Context, userID int64, data *oauth2.UserInfo) error {
	key, err := uuid.NewV7()
	if err != nil {
		return err
	}

	identity := &model.Identity{
		Key:            key,
		UserID:         userID,
		Provider:       data.Provider,
		ProviderUserID: data.ProviderUserID,
		Email:          data.Email,
	}

	return s.Repository.Transaction(ctx, func(ctx context.Context) error {
		if s.Repository.Identity.Exists(ctx, "user_id = ? AND provider = ?", userID, identity.Provider) {
			return appError(http.StatusConflict, "account already linked with provider")
		}

		if _, err = s.Repository.Identity.Create(ctx, identity); err != nil {
			return err
		}

		return s.verifyEmail(ctx, userID, *data.Email)
	})
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

	avatarURL, err := s.SyncAvatar(c, user.Key.String(), data.AvatarURL)
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

	verification, err := newVerification(data)
	if err != nil {
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
			if verification == nil {
				return nil
			}

			verification.UserID = user.ID
			if _, err = s.Repository.Verification.Create(ctx, verification); err != nil {
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

func (s *Service) SyncAvatar(c echo.Context, seed, sourceURL string) (*string, error) {
	if sourceURL == "" {
		return nil, nil
	}

	url, err := s.UploadAvatarFromURL(c.Request().Context(), seed, sourceURL)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
