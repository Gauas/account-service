package service

import (
	"fmt"
	"reflect"
)

func Validate(values ...interface{}) error {
	for _, value := range values {
		if reflect.ValueOf(value).IsZero() {
			return fmt.Errorf("invalid value %v", value)
		}
	}

	return nil
}
