package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/google/uuid"
)

type postgresAcademicUnitRepository struct {
	db *sql.DB
}

func NewPostgresAcademicUnitRepository(db *sql.DB) repository.AcademicUnitRepository {
	return &postgresAcademicUnitRepository{db: db}
}

func (r *postgresAcademicUnitRepository) Create(ctx context.Context, unit *entities.AcademicUnit) error {
	query := `INSERT INTO academic_units (id, parent_unit_id, school_id, type, name, code, description, level, academic_year, metadata, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.ExecContext(ctx, query,
		unit.ID, unit.ParentUnitID, unit.SchoolID, unit.Type, unit.Name, unit.Code,
		unit.Description, unit.Level, unit.AcademicYear, unit.Metadata, unit.IsActive,
		unit.CreatedAt, unit.UpdatedAt,
	)
	return err
}

func (r *postgresAcademicUnitRepository) FindByID(ctx context.Context, id uuid.UUID, includeDeleted bool) (*entities.AcademicUnit, error) {
	query := `SELECT id, parent_unit_id, school_id, type, name, code, description, level, academic_year, metadata, is_active, created_at, updated_at, deleted_at
		FROM academic_units WHERE id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}
	return r.scanOne(ctx, query, id)
}

func (r *postgresAcademicUnitRepository) FindBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (*entities.AcademicUnit, error) {
	query := `SELECT id, parent_unit_id, school_id, type, name, code, description, level, academic_year, metadata, is_active, created_at, updated_at, deleted_at
		FROM academic_units WHERE school_id = $1 AND code = $2 AND deleted_at IS NULL`
	return r.scanOne(ctx, query, schoolID, code)
}

func (r *postgresAcademicUnitRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	query := `SELECT id, parent_unit_id, school_id, type, name, code, description, level, academic_year, metadata, is_active, created_at, updated_at, deleted_at
		FROM academic_units WHERE school_id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}
	query += " ORDER BY name"
	return r.scanMany(ctx, query, schoolID)
}

func (r *postgresAcademicUnitRepository) FindByParentID(ctx context.Context, parentID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	query := `SELECT id, parent_unit_id, school_id, type, name, code, description, level, academic_year, metadata, is_active, created_at, updated_at, deleted_at
		FROM academic_units WHERE parent_unit_id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}
	query += " ORDER BY name"
	return r.scanMany(ctx, query, parentID)
}

func (r *postgresAcademicUnitRepository) FindRootUnits(ctx context.Context, schoolID uuid.UUID) ([]*entities.AcademicUnit, error) {
	query := `SELECT id, parent_unit_id, school_id, type, name, code, description, level, academic_year, metadata, is_active, created_at, updated_at, deleted_at
		FROM academic_units WHERE school_id = $1 AND parent_unit_id IS NULL AND deleted_at IS NULL ORDER BY name`
	return r.scanMany(ctx, query, schoolID)
}

func (r *postgresAcademicUnitRepository) FindByType(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	query := `SELECT id, parent_unit_id, school_id, type, name, code, description, level, academic_year, metadata, is_active, created_at, updated_at, deleted_at
		FROM academic_units WHERE school_id = $1 AND type = $2`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}
	query += " ORDER BY name"
	return r.scanMany(ctx, query, schoolID, unitType)
}

func (r *postgresAcademicUnitRepository) Update(ctx context.Context, unit *entities.AcademicUnit) error {
	query := `UPDATE academic_units SET parent_unit_id = $1, name = $2, code = $3, description = $4, level = $5, 
		academic_year = $6, metadata = $7, is_active = $8, updated_at = $9 WHERE id = $10 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query,
		unit.ParentUnitID, unit.Name, unit.Code, unit.Description, unit.Level,
		unit.AcademicYear, unit.Metadata, unit.IsActive, unit.UpdatedAt, unit.ID,
	)
	return err
}

func (r *postgresAcademicUnitRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE academic_units SET deleted_at = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	return err
}

func (r *postgresAcademicUnitRepository) Restore(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE academic_units SET deleted_at = NULL, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *postgresAcademicUnitRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM academic_units WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresAcademicUnitRepository) GetHierarchyPath(ctx context.Context, id uuid.UUID) ([]*entities.AcademicUnit, error) {
	// Implementación recursiva simple - puede optimizarse con CTE
	var path []*entities.AcademicUnit
	currentID := id
	
	for {
		unit, err := r.FindByID(ctx, currentID, false)
		if err != nil || unit == nil {
			break
		}
		path = append([]*entities.AcademicUnit{unit}, path...)
		if unit.ParentUnitID == nil {
			break
		}
		currentID = *unit.ParentUnitID
	}
	
	return path, nil
}

func (r *postgresAcademicUnitRepository) ExistsBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM academic_units WHERE school_id = $1 AND code = $2 AND deleted_at IS NULL)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, schoolID, code).Scan(&exists)
	return exists, err
}

// Helper methods

func (r *postgresAcademicUnitRepository) scanOne(ctx context.Context, query string, args ...interface{}) (*entities.AcademicUnit, error) {
	unit := &entities.AcademicUnit{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&unit.ID, &unit.ParentUnitID, &unit.SchoolID, &unit.Type, &unit.Name, &unit.Code,
		&unit.Description, &unit.Level, &unit.AcademicYear, &unit.Metadata, &unit.IsActive,
		&unit.CreatedAt, &unit.UpdatedAt, &unit.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("academic_unit")
	}
	if err != nil {
		return nil, err
	}
	return unit, nil
}

func (r *postgresAcademicUnitRepository) scanMany(ctx context.Context, query string, args ...interface{}) ([]*entities.AcademicUnit, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*entities.AcademicUnit
	for rows.Next() {
		unit := &entities.AcademicUnit{}
		err := rows.Scan(
			&unit.ID, &unit.ParentUnitID, &unit.SchoolID, &unit.Type, &unit.Name, &unit.Code,
			&unit.Description, &unit.Level, &unit.AcademicYear, &unit.Metadata, &unit.IsActive,
			&unit.CreatedAt, &unit.UpdatedAt, &unit.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, rows.Err()
}
