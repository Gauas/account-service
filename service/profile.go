package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gauas/account-service/dto/request"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/supports"
	"gorm.io/gorm"
)

func (s *Service) UpdateProfile(ctx context.Context, userKey string, req request.UpdateProfileRequest) error {
	user, err := s.Repository.User.Take(ctx, "key = ?", userKey)
	if err != nil {
		return err
	}

	if err := supports.Fill(user, req); err != nil {
		return err
	}

	err = s.Repository.User.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProfileByKey(ctx context.Context, key string) (*model.User, error) {
	if key == "" {
		return nil, appError(http.StatusBadRequest, "key is required")
	}

	user, err := s.Repository.User.Take(ctx, "key = ?", key)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
