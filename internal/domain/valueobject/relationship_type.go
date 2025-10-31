package valueobject

import (
	"fmt"

	"github.com/edugo/shared/pkg/errors"
)

// RelationshipType representa el tipo de relación entre guardian y estudiante
type RelationshipType string

const (
	RelationshipTypeParent    RelationshipType = "parent"
	RelationshipTypeGuardian  RelationshipType = "guardian"
	RelationshipTypeRelative  RelationshipType = "relative"
	RelationshipTypeOther     RelationshipType = "other"
)

// NewRelationshipType crea y valida un nuevo RelationshipType
func NewRelationshipType(value string) (RelationshipType, error) {
	rt := RelationshipType(value)
	if !rt.IsValid() {
		return "", errors.NewValidationError(
			fmt.Sprintf("invalid relationship type: %s", value),
		).WithField("allowed_values", AllRelationshipTypes())
	}
	return rt, nil
}

// IsValid verifica si el tipo de relación es válido
func (r RelationshipType) IsValid() bool {
	switch r {
	case RelationshipTypeParent, RelationshipTypeGuardian, RelationshipTypeRelative, RelationshipTypeOther:
		return true
	}
	return false
}

// String retorna la representación en string
func (r RelationshipType) String() string {
	return string(r)
}

// AllRelationshipTypes retorna todos los tipos válidos
func AllRelationshipTypes() []string {
	return []string{
		string(RelationshipTypeParent),
		string(RelationshipTypeGuardian),
		string(RelationshipTypeRelative),
		string(RelationshipTypeOther),
	}
}
