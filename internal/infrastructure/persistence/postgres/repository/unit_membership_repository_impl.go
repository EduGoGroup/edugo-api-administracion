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
	query := `
		INSERT INTO unit_membership (id, unit_id, user_id, role, valid_from, valid_until, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	var metadataJSON []byte
	if len(membership.Metadata()) > 0 {
		var err error
		metadataJSON, err = json.Marshal(membership.Metadata())
		if err != nil {
			return errors.NewDatabaseError("marshal metadata", err)
		}
	}

	_, err := r.db.ExecContext(ctx, query,
		membership.ID().String(),
		membership.UnitID().String(),
		membership.UserID().String(),
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
		SELECT id, unit_id, user_id, role, valid_from, valid_until, metadata, created_at, updated_at
		FROM unit_membership
		WHERE id = $1
	`

	return r.scanOneMembership(ctx, query, id.String())
}

func (r *postgresUnitMembershipRepository) FindByUnitID(ctx context.Context, unitID valueobject.UnitID, activeOnly bool) ([]*entity.UnitMembership, error) {
	query := `
		SELECT id, unit_id, user_id, role, valid_from, valid_until, metadata, created_at, updated_at
		FROM unit_membership
		WHERE unit_id = $1
	`

	if activeOnly {
		query += " AND (valid_until IS NULL OR valid_until > CURRENT_TIMESTAMP)"
	}

	query += " ORDER BY role, valid_from DESC"

	return r.scanMemberships(ctx, query, unitID.String())
}

func (r *postgresUnitMembershipRepository) FindByUserID(ctx context.Context, userID valueobject.UserID, activeOnly bool) ([]*entity.UnitMembership, error) {
	query := `
		SELECT id, unit_id, user_id, role, valid_from, valid_until, metadata, created_at, updated_at
		FROM unit_membership
		WHERE user_id = $1
	`

	if activeOnly {
		query += " AND (valid_until IS NULL OR valid_until > CURRENT_TIMESTAMP)"
	}

	query += " ORDER BY valid_from DESC"

	return r.scanMemberships(ctx, query, userID.String())
}

func (r *postgresUnitMembershipRepository) FindByUnitAndUser(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) (*entity.UnitMembership, error) {
	query := `
		SELECT id, unit_id, user_id, role, valid_from, valid_until, metadata, created_at, updated_at
		FROM unit_membership
		WHERE unit_id = $1 AND user_id = $2
		ORDER BY valid_from DESC
		LIMIT 1
	`

	return r.scanOneMembership(ctx, query, unitID.String(), userID.String())
}

func (r *postgresUnitMembershipRepository) FindByRole(ctx context.Context, unitID valueobject.UnitID, role valueobject.MembershipRole, activeOnly bool) ([]*entity.UnitMembership, error) {
	query := `
		SELECT id, unit_id, user_id, role, valid_from, valid_until, metadata, created_at, updated_at
		FROM unit_membership
		WHERE unit_id = $1 AND role = $2
	`

	if activeOnly {
		query += " AND (valid_until IS NULL OR valid_until > CURRENT_TIMESTAMP)"
	}

	query += " ORDER BY valid_from DESC"

	return r.scanMemberships(ctx, query, unitID.String(), role.String())
}

func (r *postgresUnitMembershipRepository) FindActiveAt(ctx context.Context, unitID valueobject.UnitID, at time.Time) ([]*entity.UnitMembership, error) {
	query := `
		SELECT id, unit_id, user_id, role, valid_from, valid_until, metadata, created_at, updated_at
		FROM unit_membership
		WHERE unit_id = $1
		  AND valid_from <= $2
		  AND (valid_until IS NULL OR valid_until >= $2)
		ORDER BY role, valid_from DESC
	`

	return r.scanMemberships(ctx, query, unitID.String(), at)
}

func (r *postgresUnitMembershipRepository) Update(ctx context.Context, membership *entity.UnitMembership) error {
	query := `
		UPDATE unit_membership
		SET role = $1, valid_until = $2, metadata = $3, updated_at = $4
		WHERE id = $5
	`

	var metadataJSON []byte
	if len(membership.Metadata()) > 0 {
		var err error
		metadataJSON, err = json.Marshal(membership.Metadata())
		if err != nil {
			return errors.NewDatabaseError("marshal metadata", err)
		}
	}

	_, err := r.db.ExecContext(ctx, query,
		membership.Role().String(),
		membership.ValidUntil(),
		metadataJSON,
		membership.UpdatedAt(),
		membership.ID().String(),
	)

	return err
}

func (r *postgresUnitMembershipRepository) Delete(ctx context.Context, id valueobject.MembershipID) error {
	query := `DELETE FROM unit_membership WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *postgresUnitMembershipRepository) ExistsByUnitAndUser(ctx context.Context, unitID valueobject.UnitID, userID valueobject.UserID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM unit_membership 
			WHERE unit_id = $1 AND user_id = $2 
			  AND (valid_until IS NULL OR valid_until > CURRENT_TIMESTAMP)
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, unitID.String(), userID.String()).Scan(&exists)
	return exists, err
}

func (r *postgresUnitMembershipRepository) CountByUnit(ctx context.Context, unitID valueobject.UnitID, activeOnly bool) (int, error) {
	query := `SELECT COUNT(*) FROM unit_membership WHERE unit_id = $1`

	if activeOnly {
		query += " AND (valid_until IS NULL OR valid_until > CURRENT_TIMESTAMP)"
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, unitID.String()).Scan(&count)
	return count, err
}

func (r *postgresUnitMembershipRepository) CountByRole(ctx context.Context, unitID valueobject.UnitID, role valueobject.MembershipRole) (int, error) {
	query := `
		SELECT COUNT(*) FROM unit_membership 
		WHERE unit_id = $1 AND role = $2
		  AND (valid_until IS NULL OR valid_until > CURRENT_TIMESTAMP)
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
		userIDStr    string
		role         string
		validFrom    sql.NullTime
		validUntil   sql.NullTime
		metadataJSON []byte
		createdAt    sql.NullTime
		updatedAt    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&idStr, &unitIDStr, &userIDStr, &role, &validFrom, &validUntil,
		&metadataJSON, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.buildMembership(idStr, unitIDStr, userIDStr, role, validFrom, validUntil, metadataJSON, createdAt, updatedAt)
}

// Helper: escanear múltiples membresías
func (r *postgresUnitMembershipRepository) scanMemberships(ctx context.Context, query string, args ...interface{}) ([]*entity.UnitMembership, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []*entity.UnitMembership
	for rows.Next() {
		var (
			idStr        string
			unitIDStr    string
			userIDStr    string
			role         string
			validFrom    sql.NullTime
			validUntil   sql.NullTime
			metadataJSON []byte
			createdAt    sql.NullTime
			updatedAt    sql.NullTime
		)

		if err := rows.Scan(&idStr, &unitIDStr, &userIDStr, &role, &validFrom, &validUntil, &metadataJSON, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		membership, err := r.buildMembership(idStr, unitIDStr, userIDStr, role, validFrom, validUntil, metadataJSON, createdAt, updatedAt)
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
	validFrom, validUntil sql.NullTime,
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

	var validUntilPtr *time.Time
	if validUntil.Valid {
		validUntilPtr = &validUntil.Time
	}

	return entity.ReconstructUnitMembership(
		membershipID,
		unitID,
		userID,
		roleVO,
		validFrom.Time,
		validUntilPtr,
		metadata,
		createdAt.Time,
		updatedAt.Time,
	), nil
}
