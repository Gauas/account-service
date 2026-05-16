package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/supports"
)

type createTokenRequest struct {
	UserID     uuid.UUID `json:"user_id"`
	Permission string    `json:"permission"`
}

type createTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (s *Service) CreateToken(ctx context.Context, userID uuid.UUID, permission, deviceID string) (string, string, time.Time, error) {
	if deviceID == "" {
		return "", "", time.Time{}, appError(http.StatusBadRequest, "device_id is required")
	}

	body, err := json.Marshal(createTokenRequest{UserID: userID, Permission: permission})
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.AuthorizationURL+"/v1/authorization/token", bytes.NewBuffer(body))
	if err != nil {
		return "", "", time.Time{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Device-ID", deviceID)
	req.Header.Set("Secret-Key", s.config.PrivateKey)

	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("authorization service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", time.Time{}, appError(resp.StatusCode, "authorization service: "+supports.ReadBody(resp.Body))
	}

	var res struct {
		Data createTokenResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", "", time.Time{}, err
	}

	expiry := time.Now().Add(time.Duration(res.Data.ExpiresIn) * time.Second)
	return res.Data.AccessToken, res.Data.RefreshToken, expiry, nil
}

func (s *Service) RevokeToken(ctx context.Context, refreshToken, deviceID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, s.config.AuthorizationURL+"/v1/authorization/token", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Refresh-Token", refreshToken)
	req.Header.Set("X-Device-ID", deviceID)
	req.Header.Set("Secret-Key", s.config.PrivateKey)

	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return fmt.Errorf("authorization service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return appError(resp.StatusCode, "authorization service: "+supports.ReadBody(resp.Body))
	}

	return nil
}

func (s *Service) ValidateAccessToken(ctx context.Context, token string) error {
	url := s.config.AuthorizationURL + "/v1/authorization/token/validate?token=" + token

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Secret-Key", s.config.PrivateKey)

	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return fmt.Errorf("authorization service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return appError(http.StatusUnauthorized, "invalid or expired token")
	}

	return nil
}

type googleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (s *Service) LoginWithGoogle(ctx context.Context, googleToken string) (*model.User, error) {
	info, err := fetchGoogleUserInfo(googleToken)
	if err != nil {
		return nil, appError(http.StatusUnauthorized, "invalid Google token")
	}
	if !supports.IsEmail(info.Email) {
		return nil, appError(http.StatusBadRequest, "invalid email from Google")
	}

	found, err := s.repo.User.FindOne(ctx, "email = ?", info.Email)
	if err != nil {
		return nil, err
	}

	if found != nil {
		if info.EmailVerified {
			_ = s.ensureEmailVerified(ctx, found.UserID, info.Email)
		}
		return found, nil
	}

	user, err := s.Register(ctx, RegisterRequest{
		Email:    &info.Email,
		FullName: info.Name,
		Password: "sso-" + info.Sub,
	})
	if err != nil {
		return nil, err
	}

	if info.Picture != "" {
		go func() {
			username := supports.Val(user.Username)
			_, _ = s.UpdateAvatarFromURL(context.Background(), user.UserID, username, info.Picture)
		}()
	}

	if info.EmailVerified {
		_ = s.ensureEmailVerified(ctx, user.UserID, info.Email)
	}

	return user, nil
}

func (s *Service) ensureEmailVerified(ctx context.Context, userID uuid.UUID, email string) error {
	now := time.Now()
	return s.repo.Verification.UpdateWhere(ctx, map[string]interface{}{
		"is_verified": true,
		"verified_at": now,
	}, "user_id = ? AND method = ? AND value = ?", userID, "email", email)
}

func fetchGoogleUserInfo(token string) (*googleUserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api returned %d", resp.StatusCode)
	}

	var info googleUserInfo
	return &info, json.NewDecoder(resp.Body).Decode(&info)
}
