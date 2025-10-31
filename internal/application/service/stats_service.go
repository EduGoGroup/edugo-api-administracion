package service

import (
	"context"

	"github.com/edugo/api-administracion/internal/application/dto"
	"github.com/edugo/api-administracion/internal/domain/repository"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

// StatsService define las operaciones de negocio para estadísticas
type StatsService interface {
	GetGlobalStats(ctx context.Context) (*dto.GlobalStatsResponse, error)
}

type statsService struct {
	statsRepo repository.StatsRepository
	logger    logger.Logger
}

func NewStatsService(statsRepo repository.StatsRepository, logger logger.Logger) StatsService {
	return &statsService{
		statsRepo: statsRepo,
		logger:    logger,
	}
}

func (s *statsService) GetGlobalStats(ctx context.Context) (*dto.GlobalStatsResponse, error) {
	stats, err := s.statsRepo.GetGlobalStats(ctx)
	if err != nil {
		s.logger.Error("failed to get global stats", "error", err)
		return nil, errors.NewDatabaseError("get stats", err)
	}

	s.logger.Debug("global stats retrieved",
		"total_users", stats.TotalUsers,
		"total_schools", stats.TotalSchools,
	)

	return &dto.GlobalStatsResponse{
		TotalUsers:            stats.TotalUsers,
		TotalActiveUsers:      stats.TotalActiveUsers,
		TotalSchools:          stats.TotalSchools,
		TotalSubjects:         stats.TotalSubjects,
		TotalGuardianRelations: stats.TotalGuardianRelations,
	}, nil
}
