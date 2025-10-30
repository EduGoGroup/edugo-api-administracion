package repository

import (
	"context"
	"database/sql"

	"github.com/edugo/api-administracion/internal/domain/entity"
	"github.com/edugo/api-administracion/internal/domain/repository"
	"github.com/edugo/api-administracion/internal/domain/valueobject"
)

type postgresUnitRepository struct {
	db *sql.DB
}

func NewPostgresUnitRepository(db *sql.DB) repository.UnitRepository {
	return &postgresUnitRepository{db: db}
}

func (r *postgresUnitRepository) Create(ctx context.Context, unit *entity.Unit) error {
	query := `
		INSERT INTO units (id, school_id, parent_unit_id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	var parentUnitID interface{}
	if !unit.IsRootUnit() {
		parentUnitID = unit.ParentUnitID().String()
	}

	_, err := r.db.ExecContext(ctx, query,
		unit.ID().String(),
		unit.SchoolID().String(),
		parentUnitID,
		unit.Name(),
		unit.Description(),
		unit.IsActive(),
		unit.CreatedAt(),
		unit.UpdatedAt(),
	)

	return err
}

func (r *postgresUnitRepository) FindByID(ctx context.Context, id valueobject.UnitID) (*entity.Unit, error) {
	query := `
		SELECT id, school_id, parent_unit_id, name, description, is_active, created_at, updated_at
		FROM units
		WHERE id = $1
	`

	var (
		idStr        string
		schoolIDStr  string
		parentIDStr  sql.NullString
		name         string
		description  string
		isActive     bool
		createdAt    sql.NullTime
		updatedAt    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &schoolIDStr, &parentIDStr, &name, &description, &isActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.scanToEntity(idStr, schoolIDStr, parentIDStr, name, description, isActive, createdAt, updatedAt)
}

func (r *postgresUnitRepository) Update(ctx context.Context, unit *entity.Unit) error {
	query := `
		UPDATE units
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query,
		unit.Name(),
		unit.Description(),
		unit.UpdatedAt(),
		unit.ID().String(),
	)

	return err
}

func (r *postgresUnitRepository) Delete(ctx context.Context, id valueobject.UnitID) error {
	query := `UPDATE units SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *postgresUnitRepository) FindBySchool(ctx context.Context, schoolID valueobject.SchoolID) ([]*entity.Unit, error) {
	query := `
		SELECT id, school_id, parent_unit_id, name, description, is_active, created_at, updated_at
		FROM units
		WHERE school_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, schoolID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *postgresUnitRepository) FindChildren(ctx context.Context, parentID valueobject.UnitID) ([]*entity.Unit, error) {
	query := `
		SELECT id, school_id, parent_unit_id, name, description, is_active, created_at, updated_at
		FROM units
		WHERE parent_unit_id = $1
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query, parentID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *postgresUnitRepository) HasChildren(ctx context.Context, unitID valueobject.UnitID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM units WHERE parent_unit_id = $1)`

	var hasChildren bool
	err := r.db.QueryRowContext(ctx, query, unitID.String()).Scan(&hasChildren)
	return hasChildren, err
}

func (r *postgresUnitRepository) IsDescendantOf(ctx context.Context, unitID, ancestorID valueobject.UnitID) (bool, error) {
	// Verificar si unitID es descendiente de ancestorID (para prevenir ciclos)
	// Usa recursive CTE
	query := `
		WITH RECURSIVE ancestors AS (
			SELECT id, parent_unit_id FROM units WHERE id = $1
			UNION ALL
			SELECT u.id, u.parent_unit_id
			FROM units u
			INNER JOIN ancestors a ON u.id = a.parent_unit_id
		)
		SELECT EXISTS(SELECT 1 FROM ancestors WHERE id = $2)
	`

	var isDescendant bool
	err := r.db.QueryRowContext(ctx, query, unitID.String(), ancestorID.String()).Scan(&isDescendant)
	return isDescendant, err
}

func (r *postgresUnitRepository) AddMember(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID, role string) error {
	query := `
		INSERT INTO unit_memberships (unit_id, user_id, role, created_at)
		VALUES ($1, $2, $3, NOW())
	`

	_, err := r.db.ExecContext(ctx, query, unitID.String(), userID.String(), role)
	return err
}

func (r *postgresUnitRepository) RemoveMember(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) error {
	query := `DELETE FROM unit_memberships WHERE unit_id = $1 AND user_id = $2`

	_, err := r.db.ExecContext(ctx, query, unitID.String(), userID.String())
	return err
}

func (r *postgresUnitRepository) IsMember(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM unit_memberships WHERE unit_id = $1 AND user_id = $2)`

	var isMember bool
	err := r.db.QueryRowContext(ctx, query, unitID.String(), userID.String()).Scan(&isMember)
	return isMember, err
}

// Helpers

func (r *postgresUnitRepository) scanRows(rows *sql.Rows) ([]*entity.Unit, error) {
	var units []*entity.Unit

	for rows.Next() {
		var (
			idStr        string
			schoolIDStr  string
			parentIDStr  sql.NullString
			name         string
			description  string
			isActive     bool
			createdAt    sql.NullTime
			updatedAt    sql.NullTime
		)

		if err := rows.Scan(&idStr, &schoolIDStr, &parentIDStr, &name, &description, &isActive, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		unit, err := r.scanToEntity(idStr, schoolIDStr, parentIDStr, name, description, isActive, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}

		units = append(units, unit)
	}

	return units, rows.Err()
}

func (r *postgresUnitRepository) scanToEntity(
	idStr, schoolIDStr string,
	parentIDStr sql.NullString,
	name, description string,
	isActive bool,
	createdAt, updatedAt sql.NullTime,
) (*entity.Unit, error) {
	unitID, _ := valueobject.UnitIDFromString(idStr)
	schoolID, _ := valueobject.SchoolIDFromString(schoolIDStr)

	var parentUnitID *valueobject.UnitID
	if parentIDStr.Valid && parentIDStr.String != "" {
		pID, _ := valueobject.UnitIDFromString(parentIDStr.String)
		parentUnitID = &pID
	}

	return entity.ReconstructUnit(unitID, schoolID, parentUnitID, name, description, isActive, createdAt.Time, updatedAt.Time), nil
}
