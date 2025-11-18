package service

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// HierarchyService servicio de aplicación para operaciones jerárquicas
type HierarchyService struct {
	unitRepo      repository.AcademicUnitRepository
	schoolRepo    repository.SchoolRepository
	domainService *service.AcademicUnitDomainService
}

// NewHierarchyService constructor
func NewHierarchyService(
	unitRepo repository.AcademicUnitRepository,
	schoolRepo repository.SchoolRepository,
	domainService *service.AcademicUnitDomainService,
) *HierarchyService {
	return &HierarchyService{
		unitRepo:      unitRepo,
		schoolRepo:    schoolRepo,
		domainService: domainService,
	}
}

// CreateUnit crea una nueva unidad académica
func (s *HierarchyService) CreateUnit(
	ctx context.Context,
	parentUnitID *valueobject.UnitID,
	schoolID valueobject.SchoolID,
	unitType valueobject.UnitType,
	name string,
	code string,
	description string,
) (*entity.AcademicUnit, error) {
	// 1. Validar que la escuela existe
	_, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("school %s not found", schoolID))
	}

	// 2. Validar que el código no esté duplicado en la escuela
	exists, err := s.unitRepo.ExistsBySchoolIDAndCode(ctx, schoolID, code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.NewValidationError(fmt.Sprintf("code %s already exists in school", code))
	}

	// 3. Si tiene padre, validar que existe y no crea ciclo
	var parent *entity.AcademicUnit
	if parentUnitID != nil {
		parent, err = s.unitRepo.FindByID(ctx, *parentUnitID, false)
		if err != nil {
			return nil, errors.NewNotFoundError(fmt.Sprintf("parent unit %s not found", parentUnitID))
		}

		// Validar que el padre pertenece a la misma escuela
		if parent.SchoolID() != schoolID {
			return nil, errors.NewValidationError("parent unit must belong to the same school")
		}
	}

	// 4. Crear la unidad
	unit, err := entity.NewAcademicUnit(schoolID, unitType, name, code)
	if err != nil {
		return nil, err
	}

	if description != "" {
		unit.SetDescription(description)
	}

	// 5. Establecer padre si existe
	if parent != nil {
		if err := s.domainService.SetParent(unit, parent.ID(), parent.UnitType()); err != nil {
			return nil, err
		}
	}

	// 6. Persistir
	if err := s.unitRepo.Create(ctx, unit); err != nil {
		return nil, err
	}

	return unit, nil
}

// GetUnitTree obtiene el árbol jerárquico completo de una unidad usando ltree
func (s *HierarchyService) GetUnitTree(ctx context.Context, unitID valueobject.UnitID) (*entity.AcademicUnit, []*entity.AcademicUnit, error) {
	// 1. Obtener la unidad raíz
	root, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		return nil, nil, err
	}

	// 2. Obtener todos los descendientes usando ltree (¡Sprint-03!)
	descendants, err := s.unitRepo.FindDescendants(ctx, unitID)
	if err != nil {
		return nil, nil, err
	}

	return root, descendants, nil
}

// MoveUnit mueve una unidad a un nuevo padre (o a raíz si newParentID es nil)
func (s *HierarchyService) MoveUnit(
	ctx context.Context,
	unitID valueobject.UnitID,
	newParentID *valueobject.UnitID,
) error {
	// 1. Validar que la unidad existe
	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		return err
	}

	// 2. Si hay nuevo padre, validar
	if newParentID != nil {
		newParent, err := s.unitRepo.FindByID(ctx, *newParentID, false)
		if err != nil {
			return errors.NewNotFoundError(fmt.Sprintf("new parent unit %s not found", newParentID))
		}

		// Validar misma escuela
		if newParent.SchoolID() != unit.SchoolID() {
			return errors.NewValidationError("cannot move unit to different school")
		}

		// Validar que no crea ciclo (nuevo padre no puede ser descendiente)
		descendants, err := s.unitRepo.FindDescendants(ctx, unitID)
		if err != nil {
			return err
		}

		for _, desc := range descendants {
			if desc.ID() == *newParentID {
				return errors.NewValidationError("cannot move unit: would create a circular reference")
			}
		}
	}

	// 3. Mover usando ltree (Sprint-03!)
	return s.unitRepo.MoveSubtree(ctx, unitID, newParentID)
}

// ValidateNoCircularReference valida que mover una unidad no cree un ciclo
func (s *HierarchyService) ValidateNoCircularReference(
	ctx context.Context,
	unitID valueobject.UnitID,
	newParentID valueobject.UnitID,
) error {
	// Obtener todos los descendientes de la unidad
	descendants, err := s.unitRepo.FindDescendants(ctx, unitID)
	if err != nil {
		return err
	}

	// Verificar que el nuevo padre no sea uno de los descendientes
	for _, desc := range descendants {
		if desc.ID() == newParentID {
			return errors.NewValidationError("circular reference detected")
		}
	}

	// Verificar que no sea la misma unidad
	if unitID == newParentID {
		return errors.NewValidationError("unit cannot be its own parent")
	}

	return nil
}
