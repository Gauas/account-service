package model

import (
	"time"

	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
)

type Identity struct {
	ID  int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"-"`
	Key uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"key,omitempty"`

	User   User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	UserID int64 `gorm:"type:bigint;index;not null" json:"-"`

	Provider       types.IdentityProvider `gorm:"size:50;index;not null" json:"provider,omitempty"`
	ProviderUserID string                 `gorm:"size:255;index;not null" json:"-"`

	Email *types.Email `gorm:"size:255" json:"email,omitempty"`
	Phone *string      `gorm:"size:20" json:"phone,omitempty"`
	Hash  *string      `gorm:"size:255" json:"hash,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (Identity) TableName() string {
	return "identities"
}
