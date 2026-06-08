package service

import (
	"context"
	"fmt"
)

func (s *Service) Logout(ctx context.Context, refreshToken, deviceID string) error {
	if refreshToken == "" {
		return fmt.Errorf("no refresh token provided")
	}

	return s.Infra.AuthSDK.RevokeToken(ctx, refreshToken, deviceID)
}
