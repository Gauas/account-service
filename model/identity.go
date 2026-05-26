package model

import (
	"time"

	"github.com/gauas/account-service/model/types"
	"github.com/google/uuid"
)

type Identity struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	UserID uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id,omitempty"`

	Provider       types.IdentityProvider `gorm:"size:50;index;not null" json:"provider,omitempty"`
	ProviderUserID string                 `gorm:"size:255;index;not null" json:"provider_user_id,omitempty"`

	Email *types.Email `gorm:"size:255" json:"email,omitempty"`
	Phone *string      `gorm:"size:20" json:"phone,omitempty"`
	Hash  *string      `gorm:"size:255" json:"hash,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Identity) TableName() string {
	return "identities"
}
