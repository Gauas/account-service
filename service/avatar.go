package service

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/gauas/account-service/packages/uploader"
)

func (s *Service) UpdateAvatar(ctx context.Context, userKey string, reader io.Reader, contentType, filename string) (string, error) {
	user, err := s.Repository.User.Take(ctx, "key = ?", userKey)
	if err != nil {
		return "", err
	}

	url, err := s.UploadReader(ctx, user.Key.String(), reader, contentType, filename)
	if err != nil {
		return "", err
	}

	user.AvatarURL = &url
	if err := s.Repository.User.Update(ctx, user); err != nil {
		return "", err
	}

	return url, nil
}

func (s *Service) UploadAvatarFromURL(ctx context.Context, seed, imageURL string) (string, error) {
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
		Reader:      reader,
		Filename:    filename,
		Bucket:      uploader.AVATAR_BUCKET,
		Path:        uploader.AVATAR_PATH,
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return up.URL, nil
}
