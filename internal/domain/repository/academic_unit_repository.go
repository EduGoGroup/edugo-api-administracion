package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// AcademicUnitRepository define las operaciones de persistencia para AcademicUnit
type AcademicUnitRepository interface {
	// Create crea una nueva unidad académica
	Create(ctx context.Context, unit *entities.AcademicUnit) error

	// FindByID busca una unidad por ID (incluye eliminadas si includeDeleted=true)
	FindByID(ctx context.Context, id uuid.UUID, includeDeleted bool) (*entities.AcademicUnit, error)

	// FindBySchoolIDAndCode busca una unidad por escuela y código
	FindBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (*entities.AcademicUnit, error)

	// FindBySchoolID lista todas las unidades de una escuela
	FindBySchoolID(ctx context.Context, schoolID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error)

	// FindByParentID lista unidades hijas de una unidad padre
	FindByParentID(ctx context.Context, parentID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error)

	// FindRootUnits lista unidades raíz (sin padre) de una escuela
	FindRootUnits(ctx context.Context, schoolID uuid.UUID) ([]*entities.AcademicUnit, error)

	// FindByType lista unidades por tipo
	FindByType(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool) ([]*entities.AcademicUnit, error)

	// Update actualiza una unidad existente
	Update(ctx context.Context, unit *entities.AcademicUnit) error

	// SoftDelete marca una unidad como eliminada
	SoftDelete(ctx context.Context, id uuid.UUID) error

	// Restore restaura una unidad eliminada
	Restore(ctx context.Context, id uuid.UUID) error

	// HardDelete elimina permanentemente una unidad (cuidado con cascada)
	HardDelete(ctx context.Context, id uuid.UUID) error

	// GetHierarchyPath obtiene el path jerárquico desde raíz hasta la unidad
	GetHierarchyPath(ctx context.Context, id uuid.UUID) ([]*entities.AcademicUnit, error)

	// ExistsBySchoolIDAndCode verifica si existe una unidad con ese código en la escuela
	ExistsBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (bool, error)
}
