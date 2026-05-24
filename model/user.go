package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id,omitempty"`
	Permission string    `gorm:"size:50;index" json:"permission,omitempty"`

	FullName  *string `gorm:"size:255" json:"full_name,omitempty"`
	AvatarURL *string `gorm:"size:500" json:"avatar_url,omitempty"`

	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Gender      *string    `gorm:"size:50" json:"gender,omitempty"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`

	Identities    []Identity     `gorm:"foreignKey:UserID" json:"-"`
	Verifications []Verification `gorm:"foreignKey:UserID" json:"-"`
	MFAs          []MFA          `gorm:"foreignKey:UserID" json:"-"`
}
