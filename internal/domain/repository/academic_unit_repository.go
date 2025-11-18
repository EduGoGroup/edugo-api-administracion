package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
)

// AcademicUnitRepository define las operaciones de persistencia para AcademicUnit
type AcademicUnitRepository interface {
	// Create crea una nueva unidad académica
	Create(ctx context.Context, unit *entity.AcademicUnit) error

	// FindByID busca una unidad por ID (incluye eliminadas si includeDeleted=true)
	FindByID(ctx context.Context, id valueobject.UnitID, includeDeleted bool) (*entity.AcademicUnit, error)

	// FindBySchoolIDAndCode busca una unidad por escuela y código
	FindBySchoolIDAndCode(ctx context.Context, schoolID valueobject.SchoolID, code string) (*entity.AcademicUnit, error)

	// FindBySchoolID lista todas las unidades de una escuela
	FindBySchoolID(ctx context.Context, schoolID valueobject.SchoolID, includeDeleted bool) ([]*entity.AcademicUnit, error)

	// FindByParentID lista unidades hijas de una unidad padre
	FindByParentID(ctx context.Context, parentID valueobject.UnitID, includeDeleted bool) ([]*entity.AcademicUnit, error)

	// FindRootUnits lista unidades raíz (sin padre) de una escuela
	FindRootUnits(ctx context.Context, schoolID valueobject.SchoolID) ([]*entity.AcademicUnit, error)

	// FindByType lista unidades por tipo
	FindByType(ctx context.Context, schoolID valueobject.SchoolID, unitType valueobject.UnitType, includeDeleted bool) ([]*entity.AcademicUnit, error)

	// Update actualiza una unidad existente
	Update(ctx context.Context, unit *entity.AcademicUnit) error

	// SoftDelete marca una unidad como eliminada
	SoftDelete(ctx context.Context, id valueobject.UnitID) error

	// Restore restaura una unidad eliminada
	Restore(ctx context.Context, id valueobject.UnitID) error

	// HardDelete elimina permanentemente una unidad (cuidado con cascada)
	HardDelete(ctx context.Context, id valueobject.UnitID) error

	// GetHierarchyPath obtiene el path jerárquico desde raíz hasta la unidad
	GetHierarchyPath(ctx context.Context, id valueobject.UnitID) ([]*entity.AcademicUnit, error)

	// ExistsBySchoolIDAndCode verifica si existe una unidad con ese código en la escuela
	ExistsBySchoolIDAndCode(ctx context.Context, schoolID valueobject.SchoolID, code string) (bool, error)

	// HasChildren verifica si una unidad tiene hijos
	HasChildren(ctx context.Context, id valueobject.UnitID) (bool, error)

	// =====================================================
	// Ltree-based hierarchical query methods (Sprint-03)
	// =====================================================

	// FindByPath busca una unidad académica por su path ltree
	FindByPath(ctx context.Context, path string) (*entity.AcademicUnit, error)

	// FindChildren retorna los hijos directos de una unidad académica
	FindChildren(ctx context.Context, parentID valueobject.UnitID) ([]*entity.AcademicUnit, error)

	// FindDescendants retorna TODOS los descendientes de una unidad usando ltree
	// Incluye hijos, nietos, bisnietos, etc. (toda la jerarquía debajo)
	FindDescendants(ctx context.Context, unitID valueobject.UnitID) ([]*entity.AcademicUnit, error)

	// FindAncestors retorna TODOS los ancestros de una unidad usando ltree
	// Incluye padre, abuelo, bisabuelo, etc. (toda la jerarquía arriba)
	FindAncestors(ctx context.Context, unitID valueobject.UnitID) ([]*entity.AcademicUnit, error)

	// FindBySchoolIDAndDepth retorna unidades de una escuela a una profundidad específica
	// depth=1: unidades raíz, depth=2: hijos directos de raíz, etc.
	FindBySchoolIDAndDepth(ctx context.Context, schoolID valueobject.SchoolID, depth int) ([]*entity.AcademicUnit, error)

	// MoveSubtree mueve un subárbol completo a un nuevo padre
	// Si newParentID es nil, convierte la unidad en raíz
	MoveSubtree(ctx context.Context, unitID valueobject.UnitID, newParentID *valueobject.UnitID) error
}
