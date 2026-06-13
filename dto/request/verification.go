package request

import (
	"errors"

	"github.com/gauas/account-service/model/types"
	"github.com/gauas/account-service/supports"
)

type GenerateVerification struct {
	Type   types.VerificationMethod `json:"type" validate:"required"`
	Target string                   `json:"target" validate:"required"`
}

func (r *GenerateVerification) Validate() error {
	r.Target = string(types.Email(r.Target).Normalize())
	if err := supports.Validate.Struct(r); err != nil {
		return err
	}
	if r.Type == types.PhoneVerification {
		return errors.New("phone verification is temporarily disabled")
	}
	if r.Type == types.EmailVerification {
		return types.Email(r.Target).Validate()
	}

	return r.Type.Validate()
}

type VerifyVerification struct {
	Type   types.VerificationMethod `json:"type" validate:"required"`
	Target string                   `json:"target" validate:"required"`
	Code   types.Code               `json:"code" validate:"required"`
}

func (r *VerifyVerification) Validate() error {
	r.Target = string(types.Email(r.Target).Normalize())
	r.Code = r.Code.Normalize()
	if err := supports.Validate.Struct(r); err != nil {
		return err
	}
	if r.Type == types.PhoneVerification {
		return errors.New("phone verification is temporarily disabled")
	}
	if err := r.Type.Validate(); err != nil {
		return err
	}
	if err := types.Email(r.Target).Validate(); err != nil {
		return err
	}

	return r.Code.Validate()
}
