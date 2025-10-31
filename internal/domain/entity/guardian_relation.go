package entity

import (
	"time"

	"github.com/edugo/api-administracion/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/types"
)

// GuardianRelation representa la relación entre un guardian y un estudiante
type GuardianRelation struct {
	id               types.UUID
	guardianID       valueobject.GuardianID
	studentID        valueobject.StudentID
	relationshipType valueobject.RelationshipType
	isActive         bool
	createdAt        time.Time
	updatedAt        time.Time
	createdBy        string
}

// NewGuardianRelation crea una nueva relación guardian-estudiante
func NewGuardianRelation(
	guardianID valueobject.GuardianID,
	studentID valueobject.StudentID,
	relationshipType valueobject.RelationshipType,
	createdBy string,
) (*GuardianRelation, error) {
	// Validaciones de negocio
	if guardianID.IsZero() {
		return nil, errors.NewValidationError("guardian_id is required")
	}

	if studentID.IsZero() {
		return nil, errors.NewValidationError("student_id is required")
	}

	if !relationshipType.IsValid() {
		return nil, errors.NewValidationError("invalid relationship_type")
	}

	if guardianID.Equals(valueobject.GuardianID{}) {
		return nil, errors.NewBusinessRuleError("guardian cannot be the student")
	}

	now := time.Now()

	return &GuardianRelation{
		id:               types.NewUUID(),
		guardianID:       guardianID,
		studentID:        studentID,
		relationshipType: relationshipType,
		isActive:         true,
		createdAt:        now,
		updatedAt:        now,
		createdBy:        createdBy,
	}, nil
}

// Reconstruct reconstruye una GuardianRelation desde la base de datos
func ReconstructGuardianRelation(
	id types.UUID,
	guardianID valueobject.GuardianID,
	studentID valueobject.StudentID,
	relationshipType valueobject.RelationshipType,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
	createdBy string,
) *GuardianRelation {
	return &GuardianRelation{
		id:               id,
		guardianID:       guardianID,
		studentID:        studentID,
		relationshipType: relationshipType,
		isActive:         isActive,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
		createdBy:        createdBy,
	}
}

// Getters

func (g *GuardianRelation) ID() types.UUID {
	return g.id
}

func (g *GuardianRelation) GuardianID() valueobject.GuardianID {
	return g.guardianID
}

func (g *GuardianRelation) StudentID() valueobject.StudentID {
	return g.studentID
}

func (g *GuardianRelation) RelationshipType() valueobject.RelationshipType {
	return g.relationshipType
}

func (g *GuardianRelation) IsActive() bool {
	return g.isActive
}

func (g *GuardianRelation) CreatedAt() time.Time {
	return g.createdAt
}

func (g *GuardianRelation) UpdatedAt() time.Time {
	return g.updatedAt
}

func (g *GuardianRelation) CreatedBy() string {
	return g.createdBy
}

// Business Logic Methods

// Deactivate desactiva la relación
func (g *GuardianRelation) Deactivate() error {
	if !g.isActive {
		return errors.NewBusinessRuleError("relation is already inactive")
	}

	g.isActive = false
	g.updatedAt = time.Now()
	return nil
}

// Activate activa la relación
func (g *GuardianRelation) Activate() error {
	if g.isActive {
		return errors.NewBusinessRuleError("relation is already active")
	}

	g.isActive = true
	g.updatedAt = time.Now()
	return nil
}

// ChangeRelationshipType cambia el tipo de relación
func (g *GuardianRelation) ChangeRelationshipType(newType valueobject.RelationshipType) error {
	if !newType.IsValid() {
		return errors.NewValidationError("invalid relationship type")
	}

	if g.relationshipType == newType {
		return errors.NewBusinessRuleError("new relationship type is the same as current")
	}

	g.relationshipType = newType
	g.updatedAt = time.Now()
	return nil
}
