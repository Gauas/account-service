package repository

import (
	"github.com/gauas/account-service/model"
	"gorm.io/gorm"
)

type Registry struct {
	User         Repository[model.User]
	Verification Repository[model.UserVerification]
	MFA          Repository[model.UserMFA]
}

func New(db *gorm.DB) *Registry {
	return &Registry{
		User:         Repository[model.User]{db: db},
		Verification: Repository[model.UserVerification]{db: db},
		MFA:          Repository[model.UserMFA]{db: db},
	}
}
