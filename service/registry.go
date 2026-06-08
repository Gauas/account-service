package service

import (
	"context"
	"errors"
	"net/http"

	dtoReq "github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Service) NewAccount(c echo.Context, req dtoReq.RegisterRequest) (echo.Map, error) {
	err := error(nil)
	ctx := c.Request().Context()
	email := req.Email.Normalize()

	if err = email.Validate(); err != nil {
		return nil, appError(http.StatusBadRequest, err.Error())
	}

	user := &model.User{
		Permission: "member",
		FullName:   &req.FullName,
	}

	if user.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	if email == "" {
		return nil, appError(http.StatusBadRequest, "email is required")
	}

	identity := &model.Identity{
		Provider:       types.EmailIdentityProvider,
		ProviderUserID: string(email),
		Email:          &email,
	}

	if identity.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	ver := &model.Verification{
		Method: types.EmailVerification,
		Value:  string(email),
	}

	if ver.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	if req.Password == "" {
		return nil, appError(http.StatusBadRequest, "password is required")
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	identity.Hash = &hash
	userID := int64(0)

	err = s.Repository.Transaction(ctx, func(ctx context.Context) error {
		old, err := s.Repository.Identity.Take(ctx, "email = ?", string(email))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.SaveAccount(ctx, user, identity, ver)
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
		return s.OpenSession(c, userID)
	}

	return s.TryAuthorize(c, user)
}

func (s *Service) SaveAccount(ctx context.Context, user *model.User, identity *model.Identity, verification *model.Verification) error {
	if _, err := s.Repository.User.Create(ctx, user); err != nil {
		return err
	}

	identity.UserID = user.ID
	verification.UserID = user.ID

	if _, err := s.Repository.Identity.Create(ctx, identity); err != nil {
		return err
	}

	_, err := s.Repository.Verification.Create(ctx, verification)
	return err
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
