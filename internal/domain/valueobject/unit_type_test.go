package valueobject_test

import (
	"testing"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestUnitType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		unitType valueobject.UnitType
		expected bool
	}{
		{
			name:     "school type is valid",
			unitType: valueobject.UnitTypeSchool,
			expected: true,
		},
		{
			name:     "grade type is valid",
			unitType: valueobject.UnitTypeGrade,
			expected: true,
		},
		{
			name:     "section type is valid",
			unitType: valueobject.UnitTypeSection,
			expected: true,
		},
		{
			name:     "club type is valid",
			unitType: valueobject.UnitTypeClub,
			expected: true,
		},
		{
			name:     "department type is valid",
			unitType: valueobject.UnitTypeDepartment,
			expected: true,
		},
		{
			name:     "invalid type",
			unitType: "invalid_type",
			expected: false,
		},
		{
			name:     "empty type",
			unitType: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.unitType.IsValid())
		})
	}
}

func TestParseUnitType(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr bool
		expected    valueobject.UnitType
	}{
		{
			name:        "parse school type",
			input:       "school",
			expectedErr: false,
			expected:    valueobject.UnitTypeSchool,
		},
		{
			name:        "parse grade type",
			input:       "grade",
			expectedErr: false,
			expected:    valueobject.UnitTypeGrade,
		},
		{
			name:        "parse section type",
			input:       "section",
			expectedErr: false,
			expected:    valueobject.UnitTypeSection,
		},
		{
			name:        "parse club type",
			input:       "club",
			expectedErr: false,
			expected:    valueobject.UnitTypeClub,
		},
		{
			name:        "parse department type",
			input:       "department",
			expectedErr: false,
			expected:    valueobject.UnitTypeDepartment,
		},
		{
			name:        "parse invalid type",
			input:       "invalid",
			expectedErr: true,
			expected:    "",
		},
		{
			name:        "parse empty type",
			input:       "",
			expectedErr: true,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unitType, err := valueobject.ParseUnitType(tt.input)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, unitType)
			}
		})
	}
}

func TestUnitType_String(t *testing.T) {
	tests := []struct {
		name     string
		unitType valueobject.UnitType
		expected string
	}{
		{
			name:     "school type to string",
			unitType: valueobject.UnitTypeSchool,
			expected: "school",
		},
		{
			name:     "grade type to string",
			unitType: valueobject.UnitTypeGrade,
			expected: "grade",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.unitType.String())
		})
	}
}

func TestAllUnitTypes(t *testing.T) {
	types := valueobject.AllUnitTypes()

	assert.Len(t, types, 5)
	assert.Contains(t, types, valueobject.UnitTypeSchool)
	assert.Contains(t, types, valueobject.UnitTypeGrade)
	assert.Contains(t, types, valueobject.UnitTypeSection)
	assert.Contains(t, types, valueobject.UnitTypeClub)
	assert.Contains(t, types, valueobject.UnitTypeDepartment)
}

func TestAllUnitTypesStrings(t *testing.T) {
	types := valueobject.AllUnitTypesStrings()

	assert.Len(t, types, 5)
	assert.Contains(t, types, "school")
	assert.Contains(t, types, "grade")
	assert.Contains(t, types, "section")
	assert.Contains(t, types, "club")
	assert.Contains(t, types, "department")
}
