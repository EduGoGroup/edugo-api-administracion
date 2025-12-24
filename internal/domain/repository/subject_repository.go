package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// SubjectRepository define las operaciones de persistencia para Subject
type SubjectRepository interface {
	// Create crea una nueva materia
	Create(ctx context.Context, subject *entities.Subject) error

	// FindByID busca una materia por ID
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Subject, error)

	// Update actualiza una materia
	Update(ctx context.Context, subject *entities.Subject) error

	// Delete elimina una materia (cambia is_active a false)
	Delete(ctx context.Context, id uuid.UUID) error

	// List lista todas las materias activas
	List(ctx context.Context) ([]*entities.Subject, error)

	// FindBySchoolID lista materias activas filtradas por school_id
	FindBySchoolID(ctx context.Context, schoolID uuid.UUID) ([]*entities.Subject, error)

	// ExistsByName verifica si existe una materia con ese nombre
	ExistsByName(ctx context.Context, name string) (bool, error)
}
