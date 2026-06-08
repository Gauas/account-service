package types

type IdentityProvider string

const (
	EmailIdentityProvider    IdentityProvider = "email"
	GoogleIdentityProvider   IdentityProvider = "google"
	FacebookIdentityProvider IdentityProvider = "facebook"
)
