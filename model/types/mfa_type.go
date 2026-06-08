package types

import "errors"

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
