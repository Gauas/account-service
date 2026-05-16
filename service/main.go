package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/repository"
	"github.com/gauas/account-service/supports"
	"github.com/google/uuid"
	"github.com/mozillazg/go-unidecode"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo   *repository.Registry
	config config.Config
	mq     *amqp.Channel
	cache  *redis.Client
}

func New(repo *repository.Registry, cfg config.Config, mq *amqp.Channel, cache *redis.Client) *Service {
	return &Service{repo: repo, config: cfg, mq: mq, cache: cache}
}

type RegisterRequest struct {
	Username    *string   `json:"username"`
	Password    string    `json:"password"`
	Email       *string   `json:"email"`
	Phone       *string   `json:"phone"`
	FullName    string    `json:"fullname"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*model.User, error) {
	if req.FullName == "" && req.Email == nil && req.Phone == nil {
		return nil, appError(http.StatusBadRequest, "fullname/email/phone is required")
	}
	if req.Password == "" {
		return nil, appError(http.StatusBadRequest, "password is required")
	}
	if req.Email != nil && !supports.IsEmail(*req.Email) {
		return nil, appError(http.StatusBadRequest, "invalid email format")
	}
	if req.Phone != nil && !supports.IsPhone(*req.Phone) {
		return nil, appError(http.StatusBadRequest, "invalid phone format")
	}

	hashed := hashPassword(req.Password)
	defaultAvatar := model.DefaultAvatarURL

	user := &model.User{
		UserID:      uuid.New(),
		Permission:  "member",
		Username:    req.Username,
		Password:    &hashed,
		Email:       req.Email,
		Phone:       req.Phone,
		FullName:    &req.FullName,
		Gender:      &req.Gender,
		DateOfBirth: &req.DateOfBirth,
		AvatarURL:   &defaultAvatar,
	}

	if req.Username == nil {
		username, err := s.generateUniqueUsername(ctx, req.FullName)
		if err != nil {
			return nil, err
		}
		user.Username = &username
	}

	tx := s.repo.User.DB().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if req.Email != nil && *req.Email != "" {
		ev := &model.UserVerification{ID: uuid.New(), UserID: user.UserID, Method: "email", Value: *req.Email}
		if err := tx.Create(ev).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if req.Phone != nil && *req.Phone != "" {
		pv := &model.UserVerification{ID: uuid.New(), UserID: user.UserID, Method: "phone", Value: *req.Phone}
		if err := tx.Create(pv).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	s.sendVerification(user, supports.Val(req.Email), req.FullName)

	return user, nil
}

func (s *Service) AuthenticateUser(ctx context.Context, identifierType, identifier, password string) (*model.User, error) {
	hashed := hashPassword(password)

	var field string
	switch strings.ToLower(identifierType) {
	case "email":
		field = "email"
	case "phone":
		field = "phone"
	case "username":
		field = "username"
	default:
		return nil, appError(http.StatusBadRequest, "invalid identifier type")
	}

	user, err := s.repo.User.FindOne(ctx, fmt.Sprintf("%s = ? AND password = ?", field), identifier, hashed)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, appError(http.StatusUnauthorized, "invalid credentials")
	}
	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.repo.User.FindOne(ctx, "user_id = ?", id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, appError(http.StatusNotFound, "user not found")
	}
	return user, nil
}

type UpdateProfileRequest struct {
	Username    *string    `json:"username"     patch:"username"`
	FullName    *string    `json:"fullname"     patch:"full_name"`
	Email       *string    `json:"email"        patch:"email"`
	Phone       *string    `json:"phone"        patch:"phone"`
	Gender      *string    `json:"gender"       patch:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth" patch:"date_of_birth"`
	FacebookURL *string    `json:"facebook_url" patch:"facebook_url"`
	GithubURL   *string    `json:"github_url"   patch:"github_url"`
	AvatarURL   *string    `json:"avatar_url"   patch:"avatar_url"`
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*model.User, error) {
	if _, err := s.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}

	data := supports.NewPatch()
	supports.Fill(data, req)

	if err := s.repo.User.DB().Model(&model.User{}).Where("user_id = ?", userID).Updates(data.Build()).Error; err != nil {
		return nil, err
	}
	return s.GetUserByID(ctx, userID)
}

func (s *Service) UpdateAvatarFromURL(ctx context.Context, userID uuid.UUID, username, imageURL string) (string, error) {
	fileBytes, contentType, err := downloadImage(imageURL)
	if err != nil {
		return "", err
	}
	return s.uploadAvatar(ctx, userID, username, fileBytes, contentType)
}

func (s *Service) UpdateAvatarFromBytes(ctx context.Context, userID uuid.UUID, username string, fileBytes []byte, contentType string) (string, error) {
	return s.uploadAvatar(ctx, userID, username, fileBytes, contentType)
}

func (s *Service) uploadAvatar(ctx context.Context, userID uuid.UUID, username string, fileBytes []byte, contentType string) (string, error) {
	ext := fileExtFromContentType(contentType, "")
	filename := fmt.Sprintf("%s.%s", avatarHash(username), ext)
	cdnURL, err := s.uploadToService(fileBytes, filename, contentType)
	if err != nil {
		return "", err
	}
	if err := s.repo.User.UpdateWhere(ctx, map[string]interface{}{"avatar_url": cdnURL}, "user_id = ?", userID); err != nil {
		return "", err
	}
	return cdnURL, nil
}

func (s *Service) issueToken(ctx context.Context, userID, email string) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)
	key := fmt.Sprintf("email_verification:%s", token)
	value := fmt.Sprintf("%s:%s", userID, email)
	if err := s.cache.Set(ctx, key, value, 24*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}
	return token, nil
}

func (s *Service) consumeToken(ctx context.Context, token string) (string, string, error) {
	key := fmt.Sprintf("email_verification:%s", token)
	value, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return "", "", appError(http.StatusBadRequest, "invalid or expired verification token")
	}
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return "", "", appError(http.StatusBadRequest, "malformed token data")
	}
	s.cache.Del(ctx, key)
	return parts[0], parts[1], nil
}

func (s *Service) SendVerificationEmail(ctx context.Context, userID uuid.UUID) error {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Email == nil || *user.Email == "" {
		return appError(http.StatusBadRequest, "user has no email address")
	}
	token, err := s.issueToken(ctx, userID.String(), *user.Email)
	if err != nil {
		return err
	}
	link := fmt.Sprintf("https://%s/v1/account/verify-email/%s", s.config.DomainName, token)
	name := supports.Val(user.FullName, "User")
	content := fmt.Sprintf("Hi %s,\n\nPlease verify your email by clicking the link below.\n\nThe link expires in 24 hours.", name)
	return s.sendEmail(ctx, "confirmation", *user.Email, name, content, link)
}

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	userIDStr, email, err := s.consumeToken(ctx, token)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return appError(http.StatusBadRequest, "invalid user id in token")
	}
	now := time.Now()
	return s.repo.Verification.UpdateWhere(ctx, map[string]interface{}{
		"is_verified": true,
		"verified_at": now,
	}, "user_id = ? AND method = ? AND value = ?", userID, "email", email)
}

func (s *Service) GetMFAByType(ctx context.Context, userID uuid.UUID, mfaType string) (*model.UserMFA, error) {
	return s.repo.MFA.FindOne(ctx, "user_id = ? AND type = ?", userID, mfaType)
}

func (s *Service) CreateMFA(ctx context.Context, mfa *model.UserMFA) (*model.UserMFA, error) {
	return s.repo.MFA.Create(ctx, mfa)
}

func (s *Service) UpdateMFA(ctx context.Context, mfa *model.UserMFA) (*model.UserMFA, error) {
	return s.repo.MFA.Update(ctx, mfa)
}

func (s *Service) generateUniqueUsername(ctx context.Context, fullName string) (string, error) {
	base := strings.ToLower(strings.ReplaceAll(unidecode.Unidecode(fullName), " ", ""))
	var existing []string
	if err := s.repo.User.Pluck(ctx, "username", &existing, "username LIKE ?", base+"%"); err != nil {
		return "", err
	}
	if len(existing) == 0 {
		return base, nil
	}
	max := -1
	baseLen := len(base)
	for _, u := range existing {
		if u == base {
			if max < 0 {
				max = 0
			}
			continue
		}
		if len(u) <= baseLen || !strings.HasPrefix(u, base) {
			continue
		}
		var n int
		if _, err := fmt.Sscanf(u[baseLen:], "%d", &n); err != nil {
			continue
		}
		if n > max {
			max = n
		}
	}
	if max < 0 {
		return base, nil
	}
	return fmt.Sprintf("%s%d", base, max+1), nil
}

type emailMessage struct {
	Type          string `json:"type"`
	Recipient     string `json:"recipient"`
	RecipientName string `json:"recipientName,omitempty"`
	Content       string `json:"content"`
	ActionUrl     string `json:"actionUrl,omitempty"`
}

func (s *Service) sendEmail(ctx context.Context, msgType, email, name, content, actionURL string) error {
	msg := emailMessage{Type: msgType, Recipient: email, RecipientName: name, Content: content, ActionUrl: actionURL}
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("service: marshal email: %w", err)
	}
	return s.mq.PublishWithContext(ctx, "email_exchange", "email."+msgType, false, false,
		amqp.Publishing{ContentType: "application/json", Body: body})
}

type uploadResponse struct {
	URL string `json:"url"`
}

func (s *Service) uploadToService(data []byte, filename, contentType string) (string, error) {
	url := s.config.UploadURL + "/v1/upload/file"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	_ = w.WriteField("bucket", "images")
	_ = w.WriteField("path", "avatar")
	h := map[string][]string{
		"Content-Disposition": {fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename)},
		"Content-Type":        {contentType},
	}
	fw, err := w.CreatePart(h)
	if err != nil {
		return "", err
	}
	fw.Write(data)
	w.Close()

	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Private-Key", s.config.PrivateKey)

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if !slices.Contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		return "", fmt.Errorf("upload service returned %d: %s", resp.StatusCode, supports.ReadBody(resp.Body))
	}
	var ur uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&ur); err != nil {
		return "", err
	}
	return ur.URL, nil
}
