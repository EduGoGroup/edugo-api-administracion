package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/validator"
)

// CreateGuardianRelationRequest representa la solicitud para crear una relación
type CreateGuardianRelationRequest struct {
	GuardianID       string `json:"guardian_id"`
	StudentID        string `json:"student_id"`
	RelationshipType string `json:"relationship_type"`
}

// Validate valida el request usando shared/validator
func (r *CreateGuardianRelationRequest) Validate() error {
	v := validator.New()

	v.Required(r.GuardianID, "guardian_id")
	v.UUID(r.GuardianID, "guardian_id")

	v.Required(r.StudentID, "student_id")
	v.UUID(r.StudentID, "student_id")

	v.Required(r.RelationshipType, "relationship_type")
	v.InSlice(r.RelationshipType, valueobject.AllRelationshipTypes(), "relationship_type")

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

// ToGuardianRelationResponse convierte una entidad a DTO de respuesta
func ToGuardianRelationResponse(relation *entity.GuardianRelation) *GuardianRelationResponse {
	return &GuardianRelationResponse{
		ID:               relation.ID().String(),
		GuardianID:       relation.GuardianID().String(),
		StudentID:        relation.StudentID().String(),
		RelationshipType: relation.RelationshipType().String(),
		IsActive:         relation.IsActive(),
		CreatedAt:        relation.CreatedAt(),
		UpdatedAt:        relation.UpdatedAt(),
		CreatedBy:        relation.CreatedBy(),
	}
}
