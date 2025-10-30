package repository

import (
	"context"

	"github.com/edugo/api-administracion/internal/domain/entity"
	"github.com/edugo/api-administracion/internal/domain/valueobject"
)

// SchoolRepository define las operaciones de persistencia para School
type SchoolRepository interface {
	// Create crea una nueva escuela
	Create(ctx context.Context, school *entity.School) error

	// FindByID busca una escuela por ID
	FindByID(ctx context.Context, id valueobject.SchoolID) (*entity.School, error)

	// FindByName busca una escuela por nombre
	FindByName(ctx context.Context, name string) (*entity.School, error)

	// Update actualiza una escuela existente
	Update(ctx context.Context, school *entity.School) error

	// Delete elimina una escuela (soft delete)
	Delete(ctx context.Context, id valueobject.SchoolID) error

	// List lista escuelas con filtros opcionales
	List(ctx context.Context, filters ListFilters) ([]*entity.School, error)

	// ExistsByName verifica si existe una escuela con ese nombre
	ExistsByName(ctx context.Context, name string) (bool, error)
}
