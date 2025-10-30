package repository

import (
	"context"

	"github.com/edugo/api-administracion/internal/domain/entity"
	"github.com/edugo/api-administracion/internal/domain/valueobject"
)

// UserRepository define las operaciones de persistencia para User
type UserRepository interface {
	// Create crea un nuevo usuario
	Create(ctx context.Context, user *entity.User) error

	// FindByID busca un usuario por ID
	FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error)

	// FindByEmail busca un usuario por email
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)

	// Update actualiza un usuario existente
	Update(ctx context.Context, user *entity.User) error

	// Delete elimina un usuario (soft delete)
	Delete(ctx context.Context, id valueobject.UserID) error

	// List lista usuarios con filtros opcionales
	List(ctx context.Context, filters ListFilters) ([]*entity.User, error)

	// ExistsByEmail verifica si existe un usuario con ese email
	ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error)
}

// ListFilters representa filtros para listar usuarios
type ListFilters struct {
	Role     *string
	IsActive *bool
	Limit    int
	Offset   int
}
