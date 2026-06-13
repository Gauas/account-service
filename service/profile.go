package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/model/types"
	"github.com/gauas/account-service/packages/httpresp"
	"gorm.io/gorm"
)

func (s *Service) UpdateProfile(ctx context.Context, userKey string, fullName *string, dob *time.Time, gender *types.Gender) error {
	user, err := s.Repository.User.Take(ctx, "key = ?", userKey)
	if err != nil {
		return err
	}

	if fullName != nil {
		user.FullName = fullName
	}
	if dob != nil {
		user.Dob = dob
	}
	if gender != nil {
		user.Gender = gender
	}

	err = s.Repository.User.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProfileByKey(ctx context.Context, key string) (*model.User, error) {
	if key == "" {
		return nil, httpresp.NewError(http.StatusBadRequest, "key is required")
	}

	user, err := s.Repository.User.Take(ctx, "key = ?", key)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httpresp.NewError(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
