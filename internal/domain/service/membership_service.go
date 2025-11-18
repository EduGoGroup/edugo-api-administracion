package service

import (
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// MembershipDomainService maneja la lógica de negocio de membresías
type MembershipDomainService struct{}

// NewMembershipDomainService crea una nueva instancia del servicio
func NewMembershipDomainService() *MembershipDomainService {
	return &MembershipDomainService{}
}

// IsActive verifica si una membresía está activa en el momento actual
func (s *MembershipDomainService) IsActive(membership *entity.UnitMembership) bool {
	return s.IsActiveAt(membership, time.Now())
}

// IsActiveAt verifica si una membresía está activa en un momento específico
func (s *MembershipDomainService) IsActiveAt(
	membership *entity.UnitMembership,
	t time.Time,
) bool {
	// Debe estar después de validFrom
	if t.Before(membership.ValidFrom()) {
		return false
	}

	// Si tiene validUntil, debe estar antes o igual
	if membership.ValidUntil() != nil && t.After(*membership.ValidUntil()) {
		return false
	}

	return true
}

// SetValidUntil establece la fecha de fin de validez de una membresía
func (s *MembershipDomainService) SetValidUntil(
	membership *entity.UnitMembership,
	validUntil time.Time,
) error {
	// Validar que validUntil es después de validFrom
	if validUntil.Before(membership.ValidFrom()) {
		return errors.NewValidationError("valid_until must be after valid_from")
	}

	// Validar que no sean iguales
	if validUntil.Equal(membership.ValidFrom()) {
		return errors.NewValidationError("valid_until cannot be equal to valid_from")
	}

	membership.SetValidUntilValue(&validUntil)
	membership.SetUpdatedAt(time.Now())
	return nil
}

// ExtendIndefinitely remueve la fecha de fin de una membresía
func (s *MembershipDomainService) ExtendIndefinitely(membership *entity.UnitMembership) {
	membership.SetValidUntilValue(nil)
	membership.SetUpdatedAt(time.Now())
}

// Expire marca una membresía como expirada (establece validUntil a ahora)
func (s *MembershipDomainService) Expire(membership *entity.UnitMembership) error {
	// Validar que no está ya expirada
	if membership.ValidUntil() != nil && time.Now().After(*membership.ValidUntil()) {
		return errors.NewBusinessRuleError("membership is already expired")
	}

	now := time.Now()
	membership.SetValidUntilValue(&now)
	membership.SetUpdatedAt(now)
	return nil
}

// ChangeRole cambia el rol de una membresía
func (s *MembershipDomainService) ChangeRole(
	membership *entity.UnitMembership,
	newRole valueobject.MembershipRole,
) error {
	// Validar que el nuevo rol es válido
	if !newRole.IsValid() {
		return errors.NewValidationError("invalid role")
	}

	// Validar que el rol es diferente
	if membership.Role() == newRole {
		return errors.NewBusinessRuleError("new role must be different from current role")
	}

	membership.SetRole(newRole)
	membership.SetUpdatedAt(time.Now())
	return nil
}

// HasPermission verifica si una membresía tiene un permiso específico
func (s *MembershipDomainService) HasPermission(
	membership *entity.UnitMembership,
	permission string,
) bool {
	// Si no está activa, no tiene permisos
	if !s.IsActive(membership) {
		return false
	}

	// Verificar permisos según rol
	switch membership.Role() {
	case valueobject.RoleAdmin:
		return true // Admin tiene todos los permisos
	case valueobject.RoleCoordinator:
		return permission == "view" || permission == "edit" || permission == "manage_members"
	case valueobject.RoleTeacher:
		return permission == "view" || permission == "edit"
	default:
		return false
	}
}

// Validate valida el estado completo de una membresía
func (s *MembershipDomainService) Validate(membership *entity.UnitMembership) error {
	if membership.UnitID().IsZero() {
		return errors.NewValidationError("unit_id is required")
	}

	if membership.UserID().IsZero() {
		return errors.NewValidationError("user_id is required")
	}

	if !membership.Role().IsValid() {
		return errors.NewValidationError("invalid role")
	}

	// Validar que validUntil es después de validFrom si existe
	if membership.ValidUntil() != nil {
		if membership.ValidUntil().Before(membership.ValidFrom()) {
			return errors.NewValidationError("valid_until must be after valid_from")
		}
	}

	return nil
}
