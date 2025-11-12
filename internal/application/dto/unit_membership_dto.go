package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
)

// CreateMembershipRequest representa la solicitud para crear una membresía
type CreateMembershipRequest struct {
	UnitID     string                 `json:"unit_id" validate:"required,uuid"`
	UserID     string                 `json:"user_id" validate:"required,uuid"`
	Role       string                 `json:"role" validate:"required,oneof=student teacher coordinator admin assistant"`
	ValidFrom  *time.Time             `json:"valid_from"`
	ValidUntil *time.Time             `json:"valid_until"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// UpdateMembershipRequest representa la solicitud para actualizar una membresía
type UpdateMembershipRequest struct {
	Role       *string                `json:"role" validate:"omitempty,oneof=student teacher coordinator admin assistant"`
	ValidUntil *time.Time             `json:"valid_until"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// MembershipResponse representa la respuesta con datos de una membresía
type MembershipResponse struct {
	ID         string                 `json:"id"`
	UnitID     string                 `json:"unit_id"`
	UserID     string                 `json:"user_id"`
	Role       string                 `json:"role"`
	ValidFrom  time.Time              `json:"valid_from"`
	ValidUntil *time.Time             `json:"valid_until,omitempty"`
	IsActive   bool                   `json:"is_active"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// MembershipWithDetailsResponse incluye información denormalizada de unidad y usuario
type MembershipWithDetailsResponse struct {
	MembershipResponse
	UnitName   string `json:"unit_name"`
	UnitType   string `json:"unit_type"`
	SchoolName string `json:"school_name,omitempty"`
}

// ToMembershipResponse convierte una entidad UnitMembership a response
func ToMembershipResponse(membership *entity.UnitMembership) MembershipResponse {
	return MembershipResponse{
		ID:         membership.ID().String(),
		UnitID:     membership.UnitID().String(),
		UserID:     membership.UserID().String(),
		Role:       membership.Role().String(),
		ValidFrom:  membership.ValidFrom(),
		ValidUntil: membership.ValidUntil(),
		IsActive:   membership.IsActive(),
		Metadata:   membership.Metadata(),
		CreatedAt:  membership.CreatedAt(),
		UpdatedAt:  membership.UpdatedAt(),
	}
}

// ToMembershipResponseList convierte una lista de entidades a lista de responses
func ToMembershipResponseList(memberships []*entity.UnitMembership) []MembershipResponse {
	responses := make([]MembershipResponse, len(memberships))
	for i, membership := range memberships {
		responses[i] = ToMembershipResponse(membership)
	}
	return responses
}
