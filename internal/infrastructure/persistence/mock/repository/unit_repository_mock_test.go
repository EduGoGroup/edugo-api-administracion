package repository

import (
	"context"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

func TestMockUnitRepository_Create(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	schoolID := uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	newUnit := &entities.Unit{
		SchoolID:    schoolID,
		Name:        "Nueva Unidad",
		Description: stringPtr("Descripción de la nueva unidad"),
		IsActive:    true,
	}

	err := repo.Create(ctx, newUnit)
	if err != nil {
		t.Fatalf("Error al crear unidad: %v", err)
	}

	if newUnit.ID == uuid.Nil {
		t.Fatal("El ID no fue generado")
	}

	// Verificar que se puede recuperar
	retrieved, err := repo.FindByID(ctx, newUnit.ID)
	if err != nil {
		t.Fatalf("Error al recuperar unidad: %v", err)
	}

	if retrieved.Name != newUnit.Name {
		t.Errorf("Nombre esperado %s, obtuvo %s", newUnit.Name, retrieved.Name)
	}
}

func TestMockUnitRepository_FindByID(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	unitID := uuid.MustParse("f1000000-0000-0000-0000-000000000011")
	unit, err := repo.FindByID(ctx, unitID)

	if err != nil {
		t.Fatalf("Error al buscar unidad: %v", err)
	}

	if unit.Name != "Departamento de Matemáticas" {
		t.Errorf("Nombre esperado 'Departamento de Matemáticas', obtuvo '%s'", unit.Name)
	}
}

func TestMockUnitRepository_FindByID_NotFound(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	notFoundID := uuid.New()
	_, err := repo.FindByID(ctx, notFoundID)

	if err == nil {
		t.Fatal("Se esperaba un error de not found")
	}
}

func TestMockUnitRepository_Update(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	unitID := uuid.MustParse("f1000000-0000-0000-0000-000000000011")
	unit, err := repo.FindByID(ctx, unitID)
	if err != nil {
		t.Fatalf("Error al buscar unidad: %v", err)
	}

	originalName := unit.Name
	unit.Name = "Departamento Actualizado"

	err = repo.Update(ctx, unit)
	if err != nil {
		t.Fatalf("Error al actualizar unidad: %v", err)
	}

	// Verificar que se actualizó
	updated, err := repo.FindByID(ctx, unitID)
	if err != nil {
		t.Fatalf("Error al recuperar unidad actualizada: %v", err)
	}

	if updated.Name != "Departamento Actualizado" {
		t.Errorf("Nombre no fue actualizado. Esperado 'Departamento Actualizado', obtuvo '%s'", updated.Name)
	}

	// Verificar que CreatedAt se preservó
	if !updated.CreatedAt.Equal(unit.CreatedAt) {
		t.Error("CreatedAt fue modificado")
	}
}

func TestMockUnitRepository_Delete(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	unitID := uuid.MustParse("f1000000-0000-0000-0000-000000000011")

	err := repo.Delete(ctx, unitID)
	if err != nil {
		t.Fatalf("Error al eliminar unidad: %v", err)
	}

	// Verificar que no se puede encontrar después de eliminar
	_, err = repo.FindByID(ctx, unitID)
	if err == nil {
		t.Fatal("Se esperaba un error de not found después de eliminar")
	}
}

func TestMockUnitRepository_List(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	schoolID := uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	units, err := repo.List(ctx, schoolID)

	if err != nil {
		t.Fatalf("Error al listar unidades: %v", err)
	}

	if len(units) != 2 {
		t.Errorf("Se esperaban 2 unidades, se obtuvieron %d", len(units))
	}

	// Verificar que todas pertenecen a la escuela
	for _, unit := range units {
		if unit.SchoolID != schoolID {
			t.Errorf("Unidad con ID %s no pertenece a la escuela", unit.ID)
		}
	}
}

func TestMockUnitRepository_List_Empty(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	emptySchoolID := uuid.New()
	units, err := repo.List(ctx, emptySchoolID)

	if err != nil {
		t.Fatalf("Error al listar unidades: %v", err)
	}

	if len(units) != 0 {
		t.Errorf("Se esperaban 0 unidades, se obtuvieron %d", len(units))
	}
}

func TestMockUnitRepository_Reset(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	// Crear una nueva unidad
	schoolID := uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	newUnit := &entities.Unit{
		SchoolID: schoolID,
		Name:     "Unidad Temporal",
		IsActive: true,
	}

	err := repo.Create(ctx, newUnit)
	if err != nil {
		t.Fatalf("Error al crear unidad: %v", err)
	}

	// Resetear
	repo.Reset()

	// Verificar que la unidad creada ya no existe
	_, err = repo.FindByID(ctx, newUnit.ID)
	if err == nil {
		t.Fatal("La unidad temporal debería haber sido eliminada con Reset()")
	}

	// Verificar que las unidades originales aún existen
	unitID := uuid.MustParse("f1000000-0000-0000-0000-000000000011")
	unit, err := repo.FindByID(ctx, unitID)
	if err != nil {
		t.Fatalf("Error al recuperar unidad después de Reset: %v", err)
	}

	if unit.Name != "Departamento de Matemáticas" {
		t.Error("Las unidades originales no se restauraron correctamente")
	}
}

func TestMockUnitRepository_CopyReturn(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	unitID := uuid.MustParse("f1000000-0000-0000-0000-000000000011")
	unit1, _ := repo.FindByID(ctx, unitID)
	unit2, _ := repo.FindByID(ctx, unitID)

	// Modificar una copia
	unit1.Name = "Nombre Modificado"

	// Verificar que la otra copia no se vio afectada
	if unit2.Name != "Departamento de Matemáticas" {
		t.Error("Las copias no son independientes")
	}
}

func TestMockUnitRepository_ThreadSafety(t *testing.T) {
	repo := NewMockUnitRepository()
	ctx := context.Background()

	// Este test verifica que el RWMutex está presente y funciona
	// al ejecutar operaciones concurrentes
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			unitID := uuid.MustParse("f1000000-0000-0000-0000-000000000011")
			_, _ = repo.FindByID(ctx, unitID)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// Helper para crear string pointer
func stringPtr(s string) *string {
	return &s
}
