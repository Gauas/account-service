package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gauas/account-service/model/types"
)

type GoogleUserInfo struct {
	Sub           string      `json:"sub"`
	Name          string      `json:"name"`
	GivenName     string      `json:"given_name"`
	FamilyName    string      `json:"family_name"`
	Picture       string      `json:"picture"`
	Email         types.Email `json:"email"`
	EmailVerified bool        `json:"email_verified"`
	Locale        string      `json:"locale"`
}

func TryGoogle(token string) (*GoogleUserInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	data, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call google api: %w", err)
	}
	defer data.Body.Close()

	if data.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api returned status: %d", data.StatusCode)
	}

	var response GoogleUserInfo
	if err := json.NewDecoder(data.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode google response: %w", err)
	}

	return &response, nil

}
