package request

import (
	"errors"
	"strings"
	"time"

	"github.com/gauas/account-service/model/types"
	"github.com/gauas/account-service/supports"
	"github.com/gauas/account-service/supports/oauth2"
)

type Login struct {
	Email    types.Email `json:"email" validate:"required"`
	Password string      `json:"password" validate:"required"`
}

func (r *Login) Validate() error {
	r.Email = r.Email.Normalize()
	if err := supports.Validate.Struct(r); err != nil {
		return err
	}

	return r.Email.Validate()
}

/// <----------------->

type OAuth2 struct {
	Provider string `json:"provider" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

func (r *OAuth2) Validate() error {
	r.Provider = strings.TrimSpace(r.Provider)
	r.Token = strings.TrimSpace(r.Token)
	if err := supports.Validate.Struct(r); err != nil {
		return err
	}

	if _, ok := oauth2.Providers[r.Provider]; !ok {
		return errors.New("unsupported oauth2 provider")
	}

	return nil
}

/// <----------------->

type Register struct {
	Email       types.Email  `json:"email" validate:"required"`
	Password    string       `json:"password" validate:"required"`
	FullName    string       `json:"full_name" validate:"omitempty,min=2,max=100"`
	Gender      types.Gender `json:"gender" validate:"omitempty"`
	DateOfBirth time.Time    `json:"date_of_birth" validate:"omitempty"`
}

func (r *Register) Validate() error {
	r.Email = r.Email.Normalize()
	if err := supports.Validate.Struct(r); err != nil {
		return err
	}

	return r.Email.Validate()
}
