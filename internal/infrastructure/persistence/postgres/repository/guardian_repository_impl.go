package repository

import (
	"context"
	"database/sql"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/types"
)

// postgresGuardianRepository implementa repository.GuardianRepository para PostgreSQL
type postgresGuardianRepository struct {
	db *sql.DB
}

// NewPostgresGuardianRepository crea un nuevo repository de PostgreSQL
func NewPostgresGuardianRepository(db *sql.DB) repository.GuardianRepository {
	return &postgresGuardianRepository{db: db}
}

// Create crea una nueva relación guardian-estudiante
func (r *postgresGuardianRepository) Create(ctx context.Context, relation *entity.GuardianRelation) error {
	query := `
		INSERT INTO guardian_relations (
			id, guardian_id, student_id, relationship_type,
			is_active, created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		relation.ID().String(),
		relation.GuardianID().String(),
		relation.StudentID().String(),
		relation.RelationshipType().String(),
		relation.IsActive(),
		relation.CreatedAt(),
		relation.UpdatedAt(),
		relation.CreatedBy(),
	)

	return err
}

// FindByID busca una relación por ID
func (r *postgresGuardianRepository) FindByID(ctx context.Context, id types.UUID) (*entity.GuardianRelation, error) {
	query := `
		SELECT id, guardian_id, student_id, relationship_type,
		       is_active, created_at, updated_at, created_by
		FROM guardian_relations
		WHERE id = $1
	`

	var (
		idStr               string
		guardianIDStr       string
		studentIDStr        string
		relationshipTypeStr string
		isActive            bool
		createdAt           sql.NullTime
		updatedAt           sql.NullTime
		createdBy           string
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr,
		&guardianIDStr,
		&studentIDStr,
		&relationshipTypeStr,
		&isActive,
		&createdAt,
		&updatedAt,
		&createdBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.scanToEntity(
		idStr,
		guardianIDStr,
		studentIDStr,
		relationshipTypeStr,
		isActive,
		createdAt,
		updatedAt,
		createdBy,
	)
}

// FindByGuardianAndStudent busca una relación por guardian y estudiante
func (r *postgresGuardianRepository) FindByGuardianAndStudent(
	ctx context.Context,
	guardianID valueobject.GuardianID,
	studentID valueobject.StudentID,
) (*entity.GuardianRelation, error) {
	query := `
		SELECT id, guardian_id, student_id, relationship_type,
		       is_active, created_at, updated_at, created_by
		FROM guardian_relations
		WHERE guardian_id = $1 AND student_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var (
		idStr               string
		guardianIDStr       string
		studentIDStr        string
		relationshipTypeStr string
		isActive            bool
		createdAt           sql.NullTime
		updatedAt           sql.NullTime
		createdBy           string
	)

	err := r.db.QueryRowContext(ctx, query, guardianID.String(), studentID.String()).Scan(
		&idStr,
		&guardianIDStr,
		&studentIDStr,
		&relationshipTypeStr,
		&isActive,
		&createdAt,
		&updatedAt,
		&createdBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.scanToEntity(
		idStr,
		guardianIDStr,
		studentIDStr,
		relationshipTypeStr,
		isActive,
		createdAt,
		updatedAt,
		createdBy,
	)
}

// FindByGuardian busca todas las relaciones de un guardian
func (r *postgresGuardianRepository) FindByGuardian(
	ctx context.Context,
	guardianID valueobject.GuardianID,
) ([]*entity.GuardianRelation, error) {
	query := `
		SELECT id, guardian_id, student_id, relationship_type,
		       is_active, created_at, updated_at, created_by
		FROM guardian_relations
		WHERE guardian_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, guardianID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// FindByStudent busca todas las relaciones de un estudiante
func (r *postgresGuardianRepository) FindByStudent(
	ctx context.Context,
	studentID valueobject.StudentID,
) ([]*entity.GuardianRelation, error) {
	query := `
		SELECT id, guardian_id, student_id, relationship_type,
		       is_active, created_at, updated_at, created_by
		FROM guardian_relations
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, studentID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// Update actualiza una relación existente
func (r *postgresGuardianRepository) Update(ctx context.Context, relation *entity.GuardianRelation) error {
	query := `
		UPDATE guardian_relations
		SET relationship_type = $1, is_active = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query,
		relation.RelationshipType().String(),
		relation.IsActive(),
		relation.UpdatedAt(),
		relation.ID().String(),
	)

	return err
}

// Delete elimina una relación (soft delete)
func (r *postgresGuardianRepository) Delete(ctx context.Context, id types.UUID) error {
	query := `
		UPDATE guardian_relations
		SET is_active = false, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

// ExistsActiveRelation verifica si existe una relación activa
func (r *postgresGuardianRepository) ExistsActiveRelation(
	ctx context.Context,
	guardianID valueobject.GuardianID,
	studentID valueobject.StudentID,
) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM guardian_relations
			WHERE guardian_id = $1 AND student_id = $2 AND is_active = true
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, guardianID.String(), studentID.String()).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Helper methods

func (r *postgresGuardianRepository) scanRows(rows *sql.Rows) ([]*entity.GuardianRelation, error) {
	var relations []*entity.GuardianRelation

	for rows.Next() {
		var (
			idStr               string
			guardianIDStr       string
			studentIDStr        string
			relationshipTypeStr string
			isActive            bool
			createdAt           sql.NullTime
			updatedAt           sql.NullTime
			createdBy           string
		)

		err := rows.Scan(
			&idStr,
			&guardianIDStr,
			&studentIDStr,
			&relationshipTypeStr,
			&isActive,
			&createdAt,
			&updatedAt,
			&createdBy,
		)
		if err != nil {
			return nil, err
		}

		relation, err := r.scanToEntity(
			idStr,
			guardianIDStr,
			studentIDStr,
			relationshipTypeStr,
			isActive,
			createdAt,
			updatedAt,
			createdBy,
		)
		if err != nil {
			return nil, err
		}

		relations = append(relations, relation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return relations, nil
}

func (r *postgresGuardianRepository) scanToEntity(
	idStr string,
	guardianIDStr string,
	studentIDStr string,
	relationshipTypeStr string,
	isActive bool,
	createdAt sql.NullTime,
	updatedAt sql.NullTime,
	createdBy string,
) (*entity.GuardianRelation, error) {
	// Parsear UUIDs
	id, err := types.ParseUUID(idStr)
	if err != nil {
		return nil, err
	}

	guardianID, err := valueobject.GuardianIDFromString(guardianIDStr)
	if err != nil {
		return nil, err
	}

	studentID, err := valueobject.StudentIDFromString(studentIDStr)
	if err != nil {
		return nil, err
	}

	relationshipType, err := valueobject.NewRelationshipType(relationshipTypeStr)
	if err != nil {
		return nil, err
	}

	// Reconstruir entidad
	return entity.ReconstructGuardianRelation(
		id,
		guardianID,
		studentID,
		relationshipType,
		isActive,
		createdAt.Time,
		updatedAt.Time,
		createdBy,
	), nil
}
