package valueobject

import (
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// MembershipRole representa el rol de un usuario en una unidad académica
type MembershipRole string

const (
	// RoleStudent representa un estudiante
	RoleStudent MembershipRole = "student"
	// RoleTeacher representa un profesor
	RoleTeacher MembershipRole = "teacher"
	// RoleCoordinator representa un coordinador
	RoleCoordinator MembershipRole = "coordinator"
	// RoleAdmin representa un administrador
	RoleAdmin MembershipRole = "admin"
	// RoleAssistant representa un asistente
	RoleAssistant MembershipRole = "assistant"
)

// NewMembershipRole crea un nuevo MembershipRole validando que sea válido
func NewMembershipRole(value string) (MembershipRole, error) {
	role := MembershipRole(value)
	if !role.IsValid() {
		return "", errors.NewValidationError("invalid membership role: " + value)
	}
	return role, nil
}

// IsValid verifica si el rol es válido
func (mr MembershipRole) IsValid() bool {
	switch mr {
	case RoleStudent, RoleTeacher, RoleCoordinator, RoleAdmin, RoleAssistant:
		return true
	default:
		return false
	}
}

// String retorna la representación en string
func (mr MembershipRole) String() string {
	return string(mr)
}

// HasPermission verifica si el rol tiene un permiso específico
func (mr MembershipRole) HasPermission(permission string) bool {
	// Admin tiene todos los permisos
	if mr == RoleAdmin {
		return true
	}

	// Coordinador tiene permisos de gestión
	if mr == RoleCoordinator {
		switch permission {
		case "manage_unit", "view_members", "add_members":
			return true
		}
	}

	// Profesor puede ver miembros
	if mr == RoleTeacher {
		switch permission {
		case "view_members":
			return true
		}
	}

	return false
}

// IsTeachingRole verifica si es un rol de enseñanza
func (mr MembershipRole) IsTeachingRole() bool {
	return mr == RoleTeacher || mr == RoleCoordinator || mr == RoleAssistant
}

// IsAdministrativeRole verifica si es un rol administrativo
func (mr MembershipRole) IsAdministrativeRole() bool {
	return mr == RoleAdmin || mr == RoleCoordinator
}
