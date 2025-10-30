package repository

import (
	"context"

	"github.com/edugo/api-administracion/internal/domain/entity"
	"github.com/edugo/api-administracion/internal/domain/valueobject"
)

// UnitRepository define las operaciones de persistencia para Unit
type UnitRepository interface {
	// Create crea una nueva unidad
	Create(ctx context.Context, unit *entity.Unit) error

	// FindByID busca una unidad por ID
	FindByID(ctx context.Context, id valueobject.UnitID) (*entity.Unit, error)

	// Update actualiza una unidad existente
	Update(ctx context.Context, unit *entity.Unit) error

	// Delete elimina una unidad (soft delete)
	Delete(ctx context.Context, id valueobject.UnitID) error

	// FindBySchool busca unidades de una escuela
	FindBySchool(ctx context.Context, schoolID valueobject.SchoolID) ([]*entity.Unit, error)

	// FindChildren busca unidades hijas de una unidad padre
	FindChildren(ctx context.Context, parentID valueobject.UnitID) ([]*entity.Unit, error)

	// HasChildren verifica si una unidad tiene hijos
	HasChildren(ctx context.Context, unitID valueobject.UnitID) (bool, error)

	// IsDescendantOf verifica si una unidad es descendiente de otra (para evitar ciclos)
	IsDescendantOf(ctx context.Context, unitID, ancestorID valueobject.UnitID) (bool, error)

	// AddMember agrega un miembro a una unidad
	AddMember(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID, role string) error

	// RemoveMember elimina un miembro de una unidad
	RemoveMember(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) error

	// IsMember verifica si un usuario es miembro de una unidad
	IsMember(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) (bool, error)
}
