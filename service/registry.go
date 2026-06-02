package service

import (
	"context"
	"net/http"
	"strings"

	dtoReq "github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Service) NewAccount(c echo.Context, req dtoReq.RegisterRequest) (echo.Map, error) {
	err := error(nil)
	ctx := c.Request().Context()

	if err = req.Email.Validate(); err != nil {
		return nil, appError(http.StatusBadRequest, err.Error())
	}

	user := &model.User{
		Permission: "member",
		FullName:   &req.FullName,
	}

	if user.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	if req.Email == "" {
		return nil, appError(http.StatusBadRequest, "email is required")
	}

	identity := &model.Identity{
		Provider:       types.EmailIdentityProvider,
		ProviderUserID: strings.ToLower(strings.TrimSpace(string(req.Email))),
		Email:          &req.Email,
	}

	if identity.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	verification := &model.Verification{
		Method: types.EmailVerification,
		Value:  string(req.Email),
	}

	if verification.Key, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	if req.Password == "" {
		return nil, appError(http.StatusBadRequest, "password is required")
	}

	hashed, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	identity.Hash = &hashed

	err = s.Repository.Transaction(ctx, func(ctx context.Context) error {
		if s.Repository.Identity.Exists(ctx, "email = ?", string(req.Email)) {
			return appError(http.StatusConflict, "account already exists")
		}

		if _, err = s.Repository.User.Create(ctx, user); err != nil {
			return err
		}

		identity.UserID = user.ID
		verification.UserID = user.ID

		if _, err = s.Repository.Identity.Create(ctx, identity); err != nil {
			return err
		}

		if _, err = s.Repository.Verification.Create(ctx, verification); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.TryAuthorize(c, user)
}

