package dto

import "github.com/gauas/account-service/model/types"

type LoginRequest struct {
	Email    types.Email `json:"email"`
	Password string      `json:"password"`
}
