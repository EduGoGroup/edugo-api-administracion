package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	mockData "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/mock/data"
)

// MockStatsRepository es una implementación en memoria del StatsRepository para testing
// No usa sync.RWMutex ya que es stateless y calcula dinámicamente desde otros datos mock
type MockStatsRepository struct{}

// NewMockStatsRepository crea una nueva instancia de MockStatsRepository
func NewMockStatsRepository() repository.StatsRepository {
	return &MockStatsRepository{}
}

// GetGlobalStats obtiene estadísticas globales del sistema
// Calcula dinámicamente desde los datos mock disponibles
func (r *MockStatsRepository) GetGlobalStats(ctx context.Context) (repository.GlobalStats, error) {
	// Obtener datos desde mockData
	schools := mockData.GetSchools()
	users := mockData.GetUsers()

	// Contar usuarios activos
	totalActiveUsers := 0

	for _, user := range users {
		// Excluir usuarios eliminados
		if user.DeletedAt != nil {
			continue
		}

		// Contar usuarios activos
		if user.IsActive {
			totalActiveUsers++
		}
	}

	// Crear estructura de estadísticas globales
	stats := repository.GlobalStats{
		TotalUsers:             len(users),
		TotalActiveUsers:       totalActiveUsers,
		TotalSchools:           len(schools),
		TotalSubjects:          len(mockData.GetSubjects()),
		TotalGuardianRelations: len(mockData.GetGuardianRelations()),
	}

	return stats, nil
}
