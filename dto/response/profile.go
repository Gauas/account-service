package response

type ProfileResponse struct {
	Key         string `json:"key"`
	FullName    string `json:"full_name"`
	Dob         string `json:"dob"`
	Gender      string `json:"gender"`
	AvatarURL   string `json:"avatar_url"`
	IsOnboarded bool   `json:"is_onboarded"`
	Permission  string `json:"permission"`
}
