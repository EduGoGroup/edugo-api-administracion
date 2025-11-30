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

// MockSchoolRepository implementa repository.SchoolRepository para testing
type MockSchoolRepository struct {
	mu      sync.RWMutex
	schools map[uuid.UUID]*entities.School
}

// NewMockSchoolRepository crea una nueva instancia de MockSchoolRepository
func NewMockSchoolRepository() repository.SchoolRepository {
	repo := &MockSchoolRepository{
		schools: make(map[uuid.UUID]*entities.School),
	}

	// Pre-cargar datos desde mockData
	for _, school := range mockData.GetSchools() {
		schoolCopy := *school
		repo.schools[school.ID] = &schoolCopy
	}

	return repo
}

// Create crea una nueva escuela
func (r *MockSchoolRepository) Create(ctx context.Context, school *entities.School) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validar duplicado por Code
	for _, s := range r.schools {
		if s.DeletedAt == nil && s.Code == school.Code {
			return errors.NewConflictError("school with this code already exists")
		}
	}

	// Validar duplicado por Name
	for _, s := range r.schools {
		if s.DeletedAt == nil && s.Name == school.Name {
			return errors.NewConflictError("school with this name already exists")
		}
	}

	// Generar ID si no existe
	if school.ID == uuid.Nil {
		school.ID = uuid.New()
	}

	// Establecer timestamps
	now := time.Now()
	school.CreatedAt = now
	school.UpdatedAt = now

	// Almacenar copia
	schoolCopy := *school
	r.schools[school.ID] = &schoolCopy

	return nil
}

// FindByID busca una escuela por ID
func (r *MockSchoolRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	school, exists := r.schools[id]
	if !exists || school.DeletedAt != nil {
		return nil, errors.NewNotFoundError("school not found")
	}

	// Retornar copia
	schoolCopy := *school
	return &schoolCopy, nil
}

// FindByCode busca una escuela por código
func (r *MockSchoolRepository) FindByCode(ctx context.Context, code string) (*entities.School, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, school := range r.schools {
		if school.DeletedAt == nil && school.Code == code {
			// Retornar copia
			schoolCopy := *school
			return &schoolCopy, nil
		}
	}

	return nil, errors.NewNotFoundError("school not found")
}

// FindByName busca una escuela por nombre
func (r *MockSchoolRepository) FindByName(ctx context.Context, name string) (*entities.School, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, school := range r.schools {
		if school.DeletedAt == nil && school.Name == name {
			// Retornar copia
			schoolCopy := *school
			return &schoolCopy, nil
		}
	}

	return nil, errors.NewNotFoundError("school not found")
}

// Update actualiza una escuela existente
func (r *MockSchoolRepository) Update(ctx context.Context, school *entities.School) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verificar que la escuela existe
	existing, exists := r.schools[school.ID]
	if !exists || existing.DeletedAt != nil {
		return errors.NewNotFoundError("school not found")
	}

	// Validar duplicado por Code (excepto el mismo registro)
	for _, s := range r.schools {
		if s.DeletedAt == nil && s.ID != school.ID && s.Code == school.Code {
			return errors.NewConflictError("school with this code already exists")
		}
	}

	// Validar duplicado por Name (excepto el mismo registro)
	for _, s := range r.schools {
		if s.DeletedAt == nil && s.ID != school.ID && s.Name == school.Name {
			return errors.NewConflictError("school with this name already exists")
		}
	}

	// Actualizar timestamp
	school.UpdatedAt = time.Now()

	// Preservar CreatedAt original
	school.CreatedAt = existing.CreatedAt

	// Almacenar copia actualizada
	schoolCopy := *school
	r.schools[school.ID] = &schoolCopy

	return nil
}

// Delete elimina una escuela (soft delete)
func (r *MockSchoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	school, exists := r.schools[id]
	if !exists || school.DeletedAt != nil {
		return errors.NewNotFoundError("school not found")
	}

	// Soft delete
	now := time.Now()
	school.DeletedAt = &now
	school.UpdatedAt = now

	return nil
}

// List lista todas las escuelas activas con filtros opcionales
func (r *MockSchoolRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entities.School, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.School

	// Filtrar escuelas
	for _, school := range r.schools {
		// Excluir escuelas eliminadas
		if school.DeletedAt != nil {
			continue
		}

		// Agregar copia de la escuela
		schoolCopy := *school
		result = append(result, &schoolCopy)
	}

	// Aplicar offset
	if filters.Offset > 0 {
		if filters.Offset >= len(result) {
			return []*entities.School{}, nil
		}
		result = result[filters.Offset:]
	}

	// Aplicar limit
	if filters.Limit > 0 && filters.Limit < len(result) {
		result = result[:filters.Limit]
	}

	return result, nil
}

// ExistsByName verifica si existe una escuela con el nombre dado
func (r *MockSchoolRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, school := range r.schools {
		if school.DeletedAt == nil && school.Name == name {
			return true, nil
		}
	}

	return false, nil
}

// ExistsByCode verifica si existe una escuela con el código dado
func (r *MockSchoolRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, school := range r.schools {
		if school.DeletedAt == nil && school.Code == code {
			return true, nil
		}
	}

	return false, nil
}

// Reset limpia todos los datos y recarga los datos mock (útil para testing)
func (r *MockSchoolRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Limpiar mapa
	r.schools = make(map[uuid.UUID]*entities.School)

	// Recargar datos desde mockData
	for _, school := range mockData.GetSchools() {
		schoolCopy := *school
		r.schools[school.ID] = &schoolCopy
	}
}
