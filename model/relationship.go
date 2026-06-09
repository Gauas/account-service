package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RelationshipStatus string

const (
	RelationshipStatusActive  RelationshipStatus = "active"
	RelationshipStatusPending RelationshipStatus = "pending"
)

type Relationship struct {
	ID  int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"-"`
	Key uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"key,omitempty"`

	Actor   User  `gorm:"foreignKey:ActorID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ActorID int64 `gorm:"type:bigint;not null;index" json:"-"`

	Partner   User  `gorm:"foreignKey:PartnerID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	PartnerID int64 `gorm:"type:bigint;not null;index" json:"-"`

	Status RelationshipStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (Relationship) TableName() string {
	return "relationships"
}

func (r *Relationship) BeforeCreate(_ *gorm.DB) error {
	if r.Key != uuid.Nil {
		return nil
	}

	key, err := uuid.NewV7()
	if err != nil {
		return err
	}

	r.Key = key
	return nil
}
