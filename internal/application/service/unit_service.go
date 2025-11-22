package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

type UnitService interface {
	CreateUnit(ctx context.Context, req dto.CreateUnitRequest) (*dto.UnitResponse, error)
	UpdateUnit(ctx context.Context, id string, req dto.UpdateUnitRequest) (*dto.UnitResponse, error)
	GetUnit(ctx context.Context, id string) (*dto.UnitResponse, error)
}

type unitService struct {
	unitRepo repository.UnitRepository
	logger   logger.Logger
}

func NewUnitService(unitRepo repository.UnitRepository, logger logger.Logger) UnitService {
	return &unitService{unitRepo: unitRepo, logger: logger}
}

func (s *unitService) CreateUnit(ctx context.Context, req dto.CreateUnitRequest) (*dto.UnitResponse, error) {
	// Parse IDs
	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school_id")
	}

	var parentID *uuid.UUID
	if req.ParentUnitID != nil {
		pid, err := uuid.Parse(*req.ParentUnitID)
		if err != nil {
			return nil, errors.NewValidationError("invalid parent_unit_id")
		}
		parentID = &pid
	}

	// Validaciones (lógica movida del entity)
	if req.Name == "" || len(req.Name) < 2 {
		return nil, errors.NewValidationError("name must be at least 2 characters")
	}

	// Crear entidad
	now := time.Now()
	desc := &req.Description
	unit := &entities.Unit{
		ID:           uuid.New(),
		SchoolID:     schoolID,
		ParentUnitID: parentID,
		Name:         req.Name,
		Description:  desc,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.unitRepo.Create(ctx, unit); err != nil {
		s.logger.Error("failed to create unit", "error", err)
		return nil, errors.NewDatabaseError("create unit", err)
	}

	s.logger.Info("unit created", "id", unit.ID.String())
	response := dto.ToUnitResponse(unit)
	return &response, nil
}

func (s *unitService) UpdateUnit(ctx context.Context, id string, req dto.UpdateUnitRequest) (*dto.UnitResponse, error) {
	unitID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	unit, err := s.unitRepo.FindByID(ctx, unitID)
	if err != nil || unit == nil {
		return nil, errors.NewNotFoundError("unit")
	}

	// Actualizar campos
	if req.Name != nil && *req.Name != "" {
		if len(*req.Name) < 2 {
			return nil, errors.NewValidationError("name must be at least 2 characters")
		}
		unit.Name = *req.Name
	}

	if req.Description != nil {
		unit.Description = req.Description
	}

	unit.UpdatedAt = time.Now()

	if err := s.unitRepo.Update(ctx, unit); err != nil {
		s.logger.Error("failed to update unit", "error", err)
		return nil, errors.NewDatabaseError("update unit", err)
	}

	response := dto.ToUnitResponse(unit)
	return &response, nil
}

func (s *unitService) GetUnit(ctx context.Context, id string) (*dto.UnitResponse, error) {
	unitID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	unit, err := s.unitRepo.FindByID(ctx, unitID)
	if err != nil {
		s.logger.Error("failed to find unit", "error", err)
		return nil, errors.NewDatabaseError("find unit", err)
	}
	if unit == nil {
		return nil, errors.NewNotFoundError("unit")
	}

	response := dto.ToUnitResponse(unit)
	return &response, nil
}
