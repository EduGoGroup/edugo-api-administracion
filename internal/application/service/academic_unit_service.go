package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

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

func (s *academicUnitService) CreateUnit(ctx context.Context, schoolID string, req dto.CreateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	// Parse IDs
	schoolUUID, err := uuid.Parse(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	// Verificar escuela existe
	school, err := s.schoolRepo.FindByID(ctx, schoolUUID)
	if err != nil {
		return nil, errors.NewDatabaseError("find school", err)
	}
	if school == nil {
		return nil, errors.NewNotFoundError("school")
	}

	// Verificar código único
	if req.Code != "" {
		exists, err := s.unitRepo.ExistsBySchoolIDAndCode(ctx, schoolUUID, req.Code)
		if err != nil {
			return nil, errors.NewDatabaseError("check unit code", err)
		}
		if exists {
			return nil, errors.NewAlreadyExistsError("academic unit with code").WithField("code", req.Code)
		}
	}

	// Validar padre si existe
	var parentUUID *uuid.UUID
	if req.ParentUnitID != nil {
		pid, err := uuid.Parse(*req.ParentUnitID)
		if err != nil {
			return nil, errors.NewValidationError("invalid parent_unit_id")
		}
		parent, err := s.unitRepo.FindByID(ctx, pid, false)
		if err != nil {
			return nil, errors.NewDatabaseError("find parent unit", err)
		}
		if parent == nil {
			return nil, errors.NewNotFoundError("parent unit")
		}
		parentUUID = &pid
	}

	// Validar tipo de unidad usando value object
	if _, err := valueobject.ParseUnitType(req.Type); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	// Crear unidad (lógica de validación movida aquí del entity)
	if req.DisplayName == "" {
		return nil, errors.NewValidationError("display_name is required")
	}
	if len(req.DisplayName) < 3 {
		return nil, errors.NewValidationError("display_name must be at least 3 characters")
	}

	now := time.Now()
	unit := &entities.AcademicUnit{
		ID:           uuid.New(),
		ParentUnitID: parentUUID,
		SchoolID:     schoolUUID,
		Name:         req.DisplayName,
		Code:         req.Code,
		Type:         req.Type,
		Description:  &req.Description,
		Level:        nil, // TODO: agregar si se necesita
		AcademicYear: 0,   // TODO: agregar si se necesita
		Metadata:     []byte("{}"),
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    nil,
	}

	// Persistir
	if err := s.unitRepo.Create(ctx, unit); err != nil {
		s.logger.Error("failed to create unit", "error", err, "name", req.DisplayName)
		return nil, errors.NewDatabaseError("create unit", err)
	}

	s.logger.Info("unit created successfully", "id", unit.ID.String(), "name", req.DisplayName)

	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

func (s *academicUnitService) GetUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error) {
	unitID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		s.logger.Error("failed to find unit", "error", err, "id", id)
		// Propagar AppError directamente (ej: NotFoundError)
		if _, ok := errors.GetAppError(err); ok {
			return nil, err
		}
		return nil, errors.NewDatabaseError("find unit", err)
	}
	if unit == nil {
		return nil, errors.NewNotFoundError("academic unit")
	}

	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

func (s *academicUnitService) GetUnitTree(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error) {
	schoolUUID, err := uuid.Parse(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	units, err := s.unitRepo.FindBySchoolID(ctx, schoolUUID, false)
	if err != nil {
		return nil, errors.NewDatabaseError("find units", err)
	}

	return dto.BuildUnitTree(units), nil
}

func (s *academicUnitService) ListUnitsBySchool(ctx context.Context, schoolID string, includeDeleted bool) ([]dto.AcademicUnitResponse, error) {
	schoolUUID, err := uuid.Parse(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	units, err := s.unitRepo.FindBySchoolID(ctx, schoolUUID, includeDeleted)
	if err != nil {
		return nil, errors.NewDatabaseError("find units", err)
	}

	responses := make([]dto.AcademicUnitResponse, len(units))
	for i, unit := range units {
		responses[i] = dto.ToAcademicUnitResponse(unit)
	}
	return responses, nil
}

func (s *academicUnitService) ListUnitsByType(ctx context.Context, schoolID string, unitType string) ([]dto.AcademicUnitResponse, error) {
	schoolUUID, err := uuid.Parse(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	// Validar tipo de unidad usando value object
	if _, err := valueobject.ParseUnitType(unitType); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	units, err := s.unitRepo.FindByType(ctx, schoolUUID, unitType, false)
	if err != nil {
		return nil, errors.NewDatabaseError("find units", err)
	}

	responses := make([]dto.AcademicUnitResponse, len(units))
	for i, unit := range units {
		responses[i] = dto.ToAcademicUnitResponse(unit)
	}
	return responses, nil
}

func (s *academicUnitService) UpdateUnit(ctx context.Context, id string, req dto.UpdateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	unitID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil || unit == nil {
		return nil, errors.NewNotFoundError("academic unit")
	}

	// Actualizar campos (lógica movida del entity)
	if req.DisplayName != nil {
		if len(*req.DisplayName) < 3 {
			return nil, errors.NewValidationError("display_name must be at least 3 characters")
		}
		unit.Name = *req.DisplayName
	}

	if req.Description != nil {
		unit.Description = req.Description
	}

	if req.ParentUnitID != nil {
		pid, err := uuid.Parse(*req.ParentUnitID)
		if err != nil {
			return nil, errors.NewValidationError("invalid parent_unit_id")
		}
		if pid == unitID {
			return nil, errors.NewBusinessRuleError("unit cannot be its own parent")
		}
		unit.ParentUnitID = &pid
	}

	unit.UpdatedAt = time.Now()

	if err := s.unitRepo.Update(ctx, unit); err != nil {
		s.logger.Error("failed to update unit", "error", err)
		return nil, errors.NewDatabaseError("update unit", err)
	}

	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

func (s *academicUnitService) DeleteUnit(ctx context.Context, id string) error {
	unitID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid unit ID")
	}

	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil || unit == nil {
		return errors.NewNotFoundError("academic unit")
	}

	if err := s.unitRepo.SoftDelete(ctx, unitID); err != nil {
		s.logger.Error("failed to delete unit", "error", err, "id", id)
		return errors.NewDatabaseError("delete unit", err)
	}

	s.logger.Info("unit deleted", "id", id)
	return nil
}

func (s *academicUnitService) RestoreUnit(ctx context.Context, id string) error {
	unitID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid unit ID")
	}

	if err := s.unitRepo.Restore(ctx, unitID); err != nil {
		s.logger.Error("failed to restore unit", "error", err, "id", id)
		return errors.NewDatabaseError("restore unit", err)
	}

	s.logger.Info("unit restored", "id", id)
	return nil
}

func (s *academicUnitService) GetHierarchyPath(ctx context.Context, id string) ([]dto.AcademicUnitResponse, error) {
	unitID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	units, err := s.unitRepo.GetHierarchyPath(ctx, unitID)
	if err != nil {
		return nil, errors.NewDatabaseError("get hierarchy path", err)
	}

	responses := make([]dto.AcademicUnitResponse, len(units))
	for i, unit := range units {
		responses[i] = dto.ToAcademicUnitResponse(unit)
	}
	return responses, nil
}
