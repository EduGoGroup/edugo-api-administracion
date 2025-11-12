package entity

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// School representa una escuela/institución educativa en la jerarquía académica
type School struct {
	id           valueobject.SchoolID
	name         string
	code         string
	address      string
	contactEmail *valueobject.Email
	contactPhone string
	metadata     map[string]interface{}
	createdAt    time.Time
	updatedAt    time.Time
}

// NewSchool crea una nueva escuela con validaciones de negocio
func NewSchool(name, code, address string) (*School, error) {
	// Validaciones de negocio
	if name == "" {
		return nil, errors.NewValidationError("name is required")
	}

	if len(name) < 3 {
		return nil, errors.NewValidationError("name must be at least 3 characters")
	}

	if code == "" {
		return nil, errors.NewValidationError("code is required")
	}

	if len(code) < 3 {
		return nil, errors.NewValidationError("code must be at least 3 characters")
	}

	now := time.Now()

	return &School{
		id:        valueobject.NewSchoolID(),
		name:      name,
		code:      code,
		address:   address,
		metadata:  make(map[string]interface{}),
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructSchool reconstruye una School desde la base de datos
func ReconstructSchool(
	id valueobject.SchoolID,
	name string,
	code string,
	address string,
	contactEmail *valueobject.Email,
	contactPhone string,
	metadata map[string]interface{},
	createdAt time.Time,
	updatedAt time.Time,
) *School {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return &School{
		id:           id,
		name:         name,
		code:         code,
		address:      address,
		contactEmail: contactEmail,
		contactPhone: contactPhone,
		metadata:     metadata,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters

func (s *School) ID() valueobject.SchoolID {
	return s.id
}

func (s *School) Name() string {
	return s.name
}

func (s *School) Code() string {
	return s.code
}

func (s *School) Address() string {
	return s.address
}

func (s *School) ContactEmail() *valueobject.Email {
	return s.contactEmail
}

func (s *School) ContactPhone() string {
	return s.contactPhone
}

func (s *School) Metadata() map[string]interface{} {
	// Retornar copia para evitar mutaciones externas
	copy := make(map[string]interface{})
	for k, v := range s.metadata {
		copy[k] = v
	}
	return copy
}

func (s *School) CreatedAt() time.Time {
	return s.createdAt
}

func (s *School) UpdatedAt() time.Time {
	return s.updatedAt
}

// Business Logic Methods

// UpdateInfo actualiza la información básica de la escuela
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

// UpdateContactInfo actualiza la información de contacto
func (s *School) UpdateContactInfo(email *valueobject.Email, phone string) error {
	if email == nil && phone == "" {
		return errors.NewValidationError("at least one contact field must be provided")
	}

	if email != nil {
		s.contactEmail = email
	}

	if phone != "" {
		s.contactPhone = phone
	}

	s.updatedAt = time.Now()
	return nil
}

// SetMetadata establece un valor en el metadata
func (s *School) SetMetadata(key string, value interface{}) {
	if s.metadata == nil {
		s.metadata = make(map[string]interface{})
	}
	s.metadata[key] = value
	s.updatedAt = time.Now()
}

// GetMetadata obtiene un valor del metadata
func (s *School) GetMetadata(key string) (interface{}, bool) {
	if s.metadata == nil {
		return nil, false
	}
	val, exists := s.metadata[key]
	return val, exists
}

// Validate valida el estado completo de la entidad
func (s *School) Validate() error {
	if s.name == "" {
		return errors.NewValidationError("name is required")
	}

	if len(s.name) < 3 {
		return errors.NewValidationError("name must be at least 3 characters")
	}

	if s.code == "" {
		return errors.NewValidationError("code is required")
	}

	if len(s.code) < 3 {
		return errors.NewValidationError("code must be at least 3 characters")
	}

	return nil
}
