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

type postgresUnitMembershipRepository struct {
	db *sql.DB
}

func NewPostgresUnitMembershipRepository(db *sql.DB) repository.UnitMembershipRepository {
	return &postgresUnitMembershipRepository{db: db}
}

func (r *postgresUnitMembershipRepository) Create(ctx context.Context, membership *entity.UnitMembership) error {
	// Obtener school_id de la academic_unit
	var schoolID string
	err := r.db.QueryRowContext(ctx, "SELECT school_id FROM academic_units WHERE id = $1", membership.UnitID().String()).Scan(&schoolID)
	if err != nil {
		return errors.NewDatabaseError("obtener school_id de academic_unit", err)
	}

	query := `
		INSERT INTO memberships (id, academic_unit_id, user_id, school_id, role, enrolled_at, withdrawn_at, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Serializar metadata (siempre enviar al menos {} para JSONB)
	metadata := membership.Metadata()
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return errors.NewDatabaseError("marshal metadata", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		membership.ID().String(),
		membership.UnitID().String(),
		membership.UserID().String(),
		schoolID,
		membership.Role().String(),
		membership.ValidFrom(),
		membership.ValidUntil(),
		metadataJSON,
		membership.CreatedAt(),
		membership.UpdatedAt(),
	)

	return err
}

func (r *postgresUnitMembershipRepository) FindByID(ctx context.Context, id valueobject.MembershipID) (*entity.UnitMembership, error) {
	query := `
		SELECT id, academic_unit_id, school_id, user_id, role, enrolled_at, withdrawn_at, metadata, created_at, updated_at
		FROM memberships
		WHERE id = $1
	`

	return r.scanOneMembership(ctx, query, id.String())
}

func (r *postgresUnitMembershipRepository) FindByUnitID(ctx context.Context, unitID valueobject.UnitID, activeOnly bool) ([]*entity.UnitMembership, error) {
	query := `
		SELECT id, academic_unit_id, school_id, user_id, role, enrolled_at, withdrawn_at, metadata, created_at, updated_at
		FROM memberships
		WHERE academic_unit_id = $1
	`

	if activeOnly {
		query += " AND is_active = true AND (withdrawn_at IS NULL OR withdrawn_at > CURRENT_TIMESTAMP)"
	}

	query += " ORDER BY role, enrolled_at DESC"

	return r.scanMemberships(ctx, query, unitID.String())
}

func (r *postgresUnitMembershipRepository) FindByUserID(ctx context.Context, userID valueobject.UserID, activeOnly bool) ([]*entity.UnitMembership, error) {
	query := `
		SELECT id, academic_unit_id, school_id, user_id, role, enrolled_at, withdrawn_at, metadata, created_at, updated_at
		FROM memberships
		WHERE user_id = $1
	`

	if activeOnly {
		query += " AND is_active = true AND (withdrawn_at IS NULL OR withdrawn_at > CURRENT_TIMESTAMP)"
	}

	query += " ORDER BY enrolled_at DESC"

	return r.scanMemberships(ctx, query, userID.String())
}

func (r *postgresUnitMembershipRepository) FindByUnitAndUser(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) (*entity.UnitMembership, error) {
	query := `
		SELECT id, academic_unit_id, school_id, user_id, role, enrolled_at, withdrawn_at, metadata, created_at, updated_at
		FROM memberships
		WHERE academic_unit_id = $1 AND user_id = $2
		ORDER BY enrolled_at DESC
		LIMIT 1
	`

	return r.scanOneMembership(ctx, query, unitID.String(), userID.String())
}

func (r *postgresUnitMembershipRepository) FindByRole(ctx context.Context, unitID valueobject.UnitID, role valueobject.MembershipRole, activeOnly bool) ([]*entity.UnitMembership, error) {
	query := `
		SELECT id, academic_unit_id, school_id, user_id, role, enrolled_at, withdrawn_at, metadata, created_at, updated_at
		FROM memberships
		WHERE academic_unit_id = $1 AND role = $2
	`

	if activeOnly {
		query += " AND is_active = true AND (withdrawn_at IS NULL OR withdrawn_at > CURRENT_TIMESTAMP)"
	}

	query += " ORDER BY enrolled_at DESC"

	return r.scanMemberships(ctx, query, unitID.String(), role.String())
}

func (r *postgresUnitMembershipRepository) FindActiveAt(ctx context.Context, unitID valueobject.UnitID, at time.Time) ([]*entity.UnitMembership, error) {
	query := `
		SELECT id, academic_unit_id, school_id, user_id, role, enrolled_at, withdrawn_at, metadata, created_at, updated_at
		FROM memberships
		WHERE academic_unit_id = $1
		  AND enrolled_at <= $2
		  AND (withdrawn_at IS NULL OR withdrawn_at >= $2)
		ORDER BY role, enrolled_at DESC
	`

	return r.scanMemberships(ctx, query, unitID.String(), at)
}

func (r *postgresUnitMembershipRepository) Update(ctx context.Context, membership *entity.UnitMembership) error {
	query := `
		UPDATE memberships
		SET role = $1, valid_until = $2, metadata = $3, updated_at = $4
		WHERE id = $5
	`

	// Serializar metadata (siempre enviar al menos {} para JSONB)
	metadata := membership.Metadata()
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return errors.NewDatabaseError("marshal metadata", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		membership.Role().String(),
		membership.ValidUntil(),
		metadataJSON,
		membership.UpdatedAt(),
		membership.ID().String(),
	)

	return err
}

func (r *postgresUnitMembershipRepository) Delete(ctx context.Context, id valueobject.MembershipID) error {
	query := `DELETE FROM memberships WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *postgresUnitMembershipRepository) ExistsByUnitAndUser(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM memberships
			WHERE academic_unit_id = $1 AND user_id = $2
			  AND is_active = true AND (withdrawn_at IS NULL OR withdrawn_at > CURRENT_TIMESTAMP)
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, unitID.String(), userID.String()).Scan(&exists)
	return exists, err
}

func (r *postgresUnitMembershipRepository) CountByUnit(ctx context.Context, unitID valueobject.UnitID, activeOnly bool) (int, error) {
	query := `SELECT COUNT(*) FROM memberships WHERE academic_unit_id = $1`

	if activeOnly {
		query += " AND is_active = true AND (withdrawn_at IS NULL OR withdrawn_at > CURRENT_TIMESTAMP)"
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, unitID.String()).Scan(&count)
	return count, err
}

func (r *postgresUnitMembershipRepository) CountByRole(ctx context.Context, unitID valueobject.UnitID, role valueobject.MembershipRole) (int, error) {
	query := `
		SELECT COUNT(*) FROM memberships
		WHERE academic_unit_id = $1 AND role = $2
		  AND is_active = true AND (withdrawn_at IS NULL OR withdrawn_at > CURRENT_TIMESTAMP)
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, unitID.String(), role.String()).Scan(&count)
	return count, err
}

// Helper: escanear una membresía
func (r *postgresUnitMembershipRepository) scanOneMembership(ctx context.Context, query string, args ...interface{}) (*entity.UnitMembership, error) {
	var (
		idStr        string
		unitIDStr    string
		schoolIDStr  string
		userIDStr    string
		role         string
		enrolledAt   sql.NullTime
		withdrawnAt  sql.NullTime
		metadataJSON []byte
		createdAt    sql.NullTime
		updatedAt    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&idStr, &unitIDStr, &schoolIDStr, &userIDStr, &role, &enrolledAt, &withdrawnAt,
		&metadataJSON, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.buildMembership(idStr, unitIDStr, userIDStr, role, enrolledAt, withdrawnAt, metadataJSON, createdAt, updatedAt)
}

// Helper: escanear múltiples membresías
func (r *postgresUnitMembershipRepository) scanMemberships(ctx context.Context, query string, args ...interface{}) ([]*entity.UnitMembership, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var memberships []*entity.UnitMembership
	for rows.Next() {
		var (
			idStr        string
			unitIDStr    string
			schoolIDStr  string
			userIDStr    string
			role         string
			enrolledAt   sql.NullTime
			withdrawnAt  sql.NullTime
			metadataJSON []byte
			createdAt    sql.NullTime
			updatedAt    sql.NullTime
		)

		if err := rows.Scan(&idStr, &unitIDStr, &schoolIDStr, &userIDStr, &role, &enrolledAt, &withdrawnAt, &metadataJSON, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		membership, err := r.buildMembership(idStr, unitIDStr, userIDStr, role, enrolledAt, withdrawnAt, metadataJSON, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	return memberships, rows.Err()
}

// Helper: construir entidad de membresía
func (r *postgresUnitMembershipRepository) buildMembership(
	idStr, unitIDStr, userIDStr, role string,
	enrolledAt, withdrawnAt sql.NullTime,
	metadataJSON []byte,
	createdAt, updatedAt sql.NullTime,
) (*entity.UnitMembership, error) {
	membershipID, err := valueobject.MembershipIDFromString(idStr)
	if err != nil {
		return nil, errors.NewDatabaseError("parse membership ID", err)
	}

	unitID, err := valueobject.UnitIDFromString(unitIDStr)
	if err != nil {
		return nil, errors.NewDatabaseError("parse unit ID", err)
	}

	userID, err := valueobject.UserIDFromString(userIDStr)
	if err != nil {
		return nil, errors.NewDatabaseError("parse user ID", err)
	}

	roleVO, err := valueobject.NewMembershipRole(role)
	if err != nil {
		return nil, errors.NewDatabaseError("parse role", err)
	}

	var metadata map[string]interface{}
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			return nil, errors.NewDatabaseError("unmarshal metadata", err)
		}
	}

	var withdrawnAtPtr *time.Time
	if withdrawnAt.Valid {
		withdrawnAtPtr = &withdrawnAt.Time
	}

	return entity.ReconstructUnitMembership(
		membershipID,
		unitID,
		userID,
		roleVO,
		enrolledAt.Time,
		withdrawnAtPtr,
		metadata,
		createdAt.Time,
		updatedAt.Time,
	), nil
}
