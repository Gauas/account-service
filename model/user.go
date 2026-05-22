package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID      uuid.UUID      `gorm:"type:uuid;primaryKey"                      json:"user_id,omitempty"`
	Permission  string         `gorm:"index:idx_username_permission"              json:"permission,omitempty"`
	Password    *string        `gorm:"size:255"                                   json:"password,omitempty"`
	Email       *Email         `gorm:"unique"                                     json:"email,omitempty"`
	Phone       *Phone         `gorm:"size:15"                                    json:"phone,omitempty"`
	FullName    *string        `                                                  json:"fullname,omitempty"`
	Gender      *Gender        `                                                  json:"gender,omitempty"`
	DateOfBirth *time.Time     `                                                  json:"date_of_birth,omitempty"`
	FacebookURL *string        `gorm:"unique"                                     json:"facebook_url,omitempty"`
	GithubURL   *string        `gorm:"unique"                                     json:"github_url,omitempty"`
	AvatarURL   *string        `                                                  json:"avatar_url,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                                      json:"deleted_at,omitempty"`

	Verifications []UserVerification `gorm:"foreignKey:UserID" json:"-"`
	MFAs          []UserMFA          `gorm:"foreignKey:UserID" json:"-"`
}
