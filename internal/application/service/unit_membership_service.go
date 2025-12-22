package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

type UnitMembershipService interface {
	CreateMembership(ctx context.Context, req dto.CreateMembershipRequest) (*dto.MembershipResponse, error)
	GetMembership(ctx context.Context, id string) (*dto.MembershipResponse, error)
	ListMembershipsByUnit(ctx context.Context, unitID string, activeOnly bool) ([]dto.MembershipResponse, error)
	ListMembershipsByUser(ctx context.Context, userID string, activeOnly bool) ([]dto.MembershipResponse, error)
	ListMembershipsByRole(ctx context.Context, unitID string, role string, activeOnly bool) ([]dto.MembershipResponse, error)
	UpdateMembership(ctx context.Context, id string, req dto.UpdateMembershipRequest) (*dto.MembershipResponse, error)
	ExpireMembership(ctx context.Context, id string) error
	DeleteMembership(ctx context.Context, id string) error
}

type unitMembershipService struct {
	membershipRepo repository.UnitMembershipRepository
	unitRepo       repository.AcademicUnitRepository
	logger         logger.Logger
}

func NewUnitMembershipService(
	membershipRepo repository.UnitMembershipRepository,
	unitRepo repository.AcademicUnitRepository,
	logger logger.Logger,
) UnitMembershipService {
	return &unitMembershipService{
		membershipRepo: membershipRepo,
		unitRepo:       unitRepo,
		logger:         logger,
	}
}

func (s *unitMembershipService) CreateMembership(ctx context.Context, req dto.CreateMembershipRequest) (*dto.MembershipResponse, error) {
	// Parse IDs
	unitID, err := uuid.Parse(req.UnitID)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit_id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id")
	}

	// Validar que la unidad existe
	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		return nil, errors.NewDatabaseError("find unit", err)
	}
	if unit == nil {
		return nil, errors.NewNotFoundError("academic unit")
	}

	// Verificar que no existe membresía activa
	exists, err := s.membershipRepo.ExistsByUnitAndUser(ctx, unitID, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("check membership", err)
	}
	if exists {
		return nil, errors.NewAlreadyExistsError("active membership for this user and unit")
	}

	// Validar role usando value object
	if _, err := valueobject.ParseMembershipRole(req.Role); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	// Crear entidad
	now := time.Now()
	enrolledAt := now
	if req.ValidFrom != nil {
		enrolledAt = *req.ValidFrom
	}

	membership := &entities.Membership{
		ID:             uuid.New(),
		UserID:         userID,
		SchoolID:       unit.SchoolID,
		AcademicUnitID: &unitID,
		Role:           req.Role,
		Metadata:       []byte("{}"),
		IsActive:       true,
		EnrolledAt:     enrolledAt,
		WithdrawnAt:    req.ValidUntil,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Persistir
	if err := s.membershipRepo.Create(ctx, membership); err != nil {
		s.logger.Error("failed to create membership", "error", err)
		return nil, errors.NewDatabaseError("create membership", err)
	}

	s.logger.Info("membership created", "id", membership.ID.String())
	response := dto.ToMembershipResponse(membership)
	return &response, nil
}

func (s *unitMembershipService) GetMembership(ctx context.Context, id string) (*dto.MembershipResponse, error) {
	membershipID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid membership ID")
	}

	membership, err := s.membershipRepo.FindByID(ctx, membershipID)
	if err != nil {
		s.logger.Error("failed to find membership", "error", err)
		return nil, errors.NewDatabaseError("find membership", err)
	}
	if membership == nil {
		return nil, errors.NewNotFoundError("membership")
	}

	response := dto.ToMembershipResponse(membership)
	return &response, nil
}

func (s *unitMembershipService) ListMembershipsByUnit(ctx context.Context, unitID string, activeOnly bool) ([]dto.MembershipResponse, error) {
	uid, err := uuid.Parse(unitID)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	memberships, err := s.membershipRepo.FindByUnit(ctx, uid)
	if err != nil {
		return nil, errors.NewDatabaseError("find memberships", err)
	}

	// Filtrar por activeOnly si es necesario
	if activeOnly {
		memberships = filterActiveMemberships(memberships)
	}

	responses := make([]dto.MembershipResponse, len(memberships))
	for i, m := range memberships {
		responses[i] = dto.ToMembershipResponse(m)
	}
	return responses, nil
}

func (s *unitMembershipService) ListMembershipsByUser(ctx context.Context, userID string, activeOnly bool) ([]dto.MembershipResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.NewValidationError("invalid user ID")
	}

	memberships, err := s.membershipRepo.FindByUser(ctx, uid)
	if err != nil {
		return nil, errors.NewDatabaseError("find memberships", err)
	}

	// Filtrar por activeOnly si es necesario
	if activeOnly {
		memberships = filterActiveMemberships(memberships)
	}

	responses := make([]dto.MembershipResponse, len(memberships))
	for i, m := range memberships {
		responses[i] = dto.ToMembershipResponse(m)
	}
	return responses, nil
}

func (s *unitMembershipService) ListMembershipsByRole(ctx context.Context, unitID string, role string, activeOnly bool) ([]dto.MembershipResponse, error) {
	uid, err := uuid.Parse(unitID)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}

	// Validar rol usando value object
	if _, err := valueobject.ParseMembershipRole(role); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	// Usar el método del repositorio que filtra a nivel de base de datos
	memberships, err := s.membershipRepo.FindByUnitAndRole(ctx, uid, role, activeOnly)
	if err != nil {
		return nil, errors.NewDatabaseError("find memberships by role", err)
	}

	responses := make([]dto.MembershipResponse, len(memberships))
	for i, m := range memberships {
		responses[i] = dto.ToMembershipResponse(m)
	}
	return responses, nil
}

func (s *unitMembershipService) UpdateMembership(ctx context.Context, id string, req dto.UpdateMembershipRequest) (*dto.MembershipResponse, error) {
	membershipID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid membership ID")
	}

	membership, err := s.membershipRepo.FindByID(ctx, membershipID)
	if err != nil || membership == nil {
		return nil, errors.NewNotFoundError("membership")
	}

	// Actualizar campos
	if req.Role != nil {
		// Validar rol usando value object
		if _, err := valueobject.ParseMembershipRole(*req.Role); err != nil {
			return nil, errors.NewValidationError(err.Error())
		}
		membership.Role = *req.Role
	}
	if req.ValidUntil != nil {
		membership.WithdrawnAt = req.ValidUntil
	}

	membership.UpdatedAt = time.Now()

	if err := s.membershipRepo.Update(ctx, membership); err != nil {
		s.logger.Error("failed to update membership", "error", err)
		return nil, errors.NewDatabaseError("update membership", err)
	}

	response := dto.ToMembershipResponse(membership)
	return &response, nil
}

func (s *unitMembershipService) ExpireMembership(ctx context.Context, id string) error {
	membershipID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid membership ID")
	}

	membership, err := s.membershipRepo.FindByID(ctx, membershipID)
	if err != nil || membership == nil {
		return errors.NewNotFoundError("membership")
	}

	now := time.Now()
	membership.WithdrawnAt = &now
	membership.UpdatedAt = now

	if err := s.membershipRepo.Update(ctx, membership); err != nil {
		s.logger.Error("failed to expire membership", "error", err)
		return errors.NewDatabaseError("expire membership", err)
	}

	return nil
}

func (s *unitMembershipService) DeleteMembership(ctx context.Context, id string) error {
	membershipID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid membership ID")
	}

	if err := s.membershipRepo.Delete(ctx, membershipID); err != nil {
		s.logger.Error("failed to delete membership", "error", err)
		return errors.NewDatabaseError("delete membership", err)
	}

	return nil
}

// filterActiveMemberships filtra membresías para retornar solo las activas
func filterActiveMemberships(memberships []*entities.Membership) []*entities.Membership {
	result := make([]*entities.Membership, 0, len(memberships))
	for _, m := range memberships {
		if m.IsActive && m.WithdrawnAt == nil {
			result = append(result, m)
		}
	}
	return result
}
