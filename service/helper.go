package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/packages/response"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func appError(code int, msg string) error {
	return response.NewError(code, msg)
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

func (s *Service) TryAuthorize(c echo.Context, user *model.User) (echo.Map, error) {
	tokens, err := s.Infra.AuthSDK.CreateToken(c.Request().Context(), user.Key, user.Permission, middlewares.DeviceID(c.Request().Context()))
	if err != nil {
		return nil, err
	}

	s.SetCookie(c, "access_token", tokens.AccessToken, time.Until(tokens.ExpiresAt))
	s.SetCookie(c, "refresh_token", tokens.RefreshToken, time.Until(tokens.RefreshExpiresAt))

	return echo.Map{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    tokens.ExpiresIn,
	}, nil
}

func (s *Service) SetCookie(c echo.Context, name string, value string, ttl time.Duration) {
	http.SetCookie(c.Response().Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   s.Config.Cookie.DomainName,
		Secure:   s.Config.Cookie.Secure,
		HttpOnly: s.Config.Cookie.HttpOnly,
		SameSite: s.Config.Cookie.SameSite,
		MaxAge:   int(ttl.Seconds()),
		Expires:  time.Now().Add(ttl),
	})
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
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = http.DetectContentType(data)
	}
	return data, ct, nil
}
