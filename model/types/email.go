package types

import (
	"errors"
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
