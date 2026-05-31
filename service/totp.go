package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/dto/response"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/model/types"
	"github.com/gauas/account-service/packages/mfa"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Service) GenerateTOTP(c echo.Context) (*response.TOTPSetupResponse, error) {
	ctx := c.Request().Context()

	user, err := s.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	identity, err := s.Repository.Identity.Take(ctx, "user_id = ? AND email IS NOT NULL", user.ID)
	if err != nil {
		return nil, err
	}
	if identity.Email == nil || *identity.Email == "" {
		return nil, appError(http.StatusBadRequest, "user has no email")
	}

	account := mfa.AccountName(user.Key.String())
	secret, qrURL, err := mfa.BuildKey(account)
	if err != nil {
		return nil, appError(http.StatusInternalServerError, "failed to generate totp key")
	}

	if err := s.upsertTOTP(ctx, user.ID, secret); err != nil {
		return nil, err
	}

	return &response.TOTPSetupResponse{
		Email:   string(*identity.Email),
		QRURL:   qrURL,
		Secret:  secret,
		Account: account,
		Issuer:  mfa.ISSUER,
	}, nil
}

func (s *Service) EnableTOTP(c echo.Context, req request.EnableTOTPRequest) error {
	ctx := c.Request().Context()

	user, err := s.CurrentUser(ctx)
	if err != nil {
		return err
	}
	if req.OTPCode == "" {
		return appError(http.StatusBadRequest, "otp_code is required")
	}

	totp, err := s.getTOTP(ctx, user.ID)
	if err != nil {
		return err
	}
	if totp == nil || totp.Secret == nil {
		return appError(http.StatusBadRequest, "no totp setup found")
	}
	if totp.Enabled {
		return appError(http.StatusConflict, "totp already enabled")
	}
	if !mfa.Verify(req.OTPCode, *totp.Secret, time.Now().UTC()) {
		return appError(http.StatusBadRequest, "invalid otp_code")
	}

	now := time.Now().UTC()
	totp.Enabled = true
	totp.VerifiedAt = &now

	return s.Repository.MFA.Update(ctx, totp)
}

func (s *Service) VerifyTOTP(c echo.Context, req request.VerifyTOTPRequest) (echo.Map, error) {
	ctx := c.Request().Context()

	user, err := s.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.OTPCode == "" {
		return nil, appError(http.StatusBadRequest, "otp_code is required")
	}

	totp, err := s.getTOTP(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if totp == nil || totp.Secret == nil || !totp.Enabled {
		return nil, appError(http.StatusBadRequest, "totp is not enabled")
	}
	if !mfa.Verify(req.OTPCode, *totp.Secret, time.Now().UTC()) {
		return nil, appError(http.StatusBadRequest, "invalid otp_code")
	}

	return s.TryAuthorize(c, user)
}

func (s *Service) getTOTP(ctx context.Context, userID int64) (*model.MFA, error) {
	totp, err := s.Repository.MFA.Take(ctx, "user_id = ? AND type = ?", userID, types.MFATypeTOTP)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return totp, nil
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
	if err != nil {
		return err
	}

	return nil
}
