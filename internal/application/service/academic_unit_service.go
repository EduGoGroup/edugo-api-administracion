package service

import (
	"context"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// AcademicUnitService define las operaciones de negocio para AcademicUnit
type AcademicUnitService interface {
	CreateUnit(ctx context.Context, schoolID string, req dto.CreateAcademicUnitRequest) (*dto.AcademicUnitResponse, error)
	GetUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error)
	GetUnitTree(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error)
	ListUnitsBySchool(ctx context.Context, schoolID string, includeDeleted bool) ([]dto.AcademicUnitResponse, error)
	ListUnitsByType(ctx context.Context, schoolID string, unitType string) ([]dto.AcademicUnitResponse, error)
	UpdateUnit(ctx context.Context, id string, req dto.UpdateAcademicUnitRequest) (*dto.AcademicUnitResponse, error)
	DeleteUnit(ctx context.Context, id string) error
	RestoreUnit(ctx context.Context, id string) error
	GetHierarchyPath(ctx context.Context, id string) ([]dto.AcademicUnitResponse, error)
}

type academicUnitService struct {
	unitRepo   repository.AcademicUnitRepository
	schoolRepo repository.SchoolRepository
	logger     logger.Logger
}

// NewAcademicUnitService crea un nuevo AcademicUnitService
func NewAcademicUnitService(
	unitRepo repository.AcademicUnitRepository,
	schoolRepo repository.SchoolRepository,
	logger logger.Logger,
) AcademicUnitService {
	return &academicUnitService{
		unitRepo:   unitRepo,
		schoolRepo: schoolRepo,
		logger:     logger,
	}
}

// CreateUnit crea una nueva unidad académica
func (s *academicUnitService) CreateUnit(ctx context.Context, schoolID string, req dto.CreateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	// 1. Validar que la escuela existe
	schoolIDVO, err := valueobject.SchoolIDFromString(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	school, err := s.schoolRepo.FindByID(ctx, schoolIDVO)
	if err != nil {
		return nil, errors.NewDatabaseError("find school", err)
	}
	if school == nil {
		return nil, errors.NewNotFoundError("school")
	}

	// 2. Validar y crear el tipo de unidad
	unitType, err := valueobject.NewUnitType(req.Type)
	if err != nil {
		return nil, err
	}

	// 3. Verificar código único si se proporciona
	if req.Code != "" {
		exists, err := s.unitRepo.ExistsBySchoolIDAndCode(ctx, schoolIDVO, req.Code)
		if err != nil {
			return nil, errors.NewDatabaseError("check unit code", err)
		}
		if exists {
			return nil, errors.NewAlreadyExistsError("academic unit with code").WithField("code", req.Code)
		}
	}

	// 4. Crear entidad
	unit, err := entity.NewAcademicUnit(schoolIDVO, unitType, req.DisplayName, req.Code)
	if err != nil {
		s.logger.Warn("failed to create unit entity", "error", err)
		return nil, err
	}

	// 5. Establecer padre si se proporciona
	if req.ParentUnitID != nil {
		parentID, err := valueobject.UnitIDFromString(*req.ParentUnitID)
		if err != nil {
			return nil, errors.NewValidationError("invalid parent_unit_id")
		}

		// Verificar que el padre existe y obtener su tipo
		parent, err := s.unitRepo.FindByID(ctx, parentID, false)
		if err != nil {
			return nil, errors.NewDatabaseError("find parent unit", err)
		}
		if parent == nil {
			return nil, errors.NewNotFoundError("parent unit")
		}

		// Validar la relación padre-hijo
		if err := unit.SetParent(parentID, parent.UnitType()); err != nil {
			return nil, err
		}
	}

	// 6. Agregar descripción y metadata
	if req.Description != "" {
		unit.UpdateInfo(req.DisplayName, req.Description)
	}

	if req.Metadata != nil {
		for key, value := range req.Metadata {
			unit.SetMetadata(key, value)
		}
	}

	// 7. Persistir
	if err := s.unitRepo.Create(ctx, unit); err != nil {
		s.logger.Error("failed to create unit", "error", err, "name", req.DisplayName)
		return nil, errors.NewDatabaseError("create unit", err)
	}

	s.logger.Info("unit created successfully", "id", unit.ID().String(), "name", req.DisplayName)

	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

// GetUnit obtiene una unidad por ID
func (s *academicUnitService) GetUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error) {
	unitID, err := valueobject.UnitIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		s.logger.Error("failed to find unit", "error", err, "id", id)
		return nil, errors.NewDatabaseError("find unit", err)
	}

	if unit == nil {
		return nil, errors.NewNotFoundError("academic unit")
	}

	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

// GetUnitTree obtiene el árbol jerárquico completo de una escuela
func (s *academicUnitService) GetUnitTree(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error) {
	schoolIDVO, err := valueobject.SchoolIDFromString(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	units, err := s.unitRepo.FindBySchoolID(ctx, schoolIDVO, false)
	if err != nil {
		s.logger.Error("failed to get units", "error", err, "school_id", schoolID)
		return nil, errors.NewDatabaseError("get units", err)
	}

	return dto.BuildUnitTree(units), nil
}

// ListUnitsBySchool lista todas las unidades de una escuela
func (s *academicUnitService) ListUnitsBySchool(ctx context.Context, schoolID string, includeDeleted bool) ([]dto.AcademicUnitResponse, error) {
	schoolIDVO, err := valueobject.SchoolIDFromString(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	units, err := s.unitRepo.FindBySchoolID(ctx, schoolIDVO, includeDeleted)
	if err != nil {
		s.logger.Error("failed to list units", "error", err, "school_id", schoolID)
		return nil, errors.NewDatabaseError("list units", err)
	}

	return dto.ToAcademicUnitResponseList(units), nil
}

// ListUnitsByType lista unidades por tipo
func (s *academicUnitService) ListUnitsByType(ctx context.Context, schoolID string, unitType string) ([]dto.AcademicUnitResponse, error) {
	schoolIDVO, err := valueobject.SchoolIDFromString(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	unitTypeVO, err := valueobject.NewUnitType(unitType)
	if err != nil {
		return nil, err
	}

	units, err := s.unitRepo.FindByType(ctx, schoolIDVO, unitTypeVO, false)
	if err != nil {
		s.logger.Error("failed to list units by type", "error", err, "type", unitType)
		return nil, errors.NewDatabaseError("list units", err)
	}

	return dto.ToAcademicUnitResponseList(units), nil
}

// UpdateUnit actualiza una unidad
func (s *academicUnitService) UpdateUnit(ctx context.Context, id string, req dto.UpdateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	unitID, err := valueobject.UnitIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		return nil, errors.NewDatabaseError("find unit", err)
	}
	if unit == nil {
		return nil, errors.NewNotFoundError("academic unit")
	}

	// Actualizar padre si se proporciona
	if req.ParentUnitID != nil {
		if *req.ParentUnitID == "" {
			unit.RemoveParent()
		} else {
			parentID, err := valueobject.UnitIDFromString(*req.ParentUnitID)
			if err != nil {
				return nil, errors.NewValidationError("invalid parent_unit_id")
			}

			parent, err := s.unitRepo.FindByID(ctx, parentID, false)
			if err != nil {
				return nil, errors.NewDatabaseError("find parent unit", err)
			}
			if parent == nil {
				return nil, errors.NewNotFoundError("parent unit")
			}

			if err := unit.SetParent(parentID, parent.UnitType()); err != nil {
				return nil, err
			}
		}
	}

	// Actualizar info si se proporciona
	if req.DisplayName != nil || req.Description != nil {
		name := unit.DisplayName()
		if req.DisplayName != nil {
			name = *req.DisplayName
		}
		desc := unit.Description()
		if req.Description != nil {
			desc = *req.Description
		}

		if err := unit.UpdateInfo(name, desc); err != nil {
			return nil, err
		}
	}

	// Actualizar metadata
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			unit.SetMetadata(key, value)
		}
	}

	if err := s.unitRepo.Update(ctx, unit); err != nil {
		s.logger.Error("failed to update unit", "error", err, "id", id)
		return nil, errors.NewDatabaseError("update unit", err)
	}

	s.logger.Info("unit updated successfully", "id", id)

	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

// DeleteUnit elimina una unidad (soft delete)
func (s *academicUnitService) DeleteUnit(ctx context.Context, id string) error {
	unitID, err := valueobject.UnitIDFromString(id)
	if err != nil {
		return errors.NewValidationError("invalid unit ID")
	}

	// Verificar que existe
	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		return errors.NewDatabaseError("find unit", err)
	}
	if unit == nil {
		return errors.NewNotFoundError("academic unit")
	}

	// Verificar que no tiene hijos activos
	hasChildren, err := s.unitRepo.HasChildren(ctx, unitID)
	if err != nil {
		return errors.NewDatabaseError("check children", err)
	}
	if hasChildren {
		return errors.NewBusinessRuleError("cannot delete unit with active children")
	}

	if err := s.unitRepo.SoftDelete(ctx, unitID); err != nil {
		s.logger.Error("failed to delete unit", "error", err, "id", id)
		return errors.NewDatabaseError("delete unit", err)
	}

	s.logger.Info("unit deleted successfully", "id", id)
	return nil
}

// RestoreUnit restaura una unidad eliminada
func (s *academicUnitService) RestoreUnit(ctx context.Context, id string) error {
	unitID, err := valueobject.UnitIDFromString(id)
	if err != nil {
		return errors.NewValidationError("invalid unit ID")
	}

	if err := s.unitRepo.Restore(ctx, unitID); err != nil {
		s.logger.Error("failed to restore unit", "error", err, "id", id)
		return errors.NewDatabaseError("restore unit", err)
	}

	s.logger.Info("unit restored successfully", "id", id)
	return nil
}

// GetHierarchyPath obtiene el path desde raíz hasta la unidad
func (s *academicUnitService) GetHierarchyPath(ctx context.Context, id string) ([]dto.AcademicUnitResponse, error) {
	unitID, err := valueobject.UnitIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	units, err := s.unitRepo.GetHierarchyPath(ctx, unitID)
	if err != nil {
		s.logger.Error("failed to get hierarchy path", "error", err, "id", id)
		return nil, errors.NewDatabaseError("get hierarchy path", err)
	}

	return dto.ToAcademicUnitResponseList(units), nil
}
