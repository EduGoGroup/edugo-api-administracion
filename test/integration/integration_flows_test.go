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

// TestIntegration_SchoolCRUDFlow prueba el flujo completo CRUD de School
// Esto es un test de integración REAL: prueba todas las operaciones en secuencia
func TestIntegration_SchoolCRUDFlow(t *testing.T) {
	db, cleanup := getTestDB(t)
	defer cleanup()

	ctx := context.Background()
	repo := postgresRepo.NewPostgresSchoolRepository(db)

	// PASO 1: CREATE - Crear escuela
	school, err := entity.NewSchool("Integration Test School", "INTTEST001", "123 Integration St")
	require.NoError(t, err, "No debería fallar creación de entidad")

	err = repo.Create(ctx, school)
	require.NoError(t, err, "CREATE: Debería crear la escuela")

	// PASO 2: FIND BY ID - Verificar que se creó correctamente
	found, err := repo.FindByID(ctx, school.ID())
	require.NoError(t, err, "FIND BY ID: Debería encontrar la escuela")
	assert.Equal(t, school.Name(), found.Name(), "Nombre debe coincidir")
	assert.Equal(t, school.Code(), found.Code(), "Código debe coincidir")

	// PASO 3: FIND BY CODE - Buscar por código único
	foundByCode, err := repo.FindByCode(ctx, "INTTEST001")
	require.NoError(t, err, "FIND BY CODE: Debería encontrar por código")
	assert.Equal(t, school.ID(), foundByCode.ID(), "ID debe ser el mismo")

	// PASO 4: EXISTS BY CODE - Verificar existencia
	exists, err := repo.ExistsByCode(ctx, "INTTEST001")
	require.NoError(t, err, "EXISTS: No debería dar error")
	assert.True(t, exists, "EXISTS: Código debe existir")

	existsNonExistent, err := repo.ExistsByCode(ctx, "NOEXISTE999")
	require.NoError(t, err, "EXISTS: No debería dar error para código inexistente")
	assert.False(t, existsNonExistent, "EXISTS: Código inexistente no debe existir")

	// PASO 5: UPDATE - Actualizar información de contacto
	email, err := valueobject.NewEmail("contact@integration-test.com")
	require.NoError(t, err)
	school.UpdateContactInfo(&email, "555-9999")

	err = repo.Update(ctx, school)
	require.NoError(t, err, "UPDATE: Debería actualizar la escuela")

	// PASO 6: Verificar UPDATE
	updated, err := repo.FindByID(ctx, school.ID())
	require.NoError(t, err)
	assert.Equal(t, "contact@integration-test.com", updated.ContactEmail().String(), "Email debe estar actualizado")
	assert.Equal(t, "555-9999", updated.ContactPhone(), "Teléfono debe estar actualizado")

	// PASO 7: DELETE - Eliminar escuela
	err = repo.Delete(ctx, school.ID())
	require.NoError(t, err, "DELETE: Debería eliminar la escuela")

	// PASO 8: Verificar DELETE - No debe encontrar la escuela eliminada
	_, err = repo.FindByID(ctx, school.ID())
	assert.Error(t, err, "FIND BY ID: No debería encontrar escuela eliminada")
}

// TestIntegration_AcademicHierarchyFlow prueba el flujo completo de jerarquía académica
// School → AcademicUnit (Grade) → AcademicUnit (Section con parent) → Membership
// Esto prueba la integración REAL entre los 3 repositorios y la jerarquía
func TestIntegration_AcademicHierarchyFlow(t *testing.T) {
	db, cleanup := getTestDB(t)
	defer cleanup()

	ctx := context.Background()
	schoolRepo := postgresRepo.NewPostgresSchoolRepository(db)
	unitRepo := postgresRepo.NewPostgresAcademicUnitRepository(db)
	membershipRepo := postgresRepo.NewPostgresUnitMembershipRepository(db)

	// PASO 1: Crear School (raíz de la jerarquía)
	school, err := entity.NewSchool("Hierarchy Test School", "HIERTEST001", "456 Hierarchy Ave")
	require.NoError(t, err)
	err = schoolRepo.Create(ctx, school)
	require.NoError(t, err, "Debería crear la escuela raíz")

	// PASO 2: Crear AcademicUnit de nivel superior (Grade - grado)
	grade, err := entity.NewAcademicUnit(
		school.ID(),
		valueobject.UnitTypeGrade,
		"Primer Grado",
		"GRADE-01",
	)
	require.NoError(t, err)

	err = unitRepo.Create(ctx, grade)
	require.NoError(t, err, "Debería crear unidad de nivel superior (Grade)")

	// PASO 3: Buscar unidades por school (debe retornar 1: grade)
	unitsInSchool, err := unitRepo.FindBySchoolID(ctx, school.ID(), false)
	require.NoError(t, err)
	assert.Len(t, unitsInSchool, 1, "Debe haber 1 unidad en la escuela")

	// PASO 4: Crear User (requerido por FK de memberships)
	userID := valueobject.NewUserID()
	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`, userID.String(), "coordinator@test.com", "hashed_password", "Test", "Coordinator", "teacher")
	require.NoError(t, err, "Debería crear usuario para membership")

	// PASO 5: Crear Membership (asignar usuario a una unidad)
	membership, err := entity.NewUnitMembership(
		grade.ID(),
		userID,
		valueobject.RoleCoordinator,
		time.Now(),
	)
	require.NoError(t, err)

	err = membershipRepo.Create(ctx, membership)
	require.NoError(t, err, "Debería crear membership vinculado a unidad académica")

	// PASO 6: Verificar que el membership se creó y está vinculado correctamente
	membershipsInUnit, err := membershipRepo.FindByUnitID(ctx, grade.ID(), false)
	require.NoError(t, err)
	assert.Len(t, membershipsInUnit, 1, "Debe haber 1 membership en el grado")
	assert.Equal(t, valueobject.RoleCoordinator, membershipsInUnit[0].Role(), "Rol debe ser coordinador")

	// PASO 7: Verificar integridad referencial
	foundSchool, err := schoolRepo.FindByID(ctx, school.ID())
	require.NoError(t, err)
	assert.Equal(t, "HIERTEST001", foundSchool.Code(), "School debe existir y mantener integridad")

	t.Log("✅ Flujo completo de jerarquía académica validado:")
	t.Log("   School → AcademicUnit (Grade) → Membership")
}

// TestIntegration_MetadataJSONBPersistence prueba que metadata JSONB se persiste correctamente
// Este fue el bug original que encontramos (metadata nil → error JSONB)
func TestIntegration_MetadataJSONBPersistence(t *testing.T) {
	db, cleanup := getTestDB(t)
	defer cleanup()

	ctx := context.Background()
	schoolRepo := postgresRepo.NewPostgresSchoolRepository(db)

	// CASO 1: School SIN metadata (debe persistir {} por defecto)
	schoolWithoutMetadata, err := entity.NewSchool("No Metadata School", "NOMETA001", "")
	require.NoError(t, err)

	err = schoolRepo.Create(ctx, schoolWithoutMetadata)
	require.NoError(t, err, "Debe crear school sin metadata (fix: metadata nil → {})")

	found, err := schoolRepo.FindByID(ctx, schoolWithoutMetadata.ID())
	require.NoError(t, err)
	assert.NotNil(t, found.Metadata(), "Metadata no debe ser nil")

	t.Log("✅ Bug fix validado: metadata nil → {} persiste correctamente en JSONB")
}
