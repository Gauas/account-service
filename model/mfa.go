package model

import (
	"time"

	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
)

type MFA struct {
	ID  int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"-"`
	Key uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"key,omitempty"`

	User   User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	UserID int64 `gorm:"type:bigint;index" json:"-"`

	Type types.MFAType `gorm:"size:30;index" json:"type,omitempty"`

	Secret  *string `gorm:"size:255" json:"secret,omitempty"`
	Enabled bool    `gorm:"default:false" json:"enabled,omitempty"`

	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (MFA) TableName() string {
	return "mfas"
}
