package dto

import (
	"time"

	"github.com/gauas/account-service/model"
)

type RegisterRequest struct {
	Email       *model.Email `json:"email"`
	Password    string       `json:"password"`
	Phone       *model.Phone `json:"phone"`
	FullName    string       `json:"fullname"`
	Gender      model.Gender `json:"gender"`
	DateOfBirth time.Time    `json:"date_of_birth"`
}
