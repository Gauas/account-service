package repository

import (
	"github.com/gauas/account-service/model"
	"gorm.io/gorm"
)

type Registry struct {
	db           *gorm.DB
	User         Repository[model.User]
	Verification Repository[model.Verification]
	MFA          Repository[model.MFA]
	Identity     Repository[model.Identity]
	Relationship Repository[model.Relationship]
}

func New(db *gorm.DB) *Registry {
	return &Registry{
		db:           db,
		User:         Repository[model.User]{db: db},
		Verification: Repository[model.Verification]{db: db},
		MFA:          Repository[model.MFA]{db: db},
		Identity:     Repository[model.Identity]{db: db},
		Relationship: Repository[model.Relationship]{db: db},
	}
}

func (r *Registry) DB() *gorm.DB {
	return r.db
}
