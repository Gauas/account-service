package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/packages/httpresp"
	"golang.org/x/crypto/bcrypt"
)

func appError(code int, msg string) error {
	return httpresp.NewError(code, msg)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (s *Service) OpenSession(ctx context.Context, user *model.User, deviceID string) (*Session, error) {
	tokens, err := s.Infra.AuthSDK.CreateToken(ctx, user.Key, user.Permission, deviceID)
	if err != nil {
		return nil, err
	}

	return &Session{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		ExpiresIn:        tokens.ExpiresIn,
		ExpiresAt:        tokens.ExpiresAt,
		RefreshExpiresAt: tokens.RefreshExpiresAt,
	}, nil
}

func avatarHash(username string) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s_%d", username, time.Now().UnixNano())))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func fileExtFromContentType(contentType, fallbackURL string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/webp":
		return "webp"
	case "image/gif":
		return "gif"
	default:
		if fallbackURL != "" {
			if idx := strings.LastIndex(fallbackURL, "."); idx >= 0 {
				return fallbackURL[idx+1:]
			}
		}
		return "jpg"
	}
}

func downloadImage(imageURL string) ([]byte, string, error) {
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Get(imageURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("image download returned %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	ct := normalizeContentType(resp.Header.Get("Content-Type"))
	if isGenericContentType(ct) {
		ct = http.DetectContentType(data)
	}
	return data, ct, nil
}

func normalizeContentType(raw string) string {
	if raw == "" {
		return ""
	}
	mediaType, _, err := mime.ParseMediaType(raw)
	if err != nil {
		return strings.TrimSpace(strings.ToLower(raw))
	}
	return strings.TrimSpace(strings.ToLower(mediaType))
}

func isGenericContentType(ct string) bool {
	if ct == "" {
		return true
	}
	return ct == "application/octet-stream"
}
