package service

import (
	"context"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/google/uuid"
)

type HierarchyService struct {
	unitRepo   repository.AcademicUnitRepository
	schoolRepo repository.SchoolRepository
}

func NewHierarchyService(
	unitRepo repository.AcademicUnitRepository,
	schoolRepo repository.SchoolRepository,
) *HierarchyService {
	return &HierarchyService{
		unitRepo:   unitRepo,
		schoolRepo: schoolRepo,
	}
}

func (s *HierarchyService) CreateUnit(
	ctx context.Context,
	parentUnitID *uuid.UUID,
	schoolID uuid.UUID,
	unitType string,
	name string,
	code string,
	description string,
) (*entities.AcademicUnit, error) {
	// Validar que la escuela existe
	_, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("school %s not found", schoolID))
	}

	// Validar que el código no esté duplicado
	exists, err := s.unitRepo.ExistsBySchoolIDAndCode(ctx, schoolID, code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.NewValidationError(fmt.Sprintf("code %s already exists in school", code))
	}

	// Si tiene padre, validar que existe
	if parentUnitID != nil {
		parent, err := s.unitRepo.FindByID(ctx, *parentUnitID, false)
		if err != nil {
			return nil, errors.NewNotFoundError(fmt.Sprintf("parent unit %s not found", parentUnitID))
		}

		// Validar que el padre pertenece a la misma escuela
		if parent.SchoolID != schoolID {
			return nil, errors.NewValidationError("parent unit must belong to the same school")
		}
	}

	// Crear la unidad
	now := time.Now()
	desc := &description
	unit := &entities.AcademicUnit{
		ID:           uuid.New(),
		ParentUnitID: parentUnitID,
		SchoolID:     schoolID,
		Name:         name,
		Code:         code,
		Type:         unitType,
		Description:  desc,
		Level:        nil,
		AcademicYear: 0,
		Metadata:     []byte("{}"),
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    nil,
	}

	// Persistir
	if err := s.unitRepo.Create(ctx, unit); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *HierarchyService) GetUnitTree(ctx context.Context, unitID uuid.UUID) (*entities.AcademicUnit, []*entities.AcademicUnit, error) {
	// Obtener la unidad raíz
	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		return nil, nil, err
	}
	if unit == nil {
		return nil, nil, errors.NewNotFoundError("unit")
	}

	// Obtener todos los descendientes (implementación simple)
	// TODO: mejorar con query recursivo o ltree
	descendants, err := s.unitRepo.FindBySchoolID(ctx, unit.SchoolID, false)
	if err != nil {
		return nil, nil, err
	}

	return unit, descendants, nil
}
