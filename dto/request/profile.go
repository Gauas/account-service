package request

import (
	"time"

	"github.com/gauas/account-service/model/types"
	"github.com/gauas/account-service/supports"
)

type UpdateProfile struct {
	FullName *string       `json:"full_name" validate:"omitempty,min=2,max=100"`
	Dob      *time.Time    `json:"dob" validate:"omitempty"`
	Gender   *types.Gender `json:"gender" validate:"omitempty"`
}

func (r *UpdateProfile) Validate() error {
	if err := supports.Validate.Struct(r); err != nil {
		return err
	}

	if r.Gender != nil {
		return r.Gender.Validate()
	}

	return nil
}
