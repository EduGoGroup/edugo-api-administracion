package valueobject

import "fmt"

// MembershipRole representa un rol válido para membresías de unidades académicas
type MembershipRole string

// Roles válidos para membresías
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

// validMembershipRoles contiene todos los roles válidos
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

// ValidMembershipRoles retorna una lista de todos los roles válidos
func ValidMembershipRoles() []MembershipRole {
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

// ValidMembershipRoleStrings retorna los roles como slice de strings
func ValidMembershipRoleStrings() []string {
	roles := ValidMembershipRoles()
	result := make([]string, len(roles))
	for i, r := range roles {
		result[i] = r.String()
	}
	return result
}
