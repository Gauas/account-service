package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gauas/account-service/model/types"
)

type googleUserResponse struct {
	Sub           string      `json:"sub"`
	Name          string      `json:"name"`
	Picture       string      `json:"picture"`
	Email         types.Email `json:"email"`
	EmailVerified bool        `json:"email_verified"`
}

func (GoogleProvider) GetUser(token string) (*UserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call google api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api returned status: %d", resp.StatusCode)
	}

	var data googleUserResponse

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode google response: %w", err)
	}

	return &UserInfo{
		Provider:       "google",
		ProviderUserID: data.Sub,
		Email:          &data.Email,
		EmailVerified:  data.EmailVerified,
		Name:           data.Name,
		AvatarURL:      data.Picture,
	}, nil
}

type GoogleProvider struct{}
