package request

import (
	"strings"

	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/supports"
)

type Relationship struct {
	Target string `json:"target" validate:"required"`
}

func (r *Relationship) Validate() error {
	r.Target = strings.TrimSpace(r.Target)
	return supports.Validate.Struct(r)
}

/// <----------------->

type ListRelationships struct {
	Status model.RelationshipStatus `validate:"omitempty,oneof=active pending"`
}

func (r *ListRelationships) Validate(rawStatus string) error {
	r.Status = model.RelationshipStatus(strings.TrimSpace(rawStatus))
	return supports.Validate.Struct(r)
}
