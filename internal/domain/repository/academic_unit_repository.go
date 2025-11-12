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
}
