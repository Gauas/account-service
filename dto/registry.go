package dto

import (
	"time"

	"github.com/gauas/account-service/model/types"
)

type RegisterRequest struct {
	Email       types.Email  `json:"email"`
	Password    string       `json:"password"`
	FullName    string       `json:"fullname"`
	Gender      types.Gender `json:"gender"`
	DateOfBirth time.Time    `json:"date_of_birth"`
}
