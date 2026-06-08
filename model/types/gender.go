package types

import (
	"errors"
	"slices"
)

type Gender string

func (g Gender) Validate() error {
	if !slices.Contains([]string{"male", "female", "other"}, string(g)) {
		return errors.New("invalid gender")
	}

	return nil
}
