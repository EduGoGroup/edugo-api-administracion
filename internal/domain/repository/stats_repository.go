package repository

import (
	"context"
)

// StatsRepository define las operaciones para obtener estadísticas
type StatsRepository interface {
	// GetGlobalStats obtiene estadísticas globales del sistema
	GetGlobalStats(ctx context.Context) (GlobalStats, error)
}

// GlobalStats representa las estadísticas globales
type GlobalStats struct {
	TotalUsers            int
	TotalActiveUsers      int
	TotalSchools          int
	TotalSubjects         int
	TotalGuardianRelations int
}
