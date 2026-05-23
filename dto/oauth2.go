package dto

type Oauth2Request struct {
	Provider string `json:"provider"`
	Token    string `json:"token"`
}
