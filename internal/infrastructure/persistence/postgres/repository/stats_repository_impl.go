package repository

import (
	"context"
	"database/sql"

	"github.com/edugo/api-administracion/internal/domain/repository"
)

type postgresStatsRepository struct {
	db *sql.DB
}

func NewPostgresStatsRepository(db *sql.DB) repository.StatsRepository {
	return &postgresStatsRepository{db: db}
}

func (r *postgresStatsRepository) GetGlobalStats(ctx context.Context) (repository.GlobalStats, error) {
	var stats repository.GlobalStats

	// Consultar estadísticas en paralelo usando una sola query
	query := `
		SELECT
			(SELECT COUNT(*) FROM users) AS total_users,
			(SELECT COUNT(*) FROM users WHERE is_active = true) AS total_active_users,
			(SELECT COUNT(*) FROM schools WHERE is_active = true) AS total_schools,
			(SELECT COUNT(*) FROM subjects WHERE is_active = true) AS total_subjects,
			(SELECT COUNT(*) FROM guardian_relations WHERE is_active = true) AS total_guardian_relations
	`

	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalUsers,
		&stats.TotalActiveUsers,
		&stats.TotalSchools,
		&stats.TotalSubjects,
		&stats.TotalGuardianRelations,
	)

	if err != nil {
		return repository.GlobalStats{}, err
	}

	return stats, nil
}
