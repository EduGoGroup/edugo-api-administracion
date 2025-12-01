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

// MockGuardianRepository es una implementación en memoria del GuardianRepository para testing
type MockGuardianRepository struct {
	mu        sync.RWMutex
	relations map[uuid.UUID]*entities.GuardianRelation
}

// NewMockGuardianRepository crea una nueva instancia de MockGuardianRepository
// Inicializa con un mapa vacío - los datos se crean vía API durante los tests
func NewMockGuardianRepository() repository.GuardianRepository {
	return &MockGuardianRepository{
		relations: make(map[uuid.UUID]*entities.GuardianRelation),
	}
}

// Create es un alias de CreateRelation para compatibilidad
func (r *MockGuardianRepository) Create(ctx context.Context, relation *entities.GuardianRelation) error {
	return r.CreateRelation(ctx, relation)
}

// CreateRelation crea una nueva relación guardian-student en el repositorio
func (r *MockGuardianRepository) CreateRelation(ctx context.Context, relation *entities.GuardianRelation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validar que guardian != student
	if relation.GuardianID == relation.StudentID {
		return errors.NewValidationError("guardian cannot be the same as student")
	}

	// Validar que no exista una relación activa con el mismo guardian y student
	for _, existingRelation := range r.relations {
		if existingRelation.GuardianID == relation.GuardianID &&
			existingRelation.StudentID == relation.StudentID &&
			existingRelation.IsActive {
			return errors.NewConflictError("active guardian relation already exists for this guardian and student")
		}
	}

	// Generar ID si no existe
	if relation.ID == uuid.Nil {
		relation.ID = uuid.New()
	}

	// Establecer timestamps
	now := time.Now()
	relation.CreatedAt = now
	relation.UpdatedAt = now

	// Guardar una copia de la relación
	relationCopy := *relation
	r.relations[relation.ID] = &relationCopy

	return nil
}

// FindByID es un alias de FindRelationByID para compatibilidad
func (r *MockGuardianRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error) {
	return r.FindRelationByID(ctx, id)
}

// FindRelationByID busca una relación por ID
func (r *MockGuardianRepository) FindRelationByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	relation, exists := r.relations[id]
	if !exists {
		return nil, errors.NewNotFoundError("guardian relation not found")
	}

	// Retornar una copia para evitar modificaciones externas
	return r.copyRelation(relation), nil
}

// FindByGuardian es un alias de FindRelationsByGuardian para compatibilidad
func (r *MockGuardianRepository) FindByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error) {
	return r.FindRelationsByGuardian(ctx, guardianID)
}

// FindRelationsByGuardian busca todas las relaciones de un guardian
func (r *MockGuardianRepository) FindRelationsByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.GuardianRelation

	for _, relation := range r.relations {
		if relation.GuardianID == guardianID {
			result = append(result, r.copyRelation(relation))
		}
	}

	return result, nil
}

// FindByStudent es un alias de FindRelationsByStudent para compatibilidad
func (r *MockGuardianRepository) FindByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error) {
	return r.FindRelationsByStudent(ctx, studentID)
}

// FindRelationsByStudent busca todas las relaciones de un estudiante
func (r *MockGuardianRepository) FindRelationsByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.GuardianRelation

	for _, relation := range r.relations {
		if relation.StudentID == studentID {
			result = append(result, r.copyRelation(relation))
		}
	}

	return result, nil
}

// Update es un alias de UpdateRelation para compatibilidad
func (r *MockGuardianRepository) Update(ctx context.Context, relation *entities.GuardianRelation) error {
	return r.UpdateRelation(ctx, relation)
}

// UpdateRelation actualiza una relación existente
func (r *MockGuardianRepository) UpdateRelation(ctx context.Context, relation *entities.GuardianRelation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingRelation, exists := r.relations[relation.ID]
	if !exists {
		return errors.NewNotFoundError("guardian relation not found")
	}

	if relation.GuardianID == relation.StudentID {
		return errors.NewValidationError("guardian cannot be the same as student")
	}

	for id, rel := range r.relations {
		if id != relation.ID &&
			rel.GuardianID == relation.GuardianID &&
			rel.StudentID == relation.StudentID &&
			rel.IsActive {
			return errors.NewConflictError("active guardian relation already exists for this guardian and student")
		}
	}

	relation.UpdatedAt = time.Now()
	relation.CreatedAt = existingRelation.CreatedAt
	if relation.CreatedBy == "" {
		relation.CreatedBy = existingRelation.CreatedBy
	}

	relationCopy := *relation
	r.relations[relation.ID] = &relationCopy

	return nil
}

// Delete es un alias de DeleteRelation para compatibilidad
func (r *MockGuardianRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteRelation(ctx, id)
}

// DeleteRelation elimina una relación
func (r *MockGuardianRepository) DeleteRelation(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.relations[id]
	if !exists {
		return errors.NewNotFoundError("guardian relation not found")
	}

	delete(r.relations, id)
	return nil
}

// ExistsActiveRelation verifica si existe una relación activa entre un guardian y un estudiante
func (r *MockGuardianRepository) ExistsActiveRelation(ctx context.Context, guardianID, studentID uuid.UUID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, relation := range r.relations {
		if relation.GuardianID == guardianID &&
			relation.StudentID == studentID &&
			relation.IsActive {
			return true, nil
		}
	}

	return false, nil
}

// Reset reinicia el repositorio a un estado vacío (útil para testing)
func (r *MockGuardianRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.relations = make(map[uuid.UUID]*entities.GuardianRelation)
}

// copyRelation crea una copia profunda de una relación
func (r *MockGuardianRepository) copyRelation(relation *entities.GuardianRelation) *entities.GuardianRelation {
	relationCopy := *relation
	return &relationCopy
}
