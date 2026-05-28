package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/uploader"
	"github.com/labstack/echo/v4"
)

func (s *Service) UpdateAvatar(c echo.Context, file *multipart.FileHeader) (string, error) {
	ctx := c.Request().Context()
	if file == nil {
		return "", appError(http.StatusBadRequest, "file is required")
	}

	user, err := s.Repository.User.Take(ctx, "key = ?", middlewares.UserID(ctx))
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", appError(http.StatusBadRequest, "failed to open uploaded file")
	}
	defer src.Close()

	url, err := s.UploadReader(ctx, user.Key.String(), src, file.Header.Get("Content-Type"), file.Filename)
	if err != nil {
		return "", err
	}

	user.AvatarURL = &url
	if err := s.Repository.User.Update(ctx, user); err != nil {
		return "", err
	}

	return url, nil
}

func avatarFilename(seed string, file *multipart.FileHeader) string {
	ext := fileExtFromContentType(file.Header.Get("Content-Type"), file.Filename)
	return fmt.Sprintf("%s.%s", avatarHash(seed), ext)
}

func (s *Service) UploadAvatarFromURL(ctx context.Context, seed, imageURL string) (string, error) {
	if imageURL == "" {
		return "", appError(http.StatusBadRequest, "image_url is required")
	}

	data, contentType, err := downloadImage(imageURL)
	if err != nil {
		return "", err
	}

	return s.UploadReader(ctx, seed, bytes.NewReader(data), contentType, imageURL)
}

func avatarFilenameFromMeta(seed, contentType, fallback string) string {
	ext := fileExtFromContentType(contentType, fallback)
	return fmt.Sprintf("%s.%s", avatarHash(seed), ext)
}

func (s *Service) UploadReader(ctx context.Context, seed string, reader io.Reader, contentType, fallbackName string) (string, error) {
	filename := avatarFilenameFromMeta(seed, contentType, fallbackName)
	up, err := uploader.Reader(ctx, s.Infra.Upload, uploader.Request{
		Reader:   reader,
		Filename: filename,
		Bucket:   uploader.AVATAR_BUCKET,
		Path:     uploader.AVATAR_PATH,
	})
	if err != nil {
		return "", err
	}

	return up.URL, nil
}
