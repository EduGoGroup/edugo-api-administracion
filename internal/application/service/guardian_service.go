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

type GuardianService interface {
	CreateGuardianRelation(ctx context.Context, req dto.CreateGuardianRelationRequest, createdBy string) (*dto.GuardianRelationResponse, error)
	GetGuardianRelation(ctx context.Context, id string) (*dto.GuardianRelationResponse, error)
	GetGuardianRelations(ctx context.Context, guardianID string) ([]*dto.GuardianRelationResponse, error)
	GetStudentGuardians(ctx context.Context, studentID string) ([]*dto.GuardianRelationResponse, error)
}

type guardianService struct {
	guardianRepo repository.GuardianRepository
	logger       logger.Logger
}

func NewGuardianService(
	guardianRepo repository.GuardianRepository,
	logger logger.Logger,
) GuardianService {
	return &guardianService{
		guardianRepo: guardianRepo,
		logger:       logger,
	}
}

func (s *guardianService) CreateGuardianRelation(
	ctx context.Context,
	req dto.CreateGuardianRelationRequest,
	createdBy string,
) (*dto.GuardianRelationResponse, error) {
	// Validar request
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	// Parsear IDs
	guardianID, err := uuid.Parse(req.GuardianID)
	if err != nil {
		return nil, errors.NewValidationError("invalid guardian_id format").WithField("guardian_id", req.GuardianID)
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return nil, errors.NewValidationError("invalid student_id format").WithField("student_id", req.StudentID)
	}

	// Verificar si ya existe una relación activa
	exists, err := s.guardianRepo.ExistsActiveRelation(ctx, guardianID, studentID)
	if err != nil {
		s.logger.Error("failed to check existing relation", "error", err)
		return nil, errors.NewDatabaseError("check relation", err)
	}

	if exists {
		return nil, errors.NewAlreadyExistsError("guardian relation").
			WithField("guardian_id", guardianID.String()).
			WithField("student_id", studentID.String())
	}

	// Crear entidad (lógica de negocio movida aquí)
	now := time.Now()
	relation := &entities.GuardianRelation{
		ID:               uuid.New(),
		GuardianID:       guardianID,
		StudentID:        studentID,
		RelationshipType: req.RelationshipType,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
		CreatedBy:        createdBy,
	}

	// Persistir
	if err := s.guardianRepo.Create(ctx, relation); err != nil {
		s.logger.Error("failed to save guardian relation", "error", err)
		return nil, errors.NewDatabaseError("create relation", err)
	}

	s.logger.Info("guardian relation created", "relation_id", relation.ID.String())

	return dto.ToGuardianRelationResponse(relation), nil
}

func (s *guardianService) GetGuardianRelation(ctx context.Context, id string) (*dto.GuardianRelationResponse, error) {
	relationID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid id format")
	}

	relation, err := s.guardianRepo.FindByID(ctx, relationID)
	if err != nil {
		s.logger.Error("failed to find relation", "error", err, "id", id)
		return nil, errors.NewDatabaseError("find relation", err)
	}

	if relation == nil {
		return nil, errors.NewNotFoundError("guardian relation").WithField("id", id)
	}

	return dto.ToGuardianRelationResponse(relation), nil
}

func (s *guardianService) GetGuardianRelations(ctx context.Context, guardianID string) ([]*dto.GuardianRelationResponse, error) {
	gid, err := uuid.Parse(guardianID)
	if err != nil {
		return nil, errors.NewValidationError("invalid guardian_id format")
	}

	relations, err := s.guardianRepo.FindByGuardian(ctx, gid)
	if err != nil {
		s.logger.Error("failed to find relations", "error", err, "guardian_id", guardianID)
		return nil, errors.NewDatabaseError("find relations", err)
	}

	responses := make([]*dto.GuardianRelationResponse, len(relations))
	for i, relation := range relations {
		responses[i] = dto.ToGuardianRelationResponse(relation)
	}

	return responses, nil
}

func (s *guardianService) GetStudentGuardians(ctx context.Context, studentID string) ([]*dto.GuardianRelationResponse, error) {
	sid, err := uuid.Parse(studentID)
	if err != nil {
		return nil, errors.NewValidationError("invalid student_id format")
	}

	relations, err := s.guardianRepo.FindByStudent(ctx, sid)
	if err != nil {
		s.logger.Error("failed to find relations", "error", err, "student_id", studentID)
		return nil, errors.NewDatabaseError("find relations", err)
	}

	responses := make([]*dto.GuardianRelationResponse, len(relations))
	for i, relation := range relations {
		responses[i] = dto.ToGuardianRelationResponse(relation)
	}

	return responses, nil
}
