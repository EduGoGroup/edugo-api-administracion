package valueobject

import "fmt"

// MembershipRole representa el rol de un usuario en una unidad académica
type MembershipRole string

const (
	RoleTeacher     MembershipRole = "teacher"
	RoleStudent     MembershipRole = "student"
	RoleGuardian    MembershipRole = "guardian"
	RoleCoordinator MembershipRole = "coordinator"
	RoleAdmin       MembershipRole = "admin"
	RoleAssistant   MembershipRole = "assistant"
	RoleDirector    MembershipRole = "director"
	RoleObserver    MembershipRole = "observer"
)

var validMembershipRoles = map[MembershipRole]bool{
	RoleTeacher:     true,
	RoleStudent:     true,
	RoleGuardian:    true,
	RoleCoordinator: true,
	RoleAdmin:       true,
	RoleAssistant:   true,
	RoleDirector:    true,
	RoleObserver:    true,
}

// IsValid verifica si el rol es válido
func (r MembershipRole) IsValid() bool {
	return validMembershipRoles[r]
}

// String retorna el rol como string
func (r MembershipRole) String() string {
	return string(r)
}

// ParseMembershipRole convierte un string a MembershipRole
func ParseMembershipRole(s string) (MembershipRole, error) {
	role := MembershipRole(s)
	if !role.IsValid() {
		return "", fmt.Errorf("invalid membership role: %s", s)
	}
	return role, nil
}

// AllMembershipRoles retorna todos los roles válidos
func AllMembershipRoles() []MembershipRole {
	return []MembershipRole{
		RoleTeacher,
		RoleStudent,
		RoleGuardian,
		RoleCoordinator,
		RoleAdmin,
		RoleAssistant,
		RoleDirector,
		RoleObserver,
	}
}

// AllMembershipRolesStrings retorna todos los roles como strings
func AllMembershipRolesStrings() []string {
	roles := AllMembershipRoles()
	result := make([]string, len(roles))
	for i, r := range roles {
		result[i] = string(r)
	}
	return result
}
