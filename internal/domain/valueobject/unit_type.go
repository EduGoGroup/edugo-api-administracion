package valueobject

import "fmt"

// UnitType representa un tipo válido de unidad académica
type UnitType string

// Tipos de unidades académicas válidos
const (
	UnitTypeSchool     UnitType = "school"
	UnitTypeGrade      UnitType = "grade"
	UnitTypeSection    UnitType = "section"
	UnitTypeClub       UnitType = "club"
	UnitTypeDepartment UnitType = "department"
	UnitTypeCourse     UnitType = "course"
	UnitTypeGroup      UnitType = "group"
)

// validUnitTypes contiene todos los tipos de unidad válidos
var validUnitTypes = map[UnitType]bool{
	UnitTypeSchool:     true,
	UnitTypeGrade:      true,
	UnitTypeSection:    true,
	UnitTypeClub:       true,
	UnitTypeDepartment: true,
	UnitTypeCourse:     true,
	UnitTypeGroup:      true,
}

// IsValid verifica si el tipo de unidad es válido
func (t UnitType) IsValid() bool {
	return validUnitTypes[t]
}

// String retorna el tipo como string
func (t UnitType) String() string {
	return string(t)
}

// ParseUnitType convierte un string a UnitType
func ParseUnitType(s string) (UnitType, error) {
	ut := UnitType(s)
	if !ut.IsValid() {
		return "", fmt.Errorf("invalid unit type: %s", s)
	}
	return ut, nil
}

// ValidUnitTypes retorna una lista de todos los tipos válidos
func ValidUnitTypes() []UnitType {
	return []UnitType{
		UnitTypeSchool,
		UnitTypeGrade,
		UnitTypeSection,
		UnitTypeClub,
		UnitTypeDepartment,
		UnitTypeCourse,
		UnitTypeGroup,
	}
}

// ValidUnitTypeStrings retorna los tipos como slice de strings
func ValidUnitTypeStrings() []string {
	types := ValidUnitTypes()
	result := make([]string, len(types))
	for i, t := range types {
		result[i] = t.String()
	}
	return result
}

// CanHaveParent indica si este tipo de unidad puede tener una unidad padre
func (t UnitType) CanHaveParent() bool {
	switch t {
	case UnitTypeSchool:
		return false // School es la raíz
	default:
		return true
	}
}

// AllowedChildTypes retorna los tipos de unidades que pueden ser hijos de este tipo
func (t UnitType) AllowedChildTypes() []UnitType {
	switch t {
	case UnitTypeSchool:
		return []UnitType{UnitTypeGrade, UnitTypeDepartment, UnitTypeClub}
	case UnitTypeGrade:
		return []UnitType{UnitTypeSection, UnitTypeCourse}
	case UnitTypeSection:
		return []UnitType{UnitTypeGroup}
	case UnitTypeDepartment:
		return []UnitType{UnitTypeCourse}
	default:
		return []UnitType{}
	}
}
