package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type postgresUnitMembershipRepository struct {
	db *sql.DB
}

func NewPostgresUnitMembershipRepository(db *sql.DB) repository.UnitMembershipRepository {
	return &postgresUnitMembershipRepository{db: db}
}

func (r *postgresUnitMembershipRepository) Create(ctx context.Context, membership *entities.Membership) error {
	query := `INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, withdrawn_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		membership.ID, membership.UserID, membership.SchoolID, membership.AcademicUnitID,
		membership.Role, membership.Metadata, membership.IsActive, membership.EnrolledAt,
		membership.WithdrawnAt, membership.CreatedAt, membership.UpdatedAt,
	)
	return err
}

func (r *postgresUnitMembershipRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error) {
	query := `SELECT id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, withdrawn_at, created_at, updated_at
		FROM memberships WHERE id = $1`
	membership := &entities.Membership{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&membership.ID, &membership.UserID, &membership.SchoolID, &membership.AcademicUnitID,
		&membership.Role, &membership.Metadata, &membership.IsActive, &membership.EnrolledAt,
		&membership.WithdrawnAt, &membership.CreatedAt, &membership.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return membership, err
}

func (r *postgresUnitMembershipRepository) FindByUserAndUnit(ctx context.Context, userID, unitID uuid.UUID) (*entities.Membership, error) {
	query := `SELECT id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, withdrawn_at, created_at, updated_at
		FROM memberships WHERE user_id = $1 AND academic_unit_id = $2 AND is_active = true`
	membership := &entities.Membership{}
	err := r.db.QueryRowContext(ctx, query, userID, unitID).Scan(
		&membership.ID, &membership.UserID, &membership.SchoolID, &membership.AcademicUnitID,
		&membership.Role, &membership.Metadata, &membership.IsActive, &membership.EnrolledAt,
		&membership.WithdrawnAt, &membership.CreatedAt, &membership.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return membership, err
}

func (r *postgresUnitMembershipRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error) {
	query := `SELECT id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, withdrawn_at, created_at, updated_at
		FROM memberships WHERE user_id = $1 AND is_active = true ORDER BY enrolled_at DESC`
	return r.scanMemberships(ctx, query, userID)
}

func (r *postgresUnitMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error) {
	query := `SELECT id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, withdrawn_at, created_at, updated_at
		FROM memberships WHERE academic_unit_id = $1 AND is_active = true ORDER BY enrolled_at DESC`
	return r.scanMemberships(ctx, query, unitID)
}

func (r *postgresUnitMembershipRepository) FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool) ([]*entities.Membership, error) {
	var query string
	if activeOnly {
		query = `SELECT id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, withdrawn_at, created_at, updated_at
			FROM memberships WHERE academic_unit_id = $1 AND role = $2 AND is_active = true AND withdrawn_at IS NULL ORDER BY enrolled_at DESC`
	} else {
		query = `SELECT id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, withdrawn_at, created_at, updated_at
			FROM memberships WHERE academic_unit_id = $1 AND role = $2 ORDER BY enrolled_at DESC`
	}
	return r.scanMembershipsWithRole(ctx, query, unitID, role)
}

func (r *postgresUnitMembershipRepository) scanMembershipsWithRole(ctx context.Context, query string, unitID uuid.UUID, role string) ([]*entities.Membership, error) {
	rows, err := r.db.QueryContext(ctx, query, unitID, role)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var memberships []*entities.Membership
	for rows.Next() {
		membership := &entities.Membership{}
		err := rows.Scan(
			&membership.ID, &membership.UserID, &membership.SchoolID, &membership.AcademicUnitID,
			&membership.Role, &membership.Metadata, &membership.IsActive, &membership.EnrolledAt,
			&membership.WithdrawnAt, &membership.CreatedAt, &membership.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, membership)
	}
	return memberships, rows.Err()
}

func (r *postgresUnitMembershipRepository) Update(ctx context.Context, membership *entities.Membership) error {
	query := `UPDATE memberships SET role = $1, metadata = $2, is_active = $3, withdrawn_at = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query,
		membership.Role, membership.Metadata, membership.IsActive,
		membership.WithdrawnAt, membership.UpdatedAt, membership.ID,
	)
	return err
}

func (r *postgresUnitMembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE memberships SET is_active = false, withdrawn_at = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	return err
}

func (r *postgresUnitMembershipRepository) scanMemberships(ctx context.Context, query string, id uuid.UUID) ([]*entities.Membership, error) {
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var memberships []*entities.Membership
	for rows.Next() {
		membership := &entities.Membership{}
		err := rows.Scan(
			&membership.ID, &membership.UserID, &membership.SchoolID, &membership.AcademicUnitID,
			&membership.Role, &membership.Metadata, &membership.IsActive, &membership.EnrolledAt,
			&membership.WithdrawnAt, &membership.CreatedAt, &membership.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, membership)
	}
	return memberships, rows.Err()
}

func (r *postgresUnitMembershipRepository) ExistsByUnitAndUser(ctx context.Context, unitID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM memberships WHERE academic_unit_id = $1 AND user_id = $2 AND is_active = true)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, unitID, userID).Scan(&exists)
	return exists, err
}

func (r *postgresUnitMembershipRepository) FindByUserAndSchool(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error) {
	query := `SELECT id, user_id, school_id, academic_unit_id, role, metadata, is_active, enrolled_at, withdrawn_at, created_at, updated_at
		FROM memberships WHERE user_id = $1 AND school_id = $2 AND is_active = true LIMIT 1`
	membership := &entities.Membership{}
	err := r.db.QueryRowContext(ctx, query, userID, schoolID).Scan(
		&membership.ID, &membership.UserID, &membership.SchoolID, &membership.AcademicUnitID,
		&membership.Role, &membership.Metadata, &membership.IsActive, &membership.EnrolledAt,
		&membership.WithdrawnAt, &membership.CreatedAt, &membership.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return membership, err
}
