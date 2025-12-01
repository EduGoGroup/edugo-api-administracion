package repository

import (
	"context"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
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
// Inicializa con un mapa vacío - los datos se crean vía API durante los tests
func NewMockUnitRepository() repository.UnitRepository {
	return &MockUnitRepository{
		units: make(map[uuid.UUID]*entities.Unit),
	}
}

// Create crea una nueva unidad en el repositorio
func (r *MockUnitRepository) Create(ctx context.Context, unit *entities.Unit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if unit.ID == uuid.Nil {
		unit.ID = uuid.New()
	}

	for _, existingUnit := range r.units {
		if existingUnit.SchoolID == unit.SchoolID &&
			existingUnit.Name == unit.Name &&
			existingUnit.IsActive {
			return errors.NewConflictError("unit with this name already exists in the school")
		}
	}

	if unit.SchoolID == uuid.Nil {
		return errors.NewValidationError("school_id is required")
	}

	if unit.ParentUnitID != nil && *unit.ParentUnitID != uuid.Nil {
		parentUnit, exists := r.units[*unit.ParentUnitID]
		if !exists || !parentUnit.IsActive {
			return errors.NewNotFoundError("parent unit not found")
		}
	}

	now := time.Now()
	unit.CreatedAt = now
	unit.UpdatedAt = now

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

	unitCopy := *unit
	return &unitCopy, nil
}

// Update actualiza una unidad existente
func (r *MockUnitRepository) Update(ctx context.Context, unit *entities.Unit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingUnit, exists := r.units[unit.ID]
	if !exists || !existingUnit.IsActive {
		return errors.NewNotFoundError("unit not found")
	}

	for id, u := range r.units {
		if u.SchoolID == unit.SchoolID &&
			u.Name == unit.Name &&
			id != unit.ID &&
			u.IsActive {
			return errors.NewConflictError("unit with this name already exists in the school")
		}
	}

	if unit.ParentUnitID != nil && *unit.ParentUnitID != uuid.Nil {
		parentUnit, exists := r.units[*unit.ParentUnitID]
		if !exists || !parentUnit.IsActive {
			return errors.NewNotFoundError("parent unit not found")
		}
	}

	unit.UpdatedAt = time.Now()
	unit.CreatedAt = existingUnit.CreatedAt

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

	unit.IsActive = false
	unit.UpdatedAt = time.Now()

	return nil
}

// List lista unidades activas por escuela
func (r *MockUnitRepository) List(ctx context.Context, schoolID uuid.UUID) ([]*entities.Unit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Unit

	for _, unit := range r.units {
		if !unit.IsActive {
			continue
		}
		if unit.SchoolID != schoolID {
			continue
		}
		unitCopy := *unit
		result = append(result, &unitCopy)
	}

	return result, nil
}

// Reset reinicia el repositorio a un estado vacío (útil para testing)
func (r *MockUnitRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.units = make(map[uuid.UUID]*entities.Unit)
}
