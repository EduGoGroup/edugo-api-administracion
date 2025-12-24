package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/validator"
)

// CreateGuardianRelationRequest representa la solicitud para crear una relación
type CreateGuardianRelationRequest struct {
	GuardianID       string `json:"guardian_id"`
	StudentID        string `json:"student_id"`
	RelationshipType string `json:"relationship_type"`
}

// UpdateGuardianRelationRequest representa la solicitud para actualizar una relación
type UpdateGuardianRelationRequest struct {
	RelationshipType *string `json:"relationship_type,omitempty"`
	IsActive         *bool   `json:"is_active,omitempty"`
}

// Validate valida el request usando shared/validator
func (r *CreateGuardianRelationRequest) Validate() error {
	v := validator.New()

	v.Required(r.GuardianID, "guardian_id")
	v.UUID(r.GuardianID, "guardian_id")

	v.Required(r.StudentID, "student_id")
	v.UUID(r.StudentID, "student_id")

	v.Required(r.RelationshipType, "relationship_type")
	validTypes := []string{"father", "mother", "grandfather", "grandmother", "uncle", "aunt", "other"}
	v.InSlice(r.RelationshipType, validTypes, "relationship_type")

	return v.GetError()
}

// GuardianRelationResponse representa la respuesta de una relación
type GuardianRelationResponse struct {
	ID               string    `json:"id"`
	GuardianID       string    `json:"guardian_id"`
	StudentID        string    `json:"student_id"`
	RelationshipType string    `json:"relationship_type"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedBy        string    `json:"created_by"`
}

// ToGuardianRelationResponse convierte una entidad de infrastructure a DTO de respuesta
func ToGuardianRelationResponse(relation *entities.GuardianRelation) *GuardianRelationResponse {
	return &GuardianRelationResponse{
		ID:               relation.ID.String(),
		GuardianID:       relation.GuardianID.String(),
		StudentID:        relation.StudentID.String(),
		RelationshipType: relation.RelationshipType,
		IsActive:         relation.IsActive,
		CreatedAt:        relation.CreatedAt,
		UpdatedAt:        relation.UpdatedAt,
		CreatedBy:        relation.CreatedBy,
	}
}
