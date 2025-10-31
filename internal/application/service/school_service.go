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

// SchoolService define las operaciones de negocio para escuelas
type SchoolService interface {
	CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error)
	GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error)
}

type schoolService struct {
	schoolRepo repository.SchoolRepository
	logger     logger.Logger
}

// NewSchoolService crea un nuevo SchoolService
func NewSchoolService(schoolRepo repository.SchoolRepository, logger logger.Logger) SchoolService {
	return &schoolService{
		schoolRepo: schoolRepo,
		logger:     logger,
	}
}

// CreateSchool crea una nueva escuela
func (s *schoolService) CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error) {
	// 1. Validar request
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	// 2. Verificar si ya existe
	exists, err := s.schoolRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		s.logger.Error("failed to check existing school", "error", err, "name", req.Name)
		return nil, errors.NewDatabaseError("check school", err)
	}

	if exists {
		return nil, errors.NewAlreadyExistsError("school").WithField("name", req.Name)
	}

	// 3. Crear entidad de dominio
	school, err := entity.NewSchool(req.Name, req.Address)
	if err != nil {
		return nil, err
	}

	// 4. Persistir
	if err := s.schoolRepo.Create(ctx, school); err != nil {
		s.logger.Error("failed to save school", "error", err, "name", req.Name)
		return nil, errors.NewDatabaseError("create school", err)
	}

	s.logger.Info("school created", "school_id", school.ID().String(), "name", school.Name())

	return dto.ToSchoolResponse(school), nil
}

// GetSchool obtiene una escuela por ID
func (s *schoolService) GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error) {
	schoolID, err := valueobject.SchoolIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid school_id format")
	}

	school, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		s.logger.Error("failed to find school", "error", err, "id", id)
		return nil, errors.NewDatabaseError("find school", err)
	}

	if school == nil {
		return nil, errors.NewNotFoundError("school").WithField("id", id)
	}

	return dto.ToSchoolResponse(school), nil
}
