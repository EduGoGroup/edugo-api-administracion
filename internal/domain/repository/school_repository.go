package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// SchoolRepository define las operaciones de persistencia para School
type SchoolRepository interface {
	// Create crea una nueva escuela
	Create(ctx context.Context, school *entities.School) error

	// FindByID busca una escuela por ID
	FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error)

	// FindByCode busca una escuela por código único
	FindByCode(ctx context.Context, code string) (*entities.School, error)

	// FindByName busca una escuela por nombre
	FindByName(ctx context.Context, name string) (*entities.School, error)

	// Update actualiza una escuela existente
	Update(ctx context.Context, school *entities.School) error

	// Delete elimina una escuela (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// List lista escuelas con filtros opcionales
	List(ctx context.Context, filters ListFilters) ([]*entities.School, error)

	// ExistsByName verifica si existe una escuela con ese nombre
	ExistsByName(ctx context.Context, name string) (bool, error)

	// ExistsByCode verifica si existe una escuela con ese código
	ExistsByCode(ctx context.Context, code string) (bool, error)
}
