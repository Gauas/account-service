package types

import (
	"errors"
	"slices"
	"strings"

	"github.com/gauas/account-service/supports"
)

type Email string

func (e Email) Normalize() Email {
	return Email(strings.ToLower(strings.TrimSpace(string(e))))
}

func (e Email) Validate() error {
	if e == "" {
		return nil
	}

	if !supports.IsEmail(string(e)) {
		return errors.New("invalid email")
	}

	return nil
}

type Phone string

func (p Phone) Validate() error {
	if !supports.IsPhone(string(p)) && p != "" {
		return errors.New("invalid phone")
	}

	return nil
}

type Gender string

func (g Gender) Validate() error {
	if !slices.Contains([]string{"male", "female", "other"}, string(g)) {
		return errors.New("invalid gender")
	}

	return nil
}

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

type IdentityProvider string

const (
	EmailIdentityProvider    IdentityProvider = "email"
	GoogleIdentityProvider   IdentityProvider = "google"
	FacebookIdentityProvider IdentityProvider = "facebook"
)

type MFAType string

const (
	MFATypeTOTP         MFAType = "totp"
	MFATypeSMSOTP       MFAType = "sms_otp"
	MFATypeEmailOTP     MFAType = "email_otp"
	MFATypePush         MFAType = "push"
	MFATypeSecurityKey  MFAType = "security_key"
	MFATypePasskey      MFAType = "passkey"
	MFATypeHOTP         MFAType = "hotp"
	MFATypeBackupCode   MFAType = "backup_code"
	MFATypeVoiceCallOTP MFAType = "voice_call_otp"
)

func (m MFAType) Validate() error {
	switch m {
	case MFATypeTOTP,
		MFATypeSMSOTP,
		MFATypeEmailOTP,
		MFATypePush,
		MFATypeSecurityKey,
		MFATypePasskey,
		MFATypeHOTP,
		MFATypeBackupCode,
		MFATypeVoiceCallOTP:
		return nil
	}

	return errors.New("invalid mfa type")
}
