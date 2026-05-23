package oauth2

import "github.com/gauas/account-service/model/types"

type UserInfo struct {
	Provider types.IdentityProvider

	ProviderUserID string

	Email         *types.Email
	EmailVerified bool

	Name      string
	AvatarURL string
}

type Provider interface {
	GetUser(token string) (*UserInfo, error)
}

var Providers = map[string]Provider{
	"google": GoogleProvider{},
}
