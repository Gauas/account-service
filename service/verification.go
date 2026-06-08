package service

import (
	"context"
	"errors"
	"time"

	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/model/types"
	"github.com/gauas/account-service/supports/oauth2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newVerification(data *oauth2.UserInfo) (*model.Verification, error) {
	if data.Email == nil || *data.Email == "" {
		return nil, nil
	}

	key, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	verification := &model.Verification{
		Key:        key,
		Method:     types.EmailVerification,
		Value:      string(*data.Email),
		IsVerified: data.EmailVerified,
	}
	if data.EmailVerified {
		now := time.Now().UTC()
		verification.VerifiedAt = &now
	}

	return verification, nil
}

func (s *Service) verifyEmail(ctx context.Context, userID int64, email types.Email) error {
	verification, err := s.Repository.Verification.Take(
		ctx,
		"user_id = ? AND method = ? AND value = ?",
		userID,
		types.EmailVerification,
		string(email),
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.createVerifiedEmail(ctx, userID, email)
	}
	if err != nil || verification.IsVerified {
		return err
	}

	now := time.Now().UTC()
	return s.Repository.Verification.WithContext(ctx).
		Model(&model.Verification{}).
		Where("id = ?", verification.ID).
		Updates(map[string]any{
			"is_verified": true,
			"verified_at": &now,
		}).Error
}

func (s *Service) createVerifiedEmail(ctx context.Context, userID int64, email types.Email) error {
	key, err := uuid.NewV7()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	_, err = s.Repository.Verification.Create(ctx, &model.Verification{
		Key:        key,
		UserID:     userID,
		Method:     types.EmailVerification,
		Value:      string(email),
		IsVerified: true,
		VerifiedAt: &now,
	})

	return err
}
