package model

import (
	"time"

	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id,omitempty"`
	Permission string    `gorm:"size:50;index" json:"permission,omitempty"`

	FullName  *string `gorm:"size:255" json:"full_name,omitempty"`
	AvatarURL *string `gorm:"size:500" json:"avatar_url,omitempty"`

	Dob    *time.Time    `json:"dob,omitempty"`
	Gender *types.Gender `gorm:"size:50" json:"gender,omitempty"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`

	Identities    []Identity     `gorm:"foreignKey:UserID" json:"-"`
	Verifications []Verification `gorm:"foreignKey:UserID" json:"-"`
	MFAs          []MFA          `gorm:"foreignKey:UserID" json:"-"`
}
