package valueobject

import "fmt"

// UnitType representa el tipo de una unidad académica
type UnitType string

const (
	UnitTypeSchool     UnitType = "school"
	UnitTypeGrade      UnitType = "grade"
	UnitTypeSection    UnitType = "section"
	UnitTypeClub       UnitType = "club"
	UnitTypeDepartment UnitType = "department"
)

var validUnitTypes = map[UnitType]bool{
	UnitTypeSchool:     true,
	UnitTypeGrade:      true,
	UnitTypeSection:    true,
	UnitTypeClub:       true,
	UnitTypeDepartment: true,
}

// IsValid verifica si el tipo es válido
func (t UnitType) IsValid() bool {
	return validUnitTypes[t]
}

// String retorna el tipo como string
func (t UnitType) String() string {
	return string(t)
}

// ParseUnitType convierte un string a UnitType
func ParseUnitType(s string) (UnitType, error) {
	unitType := UnitType(s)
	if !unitType.IsValid() {
		return "", fmt.Errorf("invalid unit type: %s", s)
	}
	return unitType, nil
}

// AllUnitTypes retorna todos los tipos válidos
func AllUnitTypes() []UnitType {
	return []UnitType{
		UnitTypeSchool,
		UnitTypeGrade,
		UnitTypeSection,
		UnitTypeClub,
		UnitTypeDepartment,
	}
}

// AllUnitTypesStrings retorna todos los tipos como strings
func AllUnitTypesStrings() []string {
	types := AllUnitTypes()
	result := make([]string, len(types))
	for i, t := range types {
		result[i] = string(t)
	}
	return result
}
