package uploader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

const (
	AVATAR_BUCKET = "gauas"
	AVATAR_PATH   = "images/avatar"
	UPLOAD_PATH   = "/v1/upload"
)

type Client struct {
	baseURL    string
	secretKey  string
	httpClient *http.Client
}

type Request struct {
	Reader      io.Reader
	Filename    string
	Bucket      string
	Path        string
	ContentType string
}

type Response struct {
	URL         string `json:"url"`
	FilePath    string `json:"file_path"`
	FileHash    string `json:"file_hash"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Duplicated  bool   `json:"duplicated"`
}

type envelope struct {
	Status int             `json:"status"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
}

func New(baseURL, secretKey string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		secretKey:  strings.TrimSpace(secretKey),
		httpClient: http.DefaultClient,
	}
}

func Reader(ctx context.Context, client *Client, req Request) (*Response, error) {
	body, contentType, err := multipartBody(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, client.baseURL+UPLOAD_PATH, &body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", contentType)
	httpReq.Header.Set("Secret-Key", client.secretKey)

	httpRes, err := client.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	raw, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode < http.StatusOK || httpRes.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("upload service returned %d: %s", httpRes.StatusCode, string(raw))
	}

	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, err
	}
	if env.Error != "" {
		return nil, fmt.Errorf("upload service error: %s", env.Error)
	}

	var out Response
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func multipartBody(req Request) (bytes.Buffer, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("bucket", req.Bucket); err != nil {
		return bytes.Buffer{}, "", err
	}
	if req.Path != "" {
		if err := writer.WriteField("path", req.Path); err != nil {
			return bytes.Buffer{}, "", err
		}
	}

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, req.Filename))
	if req.ContentType != "" {
		partHeader.Set("Content-Type", req.ContentType)
	}

	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return bytes.Buffer{}, "", err
	}
	if _, err := io.Copy(part, req.Reader); err != nil {
		return bytes.Buffer{}, "", err
	}
	if err := writer.Close(); err != nil {
		return bytes.Buffer{}, "", err
	}

	return body, writer.FormDataContentType(), nil
}
