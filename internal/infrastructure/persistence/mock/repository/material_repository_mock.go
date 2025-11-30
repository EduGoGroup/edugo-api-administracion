package repository

import (
	"context"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	mockData "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/mock/data"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/google/uuid"
)

// MockMaterialRepository implementa repository.MaterialRepository para testing
type MockMaterialRepository struct {
	mu        sync.RWMutex
	materials map[uuid.UUID]*entities.Material
}

// NewMockMaterialRepository crea una nueva instancia de MockMaterialRepository
// Pre-carga los materiales desde mockData.GetMaterials()
func NewMockMaterialRepository() repository.MaterialRepository {
	repo := &MockMaterialRepository{
		materials: make(map[uuid.UUID]*entities.Material),
	}

	// Pre-cargar datos desde mockData
	for _, material := range mockData.GetMaterials() {
		materialCopy := *material
		repo.materials[material.ID] = &materialCopy
	}

	return repo
}

// Delete elimina un material (soft delete - marca como DeletedAt)
func (r *MockMaterialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	material, exists := r.materials[id]
	if !exists || material.DeletedAt != nil {
		return errors.NewNotFoundError("material not found")
	}

	// Soft delete
	now := time.Now()
	material.DeletedAt = &now
	material.UpdatedAt = now

	return nil
}

// Exists verifica si un material existe (no eliminado)
func (r *MockMaterialRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	material, exists := r.materials[id]
	if !exists || material.DeletedAt != nil {
		return false, nil
	}

	return true, nil
}

// Reset limpia todos los datos y recarga los datos mock (útil para testing)
func (r *MockMaterialRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Limpiar mapa
	r.materials = make(map[uuid.UUID]*entities.Material)

	// Recargar datos desde mockData
	for _, material := range mockData.GetMaterials() {
		materialCopy := *material
		r.materials[material.ID] = &materialCopy
	}
}
