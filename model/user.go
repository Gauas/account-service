package model

import (
	"strings"
	"time"

	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const DefaultAvatarURL = "https://cdn.gauas.com/images/avatar/default_image.jpg"

type User struct {
	ID  int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"-"`
	Key uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"key"`

	Permission string `gorm:"size:50;index" json:"permission"`

	FullName  *string `gorm:"size:255" json:"full_name"`
	AvatarURL *string `gorm:"size:500" json:"avatar_url" default:"'https://cdn.gauas.com/images/avatar/default_image.jpg'"`

	Dob    *time.Time    `json:"dob"`
	Gender *types.Gender `gorm:"size:50" json:"gender"`

	IsOnboarded bool `json:"is_onboarded" gorm:"-"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`

	Identities    []Identity     `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Verifications []Verification `gorm:"foreignKey:UserID;references:ID" json:"-"`
	MFAs          []MFA          `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (u *User) AfterFind(_ *gorm.DB) error {
	u.IsOnboarded = u.HasCompletedProfile()
	return nil
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.AvatarURL == nil || strings.TrimSpace(*u.AvatarURL) == "" {
		defaultAvatar := DefaultAvatarURL
		u.AvatarURL = &defaultAvatar
	}
	return nil
}

func (u *User) HasCompletedProfile() bool {
	if u == nil {
		return false
	}

	return u.FullName != nil && strings.TrimSpace(*u.FullName) != "" &&
		u.AvatarURL != nil && strings.TrimSpace(*u.AvatarURL) != "" &&
		u.Dob != nil && !u.Dob.IsZero() &&
		u.Gender != nil && strings.TrimSpace(string(*u.Gender)) != ""
}
