package request

import (
	"time"

	"github.com/gauas/account-service/model/types"
)

type UpdateProfileRequest struct {
	FullName *string       `json:"full_name"`
	Dob      *time.Time    `json:"dob"`
	Gender   *types.Gender `json:"gender"`
}
