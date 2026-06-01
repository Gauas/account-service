package model

import (
	"time"

	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
)

type Verification struct {
	ID  int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"-"`
	Key uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"key,omitempty"`

	User   User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	UserID int64 `gorm:"type:bigint;index" json:"-"`

	Method types.VerificationMethod `gorm:"size:20;index" json:"method,omitempty"`
	Value  string                   `gorm:"size:255" json:"value,omitempty"`

	IsVerified bool       `gorm:"default:false" json:"is_verified,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (Verification) TableName() string {
	return "verifications"
}
