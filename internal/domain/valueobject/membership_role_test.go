package valueobject_test

import (
	"testing"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestMembershipRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     valueobject.MembershipRole
		expected bool
	}{
		{
			name:     "teacher role is valid",
			role:     valueobject.RoleTeacher,
			expected: true,
		},
		{
			name:     "student role is valid",
			role:     valueobject.RoleStudent,
			expected: true,
		},
		{
			name:     "guardian role is valid",
			role:     valueobject.RoleGuardian,
			expected: true,
		},
		{
			name:     "coordinator role is valid",
			role:     valueobject.RoleCoordinator,
			expected: true,
		},
		{
			name:     "admin role is valid",
			role:     valueobject.RoleAdmin,
			expected: true,
		},
		{
			name:     "assistant role is valid",
			role:     valueobject.RoleAssistant,
			expected: true,
		},
		{
			name:     "director role is valid",
			role:     valueobject.RoleDirector,
			expected: true,
		},
		{
			name:     "observer role is valid",
			role:     valueobject.RoleObserver,
			expected: true,
		},
		{
			name:     "invalid role",
			role:     "invalid_role",
			expected: false,
		},
		{
			name:     "empty role",
			role:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.IsValid())
		})
	}
}

func TestParseMembershipRole(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr bool
		expected    valueobject.MembershipRole
	}{
		{
			name:        "parse teacher role",
			input:       "teacher",
			expectedErr: false,
			expected:    valueobject.RoleTeacher,
		},
		{
			name:        "parse student role",
			input:       "student",
			expectedErr: false,
			expected:    valueobject.RoleStudent,
		},
		{
			name:        "parse guardian role",
			input:       "guardian",
			expectedErr: false,
			expected:    valueobject.RoleGuardian,
		},
		{
			name:        "parse coordinator role",
			input:       "coordinator",
			expectedErr: false,
			expected:    valueobject.RoleCoordinator,
		},
		{
			name:        "parse admin role",
			input:       "admin",
			expectedErr: false,
			expected:    valueobject.RoleAdmin,
		},
		{
			name:        "parse assistant role",
			input:       "assistant",
			expectedErr: false,
			expected:    valueobject.RoleAssistant,
		},
		{
			name:        "parse director role",
			input:       "director",
			expectedErr: false,
			expected:    valueobject.RoleDirector,
		},
		{
			name:        "parse observer role",
			input:       "observer",
			expectedErr: false,
			expected:    valueobject.RoleObserver,
		},
		{
			name:        "parse invalid role",
			input:       "invalid",
			expectedErr: true,
			expected:    "",
		},
		{
			name:        "parse empty role",
			input:       "",
			expectedErr: true,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := valueobject.ParseMembershipRole(tt.input)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, role)
			}
		})
	}
}

func TestMembershipRole_String(t *testing.T) {
	tests := []struct {
		name     string
		role     valueobject.MembershipRole
		expected string
	}{
		{
			name:     "teacher role to string",
			role:     valueobject.RoleTeacher,
			expected: "teacher",
		},
		{
			name:     "student role to string",
			role:     valueobject.RoleStudent,
			expected: "student",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.String())
		})
	}
}

func TestAllMembershipRoles(t *testing.T) {
	roles := valueobject.AllMembershipRoles()

	assert.Len(t, roles, 8)
	assert.Contains(t, roles, valueobject.RoleTeacher)
	assert.Contains(t, roles, valueobject.RoleStudent)
	assert.Contains(t, roles, valueobject.RoleGuardian)
	assert.Contains(t, roles, valueobject.RoleCoordinator)
	assert.Contains(t, roles, valueobject.RoleAdmin)
	assert.Contains(t, roles, valueobject.RoleAssistant)
	assert.Contains(t, roles, valueobject.RoleDirector)
	assert.Contains(t, roles, valueobject.RoleObserver)
}

func TestAllMembershipRolesStrings(t *testing.T) {
	roles := valueobject.AllMembershipRolesStrings()

	assert.Len(t, roles, 8)
	assert.Contains(t, roles, "teacher")
	assert.Contains(t, roles, "student")
	assert.Contains(t, roles, "guardian")
	assert.Contains(t, roles, "coordinator")
	assert.Contains(t, roles, "admin")
	assert.Contains(t, roles, "assistant")
	assert.Contains(t, roles, "director")
	assert.Contains(t, roles, "observer")
}
