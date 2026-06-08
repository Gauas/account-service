package types

import (
	"errors"

	"github.com/gauas/account-service/supports"
)

type Phone string

func (p Phone) Validate() error {
	if !supports.IsPhone(string(p)) && p != "" {
		return errors.New("invalid phone")
	}

	return nil
}
