package request

import (
	"strings"

	"github.com/gauas/account-service/supports"
)

type EnableTOTP struct {
	OTPCode string `json:"otp_code" validate:"required"`
}

func (r *EnableTOTP) Validate() error {
	r.OTPCode = strings.TrimSpace(r.OTPCode)
	return supports.Validate.Struct(r)
}

/// <----------------->

type VerifyTOTP struct {
	OTPCode string `json:"otp_code" validate:"required"`
}

func (r *VerifyTOTP) Validate() error {
	r.OTPCode = strings.TrimSpace(r.OTPCode)
	return supports.Validate.Struct(r)
}
