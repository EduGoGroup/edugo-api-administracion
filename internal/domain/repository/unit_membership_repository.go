package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type UnitMembershipRepository interface {
	Create(ctx context.Context, membership *entities.Membership) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error)
	FindByUserAndUnit(ctx context.Context, userID, unitID uuid.UUID) (*entities.Membership, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error)
	FindByUnit(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error)
	FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool) ([]*entities.Membership, error)
	Update(ctx context.Context, membership *entities.Membership) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByUnitAndUser(ctx context.Context, unitID, userID uuid.UUID) (bool, error)
}
