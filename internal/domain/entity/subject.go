package entity

import (
	"time"

	"github.com/edugo/api-administracion/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
)

// Subject representa una materia/asignatura
type Subject struct {
	id          valueobject.SubjectID
	name        string
	description string
	metadata    string
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

// NewSubject crea una nueva materia
func NewSubject(name, description, metadata string) (*Subject, error) {
	if name == "" {
		return nil, errors.NewValidationError("name is required")
	}

	now := time.Now()

	return &Subject{
		id:          valueobject.NewSubjectID(),
		name:        name,
		description: description,
		metadata:    metadata,
		isActive:    true,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ReconstructSubject reconstruye una Subject desde la base de datos
func ReconstructSubject(
	id valueobject.SubjectID,
	name, description, metadata string,
	isActive bool,
	createdAt, updatedAt time.Time,
) *Subject {
	return &Subject{
		id:          id,
		name:        name,
		description: description,
		metadata:    metadata,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getters

func (s *Subject) ID() valueobject.SubjectID   { return s.id }
func (s *Subject) Name() string                { return s.name }
func (s *Subject) Description() string         { return s.description }
func (s *Subject) Metadata() string            { return s.metadata }
func (s *Subject) IsActive() bool              { return s.isActive }
func (s *Subject) CreatedAt() time.Time        { return s.createdAt }
func (s *Subject) UpdatedAt() time.Time        { return s.updatedAt }

// Business Logic

func (s *Subject) UpdateInfo(name, description, metadata *string) error {
	if name != nil && *name != "" {
		s.name = *name
	}
	if description != nil {
		s.description = *description
	}
	if metadata != nil {
		s.metadata = *metadata
	}

	s.updatedAt = time.Now()
	return nil
}
