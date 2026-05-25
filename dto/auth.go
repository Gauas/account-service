package dto

import (
	"time"

	"github.com/gauas/account-service/model/types"
)

type LoginRequest struct {
	Email    types.Email `json:"email"`
	Password string      `json:"password"`
}

type Oauth2Request struct {
	Provider string `json:"provider"`
	Token    string `json:"token"`
}

type RegisterRequest struct {
	Email       types.Email  `json:"email"`
	Password    string       `json:"password"`
	FullName    string       `json:"fullname"`
	Gender      types.Gender `json:"gender"`
	DateOfBirth time.Time    `json:"date_of_birth"`
}
