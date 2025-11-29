package data

import (
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

var (
	// School IDs
	schoolPrimariaID   = uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	schoolSecundariaID = uuid.MustParse("b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22")
	schoolTecnicoID    = uuid.MustParse("b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33")

	// Academic Unit IDs - Escuela Primaria
	primerGradoID   = uuid.MustParse("c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	segundoGradoID  = uuid.MustParse("c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22")
	tercerGradoID   = uuid.MustParse("c3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33")
	seccionPrimer1A = uuid.MustParse("c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44")
	seccionPrimer1B = uuid.MustParse("c5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55")

	// Academic Unit IDs - Colegio Secundario
	primerAnioID    = uuid.MustParse("c6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66")
	segundoAnioID   = uuid.MustParse("c7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77")
	seccionAnio1S1  = uuid.MustParse("c8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88")
	seccionAnio1S2  = uuid.MustParse("c9eebc99-9c0b-4ef8-bb6d-6bb9bd380a99")

	// Academic Unit IDs - Instituto Técnico
	deptProgramacionID = uuid.MustParse("caeebc99-9c0b-4ef8-bb6d-6bb9bd380aaa")
	deptBasesDatosID   = uuid.MustParse("cbeebc99-9c0b-4ef8-bb6d-6bb9bd380abb")

	// Timestamp base
	baseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	academicUnitsMap map[uuid.UUID]*entities.AcademicUnit
)

func init() {
	academicUnitsMap = make(map[uuid.UUID]*entities.AcademicUnit)

	// ESCUELA PRIMARIA - Primer Grado
	primerGrado := &entities.AcademicUnit{
		ID:           primerGradoID,
		SchoolID:     schoolPrimariaID,
		ParentUnitID: uuid.NullUUID{Valid: false},
		Name:         "Primer Grado",
		Code:         "P-G1",
		Type:         "grade",
		Level:        1,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 60, "shift": "morning", "building": "Edificio A"}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[primerGradoID] = primerGrado

	// Sección A - Primer Grado
	seccionA1 := &entities.AcademicUnit{
		ID:           seccionPrimer1A,
		SchoolID:     schoolPrimariaID,
		ParentUnitID: uuid.NullUUID{UUID: primerGradoID, Valid: true},
		Name:         "Sección A",
		Code:         "P-G1-A",
		Type:         "section",
		Level:        1,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 30, "classroom": "101", "shift": "morning"}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[seccionPrimer1A] = seccionA1

	// Sección B - Primer Grado
	seccionB1 := &entities.AcademicUnit{
		ID:           seccionPrimer1B,
		SchoolID:     schoolPrimariaID,
		ParentUnitID: uuid.NullUUID{UUID: primerGradoID, Valid: true},
		Name:         "Sección B",
		Code:         "P-G1-B",
		Type:         "section",
		Level:        1,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 30, "classroom": "102", "shift": "morning"}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[seccionPrimer1B] = seccionB1

	// ESCUELA PRIMARIA - Segundo Grado
	segundoGrado := &entities.AcademicUnit{
		ID:           segundoGradoID,
		SchoolID:     schoolPrimariaID,
		ParentUnitID: uuid.NullUUID{Valid: false},
		Name:         "Segundo Grado",
		Code:         "P-G2",
		Type:         "grade",
		Level:        2,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 55, "shift": "morning", "building": "Edificio A"}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[segundoGradoID] = segundoGrado

	// ESCUELA PRIMARIA - Tercer Grado
	tercerGrado := &entities.AcademicUnit{
		ID:           tercerGradoID,
		SchoolID:     schoolPrimariaID,
		ParentUnitID: uuid.NullUUID{Valid: false},
		Name:         "Tercer Grado",
		Code:         "P-G3",
		Type:         "grade",
		Level:        3,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 50, "shift": "afternoon", "building": "Edificio B"}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[tercerGradoID] = tercerGrado

	// COLEGIO SECUNDARIO - Primer Año
	primerAnio := &entities.AcademicUnit{
		ID:           primerAnioID,
		SchoolID:     schoolSecundariaID,
		ParentUnitID: uuid.NullUUID{Valid: false},
		Name:         "Primer Año",
		Code:         "S-A1",
		Type:         "grade",
		Level:        1,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 80, "shift": "morning", "building": "Pabellón Central"}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[primerAnioID] = primerAnio

	// Sección 1 - Primer Año
	seccionAnio1 := &entities.AcademicUnit{
		ID:           seccionAnio1S1,
		SchoolID:     schoolSecundariaID,
		ParentUnitID: uuid.NullUUID{UUID: primerAnioID, Valid: true},
		Name:         "Sección 1",
		Code:         "S-A1-S1",
		Type:         "section",
		Level:        1,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 40, "classroom": "201", "shift": "morning", "orientation": "general"}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[seccionAnio1S1] = seccionAnio1

	// Sección 2 - Primer Año
	seccionAnio2 := &entities.AcademicUnit{
		ID:           seccionAnio1S2,
		SchoolID:     schoolSecundariaID,
		ParentUnitID: uuid.NullUUID{UUID: primerAnioID, Valid: true},
		Name:         "Sección 2",
		Code:         "S-A1-S2",
		Type:         "section",
		Level:        1,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 40, "classroom": "202", "shift": "morning", "orientation": "general"}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[seccionAnio1S2] = seccionAnio2

	// COLEGIO SECUNDARIO - Segundo Año
	segundoAnio := &entities.AcademicUnit{
		ID:           segundoAnioID,
		SchoolID:     schoolSecundariaID,
		ParentUnitID: uuid.NullUUID{Valid: false},
		Name:         "Segundo Año",
		Code:         "S-A2",
		Type:         "grade",
		Level:        2,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 75, "shift": "afternoon", "building": "Pabellón Central"}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[segundoAnioID] = segundoAnio

	// INSTITUTO TÉCNICO - Departamento Programación I
	deptProgramacion := &entities.AcademicUnit{
		ID:           deptProgramacionID,
		SchoolID:     schoolTecnicoID,
		ParentUnitID: uuid.NullUUID{Valid: false},
		Name:         "Programación I",
		Code:         "T-PROG1",
		Type:         "department",
		Level:        1,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 35, "lab": "Laboratorio 1", "hours_per_week": 6, "professor_count": 2}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[deptProgramacionID] = deptProgramacion

	// INSTITUTO TÉCNICO - Departamento Bases de Datos
	deptBasesDatos := &entities.AcademicUnit{
		ID:           deptBasesDatosID,
		SchoolID:     schoolTecnicoID,
		ParentUnitID: uuid.NullUUID{Valid: false},
		Name:         "Bases de Datos",
		Code:         "T-BD",
		Type:         "department",
		Level:        1,
		AcademicYear: 2024,
		Metadata:     []byte(`{"capacity": 35, "lab": "Laboratorio 2", "hours_per_week": 5, "professor_count": 2}`),
		CreatedAt:    baseTime,
		UpdatedAt:    baseTime,
		DeletedAt:    sql.NullTime{Valid: false},
	}
	academicUnitsMap[deptBasesDatosID] = deptBasesDatos
}

// GetAcademicUnits retorna todas las unidades académicas en un mapa
func GetAcademicUnits() map[uuid.UUID]*entities.AcademicUnit {
	return academicUnitsMap
}

// GetAcademicUnitsBySchool retorna todas las unidades académicas de una escuela específica
func GetAcademicUnitsBySchool(schoolID uuid.UUID) []*entities.AcademicUnit {
	var units []*entities.AcademicUnit
	for _, unit := range academicUnitsMap {
		if unit.SchoolID == schoolID {
			units = append(units, unit)
		}
	}
	return units
}

// GetRootAcademicUnits retorna las unidades académicas raíz (sin padre) de una escuela
func GetRootAcademicUnits(schoolID uuid.UUID) []*entities.AcademicUnit {
	var rootUnits []*entities.AcademicUnit
	for _, unit := range academicUnitsMap {
		if unit.SchoolID == schoolID && !unit.ParentUnitID.Valid {
			rootUnits = append(rootUnits, unit)
		}
	}
	return rootUnits
}

// GetChildAcademicUnits retorna las unidades académicas hijas de una unidad padre específica
func GetChildAcademicUnits(parentID uuid.UUID) []*entities.AcademicUnit {
	var childUnits []*entities.AcademicUnit
	for _, unit := range academicUnitsMap {
		if unit.ParentUnitID.Valid && unit.ParentUnitID.UUID == parentID {
			childUnits = append(childUnits, unit)
		}
	}
	return childUnits
}

// BuildHierarchyPath construye la ruta jerárquica completa desde la raíz hasta la unidad especificada
func BuildHierarchyPath(unitID uuid.UUID) []*entities.AcademicUnit {
	var path []*entities.AcademicUnit

	currentUnit, exists := academicUnitsMap[unitID]
	if !exists {
		return path
	}

	// Construir el path desde la unidad actual hacia arriba
	for currentUnit != nil {
		// Insertar al inicio para mantener el orden de raíz a hoja
		path = append([]*entities.AcademicUnit{currentUnit}, path...)

		// Buscar el padre
		if currentUnit.ParentUnitID.Valid {
			currentUnit = academicUnitsMap[currentUnit.ParentUnitID.UUID]
		} else {
			break
		}
	}

	return path
}
