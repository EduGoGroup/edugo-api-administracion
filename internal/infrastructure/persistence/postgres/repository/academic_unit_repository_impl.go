package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

type postgresAcademicUnitRepository struct {
	db *sql.DB
}

func NewPostgresAcademicUnitRepository(db *sql.DB) repository.AcademicUnitRepository {
	return &postgresAcademicUnitRepository{db: db}
}

func (r *postgresAcademicUnitRepository) Create(ctx context.Context, unit *entity.AcademicUnit) error {
	query := `
		INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code, description, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var parentID *string
	if unit.ParentUnitID() != nil {
		id := unit.ParentUnitID().String()
		parentID = &id
	}

	var metadataJSON []byte
	if len(unit.Metadata()) > 0 {
		var err error
		metadataJSON, err = json.Marshal(unit.Metadata())
		if err != nil {
			return errors.NewDatabaseError("marshal metadata", err)
		}
	}

	_, err := r.db.ExecContext(ctx, query,
		unit.ID().String(),
		parentID,
		unit.SchoolID().String(),
		unit.UnitType().String(),
		unit.DisplayName(),
		unit.Code(),
		unit.Description(),
		metadataJSON,
		unit.CreatedAt(),
		unit.UpdatedAt(),
	)

	return err
}

func (r *postgresAcademicUnitRepository) FindByID(ctx context.Context, id valueobject.UnitID, includeDeleted bool) (*entity.AcademicUnit, error) {
	query := `
		SELECT id, parent_unit_id, school_id, unit_type, display_name, code, description, metadata, created_at, updated_at, deleted_at
		FROM academic_unit
		WHERE id = $1
	`

	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	return r.scanOneUnit(ctx, query, id.String())
}

func (r *postgresAcademicUnitRepository) FindBySchoolIDAndCode(ctx context.Context, schoolID valueobject.SchoolID, code string) (*entity.AcademicUnit, error) {
	query := `
		SELECT id, parent_unit_id, school_id, unit_type, display_name, code, description, metadata, created_at, updated_at, deleted_at
		FROM academic_unit
		WHERE school_id = $1 AND code = $2 AND deleted_at IS NULL
	`

	return r.scanOneUnit(ctx, query, schoolID.String(), code)
}

func (r *postgresAcademicUnitRepository) FindBySchoolID(ctx context.Context, schoolID valueobject.SchoolID, includeDeleted bool) ([]*entity.AcademicUnit, error) {
	query := `
		SELECT id, parent_unit_id, school_id, unit_type, display_name, code, description, metadata, created_at, updated_at, deleted_at
		FROM academic_unit
		WHERE school_id = $1
	`

	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	query += " ORDER BY unit_type, display_name"

	return r.scanUnits(ctx, query, schoolID.String())
}

func (r *postgresAcademicUnitRepository) FindByParentID(ctx context.Context, parentID valueobject.UnitID, includeDeleted bool) ([]*entity.AcademicUnit, error) {
	query := `
		SELECT id, parent_unit_id, school_id, unit_type, display_name, code, description, metadata, created_at, updated_at, deleted_at
		FROM academic_unit
		WHERE parent_unit_id = $1
	`

	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	query += " ORDER BY display_name"

	return r.scanUnits(ctx, query, parentID.String())
}

func (r *postgresAcademicUnitRepository) FindRootUnits(ctx context.Context, schoolID valueobject.SchoolID) ([]*entity.AcademicUnit, error) {
	query := `
		SELECT id, parent_unit_id, school_id, unit_type, display_name, code, description, metadata, created_at, updated_at, deleted_at
		FROM academic_unit
		WHERE school_id = $1 AND parent_unit_id IS NULL AND deleted_at IS NULL
		ORDER BY unit_type, display_name
	`

	return r.scanUnits(ctx, query, schoolID.String())
}

func (r *postgresAcademicUnitRepository) FindByType(ctx context.Context, schoolID valueobject.SchoolID, unitType valueobject.UnitType, includeDeleted bool) ([]*entity.AcademicUnit, error) {
	query := `
		SELECT id, parent_unit_id, school_id, unit_type, display_name, code, description, metadata, created_at, updated_at, deleted_at
		FROM academic_unit
		WHERE school_id = $1 AND unit_type = $2
	`

	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	query += " ORDER BY display_name"

	return r.scanUnits(ctx, query, schoolID.String(), unitType.String())
}

func (r *postgresAcademicUnitRepository) Update(ctx context.Context, unit *entity.AcademicUnit) error {
	query := `
		UPDATE academic_unit
		SET parent_unit_id = $1, display_name = $2, description = $3, metadata = $4, updated_at = $5
		WHERE id = $6
	`

	var parentID *string
	if unit.ParentUnitID() != nil {
		id := unit.ParentUnitID().String()
		parentID = &id
	}

	var metadataJSON []byte
	if len(unit.Metadata()) > 0 {
		var err error
		metadataJSON, err = json.Marshal(unit.Metadata())
		if err != nil {
			return errors.NewDatabaseError("marshal metadata", err)
		}
	}

	_, err := r.db.ExecContext(ctx, query,
		parentID,
		unit.DisplayName(),
		unit.Description(),
		metadataJSON,
		unit.UpdatedAt(),
		unit.ID().String(),
	)

	return err
}

func (r *postgresAcademicUnitRepository) SoftDelete(ctx context.Context, id valueobject.UnitID) error {
	query := `UPDATE academic_unit SET deleted_at = $1, updated_at = $2 WHERE id = $3`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, id.String())
	return err
}

func (r *postgresAcademicUnitRepository) Restore(ctx context.Context, id valueobject.UnitID) error {
	query := `UPDATE academic_unit SET deleted_at = NULL, updated_at = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id.String())
	return err
}

func (r *postgresAcademicUnitRepository) HardDelete(ctx context.Context, id valueobject.UnitID) error {
	query := `DELETE FROM academic_unit WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *postgresAcademicUnitRepository) GetHierarchyPath(ctx context.Context, id valueobject.UnitID) ([]*entity.AcademicUnit, error) {
	query := `
		WITH RECURSIVE path AS (
			SELECT id, parent_unit_id, school_id, unit_type, display_name, code, description, metadata, created_at, updated_at, deleted_at, 1 as depth
			FROM academic_unit
			WHERE id = $1
			
			UNION ALL
			
			SELECT au.id, au.parent_unit_id, au.school_id, au.unit_type, au.display_name, au.code, au.description, au.metadata, au.created_at, au.updated_at, au.deleted_at, p.depth + 1
			FROM academic_unit au
			INNER JOIN path p ON au.id = p.parent_unit_id
		)
		SELECT id, parent_unit_id, school_id, unit_type, display_name, code, description, metadata, created_at, updated_at, deleted_at
		FROM path
		ORDER BY depth DESC
	`

	return r.scanUnits(ctx, query, id.String())
}

func (r *postgresAcademicUnitRepository) ExistsBySchoolIDAndCode(ctx context.Context, schoolID valueobject.SchoolID, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM academic_unit WHERE school_id = $1 AND code = $2 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, schoolID.String(), code).Scan(&exists)
	return exists, err
}

func (r *postgresAcademicUnitRepository) HasChildren(ctx context.Context, id valueobject.UnitID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM academic_unit WHERE parent_unit_id = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&exists)
	return exists, err
}

// Helper: escanear una unidad
func (r *postgresAcademicUnitRepository) scanOneUnit(ctx context.Context, query string, args ...interface{}) (*entity.AcademicUnit, error) {
	var (
		idStr        string
		parentIDStr  sql.NullString
		schoolIDStr  string
		unitType     string
		displayName  string
		code         string
		description  string
		metadataJSON []byte
		createdAt    sql.NullTime
		updatedAt    sql.NullTime
		deletedAt    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&idStr, &parentIDStr, &schoolIDStr, &unitType, &displayName, &code, &description,
		&metadataJSON, &createdAt, &updatedAt, &deletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.buildUnit(idStr, parentIDStr, schoolIDStr, unitType, displayName, code, description, metadataJSON, createdAt, updatedAt, deletedAt)
}

// Helper: escanear múltiples unidades
func (r *postgresAcademicUnitRepository) scanUnits(ctx context.Context, query string, args ...interface{}) ([]*entity.AcademicUnit, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*entity.AcademicUnit
	for rows.Next() {
		var (
			idStr        string
			parentIDStr  sql.NullString
			schoolIDStr  string
			unitType     string
			displayName  string
			code         string
			description  string
			metadataJSON []byte
			createdAt    sql.NullTime
			updatedAt    sql.NullTime
			deletedAt    sql.NullTime
		)

		if err := rows.Scan(&idStr, &parentIDStr, &schoolIDStr, &unitType, &displayName, &code, &description, &metadataJSON, &createdAt, &updatedAt, &deletedAt); err != nil {
			return nil, err
		}

		unit, err := r.buildUnit(idStr, parentIDStr, schoolIDStr, unitType, displayName, code, description, metadataJSON, createdAt, updatedAt, deletedAt)
		if err != nil {
			return nil, err
		}

		units = append(units, unit)
	}

	return units, rows.Err()
}

// Helper: construir entidad desde campos escaneados
func (r *postgresAcademicUnitRepository) buildUnit(
	idStr string,
	parentIDStr sql.NullString,
	schoolIDStr string,
	unitType string,
	displayName string,
	code string,
	description string,
	metadataJSON []byte,
	createdAt, updatedAt, deletedAt sql.NullTime,
) (*entity.AcademicUnit, error) {
	unitID, err := valueobject.UnitIDFromString(idStr)
	if err != nil {
		return nil, errors.NewDatabaseError("parse unit ID", err)
	}

	schoolID, err := valueobject.SchoolIDFromString(schoolIDStr)
	if err != nil {
		return nil, errors.NewDatabaseError("parse school ID", err)
	}

	var parentID *valueobject.UnitID
	if parentIDStr.Valid && parentIDStr.String != "" {
		pID, err := valueobject.UnitIDFromString(parentIDStr.String)
		if err != nil {
			return nil, errors.NewDatabaseError("parse parent unit ID", err)
		}
		parentID = &pID
	}

	unitTypeVO, err := valueobject.NewUnitType(unitType)
	if err != nil {
		return nil, errors.NewDatabaseError("parse unit type", err)
	}

	var metadata map[string]interface{}
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			return nil, errors.NewDatabaseError("unmarshal metadata", err)
		}
	}

	var deletedAtPtr *time.Time
	if deletedAt.Valid {
		deletedAtPtr = &deletedAt.Time
	}

	return entity.ReconstructAcademicUnit(
		unitID,
		parentID,
		schoolID,
		unitTypeVO,
		displayName,
		code,
		description,
		metadata,
		createdAt.Time,
		updatedAt.Time,
		deletedAtPtr,
	), nil
}
