package service

import (
	"fmt"
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

	filename := avatarFilename(user.Key.String(), file)
	up, err := uploader.Upload(ctx, s.Infra.Upload, uploader.Request{
		Reader:   src,
		Filename: filename,
		Bucket:   uploader.AVATAR_BUCKET,
		Path:     uploader.AVATAR_PATH,
	})
	if err != nil {
		return "", err
	}

	user.AvatarURL = &up.URL
	if err := s.Repository.User.Update(ctx, user); err != nil {
		return "", err
	}

	return up.URL, nil
}

func avatarFilename(seed string, file *multipart.FileHeader) string {
	ext := fileExtFromContentType(file.Header.Get("Content-Type"), file.Filename)
	return fmt.Sprintf("%s.%s", avatarHash(seed), ext)
}
