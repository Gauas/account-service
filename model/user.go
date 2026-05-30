package model

import (
	"time"

	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID  int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"-"`
	Key uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"key"`

	Permission string `gorm:"size:50;index" json:"permission"`

	FullName  *string `gorm:"size:255" json:"full_name"`
	AvatarURL *string `gorm:"size:500" json:"avatar_url" default:"'https://cdn.gauas.com/images/avatar/default.jpg'"`

	Dob    *time.Time    `json:"dob"`
	Gender *types.Gender `gorm:"size:50" json:"gender"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`

	Identities    []Identity     `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Verifications []Verification `gorm:"foreignKey:UserID;references:ID" json:"-"`
	MFAs          []MFA          `gorm:"foreignKey:UserID;references:ID" json:"-"`
}
