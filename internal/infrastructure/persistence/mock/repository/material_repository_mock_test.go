package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestMockMaterialRepository_Delete(t *testing.T) {
	repo := NewMockMaterialRepository()
	ctx := context.Background()

	// Material existente
	materialID := uuid.MustParse("f1aabc99-9c0b-4ef8-bb6d-6bb9bd380f11")

	// Verificar que existe antes de eliminar
	exists, err := repo.Exists(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al verificar existencia: %v", err)
	}
	if !exists {
		t.Fatal("Material debería existir antes de eliminar")
	}

	// Eliminar
	err = repo.Delete(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al eliminar material: %v", err)
	}

	// Verificar que ya no existe
	exists, err = repo.Exists(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al verificar existencia después de eliminar: %v", err)
	}
	if exists {
		t.Fatal("Material no debería existir después de eliminar")
	}
}

func TestMockMaterialRepository_Delete_NotFound(t *testing.T) {
	repo := NewMockMaterialRepository()
	ctx := context.Background()

	notFoundID := uuid.New()

	err := repo.Delete(ctx, notFoundID)
	if err == nil {
		t.Fatal("Se esperaba un error de not found")
	}
}

func TestMockMaterialRepository_Delete_AlreadyDeleted(t *testing.T) {
	repo := NewMockMaterialRepository()
	ctx := context.Background()

	materialID := uuid.MustParse("f1aabc99-9c0b-4ef8-bb6d-6bb9bd380f11")

	// Eliminar una primera vez
	err := repo.Delete(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al eliminar material la primera vez: %v", err)
	}

	// Intentar eliminar nuevamente
	err = repo.Delete(ctx, materialID)
	if err == nil {
		t.Fatal("Se esperaba un error al intentar eliminar un material ya eliminado")
	}
}

func TestMockMaterialRepository_Exists(t *testing.T) {
	repo := NewMockMaterialRepository()
	ctx := context.Background()

	// Material existente
	materialID := uuid.MustParse("f1aabc99-9c0b-4ef8-bb6d-6bb9bd380f11")

	exists, err := repo.Exists(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al verificar existencia: %v", err)
	}

	if !exists {
		t.Fatal("Material debería existir")
	}
}

func TestMockMaterialRepository_Exists_NotFound(t *testing.T) {
	repo := NewMockMaterialRepository()
	ctx := context.Background()

	notFoundID := uuid.New()

	exists, err := repo.Exists(ctx, notFoundID)
	if err != nil {
		t.Fatalf("Error al verificar existencia: %v", err)
	}

	if exists {
		t.Fatal("Material no debería existir")
	}
}

func TestMockMaterialRepository_Exists_Deleted(t *testing.T) {
	repo := NewMockMaterialRepository()
	ctx := context.Background()

	materialID := uuid.MustParse("f1aabc99-9c0b-4ef8-bb6d-6bb9bd380f11")

	// Eliminar el material
	err := repo.Delete(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al eliminar material: %v", err)
	}

	// Verificar que no existe después de eliminar
	exists, err := repo.Exists(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al verificar existencia: %v", err)
	}

	if exists {
		t.Fatal("Material eliminado no debería existir")
	}
}

func TestMockMaterialRepository_Reset(t *testing.T) {
	repo := NewMockMaterialRepository()
	mockRepo, ok := repo.(*MockMaterialRepository)
	if !ok {
		t.Fatal("No se pudo hacer cast a *MockMaterialRepository")
	}
	ctx := context.Background()

	materialID := uuid.MustParse("f1aabc99-9c0b-4ef8-bb6d-6bb9bd380f11")

	// Eliminar un material
	err := repo.Delete(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al eliminar material: %v", err)
	}

	// Verificar que está eliminado
	exists, err := repo.Exists(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al verificar existencia: %v", err)
	}
	if exists {
		t.Fatal("Material debería estar eliminado")
	}

	// Resetear
	mockRepo.Reset()

	// Verificar que el material vuelve a existir
	exists, err = repo.Exists(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al verificar existencia después de reset: %v", err)
	}

	if !exists {
		t.Fatal("Material debería existir después de reset")
	}
}

func TestMockMaterialRepository_ThreadSafety(t *testing.T) {
	repo := NewMockMaterialRepository()
	ctx := context.Background()

	// Este test verifica que el RWMutex está presente y funciona
	// al ejecutar operaciones concurrentes
	done := make(chan bool)

	materialID := uuid.MustParse("f1aabc99-9c0b-4ef8-bb6d-6bb9bd380f11")

	// Ejecutar múltiples operaciones de lectura concurrentes
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = repo.Exists(ctx, materialID)
			done <- true
		}()
	}

	// Ejecutar una operación de escritura concurrente
	for i := 0; i < 2; i++ {
		go func() {
			mockRepo := repo.(*MockMaterialRepository)
			mockRepo.Reset()
			done <- true
		}()
	}

	// Esperar a que todas las goroutines terminen
	for i := 0; i < 12; i++ {
		<-done
	}
}

func TestMockMaterialRepository_MultipleExistsChecks(t *testing.T) {
	repo := NewMockMaterialRepository()
	ctx := context.Background()

	materialIDs := []uuid.UUID{
		uuid.MustParse("f1aabc99-9c0b-4ef8-bb6d-6bb9bd380f11"),
		uuid.MustParse("f2aabc99-9c0b-4ef8-bb6d-6bb9bd380f22"),
		uuid.MustParse("f3aabc99-9c0b-4ef8-bb6d-6bb9bd380f33"),
		uuid.MustParse("f4aabc99-9c0b-4ef8-bb6d-6bb9bd380f44"),
		uuid.MustParse("f5aabc99-9c0b-4ef8-bb6d-6bb9bd380f55"),
		uuid.MustParse("f6aabc99-9c0b-4ef8-bb6d-6bb9bd380f66"),
	}

	// Verificar que todos los materiales existen
	for _, materialID := range materialIDs {
		exists, err := repo.Exists(ctx, materialID)
		if err != nil {
			t.Fatalf("Error al verificar existencia de material %s: %v", materialID, err)
		}
		if !exists {
			t.Fatalf("Material %s debería existir", materialID)
		}
	}
}

func TestMockMaterialRepository_DeleteAndExistsConsistency(t *testing.T) {
	repo := NewMockMaterialRepository()
	ctx := context.Background()

	materialID := uuid.MustParse("f1aabc99-9c0b-4ef8-bb6d-6bb9bd380f11")

	// Estado inicial: debe existir
	exists, err := repo.Exists(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al verificar existencia inicial: %v", err)
	}
	if !exists {
		t.Fatal("Material debería existir inicialmente")
	}

	// Eliminar
	err = repo.Delete(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al eliminar material: %v", err)
	}

	// Verificar consistencia: Exists debe retornar false
	exists, err = repo.Exists(ctx, materialID)
	if err != nil {
		t.Fatalf("Error al verificar existencia después de eliminar: %v", err)
	}
	if exists {
		t.Fatal("Material no debería existir después de eliminación (inconsistencia)")
	}

	// Intentar eliminar nuevamente debe fallar
	err = repo.Delete(ctx, materialID)
	if err == nil {
		t.Fatal("Eliminar un material ya eliminado debería fallar")
	}
}
