package repository

import (
	"context"
	"database/sql"
	"time"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type postgresUnitRepository struct {
	db *sql.DB
}

func NewPostgresUnitRepository(db *sql.DB) repository.UnitRepository {
	return &postgresUnitRepository{db: db}
}

func (r *postgresUnitRepository) Create(ctx context.Context, unit *entities.Unit) error {
	query := `INSERT INTO units (id, school_id, parent_unit_id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		unit.ID, unit.SchoolID, unit.ParentUnitID, unit.Name, unit.Description,
		unit.IsActive, unit.CreatedAt, unit.UpdatedAt,
	)
	return err
}

func (r *postgresUnitRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Unit, error) {
	query := `SELECT id, school_id, parent_unit_id, name, description, is_active, created_at, updated_at
		FROM units WHERE id = $1`
	unit := &entities.Unit{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&unit.ID, &unit.SchoolID, &unit.ParentUnitID, &unit.Name, &unit.Description,
		&unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return unit, err
}

func (r *postgresUnitRepository) Update(ctx context.Context, unit *entities.Unit) error {
	query := `UPDATE units SET name = $1, description = $2, parent_unit_id = $3, is_active = $4, updated_at = $5
		WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query,
		unit.Name, unit.Description, unit.ParentUnitID, unit.IsActive, unit.UpdatedAt, unit.ID,
	)
	return err
}

func (r *postgresUnitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE units SET is_active = false, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *postgresUnitRepository) List(ctx context.Context, schoolID uuid.UUID) ([]*entities.Unit, error) {
	query := `SELECT id, school_id, parent_unit_id, name, description, is_active, created_at, updated_at
		FROM units WHERE school_id = $1 AND is_active = true ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*entities.Unit
	for rows.Next() {
		unit := &entities.Unit{}
		err := rows.Scan(&unit.ID, &unit.SchoolID, &unit.ParentUnitID, &unit.Name, &unit.Description,
			&unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, rows.Err()
}
