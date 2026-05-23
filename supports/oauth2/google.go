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

//
//func TryGoogle(token string) (*GoogleUserInfo, error) {
//	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create request: %w", err)
//	}
//	req.Header.Set("Authorization", "Bearer "+token)
//
//	client := &http.Client{}
//	data, err := client.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("failed to call google api: %w", err)
//	}
//	defer data.Body.Close()
//
//	if data.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("google api returned status: %d", data.StatusCode)
//	}
//
//	var response GoogleUserInfo
//	if err := json.NewDecoder(data.Body).Decode(&response); err != nil {
//		return nil, fmt.Errorf("failed to decode google response: %w", err)
//	}
//
//	return &response, nil
//
//}

type GoogleProvider struct{}
