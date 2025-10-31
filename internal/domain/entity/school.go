package entity

import (
	"time"

	"github.com/edugo/api-administracion/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
)

// School representa una escuela/institución educativa
type School struct {
	id        valueobject.SchoolID
	name      string
	address   string
	isActive  bool
	createdAt time.Time
	updatedAt time.Time
}

// NewSchool crea una nueva escuela con validaciones de negocio
func NewSchool(name, address string) (*School, error) {
	// Validaciones de negocio
	if name == "" {
		return nil, errors.NewValidationError("name is required")
	}

	if len(name) < 3 {
		return nil, errors.NewValidationError("name must be at least 3 characters")
	}

	if address == "" {
		return nil, errors.NewValidationError("address is required")
	}

	now := time.Now()

	return &School{
		id:        valueobject.NewSchoolID(),
		name:      name,
		address:   address,
		isActive:  true,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructSchool reconstruye una School desde la base de datos
func ReconstructSchool(
	id valueobject.SchoolID,
	name string,
	address string,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *School {
	return &School{
		id:        id,
		name:      name,
		address:   address,
		isActive:  isActive,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// Getters

func (s *School) ID() valueobject.SchoolID {
	return s.id
}

func (s *School) Name() string {
	return s.name
}

func (s *School) Address() string {
	return s.address
}

func (s *School) IsActive() bool {
	return s.isActive
}

func (s *School) CreatedAt() time.Time {
	return s.createdAt
}

func (s *School) UpdatedAt() time.Time {
	return s.updatedAt
}

// Business Logic Methods

// UpdateInfo actualiza la información de la escuela
func (s *School) UpdateInfo(name, address string) error {
	if name == "" && address == "" {
		return errors.NewValidationError("at least one field must be provided")
	}

	if name != "" {
		if len(name) < 3 {
			return errors.NewValidationError("name must be at least 3 characters")
		}
		s.name = name
	}

	if address != "" {
		s.address = address
	}

	s.updatedAt = time.Now()
	return nil
}

// Deactivate desactiva la escuela
func (s *School) Deactivate() error {
	if !s.isActive {
		return errors.NewBusinessRuleError("school is already inactive")
	}

	s.isActive = false
	s.updatedAt = time.Now()
	return nil
}

// Activate activa la escuela
func (s *School) Activate() error {
	if s.isActive {
		return errors.NewBusinessRuleError("school is already active")
	}

	s.isActive = true
	s.updatedAt = time.Now()
	return nil
}
