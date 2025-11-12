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

// SchoolService define las operaciones de negocio para School
type SchoolService interface {
	CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error)
	GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error)
	GetSchoolByCode(ctx context.Context, code string) (*dto.SchoolResponse, error)
	UpdateSchool(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error)
	ListSchools(ctx context.Context) ([]dto.SchoolResponse, error)
	DeleteSchool(ctx context.Context, id string) error
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
	// 1. Verificar si ya existe por código
	exists, err := s.schoolRepo.ExistsByCode(ctx, req.Code)
	if err != nil {
		s.logger.Error("failed to check existing school", "error", err, "code", req.Code)
		return nil, errors.NewDatabaseError("check school", err)
	}

	if exists {
		return nil, errors.NewAlreadyExistsError("school").WithField("code", req.Code)
	}

	// 2. Crear entidad de dominio
	school, err := entity.NewSchool(req.Name, req.Code, req.Address)
	if err != nil {
		s.logger.Warn("failed to create school entity", "error", err)
		return nil, err
	}

	// 3. Agregar contacto si se proporciona
	if req.ContactEmail != "" || req.ContactPhone != "" {
		var email *valueobject.Email
		if req.ContactEmail != "" {
			e, err := valueobject.NewEmail(req.ContactEmail)
			if err != nil {
				return nil, err
			}
			email = &e
		}

		if err := school.UpdateContactInfo(email, req.ContactPhone); err != nil {
			return nil, err
		}
	}

	// 4. Agregar metadata si se proporciona
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			school.SetMetadata(key, value)
		}
	}

	// 5. Persistir
	if err := s.schoolRepo.Create(ctx, school); err != nil {
		s.logger.Error("failed to create school", "error", err, "name", req.Name)
		return nil, errors.NewDatabaseError("create school", err)
	}

	s.logger.Info("school created successfully", "id", school.ID().String(), "name", req.Name)

	response := dto.ToSchoolResponse(school)
	return &response, nil
}

// GetSchool obtiene una escuela por ID
func (s *schoolService) GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error) {
	schoolID, err := valueobject.SchoolIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	school, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		s.logger.Error("failed to find school", "error", err, "id", id)
		return nil, errors.NewDatabaseError("find school", err)
	}

	if school == nil {
		return nil, errors.NewNotFoundError("school")
	}

	response := dto.ToSchoolResponse(school)
	return &response, nil
}

// GetSchoolByCode obtiene una escuela por código
func (s *schoolService) GetSchoolByCode(ctx context.Context, code string) (*dto.SchoolResponse, error) {
	school, err := s.schoolRepo.FindByCode(ctx, code)
	if err != nil {
		s.logger.Error("failed to find school by code", "error", err, "code", code)
		return nil, errors.NewDatabaseError("find school", err)
	}

	if school == nil {
		return nil, errors.NewNotFoundError("school")
	}

	response := dto.ToSchoolResponse(school)
	return &response, nil
}

// UpdateSchool actualiza una escuela
func (s *schoolService) UpdateSchool(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error) {
	schoolID, err := valueobject.SchoolIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	school, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		return nil, errors.NewDatabaseError("find school", err)
	}

	if school == nil {
		return nil, errors.NewNotFoundError("school")
	}

	// Actualizar info básica si se proporciona
	if req.Name != nil || req.Address != nil {
		name := ""
		if req.Name != nil {
			name = *req.Name
		}
		address := ""
		if req.Address != nil {
			address = *req.Address
		}

		if err := school.UpdateInfo(name, address); err != nil {
			return nil, err
		}
	}

	// Actualizar contacto si se proporciona
	if req.ContactEmail != nil || req.ContactPhone != nil {
		var email *valueobject.Email
		if req.ContactEmail != nil && *req.ContactEmail != "" {
			e, err := valueobject.NewEmail(*req.ContactEmail)
			if err != nil {
				return nil, err
			}
			email = &e
		}

		phone := ""
		if req.ContactPhone != nil {
			phone = *req.ContactPhone
		}

		if err := school.UpdateContactInfo(email, phone); err != nil {
			return nil, err
		}
	}

	// Actualizar metadata si se proporciona
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			school.SetMetadata(key, value)
		}
	}

	if err := s.schoolRepo.Update(ctx, school); err != nil {
		s.logger.Error("failed to update school", "error", err, "id", id)
		return nil, errors.NewDatabaseError("update school", err)
	}

	s.logger.Info("school updated successfully", "id", id)

	response := dto.ToSchoolResponse(school)
	return &response, nil
}

// ListSchools lista todas las escuelas
func (s *schoolService) ListSchools(ctx context.Context) ([]dto.SchoolResponse, error) {
	schools, err := s.schoolRepo.List(ctx, repository.ListFilters{})
	if err != nil {
		s.logger.Error("failed to list schools", "error", err)
		return nil, errors.NewDatabaseError("list schools", err)
	}

	return dto.ToSchoolResponseList(schools), nil
}

// DeleteSchool elimina una escuela
// DeleteSchool elimina una escuela
func (s *schoolService) DeleteSchool(ctx context.Context, id string) error {
	schoolID, err := valueobject.SchoolIDFromString(id)
	if err != nil {
		return errors.NewValidationError("invalid school ID")
	}

	// Verificar que la escuela existe antes de eliminar
	school, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		s.logger.Error("failed to find school", "error", err, "id", id)
		return errors.NewDatabaseError("find school", err)
	}

	if school == nil {
		return errors.NewNotFoundError("school")
	}

	if err := s.schoolRepo.Delete(ctx, schoolID); err != nil {
		s.logger.Error("failed to delete school", "error", err, "id", id)
		return errors.NewDatabaseError("delete school", err)
	}

	s.logger.Info("school deleted successfully", "id", id)
	return nil
}
