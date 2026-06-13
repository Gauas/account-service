package types

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type Code string

func (c Code) Normalize() Code {
	return Code(strings.TrimSpace(string(c)))
}

func (c Code) Validate() error {
	if c == "" {
		return errors.New("verification code is required")
	}
	if len(c) != 6 {
		return errors.New("invalid verification code")
	}
	for _, ch := range c {
		if ch < '0' || ch > '9' {
			return errors.New("invalid verification code")
		}
	}

	return nil
}

func NewCode() (Code, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	return Code(fmt.Sprintf("%06d", n.Int64())), nil
}
