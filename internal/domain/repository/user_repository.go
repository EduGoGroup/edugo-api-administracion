package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// UserRepository define las operaciones de persistencia para User
type UserRepository interface {
	// Create crea un nuevo usuario
	Create(ctx context.Context, user *entities.User) error

	// FindByID busca un usuario por ID
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)

	// FindByEmail busca un usuario por email
	FindByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update actualiza un usuario existente
	Update(ctx context.Context, user *entities.User) error

	// Delete elimina un usuario (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// List lista usuarios con filtros opcionales
	List(ctx context.Context, filters ListFilters) ([]*entities.User, error)

	// ExistsByEmail verifica si existe un usuario con ese email
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// ListFilters representa filtros para listar usuarios
type ListFilters struct {
	Role     *string
	IsActive *bool
	Limit    int
	Offset   int
}
