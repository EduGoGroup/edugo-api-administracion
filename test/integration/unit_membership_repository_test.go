//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	postgresRepo "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/postgres/repository"
)

func TestMembershipRepository_Create(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	schoolRepo := postgresRepo.NewPostgresSchoolRepository(db)
	unitRepo := postgresRepo.NewPostgresAcademicUnitRepository(db)
	membershipRepo := postgresRepo.NewPostgresUnitMembershipRepository(db)

	// Setup: School + Unit
	school, _ := entity.NewSchool("Test School", "MEM001", "")
	schoolRepo.Create(ctx, school)

	unitType, _ := valueobject.NewUnitType("section")
	unit, _ := entity.NewAcademicUnit(school.ID(), unitType, "Section A", "", "")
	unitRepo.Create(ctx, unit)

	// Crear membership
	userID, _ := valueobject.UserIDFromString("user-123")
	role, _ := valueobject.NewMembershipRole("student")
	membership, err := entity.NewUnitMembership(unit.ID(), userID, role, time.Now(), nil)
	require.NoError(t, err)

	err = membershipRepo.Create(ctx, membership)
	assert.NoError(t, err)

	// Verificar
	found, err := membershipRepo.FindByID(ctx, membership.ID())
	assert.NoError(t, err)
	assert.Equal(t, membership.Role(), found.Role())
}

func TestMembershipRepository_FindByUnit(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	schoolRepo := postgresRepo.NewPostgresSchoolRepository(db)
	unitRepo := postgresRepo.NewPostgresAcademicUnitRepository(db)
	membershipRepo := postgresRepo.NewPostgresUnitMembershipRepository(db)

	// Setup
	school, _ := entity.NewSchool("Test School", "MEM002", "")
	schoolRepo.Create(ctx, school)

	unitType, _ := valueobject.NewUnitType("section")
	unit, _ := entity.NewAcademicUnit(school.ID(), unitType, "Section B", "", "")
	unitRepo.Create(ctx, unit)

	// Crear 2 membresías
	userID1, _ := valueobject.UserIDFromString("user-1")
	userID2, _ := valueobject.UserIDFromString("user-2")
	role, _ := valueobject.NewMembershipRole("student")

	m1, _ := entity.NewUnitMembership(unit.ID(), userID1, role, time.Now(), nil)
	m2, _ := entity.NewUnitMembership(unit.ID(), userID2, role, time.Now(), nil)

	membershipRepo.Create(ctx, m1)
	membershipRepo.Create(ctx, m2)

	// Buscar por unidad
	memberships, err := membershipRepo.FindByUnit(ctx, unit.ID(), true)
	assert.NoError(t, err)
	assert.Len(t, memberships, 2)
}

func TestMembershipRepository_FindActiveAt(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	schoolRepo := postgresRepo.NewPostgresSchoolRepository(db)
	unitRepo := postgresRepo.NewPostgresAcademicUnitRepository(db)
	membershipRepo := postgresRepo.NewPostgresUnitMembershipRepository(db)

	// Setup
	school, _ := entity.NewSchool("Test School", "MEM003", "")
	schoolRepo.Create(ctx, school)

	unitType, _ := valueobject.NewUnitType("section")
	unit, _ := entity.NewAcademicUnit(school.ID(), unitType, "Section C", "", "")
	unitRepo.Create(ctx, unit)

	// Crear membership activa
	userID, _ := valueobject.UserIDFromString("user-active")
	role, _ := valueobject.NewMembershipRole("student")
	validUntil := time.Now().AddDate(1, 0, 0) // Válida por 1 año

	membership, _ := entity.NewUnitMembership(unit.ID(), userID, role, time.Now(), &validUntil)
	membershipRepo.Create(ctx, membership)

	// Buscar activas
	active, err := membershipRepo.FindActiveAt(ctx, unit.ID(), time.Now())
	assert.NoError(t, err)
	assert.Len(t, active, 1)
}
