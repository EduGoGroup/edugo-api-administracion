package dto

import (
	"time"

	"github.com/edugo/api-administracion/internal/domain/entity"
	"github.com/edugo/shared/pkg/validator"
)

// CreateSubjectRequest solicitud para crear materia
type CreateSubjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Metadata    string `json:"metadata"`
}

func (r *CreateSubjectRequest) Validate() error {
	v := validator.New()
	v.Required(r.Name, "name")
	v.MinLength(r.Name, 2, "name")
	v.MaxLength(r.Name, 100, "name")
	return v.GetError()
}

// UpdateSubjectRequest solicitud para actualizar materia
type UpdateSubjectRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Metadata    *string `json:"metadata,omitempty"`
}

func (r *UpdateSubjectRequest) Validate() error {
	v := validator.New()
	if r.Name != nil {
		v.MinLength(*r.Name, 2, "name")
		v.MaxLength(*r.Name, 100, "name")
	}
	return v.GetError()
}

// SubjectResponse respuesta de materia
type SubjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Metadata    string    `json:"metadata"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToSubjectResponse(subject *entity.Subject) *SubjectResponse {
	return &SubjectResponse{
		ID:          subject.ID().String(),
		Name:        subject.Name(),
		Description: subject.Description(),
		Metadata:    subject.Metadata(),
		IsActive:    subject.IsActive(),
		CreatedAt:   subject.CreatedAt(),
		UpdatedAt:   subject.UpdatedAt(),
	}
}
