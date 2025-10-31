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

// GuardianService define las operaciones de negocio para guardians
type GuardianService interface {
	// CreateGuardianRelation crea una nueva relación guardian-estudiante
	CreateGuardianRelation(ctx context.Context, req dto.CreateGuardianRelationRequest, createdBy string) (*dto.GuardianRelationResponse, error)

	// GetGuardianRelation obtiene una relación por ID
	GetGuardianRelation(ctx context.Context, id string) (*dto.GuardianRelationResponse, error)

	// GetGuardianRelations obtiene todas las relaciones de un guardian
	GetGuardianRelations(ctx context.Context, guardianID string) ([]*dto.GuardianRelationResponse, error)

	// GetStudentGuardians obtiene todos los guardians de un estudiante
	GetStudentGuardians(ctx context.Context, studentID string) ([]*dto.GuardianRelationResponse, error)
}

// guardianService implementa GuardianService
type guardianService struct {
	guardianRepo repository.GuardianRepository
	logger       logger.Logger
}

// NewGuardianService crea un nuevo GuardianService
func NewGuardianService(
	guardianRepo repository.GuardianRepository,
	logger logger.Logger,
) GuardianService {
	return &guardianService{
		guardianRepo: guardianRepo,
		logger:       logger,
	}
}

// CreateGuardianRelation implementa la creación de una relación
func (s *guardianService) CreateGuardianRelation(
	ctx context.Context,
	req dto.CreateGuardianRelationRequest,
	createdBy string,
) (*dto.GuardianRelationResponse, error) {
	// 1. Validar request
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	// 2. Convertir IDs a value objects
	guardianID, err := valueobject.GuardianIDFromString(req.GuardianID)
	if err != nil {
		return nil, errors.NewValidationError("invalid guardian_id format").
			WithField("guardian_id", req.GuardianID)
	}

	studentID, err := valueobject.StudentIDFromString(req.StudentID)
	if err != nil {
		return nil, errors.NewValidationError("invalid student_id format").
			WithField("student_id", req.StudentID)
	}

	// 3. Crear relationship type
	relationshipType, err := valueobject.NewRelationshipType(req.RelationshipType)
	if err != nil {
		return nil, err
	}

	// 4. Verificar si ya existe una relación activa
	exists, err := s.guardianRepo.ExistsActiveRelation(ctx, guardianID, studentID)
	if err != nil {
		s.logger.Error("failed to check existing relation",
			"error", err,
			"guardian_id", guardianID.String(),
			"student_id", studentID.String(),
		)
		return nil, errors.NewDatabaseError("check relation", err)
	}

	if exists {
		return nil, errors.NewAlreadyExistsError("guardian relation").
			WithField("guardian_id", guardianID.String()).
			WithField("student_id", studentID.String())
	}

	// 5. Crear entidad de dominio
	relation, err := entity.NewGuardianRelation(
		guardianID,
		studentID,
		relationshipType,
		createdBy,
	)
	if err != nil {
		s.logger.Warn("failed to create guardian relation entity", "error", err)
		return nil, err
	}

	// 6. Persistir en repositorio
	if err := s.guardianRepo.Create(ctx, relation); err != nil {
		s.logger.Error("failed to save guardian relation",
			"error", err,
			"guardian_id", guardianID.String(),
			"student_id", studentID.String(),
		)
		return nil, errors.NewDatabaseError("create relation", err)
	}

	s.logger.Info("guardian relation created",
		"relation_id", relation.ID().String(),
		"guardian_id", guardianID.String(),
		"student_id", studentID.String(),
		"relationship_type", relationshipType.String(),
	)

	// 7. Retornar DTO de respuesta
	return dto.ToGuardianRelationResponse(relation), nil
}

// GetGuardianRelation obtiene una relación por ID
func (s *guardianService) GetGuardianRelation(ctx context.Context, id string) (*dto.GuardianRelationResponse, error) {
	// Parsear UUID
	uuid, err := valueobject.GuardianIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid id format")
	}

	// Buscar en repositorio
	relation, err := s.guardianRepo.FindByID(ctx, uuid.UUID())
	if err != nil {
		s.logger.Error("failed to find relation", "error", err, "id", id)
		return nil, errors.NewDatabaseError("find relation", err)
	}

	if relation == nil {
		return nil, errors.NewNotFoundError("guardian relation").
			WithField("id", id)
	}

	return dto.ToGuardianRelationResponse(relation), nil
}

// GetGuardianRelations obtiene todas las relaciones de un guardian
func (s *guardianService) GetGuardianRelations(ctx context.Context, guardianID string) ([]*dto.GuardianRelationResponse, error) {
	// Parsear UUID
	gid, err := valueobject.GuardianIDFromString(guardianID)
	if err != nil {
		return nil, errors.NewValidationError("invalid guardian_id format")
	}

	// Buscar en repositorio
	relations, err := s.guardianRepo.FindByGuardian(ctx, gid)
	if err != nil {
		s.logger.Error("failed to find relations", "error", err, "guardian_id", guardianID)
		return nil, errors.NewDatabaseError("find relations", err)
	}

	// Convertir a DTOs
	responses := make([]*dto.GuardianRelationResponse, len(relations))
	for i, relation := range relations {
		responses[i] = dto.ToGuardianRelationResponse(relation)
	}

	return responses, nil
}

// GetStudentGuardians obtiene todos los guardians de un estudiante
func (s *guardianService) GetStudentGuardians(ctx context.Context, studentID string) ([]*dto.GuardianRelationResponse, error) {
	// Parsear UUID
	sid, err := valueobject.StudentIDFromString(studentID)
	if err != nil {
		return nil, errors.NewValidationError("invalid student_id format")
	}

	// Buscar en repositorio
	relations, err := s.guardianRepo.FindByStudent(ctx, sid)
	if err != nil {
		s.logger.Error("failed to find relations", "error", err, "student_id", studentID)
		return nil, errors.NewDatabaseError("find relations", err)
	}

	// Convertir a DTOs
	responses := make([]*dto.GuardianRelationResponse, len(relations))
	for i, relation := range relations {
		responses[i] = dto.ToGuardianRelationResponse(relation)
	}

	return responses, nil
}
