//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	postgresRepo "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/postgres/repository"
)

func TestAcademicUnitRepository_Create(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	schoolRepo := postgresRepo.NewPostgresSchoolRepository(db)
	unitRepo := postgresRepo.NewPostgresAcademicUnitRepository(db)

	// Crear escuela primero
	school, _ := entity.NewSchool("Test School", "UNIT001", "")
	schoolRepo.Create(ctx, school)

	// Crear unidad
	unitType, _ := valueobject.NewUnitType("grade")
	unit, err := entity.NewAcademicUnit(school.ID(), unitType, "1st Grade", "G01", "")
	require.NoError(t, err)

	err = unitRepo.Create(ctx, unit)
	assert.NoError(t, err)

	// Verificar
	found, err := unitRepo.FindByID(ctx, unit.ID())
	assert.NoError(t, err)
	assert.Equal(t, unit.DisplayName(), found.DisplayName())
}

func TestAcademicUnitRepository_FindBySchool(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	schoolRepo := postgresRepo.NewPostgresSchoolRepository(db)
	unitRepo := postgresRepo.NewPostgresAcademicUnitRepository(db)

	// Crear escuela
	school, _ := entity.NewSchool("Test School", "UNIT002", "")
	schoolRepo.Create(ctx, school)

	// Crear 3 unidades
	unitType, _ := valueobject.NewUnitType("grade")
	unit1, _ := entity.NewAcademicUnit(school.ID(), unitType, "Grade 1", "", "")
	unit2, _ := entity.NewAcademicUnit(school.ID(), unitType, "Grade 2", "", "")
	unit3, _ := entity.NewAcademicUnit(school.ID(), unitType, "Grade 3", "", "")

	unitRepo.Create(ctx, unit1)
	unitRepo.Create(ctx, unit2)
	unitRepo.Create(ctx, unit3)

	// Buscar por escuela
	units, err := unitRepo.FindBySchool(ctx, school.ID(), false)
	assert.NoError(t, err)
	assert.Len(t, units, 3)
}

func TestAcademicUnitRepository_SoftDelete(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	schoolRepo := postgresRepo.NewPostgresSchoolRepository(db)
	unitRepo := postgresRepo.NewPostgresAcademicUnitRepository(db)

	// Crear escuela y unidad
	school, _ := entity.NewSchool("Test School", "UNIT003", "")
	schoolRepo.Create(ctx, school)

	unitType, _ := valueobject.NewUnitType("grade")
	unit, _ := entity.NewAcademicUnit(school.ID(), unitType, "To Delete", "", "")
	unitRepo.Create(ctx, unit)

	// Soft delete
	err := unitRepo.Delete(ctx, unit.ID())
	assert.NoError(t, err)

	// No aparece en búsquedas normales
	units, _ := unitRepo.FindBySchool(ctx, school.ID(), false)
	assert.Len(t, units, 0)

	// Sí aparece con includeDeleted
	unitsDeleted, _ := unitRepo.FindBySchool(ctx, school.ID(), true)
	assert.Len(t, unitsDeleted, 1)
}

func TestAcademicUnitRepository_GetHierarchyPath(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	schoolRepo := postgresRepo.NewPostgresSchoolRepository(db)
	unitRepo := postgresRepo.NewPostgresAcademicUnitRepository(db)

	// Crear escuela
	school, _ := entity.NewSchool("Test School", "UNIT004", "")
	schoolRepo.Create(ctx, school)

	// Crear jerarquía: Grade -> Section
	gradeType, _ := valueobject.NewUnitType("grade")
	grade, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 1", "", "")
	unitRepo.Create(ctx, grade)

	sectionType, _ := valueobject.NewUnitType("section")
	section, _ := entity.NewAcademicUnit(school.ID(), sectionType, "Section A", "", "")
	section.SetParent(&grade)
	unitRepo.Create(ctx, section)

	// Obtener path jerárquico
	path, err := unitRepo.GetHierarchyPath(ctx, section.ID())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(path), 2) // Al menos grade y section
}
