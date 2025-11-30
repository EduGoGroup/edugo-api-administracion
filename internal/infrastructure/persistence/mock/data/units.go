package data

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

var (
	// UUIDs fijos para las unidades organizacionales de prueba
	unitID1 = uuid.MustParse("f1000000-0000-0000-0000-000000000011") // Departamento de Matemáticas
	unitID2 = uuid.MustParse("f2000000-0000-0000-0000-000000000022") // Departamento de Ciencias
	unitID3 = uuid.MustParse("f3000000-0000-0000-0000-000000000033") // Coordinación Académica
	unitID4 = uuid.MustParse("f4000000-0000-0000-0000-000000000044") // Grupo de Docentes
)

// GetUnits retorna un map con las 4 unidades organizacionales de prueba
func GetUnits() map[uuid.UUID]*entities.Unit {
	units := make(map[uuid.UUID]*entities.Unit)

	// Departamento de Matemáticas - Escuela Primaria (raíz)
	units[unitID1] = &entities.Unit{
		ID:           unitID1,
		SchoolID:     schoolID1, // Escuela Primaria Demo
		ParentUnitID: nil,       // Es una unidad raíz
		Name:         "Departamento de Matemáticas",
		Description:  strPtr("Departamento responsable de la enseñanza de matemáticas en nivel primario"),
		IsActive:     true,
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
	}

	// Departamento de Ciencias - Escuela Primaria (raíz)
	units[unitID2] = &entities.Unit{
		ID:           unitID2,
		SchoolID:     schoolID1, // Escuela Primaria Demo
		ParentUnitID: nil,       // Es una unidad raíz
		Name:         "Departamento de Ciencias",
		Description:  strPtr("Departamento responsable de la enseñanza de ciencias naturales en nivel primario"),
		IsActive:     true,
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
	}

	// Coordinación Académica - Colegio Secundario (raíz)
	units[unitID3] = &entities.Unit{
		ID:           unitID3,
		SchoolID:     schoolID2, // Colegio Secundario Demo
		ParentUnitID: nil,       // Es una unidad raíz
		Name:         "Coordinación Académica",
		Description:  strPtr("Coordinación responsable de la gestión académica en nivel secundario"),
		IsActive:     true,
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
	}

	// Grupo de Docentes - Colegio Secundario (sub-unidad de Coordinación Académica)
	units[unitID4] = &entities.Unit{
		ID:           unitID4,
		SchoolID:     schoolID2,   // Colegio Secundario Demo
		ParentUnitID: &unitID3,    // Padre: Coordinación Académica
		Name:         "Grupo de Docentes",
		Description:  strPtr("Grupo de docentes pertenecientes a la coordinación académica"),
		IsActive:     true,
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
	}

	return units
}

// GetUnitsBySchool retorna todas las unidades asociadas a una escuela
func GetUnitsBySchool(schoolID uuid.UUID) []*entities.Unit {
	units := GetUnits()
	var result []*entities.Unit

	for _, unit := range units {
		if unit.SchoolID == schoolID {
			result = append(result, unit)
		}
	}

	return result
}

// GetRootUnits retorna solo las unidades raíz (sin padre) de una escuela
func GetRootUnits(schoolID uuid.UUID) []*entities.Unit {
	units := GetUnits()
	var result []*entities.Unit

	for _, unit := range units {
		if unit.SchoolID == schoolID && unit.ParentUnitID == nil {
			result = append(result, unit)
		}
	}

	return result
}

// GetUnitByID retorna una unidad por su ID
func GetUnitByID(id uuid.UUID) *entities.Unit {
	units := GetUnits()
	return units[id]
}

// GetUnitByName retorna una unidad por su nombre dentro de una escuela
func GetUnitByName(schoolID uuid.UUID, name string) *entities.Unit {
	units := GetUnits()
	for _, unit := range units {
		if unit.SchoolID == schoolID && unit.Name == name {
			return unit
		}
	}
	return nil
}

// GetUnitChildren retorna todas las sub-unidades de una unidad padre
func GetUnitChildren(parentUnitID uuid.UUID) []*entities.Unit {
	units := GetUnits()
	var result []*entities.Unit

	for _, unit := range units {
		if unit.ParentUnitID != nil && *unit.ParentUnitID == parentUnitID {
			result = append(result, unit)
		}
	}

	return result
}
