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
	unit, err := entity.NewAcademicUnit(school.ID(), unitType, "1st Grade", "G01")
	require.NoError(t, err)

	err = unitRepo.Create(ctx, unit)
	assert.NoError(t, err)

	// Verificar
	found, err := unitRepo.FindByID(ctx, unit.ID(), false)
	assert.NoError(t, err)
	assert.Equal(t, unit.DisplayName(), found.DisplayName())
}

func TestAcademicUnitRepository_FindBySchoolID(t *testing.T) {
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
	unit1, _ := entity.NewAcademicUnit(school.ID(), unitType, "Grade 1", "")
	unit2, _ := entity.NewAcademicUnit(school.ID(), unitType, "Grade 2", "")
	unit3, _ := entity.NewAcademicUnit(school.ID(), unitType, "Grade 3", "")

	unitRepo.Create(ctx, unit1)
	unitRepo.Create(ctx, unit2)
	unitRepo.Create(ctx, unit3)

	// Buscar por escuela
	units, err := unitRepo.FindBySchoolID(ctx, school.ID(), false)
	assert.NoError(t, err)
	assert.Len(t, units, 3)
}
