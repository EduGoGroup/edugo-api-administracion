package valueobject

import (
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// UnitType representa el tipo de unidad académica
type UnitType string

const (
	// UnitTypeSchool representa el nivel raíz de la escuela
	UnitTypeSchool UnitType = "school"
	// UnitTypeGrade representa un grado académico
	UnitTypeGrade UnitType = "grade"
	// UnitTypeSection representa una sección de grado
	UnitTypeSection UnitType = "section"
	// UnitTypeClub representa un club extracurricular
	UnitTypeClub UnitType = "club"
	// UnitTypeDepartment representa un departamento administrativo
	UnitTypeDepartment UnitType = "department"
)

// NewUnitType crea un nuevo UnitType validando que sea un valor válido
func NewUnitType(value string) (UnitType, error) {
	ut := UnitType(value)
	if !ut.IsValid() {
		return "", errors.NewValidationError("invalid unit type: " + value)
	}
	return ut, nil
}

// IsValid verifica si el tipo es válido
func (ut UnitType) IsValid() bool {
	switch ut {
	case UnitTypeSchool, UnitTypeGrade, UnitTypeSection, UnitTypeClub, UnitTypeDepartment:
		return true
	default:
		return false
	}
}

// String retorna la representación en string
func (ut UnitType) String() string {
	return string(ut)
}

// CanHaveChildren determina si este tipo puede tener hijos
func (ut UnitType) CanHaveChildren() bool {
	switch ut {
	case UnitTypeSchool, UnitTypeGrade:
		return true
	case UnitTypeSection, UnitTypeClub, UnitTypeDepartment:
		return false
	default:
		return false
	}
}

// AllowedChildTypes retorna los tipos de hijos permitidos
func (ut UnitType) AllowedChildTypes() []UnitType {
	switch ut {
	case UnitTypeSchool:
		return []UnitType{UnitTypeGrade, UnitTypeClub, UnitTypeDepartment}
	case UnitTypeGrade:
		return []UnitType{UnitTypeSection}
	default:
		return []UnitType{}
	}
}
