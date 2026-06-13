package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/supports/oauth2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Service) TryOAuth2(ctx context.Context, providerName string, token string, deviceID string) (*Session, error) {
	provider := oauth2.Providers[providerName]
	data, err := provider.GetUser(token)

	if err != nil {
		return nil, err
	}

	if data.Email != nil {
		email := data.Email.Normalize()
		data.Email = &email
	}

	identity, err := s.Repository.Identity.Take(ctx, "provider = ? AND provider_user_id = ?", data.Provider, data.ProviderUserID)
	if err == nil {
		return s.OpenSessionByID(ctx, identity.UserID, deviceID)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if data.Email == nil || *data.Email == "" || !data.EmailVerified {
		return s.NewOAuthAccount(ctx, data, deviceID)
	}

	identity, err = s.Repository.Identity.Take(ctx, "email = ?", string(*data.Email))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.NewOAuthAccount(ctx, data, deviceID)
	}

	if err != nil {
		return nil, err
	}

	if err = s.LinkIdentify(ctx, identity.UserID, data); err != nil {
		return nil, err
	}

	return s.OpenSessionByID(ctx, identity.UserID, deviceID)
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

func (s *Service) NewOAuthAccount(ctx context.Context, data *oauth2.UserInfo, deviceID string) (*Session, error) {
	err := error(nil)

	user := &model.User{
		Permission: "member",
		FullName:   &data.Name,
	}

	if user.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	avatarURL, err := s.SyncAvatar(ctx, user.Key.String(), data.AvatarURL)
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

			if data.Email == nil || *data.Email == "" || !data.EmailVerified {
				return nil
			}

			return s.verifyEmail(ctx, user.ID, *data.Email)
		},
	)
	if err != nil {
		return nil, err
	}

	return s.OpenSession(ctx, user, deviceID)
}

func (s *Service) SyncAvatar(ctx context.Context, seed, sourceURL string) (*string, error) {
	if sourceURL == "" {
		return nil, nil
	}

	url, err := s.UploadAvatarFromURL(ctx, seed, sourceURL)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
