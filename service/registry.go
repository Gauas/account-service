package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Service) NewAccount(ctx context.Context, email types.Email, password string, fullName string, deviceID string) (*Session, error) {
	err := error(nil)
	email = email.Normalize()

	user := &model.User{
		Permission: "member",
		FullName:   &fullName,
	}

	if user.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	identity := &model.Identity{
		Provider:       types.EmailIdentityProvider,
		ProviderUserID: string(email),
		Email:          &email,
	}

	if identity.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	identity.Hash = &hash
	userID := int64(0)

	err = s.Repository.Transaction(ctx, func(ctx context.Context) error {
		old, err := s.Repository.Identity.Take(ctx, "email = ?", string(email))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.SaveAccount(ctx, user, identity)
		}

		if err != nil {
			return err
		}

		if old.Provider == types.EmailIdentityProvider {
			return appError(http.StatusConflict, "account already exists")
		}

		userID = old.UserID
		return s.LinkAccount(ctx, userID, email, hash)
	})
	if err != nil {
		return nil, err
	}

	if userID != int64(0) {
		return s.OpenSessionByID(ctx, userID, deviceID)
	}

	return s.OpenSession(ctx, user, deviceID)
}

func (s *Service) OpenSessionByID(ctx context.Context, userID int64, deviceID string) (*Session, error) {
	user, err := s.Repository.User.Take(ctx, "id = ?", userID)
	if err != nil {
		return nil, err
	}

	return s.OpenSession(ctx, user, deviceID)
}

func (s *Service) SaveAccount(ctx context.Context, user *model.User, identity *model.Identity) error {
	if _, err := s.Repository.User.Create(ctx, user); err != nil {
		return err
	}

	identity.UserID = user.ID

	if _, err := s.Repository.Identity.Create(ctx, identity); err != nil {
		return err
	}

	return nil
}

func (s *Service) LinkAccount(ctx context.Context, userID int64, email types.Email, hash string) error {
	key, err := uuid.NewV7()
	if err != nil {
		return err
	}

	if s.Repository.Identity.Exists(ctx, "user_id = ? AND provider = ?", userID, types.EmailIdentityProvider) {
		return appError(http.StatusConflict, "account already linked with email")
	}

	_, err = s.Repository.Identity.Create(ctx, &model.Identity{
		Key:            key,
		UserID:         userID,
		Provider:       types.EmailIdentityProvider,
		ProviderUserID: string(email),
		Email:          &email,
		Hash:           &hash,
	})

	return err
}
