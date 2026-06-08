package types

import "errors"

type VerificationMethod string

const (
	EmailVerification VerificationMethod = "email"
	PhoneVerification VerificationMethod = "phone"
)

func (v VerificationMethod) Validate() error {
	switch v {
	case EmailVerification, PhoneVerification:
		return nil
	}

	return errors.New("invalid verification method")
}
