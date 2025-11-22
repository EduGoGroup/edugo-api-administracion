package repository

import (
	"context"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type UnitRepository interface {
	Create(ctx context.Context, unit *entities.Unit) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Unit, error)
	Update(ctx context.Context, unit *entities.Unit) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, schoolID uuid.UUID) ([]*entities.Unit, error)
}
