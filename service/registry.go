package service

import (
	"context"

	"github.com/gauas/account-service/dto"
	"github.com/gauas/account-service/model"
	"github.com/google/uuid"
)

func (s *Service) Registry(ctx context.Context, req dto.RegisterRequest) error {
	return s.Repository.Transaction(ctx,
		func(ctx context.Context) error {
			hashed := hashPassword(req.Password)

			user := &model.User{
				UserID:      uuid.New(),
				Permission:  "member",
				Password:    &hashed,
				Email:       req.Email,
				FullName:    &req.FullName,
				Gender:      &req.Gender,
				DateOfBirth: &req.DateOfBirth,
			}

			if _, err := s.Repository.User.Create(ctx, user); err != nil {
				return err
			}

			if err := req.Email.Validate(); err != nil {
				return err
			}

			emailVerification := &model.Verification{
				ID:     uuid.New(),
				UserID: user.UserID,
				Method: model.EmailVerification,
				Value:  string(*user.Email),
			}

			if _, err := s.Repository.Verification.Create(ctx, emailVerification); err != nil {
				return err
			}

			return nil
		},
	)
}
