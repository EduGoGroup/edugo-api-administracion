package service

import (
	"context"

	"github.com/edugo/api-administracion/internal/application/dto"
	"github.com/edugo/api-administracion/internal/domain/entity"
	"github.com/edugo/api-administracion/internal/domain/repository"
	"github.com/edugo/api-administracion/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

// UnitService define las operaciones de negocio para unidades
type UnitService interface {
	CreateUnit(ctx context.Context, req dto.CreateUnitRequest) (*dto.UnitResponse, error)
	UpdateUnit(ctx context.Context, id string, req dto.UpdateUnitRequest) (*dto.UnitResponse, error)
	AssignMember(ctx context.Context, unitID string, req dto.AssignMemberRequest) error
}

type unitService struct {
	unitRepo repository.UnitRepository
	logger   logger.Logger
}

func NewUnitService(unitRepo repository.UnitRepository, logger logger.Logger) UnitService {
	return &unitService{
		unitRepo: unitRepo,
		logger:   logger,
	}
}

func (s *unitService) CreateUnit(ctx context.Context, req dto.CreateUnitRequest) (*dto.UnitResponse, error) {
	// Validar request
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	// Parsear school ID
	schoolID, err := valueobject.SchoolIDFromString(req.SchoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school_id format")
	}

	// Parsear parent unit ID (si existe)
	var parentUnitID *valueobject.UnitID
	if req.ParentUnitID != nil && *req.ParentUnitID != "" {
		pID, err := valueobject.UnitIDFromString(*req.ParentUnitID)
		if err != nil {
			return nil, errors.NewValidationError("invalid parent_unit_id format")
		}
		parentUnitID = &pID
	}

	// Crear entidad de dominio
	unit, err := entity.NewUnit(schoolID, parentUnitID, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	// Persistir
	if err := s.unitRepo.Create(ctx, unit); err != nil {
		s.logger.Error("failed to save unit", "error", err)
		return nil, errors.NewDatabaseError("create unit", err)
	}

	s.logger.Info("unit created",
		"unit_id", unit.ID().String(),
		"school_id", schoolID.String(),
		"name", unit.Name(),
	)

	return dto.ToUnitResponse(unit), nil
}

func (s *unitService) UpdateUnit(ctx context.Context, id string, req dto.UpdateUnitRequest) (*dto.UnitResponse, error) {
	// Validar request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Parsear unit ID
	unitID, err := valueobject.UnitIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit_id format")
	}

	// Buscar unidad
	unit, err := s.unitRepo.FindByID(ctx, unitID)
	if err != nil || unit == nil {
		return nil, errors.NewNotFoundError("unit")
	}

	// Actualizar
	if err := unit.UpdateInfo(req.Name, req.Description); err != nil {
		return nil, err
	}

	// Persistir
	if err := s.unitRepo.Update(ctx, unit); err != nil {
		s.logger.Error("failed to update unit", "error", err)
		return nil, errors.NewDatabaseError("update unit", err)
	}

	s.logger.Info("unit updated", "unit_id", unitID.String())

	return dto.ToUnitResponse(unit), nil
}

func (s *unitService) AssignMember(ctx context.Context, unitID string, req dto.AssignMemberRequest) error {
	// Validar request
	if err := req.Validate(); err != nil {
		return err
	}

	// Parsear IDs
	uID, err := valueobject.UnitIDFromString(unitID)
	if err != nil {
		return errors.NewValidationError("invalid unit_id format")
	}

	userID, err := valueobject.UserIDFromString(req.UserID)
	if err != nil {
		return errors.NewValidationError("invalid user_id format")
	}

	// Verificar que la unidad existe
	unit, err := s.unitRepo.FindByID(ctx, uID)
	if err != nil || unit == nil {
		return errors.NewNotFoundError("unit")
	}

	// Verificar si ya es miembro
	isMember, err := s.unitRepo.IsMember(ctx, uID, userID)
	if err != nil {
		s.logger.Error("failed to check membership", "error", err)
		return errors.NewDatabaseError("check member", err)
	}

	if isMember {
		return errors.NewAlreadyExistsError("unit member").
			WithField("unit_id", unitID).
			WithField("user_id", req.UserID)
	}

	// Agregar miembro
	if err := s.unitRepo.AddMember(ctx, uID, userID, req.Role); err != nil {
		s.logger.Error("failed to add member", "error", err)
		return errors.NewDatabaseError("add member", err)
	}

	s.logger.Info("member assigned",
		"unit_id", unitID,
		"user_id", req.UserID,
		"role", req.Role,
	)

	return nil
}
