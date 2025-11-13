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

func TestSchoolRepository_Create(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	repo := postgresRepo.NewPostgresSchoolRepository(db)

	// Crear escuela
	school, err := entity.NewSchool("Test School", "TEST001", "123 Test St")
	require.NoError(t, err)

	err = repo.Create(ctx, school)
	assert.NoError(t, err)

	// Verificar que se creó
	found, err := repo.FindByID(ctx, school.ID())
	assert.NoError(t, err)
	assert.Equal(t, school.Name(), found.Name())
	assert.Equal(t, school.Code(), found.Code())
}

func TestSchoolRepository_FindByCode(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	repo := postgresRepo.NewPostgresSchoolRepository(db)

	// Crear escuela
	school, _ := entity.NewSchool("Test School", "FINDCODE", "123 Test St")
	err := repo.Create(ctx, school)
	require.NoError(t, err)

	// Buscar por código
	found, err := repo.FindByCode(ctx, "FINDCODE")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "FINDCODE", found.Code())
}

func TestSchoolRepository_ExistsByCode(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	repo := postgresRepo.NewPostgresSchoolRepository(db)

	// No existe
	exists, err := repo.ExistsByCode(ctx, "NOEXIST")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Crear escuela
	school, _ := entity.NewSchool("Test School", "EXISTS", "123 Test St")
	repo.Create(ctx, school)

	// Ahora existe
	exists, err = repo.ExistsByCode(ctx, "EXISTS")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestSchoolRepository_Update(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	repo := postgresRepo.NewPostgresSchoolRepository(db)

	// Crear escuela
	school, _ := entity.NewSchool("Original Name", "UPDATE", "123 Test St")
	repo.Create(ctx, school)

	// Actualizar
	email, _ := valueobject.NewEmail("test@example.com")
	school.UpdateContactInfo(&email, "555-1234")

	err := repo.Update(ctx, school)
	assert.NoError(t, err)

	// Verificar cambios
	found, _ := repo.FindByID(ctx, school.ID())
	assert.Equal(t, "test@example.com", found.ContactEmail().String())
	assert.Equal(t, "555-1234", found.ContactPhone())
}

func TestSchoolRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDBWithMigrations(t)
	defer cleanup()

	ctx := context.Background()
	repo := postgresRepo.NewPostgresSchoolRepository(db)

	// Crear escuela
	school, _ := entity.NewSchool("To Delete", "DELETE", "")
	repo.Create(ctx, school)

	// Eliminar
	err := repo.Delete(ctx, school.ID())
	assert.NoError(t, err)

	// Verificar que no existe
	_, err = repo.FindByID(ctx, school.ID())
	assert.Error(t, err)
}
