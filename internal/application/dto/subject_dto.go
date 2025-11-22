package dto

import (
	"time"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

type CreateSubjectRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	Metadata    string `json:"metadata"`
}

type UpdateSubjectRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2"`
	Description *string `json:"description"`
	Metadata    *string `json:"metadata"`
}

type SubjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Metadata    string    `json:"metadata,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToSubjectResponse(subject *entities.Subject) SubjectResponse {
	desc := ""
	if subject.Description != nil {
		desc = *subject.Description
	}
	meta := ""
	if subject.Metadata != nil {
		meta = *subject.Metadata
	}
	return SubjectResponse{
		ID:          subject.ID.String(),
		Name:        subject.Name,
		Description: desc,
		Metadata:    meta,
		IsActive:    subject.IsActive,
		CreatedAt:   subject.CreatedAt,
		UpdatedAt:   subject.UpdatedAt,
	}
}
