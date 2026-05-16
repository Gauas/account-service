package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gauas/account-service/packages/response"
)

func appError(code int, msg string) error {
	return response.NewError(code, msg)
}

func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
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
