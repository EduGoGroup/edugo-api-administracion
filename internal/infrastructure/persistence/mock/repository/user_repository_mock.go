package repository

import (
	"context"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/mock/dataset"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/google/uuid"
)

// MockUserRepository es una implementación en memoria del UserRepository para testing
// Usa el dataset generado automáticamente desde SQL migrations
type MockUserRepository struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*entities.User
}

// NewMockUserRepository crea una nueva instancia de MockUserRepository
// Pre-carga los usuarios desde el dataset generado
func NewMockUserRepository() repository.UserRepository {
	users := make(map[uuid.UUID]*entities.User)

	// Pre-cargar datos desde dataset generado
	for _, user := range dataset.DB.Users.List() {
		// Hacer una copia del usuario para evitar modificaciones externas
		userCopy := *user
		users[user.ID] = &userCopy
	}

	return &MockUserRepository{
		users: users,
	}
}

// Create crea un nuevo usuario en el repositorio
func (r *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validar que el email no exista
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email && existingUser.DeletedAt == nil {
			return errors.NewConflictError("user with this email already exists")
		}
	}

	// Generar ID si no existe
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// Establecer timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Guardar una copia del usuario
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// FindByID busca un usuario por ID
func (r *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists || user.DeletedAt != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	// Retornar una copia para evitar modificaciones externas
	userCopy := *user
	return &userCopy, nil
}

// FindByEmail busca un usuario por email
func (r *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email && user.DeletedAt == nil {
			// Retornar una copia para evitar modificaciones externas
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, errors.NewNotFoundError("user not found")
}

// Update actualiza un usuario existente
func (r *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verificar que el usuario existe
	existingUser, exists := r.users[user.ID]
	if !exists || existingUser.DeletedAt != nil {
		return errors.NewNotFoundError("user not found")
	}

	// Validar que el email no esté en uso por otro usuario
	for id, u := range r.users {
		if u.Email == user.Email && id != user.ID && u.DeletedAt == nil {
			return errors.NewConflictError("user with this email already exists")
		}
	}

	// Actualizar timestamp
	user.UpdatedAt = time.Now()

	// Preservar CreatedAt original
	user.CreatedAt = existingUser.CreatedAt

	// Guardar una copia del usuario actualizado
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// Delete elimina un usuario (soft delete)
func (r *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists || user.DeletedAt != nil {
		return errors.NewNotFoundError("user not found")
	}

	// Soft delete: establecer DeletedAt
	now := time.Now()
	user.DeletedAt = &now
	user.UpdatedAt = now

	return nil
}

// List lista usuarios con filtros opcionales
func (r *MockUserRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.User

	// Filtrar usuarios
	for _, user := range r.users {
		// Excluir usuarios eliminados
		if user.DeletedAt != nil {
			continue
		}

		// Aplicar filtro de rol
		if filters.Role != nil && user.Role != *filters.Role {
			continue
		}

		// Aplicar filtro de estado activo
		if filters.IsActive != nil && user.IsActive != *filters.IsActive {
			continue
		}

		// Agregar copia del usuario
		userCopy := *user
		result = append(result, &userCopy)
	}

	// Aplicar offset
	if filters.Offset > 0 {
		if filters.Offset >= len(result) {
			return []*entities.User{}, nil
		}
		result = result[filters.Offset:]
	}

	// Aplicar limit
	if filters.Limit > 0 && filters.Limit < len(result) {
		result = result[:filters.Limit]
	}

	return result, nil
}

// ExistsByEmail verifica si existe un usuario con ese email
func (r *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email && user.DeletedAt == nil {
			return true, nil
		}
	}

	return false, nil
}

// Reset reinicia el repositorio a su estado inicial (útil para testing)
func (r *MockUserRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Recargar datos desde dataset generado
	users := make(map[uuid.UUID]*entities.User)
	for _, user := range dataset.DB.Users.List() {
		userCopy := *user
		users[user.ID] = &userCopy
	}

	r.users = users
}
