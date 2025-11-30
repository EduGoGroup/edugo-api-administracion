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

// MockSubjectRepository implementa repository.SubjectRepository para testing
type MockSubjectRepository struct {
	mu       sync.RWMutex
	subjects map[uuid.UUID]*entities.Subject
}

// NewMockSubjectRepository crea una nueva instancia de MockSubjectRepository
func NewMockSubjectRepository() repository.SubjectRepository {
	repo := &MockSubjectRepository{
		subjects: make(map[uuid.UUID]*entities.Subject),
	}

	// Pre-cargar datos desde mockData
	for _, subject := range mockData.GetSubjects() {
		subjectCopy := *subject
		repo.subjects[subject.ID] = &subjectCopy
	}

	return repo
}

// Create crea una nueva materia
func (r *MockSubjectRepository) Create(ctx context.Context, subject *entities.Subject) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validar duplicado por Name
	for _, s := range r.subjects {
		if s.Name == subject.Name {
			return errors.NewConflictError("subject with this name already exists")
		}
	}

	// Generar ID si no existe
	if subject.ID == uuid.Nil {
		subject.ID = uuid.New()
	}

	// Establecer timestamps
	now := time.Now()
	subject.CreatedAt = now
	subject.UpdatedAt = now

	// Almacenar copia
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

	// Retornar copia
	subjectCopy := *subject
	return &subjectCopy, nil
}

// Update actualiza una materia existente
func (r *MockSubjectRepository) Update(ctx context.Context, subject *entities.Subject) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verificar que la materia existe
	existing, exists := r.subjects[subject.ID]
	if !exists {
		return errors.NewNotFoundError("subject not found")
	}

	// Validar duplicado por Name (excepto el mismo registro)
	for _, s := range r.subjects {
		if s.ID != subject.ID && s.Name == subject.Name {
			return errors.NewConflictError("subject with this name already exists")
		}
	}

	// Actualizar timestamp
	subject.UpdatedAt = time.Now()

	// Preservar CreatedAt original
	subject.CreatedAt = existing.CreatedAt

	// Almacenar copia actualizada
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

	// Soft delete
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

	// Filtrar materias activas
	for _, subject := range r.subjects {
		if subject.IsActive {
			// Agregar copia de la materia
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

// Reset limpia todos los datos y recarga los datos mock (útil para testing)
func (r *MockSubjectRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Limpiar mapa
	r.subjects = make(map[uuid.UUID]*entities.Subject)

	// Recargar datos desde mockData
	for _, subject := range mockData.GetSubjects() {
		subjectCopy := *subject
		r.subjects[subject.ID] = &subjectCopy
	}
}
