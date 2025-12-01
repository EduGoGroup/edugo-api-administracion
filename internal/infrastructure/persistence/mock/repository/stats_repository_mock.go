package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
)

// MockStatsRepository es una implementación en memoria del StatsRepository para testing
// Retorna valores vacíos/cero ya que en modo mock las estadísticas se calculan
// dinámicamente desde los datos creados durante los tests
type MockStatsRepository struct{}

// NewMockStatsRepository crea una nueva instancia de MockStatsRepository
func NewMockStatsRepository() repository.StatsRepository {
	return &MockStatsRepository{}
}

// GetGlobalStats obtiene estadísticas globales del sistema
// En modo mock retorna valores cero - los tests deben verificar
// estadísticas después de crear datos específicos
func (r *MockStatsRepository) GetGlobalStats(ctx context.Context) (repository.GlobalStats, error) {
	return repository.GlobalStats{
		TotalUsers:             0,
		TotalActiveUsers:       0,
		TotalSchools:           0,
		TotalSubjects:          0,
		TotalGuardianRelations: 0,
	}, nil
}
