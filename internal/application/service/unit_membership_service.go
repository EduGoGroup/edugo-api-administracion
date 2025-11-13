package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// UnitMembershipService define las operaciones de negocio para UnitMembership
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

// NewUnitMembershipService crea un nuevo UnitMembershipService
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

// CreateMembership crea una nueva membresía
func (s *unitMembershipService) CreateMembership(ctx context.Context, req dto.CreateMembershipRequest) (*dto.MembershipResponse, error) {
	// 1. Validar que la unidad existe
	unitID, err := valueobject.UnitIDFromString(req.UnitID)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit_id")
	}

	unit, err := s.unitRepo.FindByID(ctx, unitID, false)
	if err != nil {
		return nil, errors.NewDatabaseError("find unit", err)
	}
	if unit == nil {
		return nil, errors.NewNotFoundError("academic unit")
	}

	// 2. Validar user_id
	userID, err := valueobject.UserIDFromString(req.UserID)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id")
	}

	// 3. Verificar que no existe membresía activa
	exists, err := s.membershipRepo.ExistsByUnitAndUser(ctx, unitID, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("check membership", err)
	}
	if exists {
		return nil, errors.NewAlreadyExistsError("active membership for this user and unit")
	}

	// 4. Validar rol
	role, err := valueobject.NewMembershipRole(req.Role)
	if err != nil {
		return nil, err
	}

	// 5. Crear entidad
	validFrom := time.Now()
	if req.ValidFrom != nil {
		validFrom = *req.ValidFrom
	}

	membership, err := entity.NewUnitMembership(unitID, userID, role, validFrom)
	if err != nil {
		s.logger.Warn("failed to create membership entity", "error", err)
		return nil, err
	}

	// 6. Establecer validUntil si se proporciona
	if req.ValidUntil != nil {
		if err := membership.SetValidUntil(*req.ValidUntil); err != nil {
			return nil, err
		}
	}

	// 7. Agregar metadata
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			membership.SetMetadata(key, value)
		}
	}

	// 8. Persistir
	if err := s.membershipRepo.Create(ctx, membership); err != nil {
		s.logger.Error("failed to create membership", "error", err, "unit_id", req.UnitID, "user_id", req.UserID)
		return nil, errors.NewDatabaseError("create membership", err)
	}

	s.logger.Info("membership created successfully", "id", membership.ID().String(), "role", req.Role)

	response := dto.ToMembershipResponse(membership)
	return &response, nil
}

// GetMembership obtiene una membresía por ID
func (s *unitMembershipService) GetMembership(ctx context.Context, id string) (*dto.MembershipResponse, error) {
	membershipID, err := valueobject.MembershipIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid membership ID")
	}

	membership, err := s.membershipRepo.FindByID(ctx, membershipID)
	if err != nil {
		s.logger.Error("failed to find membership", "error", err, "id", id)
		return nil, errors.NewDatabaseError("find membership", err)
	}

	if membership == nil {
		return nil, errors.NewNotFoundError("membership")
	}

	response := dto.ToMembershipResponse(membership)
	return &response, nil
}

// ListMembershipsByUnit lista membresías de una unidad
func (s *unitMembershipService) ListMembershipsByUnit(ctx context.Context, unitID string, activeOnly bool) ([]dto.MembershipResponse, error) {
	unitIDVO, err := valueobject.UnitIDFromString(unitID)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit_id")
	}

	memberships, err := s.membershipRepo.FindByUnitID(ctx, unitIDVO, activeOnly)
	if err != nil {
		s.logger.Error("failed to list memberships", "error", err, "unit_id", unitID)
		return nil, errors.NewDatabaseError("list memberships", err)
	}

	return dto.ToMembershipResponseList(memberships), nil
}

// ListMembershipsByUser lista membresías de un usuario
func (s *unitMembershipService) ListMembershipsByUser(ctx context.Context, userID string, activeOnly bool) ([]dto.MembershipResponse, error) {
	userIDVO, err := valueobject.UserIDFromString(userID)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id")
	}

	memberships, err := s.membershipRepo.FindByUserID(ctx, userIDVO, activeOnly)
	if err != nil {
		s.logger.Error("failed to list memberships", "error", err, "user_id", userID)
		return nil, errors.NewDatabaseError("list memberships", err)
	}

	return dto.ToMembershipResponseList(memberships), nil
}

// ListMembershipsByRole lista membresías por rol
func (s *unitMembershipService) ListMembershipsByRole(ctx context.Context, unitID string, role string, activeOnly bool) ([]dto.MembershipResponse, error) {
	unitIDVO, err := valueobject.UnitIDFromString(unitID)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit_id")
	}

	roleVO, err := valueobject.NewMembershipRole(role)
	if err != nil {
		return nil, err
	}

	memberships, err := s.membershipRepo.FindByRole(ctx, unitIDVO, roleVO, activeOnly)
	if err != nil {
		s.logger.Error("failed to list memberships by role", "error", err, "role", role)
		return nil, errors.NewDatabaseError("list memberships", err)
	}

	return dto.ToMembershipResponseList(memberships), nil
}

// UpdateMembership actualiza una membresía
func (s *unitMembershipService) UpdateMembership(ctx context.Context, id string, req dto.UpdateMembershipRequest) (*dto.MembershipResponse, error) {
	membershipID, err := valueobject.MembershipIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid membership ID")
	}

	membership, err := s.membershipRepo.FindByID(ctx, membershipID)
	if err != nil {
		return nil, errors.NewDatabaseError("find membership", err)
	}
	if membership == nil {
		return nil, errors.NewNotFoundError("membership")
	}

	// Cambiar rol si se proporciona
	if req.Role != nil {
		roleVO, err := valueobject.NewMembershipRole(*req.Role)
		if err != nil {
			return nil, err
		}

		if err := membership.ChangeRole(roleVO); err != nil {
			return nil, err
		}
	}

	// Actualizar validUntil
	if req.ValidUntil != nil {
		if err := membership.SetValidUntil(*req.ValidUntil); err != nil {
			return nil, err
		}
	}

	// Actualizar metadata
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			membership.SetMetadata(key, value)
		}
	}

	if err := s.membershipRepo.Update(ctx, membership); err != nil {
		s.logger.Error("failed to update membership", "error", err, "id", id)
		return nil, errors.NewDatabaseError("update membership", err)
	}

	s.logger.Info("membership updated successfully", "id", id)

	response := dto.ToMembershipResponse(membership)
	return &response, nil
}

// ExpireMembership marca una membresía como expirada
func (s *unitMembershipService) ExpireMembership(ctx context.Context, id string) error {
	membershipID, err := valueobject.MembershipIDFromString(id)
	if err != nil {
		return errors.NewValidationError("invalid membership ID")
	}

	membership, err := s.membershipRepo.FindByID(ctx, membershipID)
	if err != nil {
		return errors.NewDatabaseError("find membership", err)
	}
	if membership == nil {
		return errors.NewNotFoundError("membership")
	}

	if err := membership.Expire(); err != nil {
		return err
	}

	if err := s.membershipRepo.Update(ctx, membership); err != nil {
		s.logger.Error("failed to expire membership", "error", err, "id", id)
		return errors.NewDatabaseError("expire membership", err)
	}

	s.logger.Info("membership expired successfully", "id", id)
	return nil
}

// DeleteMembership elimina una membresía permanentemente
func (s *unitMembershipService) DeleteMembership(ctx context.Context, id string) error {
	membershipID, err := valueobject.MembershipIDFromString(id)
	if err != nil {
		return errors.NewValidationError("invalid membership ID")
	}

	// Verificar que existe
	membership, err := s.membershipRepo.FindByID(ctx, membershipID)
	if err != nil {
		return errors.NewDatabaseError("find membership", err)
	}
	if membership == nil {
		return errors.NewNotFoundError("membership")
	}

	if err := s.membershipRepo.Delete(ctx, membershipID); err != nil {
		s.logger.Error("failed to delete membership", "error", err, "id", id)
		return errors.NewDatabaseError("delete membership", err)
	}

	s.logger.Info("membership deleted successfully", "id", id)
	return nil
}
