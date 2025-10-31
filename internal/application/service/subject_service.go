package service

import (
	"context"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/logger"
)

// SubjectService define las operaciones de negocio para materias
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
	return &subjectService{
		subjectRepo: subjectRepo,
		logger:      logger,
	}
}

func (s *subjectService) CreateSubject(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	subject, err := entity.NewSubject(req.Name, req.Description, req.Metadata)
	if err != nil {
		return nil, err
	}

	if err := s.subjectRepo.Create(ctx, subject); err != nil {
		s.logger.Error("failed to save subject", "error", err)
		return nil, errors.NewDatabaseError("create subject", err)
	}

	s.logger.Info("subject created", "subject_id", subject.ID().String(), "name", subject.Name())
	return dto.ToSubjectResponse(subject), nil
}

func (s *subjectService) UpdateSubject(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	subjectID, err := valueobject.SubjectIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid subject_id format")
	}

	subject, err := s.subjectRepo.FindByID(ctx, subjectID)
	if err != nil || subject == nil {
		return nil, errors.NewNotFoundError("subject")
	}

	if err := subject.UpdateInfo(req.Name, req.Description, req.Metadata); err != nil {
		return nil, err
	}

	if err := s.subjectRepo.Update(ctx, subject); err != nil {
		s.logger.Error("failed to update subject", "error", err)
		return nil, errors.NewDatabaseError("update subject", err)
	}

	s.logger.Info("subject updated", "subject_id", subject.ID().String())
	return dto.ToSubjectResponse(subject), nil
}

func (s *subjectService) GetSubject(ctx context.Context, id string) (*dto.SubjectResponse, error) {
	subjectID, err := valueobject.SubjectIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid subject_id format")
	}

	subject, err := s.subjectRepo.FindByID(ctx, subjectID)
	if err != nil || subject == nil {
		return nil, errors.NewNotFoundError("subject")
	}

	return dto.ToSubjectResponse(subject), nil
}
