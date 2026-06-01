package uploader

import (
	"context"
	"io"

	upload "github.com/gauas/upload-service/sdk"
)

const (
	AVATAR_BUCKET = "images"
	AVATAR_PATH   = "avatar"
)

type Client interface {
	UploadFile(ctx context.Context, req upload.UploadRequest) (*upload.UploadResponse, error)
}

type Request struct {
	Reader   io.Reader
	Filename string
	Bucket   string
	Path     string
}

func Reader(ctx context.Context, client Client, req Request) (*upload.UploadResponse, error) {
	return client.UploadFile(ctx, upload.UploadRequest{
		Reader:   req.Reader,
		Filename: req.Filename,
		Bucket:   req.Bucket,
		Path:     req.Path,
	})
}
