package entity

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// UnitMembership representa la membresía de un usuario en una unidad académica
type UnitMembership struct {
	id         valueobject.MembershipID
	unitID     valueobject.UnitID
	userID     valueobject.UserID
	role       valueobject.MembershipRole
	validFrom  time.Time
	validUntil *time.Time
	metadata   map[string]interface{}
	createdAt  time.Time
	updatedAt  time.Time
}

// NewUnitMembership crea una nueva membresía
func NewUnitMembership(
	unitID valueobject.UnitID,
	userID valueobject.UserID,
	role valueobject.MembershipRole,
	validFrom time.Time,
) (*UnitMembership, error) {
	// Validaciones de negocio
	if unitID.IsZero() {
		return nil, errors.NewValidationError("unit_id is required")
	}

	if userID.IsZero() {
		return nil, errors.NewValidationError("user_id is required")
	}

	if !role.IsValid() {
		return nil, errors.NewValidationError("invalid membership role")
	}

	if validFrom.IsZero() {
		validFrom = time.Now()
	}

	now := time.Now()

	return &UnitMembership{
		id:        valueobject.NewMembershipID(),
		unitID:    unitID,
		userID:    userID,
		role:      role,
		validFrom: validFrom,
		metadata:  make(map[string]interface{}),
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructUnitMembership reconstruye una UnitMembership desde la base de datos
func ReconstructUnitMembership(
	id valueobject.MembershipID,
	unitID valueobject.UnitID,
	userID valueobject.UserID,
	role valueobject.MembershipRole,
	validFrom time.Time,
	validUntil *time.Time,
	metadata map[string]interface{},
	createdAt time.Time,
	updatedAt time.Time,
) *UnitMembership {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return &UnitMembership{
		id:         id,
		unitID:     unitID,
		userID:     userID,
		role:       role,
		validFrom:  validFrom,
		validUntil: validUntil,
		metadata:   metadata,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

// Getters

func (um *UnitMembership) ID() valueobject.MembershipID {
	return um.id
}

func (um *UnitMembership) UnitID() valueobject.UnitID {
	return um.unitID
}

func (um *UnitMembership) UserID() valueobject.UserID {
	return um.userID
}

func (um *UnitMembership) Role() valueobject.MembershipRole {
	return um.role
}

func (um *UnitMembership) ValidFrom() time.Time {
	return um.validFrom
}

func (um *UnitMembership) ValidUntil() *time.Time {
	return um.validUntil
}

func (um *UnitMembership) Metadata() map[string]interface{} {
	copy := make(map[string]interface{})
	for k, v := range um.metadata {
		copy[k] = v
	}
	return copy
}

func (um *UnitMembership) CreatedAt() time.Time {
	return um.createdAt
}

func (um *UnitMembership) UpdatedAt() time.Time {
	return um.updatedAt
}

// Setters - Para uso exclusivo de Domain Services
// ⚠️ NO usar directamente - pueden romper invariantes

func (um *UnitMembership) SetRole(role valueobject.MembershipRole) {
	um.role = role
}

func (um *UnitMembership) SetValidUntilValue(validUntil *time.Time) {
	um.validUntil = validUntil
}

func (um *UnitMembership) SetUpdatedAt(t time.Time) {
	um.updatedAt = t
}

func (um *UnitMembership) SetMetadataValue(key string, value interface{}) {
	if um.metadata == nil {
		um.metadata = make(map[string]interface{})
	}
	um.metadata[key] = value
}

// Business Logic Methods

// IsActive verifica si la membresía está activa en el momento actual
func (um *UnitMembership) IsActive() bool {
	return um.IsActiveAt(time.Now())
}

// IsActiveAt verifica si la membresía está activa en un momento específico
func (um *UnitMembership) IsActiveAt(t time.Time) bool {
	// Debe estar después de validFrom
	if t.Before(um.validFrom) {
		return false
	}

	// Si no tiene validUntil, es indefinida
	if um.validUntil == nil {
		return true
	}

	// Debe estar antes de validUntil
	return t.Before(*um.validUntil) || t.Equal(*um.validUntil)
}

// SetValidUntil establece la fecha de fin de la membresía
func (um *UnitMembership) SetValidUntil(validUntil time.Time) error {
	// ValidUntil debe ser después de validFrom
	if !validUntil.After(um.validFrom) {
		return errors.NewValidationError("valid_until must be after valid_from")
	}

	um.validUntil = &validUntil
	um.updatedAt = time.Now()
	return nil
}

// ExtendIndefinitely remueve la fecha de fin (membresía indefinida)
func (um *UnitMembership) ExtendIndefinitely() {
	um.validUntil = nil
	um.updatedAt = time.Now()
}

// Expire marca la membresía como expirada (establece validUntil a ahora)
func (um *UnitMembership) Expire() error {
	now := time.Now()

	if um.validUntil != nil && um.validUntil.Before(now) {
		return errors.NewBusinessRuleError("membership is already expired")
	}

	um.validUntil = &now
	um.updatedAt = now
	return nil
}

// ChangeRole cambia el rol de la membresía
func (um *UnitMembership) ChangeRole(newRole valueobject.MembershipRole) error {
	if !newRole.IsValid() {
		return errors.NewValidationError("invalid membership role")
	}

	if um.role == newRole {
		return errors.NewBusinessRuleError("role is already " + newRole.String())
	}

	um.role = newRole
	um.updatedAt = time.Now()
	return nil
}

// HasPermission verifica si esta membresía tiene un permiso específico
func (um *UnitMembership) HasPermission(permission string) bool {
	// Solo membresías activas tienen permisos
	if !um.IsActive() {
		return false
	}

	return um.role.HasPermission(permission)
}

// SetMetadata establece un valor en el metadata
func (um *UnitMembership) SetMetadata(key string, value interface{}) {
	if um.metadata == nil {
		um.metadata = make(map[string]interface{})
	}
	um.metadata[key] = value
	um.updatedAt = time.Now()
}

// GetMetadata obtiene un valor del metadata
func (um *UnitMembership) GetMetadata(key string) (interface{}, bool) {
	if um.metadata == nil {
		return nil, false
	}
	val, exists := um.metadata[key]
	return val, exists
}

// Validate valida el estado completo de la entidad
func (um *UnitMembership) Validate() error {
	if um.unitID.IsZero() {
		return errors.NewValidationError("unit_id is required")
	}

	if um.userID.IsZero() {
		return errors.NewValidationError("user_id is required")
	}

	if !um.role.IsValid() {
		return errors.NewValidationError("invalid membership role")
	}

	if um.validUntil != nil && !um.validUntil.After(um.validFrom) {
		return errors.NewValidationError("valid_until must be after valid_from")
	}

	return nil
}
