//go:build integration

package integration

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/postgres/repository"
)

// TestAcademicUnitRepository_FindChildren verifica que se pueden obtener los hijos directos de una unidad
func TestAcademicUnitRepository_FindChildren(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPostgresAcademicUnitRepository(db)

	// Setup: Crear jerarquía
	// Escuela -> Grado -> Sección A
	//                 -> Sección B

	school := createTestSchool(t, db, "Test School", "TS001")

	gradeType, _ := valueobject.NewUnitType("grade")
	grade, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 1", "G1")
	require.NoError(t, repo.Create(ctx, grade))

	sectionType, _ := valueobject.NewUnitType("section")
	sectionA, _ := entity.NewAcademicUnit(school.ID(), sectionType, "Section A", "G1-A")
	sectionA.SetParentID(grade.ID())
	require.NoError(t, repo.Create(ctx, sectionA))

	sectionB, _ := entity.NewAcademicUnit(school.ID(), sectionType, "Section B", "G1-B")
	sectionB.SetParentID(grade.ID())
	require.NoError(t, repo.Create(ctx, sectionB))

	// Test: Obtener hijos directos del grado
	children, err := repo.FindChildren(ctx, grade.ID())
	require.NoError(t, err)
	assert.Len(t, children, 2, "Grade should have 2 children")

	// Verificar que son las secciones correctas
	childNames := []string{children[0].DisplayName(), children[1].DisplayName()}
	assert.Contains(t, childNames, "Section A")
	assert.Contains(t, childNames, "Section B")
}

// TestAcademicUnitRepository_FindDescendants verifica que se pueden obtener TODOS los descendientes
func TestAcademicUnitRepository_FindDescendants(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPostgresAcademicUnitRepository(db)

	// Setup: Crear jerarquía profunda
	// Escuela -> Grado -> Sección -> Grupo de estudio

	school := createTestSchool(t, db, "Test School", "TS001")

	gradeType, _ := valueobject.NewUnitType("grade")
	grade, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 1", "G1")
	require.NoError(t, repo.Create(ctx, grade))

	sectionType, _ := valueobject.NewUnitType("section")
	section, _ := entity.NewAcademicUnit(school.ID(), sectionType, "Section A", "G1-A")
	section.SetParentID(grade.ID())
	require.NoError(t, repo.Create(ctx, section))

	clubType, _ := valueobject.NewUnitType("club")
	club, _ := entity.NewAcademicUnit(school.ID(), clubType, "Math Club", "G1-A-MC")
	club.SetParentID(section.ID())
	require.NoError(t, repo.Create(ctx, club))

	// Test: Obtener TODOS los descendientes del grado (debe incluir sección Y club)
	descendants, err := repo.FindDescendants(ctx, grade.ID())
	require.NoError(t, err)
	assert.Len(t, descendants, 2, "Grade should have 2 descendants (section + club)")

	// Verificar orden jerárquico (ordenado por path ltree)
	assert.Equal(t, "Section A", descendants[0].DisplayName())
	assert.Equal(t, "Math Club", descendants[1].DisplayName())
}

// TestAcademicUnitRepository_FindAncestors verifica que se pueden obtener TODOS los ancestros
func TestAcademicUnitRepository_FindAncestors(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPostgresAcademicUnitRepository(db)

	// Setup: Crear jerarquía de 3 niveles
	// Escuela -> Grado -> Sección -> Club

	school := createTestSchool(t, db, "Test School", "TS001")

	gradeType, _ := valueobject.NewUnitType("grade")
	grade, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 1", "G1")
	require.NoError(t, repo.Create(ctx, grade))

	sectionType, _ := valueobject.NewUnitType("section")
	section, _ := entity.NewAcademicUnit(school.ID(), sectionType, "Section A", "G1-A")
	section.SetParentID(grade.ID())
	require.NoError(t, repo.Create(ctx, section))

	clubType, _ := valueobject.NewUnitType("club")
	club, _ := entity.NewAcademicUnit(school.ID(), clubType, "Math Club", "G1-A-MC")
	club.SetParentID(section.ID())
	require.NoError(t, repo.Create(ctx, club))

	// Test: Obtener TODOS los ancestros del club (debe incluir sección Y grado)
	ancestors, err := repo.FindAncestors(ctx, club.ID())
	require.NoError(t, err)
	assert.Len(t, ancestors, 2, "Club should have 2 ancestors (grade + section)")

	// Verificar orden jerárquico (de raíz a hoja)
	assert.Equal(t, "Grade 1", ancestors[0].DisplayName())
	assert.Equal(t, "Section A", ancestors[1].DisplayName())
}

// TestAcademicUnitRepository_FindByPath verifica búsqueda por path ltree
func TestAcademicUnitRepository_FindByPath(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPostgresAcademicUnitRepository(db)

	// Setup: Crear unidad
	school := createTestSchool(t, db, "Test School", "TS001")

	gradeType, _ := valueobject.NewUnitType("grade")
	grade, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 1", "G1")
	require.NoError(t, repo.Create(ctx, grade))

	// El path debería ser solo el UUID del grado (es raíz), con guiones reemplazados por guiones bajos
	expectedPath := strings.ReplaceAll(grade.ID().String(), "-", "_")

	// Test: Buscar por path
	found, err := repo.FindByPath(ctx, expectedPath)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, grade.ID(), found.ID())
	assert.Equal(t, "Grade 1", found.DisplayName())
}

// TestAcademicUnitRepository_FindBySchoolIDAndDepth verifica búsqueda por profundidad
func TestAcademicUnitRepository_FindBySchoolIDAndDepth(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPostgresAcademicUnitRepository(db)

	// Setup: Crear jerarquía de múltiples niveles
	school := createTestSchool(t, db, "Test School", "TS001")

	// Nivel 1 (depth=1): 2 grados raíz
	gradeType, _ := valueobject.NewUnitType("grade")
	grade1, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 1", "G1")
	require.NoError(t, repo.Create(ctx, grade1))

	grade2, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 2", "G2")
	require.NoError(t, repo.Create(ctx, grade2))

	// Nivel 2 (depth=2): 2 secciones bajo Grade 1
	sectionType, _ := valueobject.NewUnitType("section")
	sectionA, _ := entity.NewAcademicUnit(school.ID(), sectionType, "Section A", "G1-A")
	sectionA.SetParentID(grade1.ID())
	require.NoError(t, repo.Create(ctx, sectionA))

	sectionB, _ := entity.NewAcademicUnit(school.ID(), sectionType, "Section B", "G1-B")
	sectionB.SetParentID(grade1.ID())
	require.NoError(t, repo.Create(ctx, sectionB))

	// Test: Buscar unidades de profundidad 1 (solo grados raíz)
	depth1Units, err := repo.FindBySchoolIDAndDepth(ctx, school.ID(), 1)
	require.NoError(t, err)
	assert.Len(t, depth1Units, 2, "Should find 2 units at depth 1")

	// Test: Buscar unidades de profundidad 2 (solo secciones)
	depth2Units, err := repo.FindBySchoolIDAndDepth(ctx, school.ID(), 2)
	require.NoError(t, err)
	assert.Len(t, depth2Units, 2, "Should find 2 units at depth 2")
}

// TestAcademicUnitRepository_MoveSubtree verifica mover un subárbol completo
func TestAcademicUnitRepository_MoveSubtree(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPostgresAcademicUnitRepository(db)

	// Setup: Crear dos jerarquías independientes
	school := createTestSchool(t, db, "Test School", "TS001")

	gradeType, _ := valueobject.NewUnitType("grade")
	grade1, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 1", "G1")
	require.NoError(t, repo.Create(ctx, grade1))

	grade2, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 2", "G2")
	require.NoError(t, repo.Create(ctx, grade2))

	sectionType, _ := valueobject.NewUnitType("section")
	section, _ := entity.NewAcademicUnit(school.ID(), sectionType, "Section A", "G1-A")
	section.SetParentID(grade1.ID())
	require.NoError(t, repo.Create(ctx, section))

	clubType, _ := valueobject.NewUnitType("club")
	club, _ := entity.NewAcademicUnit(school.ID(), clubType, "Math Club", "G1-A-MC")
	club.SetParentID(section.ID())
	require.NoError(t, repo.Create(ctx, club))

	// Test: Mover la sección (con su club) de Grade 1 a Grade 2
	grade2ID := grade2.ID()
	err := repo.MoveSubtree(ctx, section.ID(), &grade2ID)
	require.NoError(t, err)

	// Verificar que la sección ahora está bajo Grade 2
	updatedSection, err := repo.FindByID(ctx, section.ID(), false)
	require.NoError(t, err)
	assert.NotNil(t, updatedSection.ParentUnitID())
	assert.Equal(t, grade2.ID(), *updatedSection.ParentUnitID())

	// Verificar que el club sigue siendo descendiente de la sección
	descendants, err := repo.FindDescendants(ctx, section.ID())
	require.NoError(t, err)
	assert.Len(t, descendants, 1, "Section should still have its club as descendant")
	assert.Equal(t, club.ID(), descendants[0].ID())

	// Verificar que Grade 2 ahora tiene 2 descendientes (section + club)
	grade2Descendants, err := repo.FindDescendants(ctx, grade2.ID())
	require.NoError(t, err)
	assert.Len(t, grade2Descendants, 2, "Grade 2 should have 2 descendants after move")
}

// TestAcademicUnitRepository_MoveSubtreeToRoot verifica mover un subárbol a raíz
func TestAcademicUnitRepository_MoveSubtreeToRoot(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPostgresAcademicUnitRepository(db)

	// Setup: Crear jerarquía
	school := createTestSchool(t, db, "Test School", "TS001")

	gradeType, _ := valueobject.NewUnitType("grade")
	grade, _ := entity.NewAcademicUnit(school.ID(), gradeType, "Grade 1", "G1")
	require.NoError(t, repo.Create(ctx, grade))

	sectionType, _ := valueobject.NewUnitType("section")
	section, _ := entity.NewAcademicUnit(school.ID(), sectionType, "Section A", "G1-A")
	section.SetParentID(grade.ID())
	require.NoError(t, repo.Create(ctx, section))

	// Test: Convertir la sección en unidad raíz (sin padre)
	err := repo.MoveSubtree(ctx, section.ID(), nil)
	require.NoError(t, err)

	// Verificar que la sección ahora es raíz
	updatedSection, err := repo.FindByID(ctx, section.ID(), false)
	require.NoError(t, err)
	assert.Nil(t, updatedSection.ParentUnitID(), "Section should be root (no parent)")

	// Verificar que Grade 1 ya no tiene hijos
	children, err := repo.FindChildren(ctx, grade.ID())
	require.NoError(t, err)
	assert.Len(t, children, 0, "Grade should have no children after move to root")
}

// TestAcademicUnitRepository_LtreePerformance benchmark de queries ltree vs recursivos
func TestAcademicUnitRepository_LtreePerformance(t *testing.T) {
	// TODO_FASE2: Implementar benchmark comparativo
	// Este test debería comparar:
	// 1. FindDescendants (ltree) vs GetHierarchyPath (CTE recursivo)
	// 2. Medir tiempo de ejecución con jerarquías profundas (10+ niveles)
	// 3. Validar que ltree es significativamente más rápido (>50%)
	//
	// Ejemplo de estructura:
	// - Crear árbol de 100+ unidades con 5-6 niveles de profundidad
	// - Ejecutar FindDescendants en nodo raíz (debe retornar ~100 unidades)
	// - Comparar tiempo vs CTE recursivo
	//
	// Expectativa: ltree debería ser 2-5x más rápido para jerarquías profundas
}

// =====================================================
// Helper functions
// =====================================================

// createTestSchool crea una escuela de prueba en la base de datos
func createTestSchool(t *testing.T, db *sql.DB, name, code string) *entity.School {
	ctx := context.Background()
	schoolID := valueobject.NewSchoolID()

	// Hacer el código único agregando parte del UUID para evitar duplicados entre tests
	uniqueCode := code + "-" + schoolID.String()[:8]

	// Insertar escuela en la base de datos
	query := `INSERT INTO schools (id, name, code, address, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := db.ExecContext(ctx, query, schoolID.String(), name, uniqueCode, "Test Address")
	require.NoError(t, err)

	// Reconstruir entidad usando ReconstructSchool
	now := time.Now()
	school := entity.ReconstructSchool(
		schoolID,
		name,
		uniqueCode,
		"Test Address",
		nil, // contactEmail
		"",  // contactPhone
		nil, // metadata
		now,
		now,
	)

	return school
}
