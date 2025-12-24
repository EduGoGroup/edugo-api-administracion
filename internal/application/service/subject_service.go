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

type SubjectService interface {
	CreateSubject(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error)
	UpdateSubject(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error)
	GetSubject(ctx context.Context, id string) (*dto.SubjectResponse, error)
	ListSubjects(ctx context.Context, schoolID string) ([]dto.SubjectResponse, error)
	DeleteSubject(ctx context.Context, id string) error
}

type subjectService struct {
	subjectRepo repository.SubjectRepository
	logger      logger.Logger
}

func NewSubjectService(subjectRepo repository.SubjectRepository, logger logger.Logger) SubjectService {
	return &subjectService{subjectRepo: subjectRepo, logger: logger}
}

func (s *subjectService) CreateSubject(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	// Validar
	if req.Name == "" {
		return nil, errors.NewValidationError("name is required")
	}

	// Crear entidad
	now := time.Now()
	desc := &req.Description
	meta := &req.Metadata
	subject := &entities.Subject{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: desc,
		Metadata:    meta,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.subjectRepo.Create(ctx, subject); err != nil {
		s.logger.Error("failed to create subject", "error", err)
		return nil, errors.NewDatabaseError("create subject", err)
	}

	s.logger.Info("subject created", "id", subject.ID.String())
	response := dto.ToSubjectResponse(subject)
	return &response, nil
}

func (s *subjectService) UpdateSubject(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid subject ID")
	}

	subject, err := s.subjectRepo.FindByID(ctx, subjectID)
	if err != nil {
		s.logger.Error("database error",
			"operation", "find_subject",
			"subject_id", subjectID,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError("find subject", err)
	}
	if subject == nil {
		return nil, errors.NewNotFoundError("subject")
	}

	// Actualizar campos
	if req.Name != nil && *req.Name != "" {
		subject.Name = *req.Name
	}
	if req.Description != nil {
		subject.Description = req.Description
	}
	if req.Metadata != nil {
		subject.Metadata = req.Metadata
	}

	subject.UpdatedAt = time.Now()

	if err := s.subjectRepo.Update(ctx, subject); err != nil {
		s.logger.Error("failed to update subject", "error", err)
		return nil, errors.NewDatabaseError("update subject", err)
	}

	s.logger.Info("subject updated", "id", subject.ID.String())
	response := dto.ToSubjectResponse(subject)
	return &response, nil
}

func (s *subjectService) GetSubject(ctx context.Context, id string) (*dto.SubjectResponse, error) {
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid subject ID")
	}

	subject, err := s.subjectRepo.FindByID(ctx, subjectID)
	if err != nil {
		s.logger.Error("failed to find subject", "error", err, "id", id)
		return nil, errors.NewDatabaseError("find subject", err)
	}
	if subject == nil {
		return nil, errors.NewNotFoundError("subject")
	}

	response := dto.ToSubjectResponse(subject)
	return &response, nil
}

func (s *subjectService) ListSubjects(ctx context.Context, schoolID string) ([]dto.SubjectResponse, error) {
	var subjects []*entities.Subject
	var err error

	if schoolID != "" {
		// Validar UUID de school
		schoolUUID, parseErr := uuid.Parse(schoolID)
		if parseErr != nil {
			return nil, errors.NewValidationError("invalid school ID")
		}
		subjects, err = s.subjectRepo.FindBySchoolID(ctx, schoolUUID)
	} else {
		subjects, err = s.subjectRepo.List(ctx)
	}

	if err != nil {
		s.logger.Error("failed to list subjects", "error", err, "school_id", schoolID)
		return nil, errors.NewDatabaseError("list subjects", err)
	}

	responses := make([]dto.SubjectResponse, len(subjects))
	for i, subject := range subjects {
		responses[i] = dto.ToSubjectResponse(subject)
	}

	return responses, nil
}

func (s *subjectService) DeleteSubject(ctx context.Context, id string) error {
	subjectID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid subject ID")
	}

	// Verificar que existe antes de eliminar
	subject, err := s.subjectRepo.FindByID(ctx, subjectID)
	if err != nil {
		s.logger.Error("failed to find subject for deletion", "error", err, "id", id)
		return errors.NewDatabaseError("find subject", err)
	}
	if subject == nil {
		return errors.NewNotFoundError("subject")
	}

	if err := s.subjectRepo.Delete(ctx, subjectID); err != nil {
		s.logger.Error("failed to delete subject", "error", err, "id", id)
		return errors.NewDatabaseError("delete subject", err)
	}

	s.logger.Info("subject deleted", "id", id)
	return nil
}
