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

// MockSubjectRepository implementa repository.SubjectRepository para testing
type MockSubjectRepository struct {
	mu       sync.RWMutex
	subjects map[uuid.UUID]*entities.Subject
}

// NewMockSubjectRepository crea una nueva instancia de MockSubjectRepository
// Inicializa con un mapa vacío - los datos se crean vía API durante los tests
func NewMockSubjectRepository() repository.SubjectRepository {
	return &MockSubjectRepository{
		subjects: make(map[uuid.UUID]*entities.Subject),
	}
}

// Create crea una nueva materia
func (r *MockSubjectRepository) Create(ctx context.Context, subject *entities.Subject) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.subjects {
		if s.Name == subject.Name {
			return errors.NewConflictError("subject with this name already exists")
		}
	}

	if subject.ID == uuid.Nil {
		subject.ID = uuid.New()
	}

	now := time.Now()
	subject.CreatedAt = now
	subject.UpdatedAt = now

	subjectCopy := *subject
	r.subjects[subject.ID] = &subjectCopy

	return nil
}

// FindByID busca una materia por ID
func (r *MockSubjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Subject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	subject, exists := r.subjects[id]
	if !exists {
		return nil, errors.NewNotFoundError("subject not found")
	}

	subjectCopy := *subject
	return &subjectCopy, nil
}

// Update actualiza una materia existente
func (r *MockSubjectRepository) Update(ctx context.Context, subject *entities.Subject) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.subjects[subject.ID]
	if !exists {
		return errors.NewNotFoundError("subject not found")
	}

	for _, s := range r.subjects {
		if s.ID != subject.ID && s.Name == subject.Name {
			return errors.NewConflictError("subject with this name already exists")
		}
	}

	subject.UpdatedAt = time.Now()
	subject.CreatedAt = existing.CreatedAt

	subjectCopy := *subject
	r.subjects[subject.ID] = &subjectCopy

	return nil
}

// Delete elimina una materia (soft delete)
func (r *MockSubjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	subject, exists := r.subjects[id]
	if !exists {
		return errors.NewNotFoundError("subject not found")
	}

	now := time.Now()
	subject.IsActive = false
	subject.UpdatedAt = now

	return nil
}

// List lista todas las materias activas
func (r *MockSubjectRepository) List(ctx context.Context) ([]*entities.Subject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Subject

	for _, subject := range r.subjects {
		if subject.IsActive {
			subjectCopy := *subject
			result = append(result, &subjectCopy)
		}
	}

	return result, nil
}

// ExistsByName verifica si existe una materia con el nombre dado
func (r *MockSubjectRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, subject := range r.subjects {
		if subject.Name == name {
			return true, nil
		}
	}

	return false, nil
}

// Reset limpia todos los datos (útil para testing)
func (r *MockSubjectRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subjects = make(map[uuid.UUID]*entities.Subject)
}
