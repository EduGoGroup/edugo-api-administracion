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
	if err != nil || subject == nil {
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
