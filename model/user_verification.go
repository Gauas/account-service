package model

import (
	"time"

	"github.com/google/uuid"
)

type UserVerification struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"   json:"id,omitempty"`
	UserID     uuid.UUID  `gorm:"type:uuid;index"        json:"user_id,omitempty"`
	Method     string     `gorm:"size:20;index"          json:"method,omitempty"`
	Value      string     `gorm:"size:255"               json:"value,omitempty"`
	IsVerified bool       `gorm:"default:false"          json:"is_verified,omitempty"`
	VerifiedAt *time.Time `                              json:"verified_at,omitempty"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (UserVerification) TableName() string { return "verifications" }
