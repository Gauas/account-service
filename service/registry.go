package service

import (
	"context"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) error {
	hashed := hashPassword(req.Password)

	user := &model.User{
		UserID:      uuid.New(),
		Permission:  "member",
		Password:    &hashed,
		Email:       req.Email,
		Phone:       req.Phone,
		FullName:    &req.FullName,
		Gender:      &req.Gender,
		DateOfBirth: &req.DateOfBirth,
	}

	return s.Repository.User.DB().
		WithContext(ctx).
		Transaction(
			func(tx *gorm.DB) error {
				if err := tx.Create(user).Error; err != nil {
					return err
				}

				if err := req.Email.Validate(); err != nil {
					return err
				}

				if req.Phone.Validate() == nil {
					model.NewVerification(string(*req.Phone), model.PhoneVerification, user.UserID)
				}

				return nil
			})
}
