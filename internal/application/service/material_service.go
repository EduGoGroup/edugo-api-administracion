package service

import (
	"context"

	"github.com/edugo/api-administracion/internal/domain/repository"
	"github.com/edugo/api-administracion/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

// MaterialService define las operaciones de negocio para materiales
// Nota: Solo DELETE en API Admin, los demás endpoints están en API Mobile
type MaterialService interface {
	DeleteMaterial(ctx context.Context, id string) error
}

type materialService struct {
	materialRepo repository.MaterialRepository
	logger       logger.Logger
}

func NewMaterialService(materialRepo repository.MaterialRepository, logger logger.Logger) MaterialService {
	return &materialService{
		materialRepo: materialRepo,
		logger:       logger,
	}
}

func (s *materialService) DeleteMaterial(ctx context.Context, id string) error {
	// Validar ID
	materialID, err := valueobject.MaterialIDFromString(id)
	if err != nil {
		return errors.NewValidationError("invalid material_id format")
	}

	// Verificar que existe
	exists, err := s.materialRepo.Exists(ctx, materialID)
	if err != nil {
		s.logger.Error("failed to check material", "error", err, "id", id)
		return errors.NewDatabaseError("check material", err)
	}

	if !exists {
		return errors.NewNotFoundError("material").WithField("id", id)
	}

	// Eliminar (soft delete)
	if err := s.materialRepo.Delete(ctx, materialID); err != nil {
		s.logger.Error("failed to delete material", "error", err, "id", id)
		return errors.NewDatabaseError("delete material", err)
	}

	s.logger.Info("material deleted", "material_id", materialID.String())
	return nil
}
