package repository

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
)

// UnitMembershipRepository define las operaciones de persistencia para UnitMembership
type UnitMembershipRepository interface {
	// Create crea una nueva membresía
	Create(ctx context.Context, membership *entity.UnitMembership) error

	// FindByID busca una membresía por ID
	FindByID(ctx context.Context, id valueobject.MembershipID) (*entity.UnitMembership, error)

	// FindByUnitID lista todas las membresías de una unidad
	FindByUnitID(ctx context.Context, unitID valueobject.UnitID, activeOnly bool) ([]*entity.UnitMembership, error)

	// FindByUserID lista todas las membresías de un usuario
	FindByUserID(ctx context.Context, userID valueobject.UserID, activeOnly bool) ([]*entity.UnitMembership, error)

	// FindByUnitAndUser busca una membresía específica usuario-unidad
	FindByUnitAndUser(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) (*entity.UnitMembership, error)

	// FindByRole lista membresías por rol
	FindByRole(ctx context.Context, unitID valueobject.UnitID, role valueobject.MembershipRole, activeOnly bool) ([]*entity.UnitMembership, error)

	// FindActiveAt lista membresías activas en un momento específico
	FindActiveAt(ctx context.Context, unitID valueobject.UnitID, at time.Time) ([]*entity.UnitMembership, error)

	// Update actualiza una membresía existente
	Update(ctx context.Context, membership *entity.UnitMembership) error

	// Delete elimina una membresía permanentemente
	Delete(ctx context.Context, id valueobject.MembershipID) error

	// ExistsByUnitAndUser verifica si existe una membresía activa
	ExistsByUnitAndUser(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) (bool, error)

	// CountByUnit cuenta membresías de una unidad
	CountByUnit(ctx context.Context, unitID valueobject.UnitID, activeOnly bool) (int, error)

	// CountByRole cuenta membresías por rol
	CountByRole(ctx context.Context, unitID valueobject.UnitID, role valueobject.MembershipRole) (int, error)
}
