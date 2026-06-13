package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/model/types"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/gauas/account-service/publisher"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Service) GenerateVerification(ctx context.Context, userKey string, method types.VerificationMethod, target string) error {
	user, err := s.GetProfileByKey(ctx, userKey)
	if err != nil {
		return err
	}
	if !s.Repository.Identity.Exists(ctx, "user_id = ? AND email = ?", user.ID, target) {
		return httpresp.NewError(http.StatusForbidden, "email does not belong to user")
	}

	code, err := types.NewCode()
	if err != nil {
		return err
	}

	key, err := uuid.NewV7()
	if err != nil {
		return err
	}

	verification, err := s.Repository.Verification.Take(ctx, "user_id = ? AND method = ? AND target = ?", user.ID, method, target)
	if err == nil && verification.IsVerified {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		verification.Value = string(code)
		err = s.Repository.Verification.Update(ctx, verification)
	} else {
		verification = &model.Verification{
			Key:    key,
			UserID: user.ID,
			Method: method,
			Target: target,
			Value:  string(code),
		}
		verification, err = s.Repository.Verification.Create(ctx, verification)
	}
	if err != nil {
		return err
	}

	message := publisher.EmailMessage{
		Type: "verification",
		To:   target,
		Data: map[string]any{"code": verification.Value},
	}
	return s.Publisher.Email.Send(ctx, message)
}

func (s *Service) TryVerification(ctx context.Context, userKey string, method types.VerificationMethod, target string, code types.Code) error {
	user, err := s.GetProfileByKey(ctx, userKey)
	if err != nil {
		return err
	}
	if !s.Repository.Identity.Exists(ctx, "user_id = ? AND email = ?", user.ID, target) {
		return httpresp.NewError(http.StatusForbidden, "email does not belong to user")
	}

	verification, err := s.Repository.Verification.Take(ctx, "user_id = ? AND method = ? AND target = ?", user.ID, method, target)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return httpresp.NewError(http.StatusNotFound, "verification not found")
	}
	if err != nil {
		return err
	}
	if verification.IsVerified {
		return nil
	}
	if verification.Value != string(code) {
		return httpresp.NewError(http.StatusBadRequest, "invalid verification code")
	}

	now := time.Now().UTC()
	verification.IsVerified = true
	verification.VerifiedAt = &now

	return s.Repository.Verification.Update(ctx, verification)
}

func (s *Service) verifyEmail(ctx context.Context, userID int64, email types.Email) error {
	email = email.Normalize()

	verification, err := s.Repository.Verification.Take(ctx, "user_id = ? AND method = ? AND target = ?", userID, types.EmailVerification, string(email))
	if err == nil {
		now := time.Now().UTC()
		verification.IsVerified = true
		verification.VerifiedAt = &now

		return s.Repository.Verification.Update(ctx, verification)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	code, err := types.NewCode()
	if err != nil {
		return err
	}

	key, err := uuid.NewV7()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	verification = &model.Verification{
		Key:        key,
		UserID:     userID,
		Method:     types.EmailVerification,
		Target:     string(email),
		Value:      string(code),
		IsVerified: true,
		VerifiedAt: &now,
	}
	_, err = s.Repository.Verification.Create(ctx, verification)

	return err
}
