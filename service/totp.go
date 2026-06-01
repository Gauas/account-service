package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/model/types"
	mfapkg "github.com/gauas/account-service/packages/mfa"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TOTPSetup struct {
	QRURL   string `json:"qr_code"`
	Secret  string `json:"secret"`
	Account string `json:"account"`
	Issuer  string `json:"issuer"`
}

func (s *Service) GenerateTOTP(c echo.Context) (*TOTPSetup, error) {
	ctx := c.Request().Context()

	user, err := s.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	account := mfapkg.AccountName(user.Key.String())
	secret, qrURL, err := mfapkg.BuildKey(account)
	if err != nil {
		return nil, appError(http.StatusInternalServerError, "failed to generate totp key")
	}

	if err := s.upsertTOTP(ctx, user.ID, secret); err != nil {
		return nil, err
	}

	return &TOTPSetup{
		QRURL:   qrURL,
		Secret:  secret,
		Account: account,
		Issuer:  mfapkg.ISSUER,
	}, nil
}

func (s *Service) EnableTOTP(c echo.Context, otpCode string) error {
	ctx := c.Request().Context()

	user, err := s.CurrentUser(ctx)
	if err != nil {
		return err
	}
	if otpCode == "" {
		return appError(http.StatusBadRequest, "otp_code is required")
	}

	mfa, err := s.getTOTP(ctx, user.ID)
	if err != nil {
		return err
	}
	if mfa == nil || mfa.Secret == nil {
		return appError(http.StatusBadRequest, "no totp setup found")
	}
	if mfa.Enabled {
		return appError(http.StatusConflict, "totp already enabled")
	}
	if !mfapkg.Verify(otpCode, *mfa.Secret, time.Now().UTC()) {
		return appError(http.StatusBadRequest, "invalid otp_code")
	}

	now := time.Now().UTC()
	mfa.Enabled = true
	mfa.VerifiedAt = &now

	return s.Repository.MFA.Update(ctx, mfa)
}

func (s *Service) VerifyTOTP(c echo.Context, otpCode string) (echo.Map, error) {
	ctx := c.Request().Context()

	user, err := s.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if otpCode == "" {
		return nil, appError(http.StatusBadRequest, "otp_code is required")
	}

	mfa, err := s.getTOTP(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if mfa == nil || mfa.Secret == nil || !mfa.Enabled {
		return nil, appError(http.StatusBadRequest, "totp is not enabled")
	}
	if !mfapkg.Verify(otpCode, *mfa.Secret, time.Now().UTC()) {
		return nil, appError(http.StatusBadRequest, "invalid otp_code")
	}

	return s.TryAuthorize(c, user)
}

func (s *Service) getTOTP(ctx context.Context, userID int64) (*model.MFA, error) {
	mfa, err := s.Repository.MFA.Take(ctx, "user_id = ? AND type = ?", userID, types.MFATypeTOTP)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mfa, nil
}

func (s *Service) upsertTOTP(ctx context.Context, userID int64, secret string) error {
	mfa, err := s.getTOTP(ctx, userID)
	if err != nil {
		return err
	}
	if mfa == nil {
		return s.createTOTP(ctx, userID, secret)
	}
	if mfa.Enabled {
		return appError(http.StatusConflict, "totp already enabled")
	}

	mfa.Secret = &secret
	return s.Repository.MFA.Update(ctx, mfa)
}

func (s *Service) createTOTP(ctx context.Context, userID int64, secret string) error {
	key, err := uuid.NewV7()
	if err != nil {
		return err
	}

	_, err = s.Repository.MFA.Create(ctx, &model.MFA{
		Key:    key,
		UserID: userID,
		Type:   types.MFATypeTOTP,
		Secret: &secret,
	})
	return err
}
