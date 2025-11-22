package dto

import (
	"time"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

type CreateMembershipRequest struct {
	UnitID     string     `json:"unit_id" validate:"required,uuid"`
	UserID     string     `json:"user_id" validate:"required,uuid"`
	Role       string     `json:"role" validate:"required"`
	ValidFrom  *time.Time `json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until"`
}

type UpdateMembershipRequest struct {
	Role       *string    `json:"role"`
	ValidUntil *time.Time `json:"valid_until"`
}

type MembershipResponse struct {
	ID         string     `json:"id"`
	UnitID     string     `json:"unit_id"`
	UserID     string     `json:"user_id"`
	Role       string     `json:"role"`
	EnrolledAt time.Time  `json:"enrolled_at"`
	WithdrawnAt *time.Time `json:"withdrawn_at,omitempty"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func ToMembershipResponse(m *entities.Membership) MembershipResponse {
	var unitID string
	if m.AcademicUnitID != nil {
		unitID = m.AcademicUnitID.String()
	}
	return MembershipResponse{
		ID:          m.ID.String(),
		UnitID:      unitID,
		UserID:      m.UserID.String(),
		Role:        m.Role,
		EnrolledAt:  m.EnrolledAt,
		WithdrawnAt: m.WithdrawnAt,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
