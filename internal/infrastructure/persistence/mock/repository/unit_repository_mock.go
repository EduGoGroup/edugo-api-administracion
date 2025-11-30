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

// MockUnitRepository es una implementación en memoria del UnitRepository para testing
type MockUnitRepository struct {
	mu    sync.RWMutex
	units map[uuid.UUID]*entities.Unit
}

// NewMockUnitRepository crea una nueva instancia de MockUnitRepository
// Pre-carga las unidades desde mockData.GetUnits()
func NewMockUnitRepository() repository.UnitRepository {
	units := make(map[uuid.UUID]*entities.Unit)

	// Pre-cargar datos desde mockData
	mockUnits := mockData.GetUnits()
	for id, unit := range mockUnits {
		// Hacer una copia de la unidad para evitar modificaciones externas
		unitCopy := *unit
		units[id] = &unitCopy
	}

	return &MockUnitRepository{
		units: units,
	}
}

// Create crea una nueva unidad en el repositorio
func (r *MockUnitRepository) Create(ctx context.Context, unit *entities.Unit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generar ID si no existe
	if unit.ID == uuid.Nil {
		unit.ID = uuid.New()
	}

	// Validar que el nombre no exista en la misma escuela
	for _, existingUnit := range r.units {
		if existingUnit.SchoolID == unit.SchoolID &&
			existingUnit.Name == unit.Name &&
			existingUnit.IsActive {
			return errors.NewConflictError("unit with this name already exists in the school")
		}
	}

	// Validar que la escuela existe (opcional pero recomendado)
	if unit.SchoolID == uuid.Nil {
		return errors.NewValidationError("school_id is required")
	}

	// Validar que la unidad padre existe si se proporciona
	if unit.ParentUnitID != nil && *unit.ParentUnitID != uuid.Nil {
		parentUnit, exists := r.units[*unit.ParentUnitID]
		if !exists || !parentUnit.IsActive {
			return errors.NewNotFoundError("parent unit not found")
		}
	}

	// Establecer timestamps
	now := time.Now()
	unit.CreatedAt = now
	unit.UpdatedAt = now

	// Guardar una copia de la unidad
	unitCopy := *unit
	r.units[unit.ID] = &unitCopy

	return nil
}

// FindByID busca una unidad por ID
func (r *MockUnitRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Unit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	unit, exists := r.units[id]
	if !exists || !unit.IsActive {
		return nil, errors.NewNotFoundError("unit not found")
	}

	// Retornar una copia para evitar modificaciones externas
	unitCopy := *unit
	return &unitCopy, nil
}

// Update actualiza una unidad existente
func (r *MockUnitRepository) Update(ctx context.Context, unit *entities.Unit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verificar que la unidad existe
	existingUnit, exists := r.units[unit.ID]
	if !exists || !existingUnit.IsActive {
		return errors.NewNotFoundError("unit not found")
	}

	// Validar que el nombre no esté en uso por otra unidad en la misma escuela
	for id, u := range r.units {
		if u.SchoolID == unit.SchoolID &&
			u.Name == unit.Name &&
			id != unit.ID &&
			u.IsActive {
			return errors.NewConflictError("unit with this name already exists in the school")
		}
	}

	// Validar que la unidad padre existe si se proporciona
	if unit.ParentUnitID != nil && *unit.ParentUnitID != uuid.Nil {
		parentUnit, exists := r.units[*unit.ParentUnitID]
		if !exists || !parentUnit.IsActive {
			return errors.NewNotFoundError("parent unit not found")
		}
	}

	// Actualizar timestamp
	unit.UpdatedAt = time.Now()

	// Preservar CreatedAt original
	unit.CreatedAt = existingUnit.CreatedAt

	// Guardar una copia de la unidad actualizada
	unitCopy := *unit
	r.units[unit.ID] = &unitCopy

	return nil
}

// Delete desactiva una unidad (marca IsActive como false)
func (r *MockUnitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	unit, exists := r.units[id]
	if !exists || !unit.IsActive {
		return errors.NewNotFoundError("unit not found")
	}

	// Desactivar la unidad
	unit.IsActive = false
	unit.UpdatedAt = time.Now()

	return nil
}

// List lista unidades activas por escuela
func (r *MockUnitRepository) List(ctx context.Context, schoolID uuid.UUID) ([]*entities.Unit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Unit

	// Filtrar unidades
	for _, unit := range r.units {
		// Incluir solo unidades activas
		if !unit.IsActive {
			continue
		}

		// Filtrar por escuela
		if unit.SchoolID != schoolID {
			continue
		}

		// Agregar copia de la unidad
		unitCopy := *unit
		result = append(result, &unitCopy)
	}

	return result, nil
}

// Reset reinicia el repositorio a su estado inicial (útil para testing)
func (r *MockUnitRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Recargar datos desde mockData
	units := make(map[uuid.UUID]*entities.Unit)
	mockUnits := mockData.GetUnits()
	for id, unit := range mockUnits {
		unitCopy := *unit
		units[id] = &unitCopy
	}

	r.units = units
}
