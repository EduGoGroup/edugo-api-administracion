package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type postgresGuardianRepository struct {
	db *sql.DB
}

func NewPostgresGuardianRepository(db *sql.DB) repository.GuardianRepository {
	return &postgresGuardianRepository{db: db}
}

func (r *postgresGuardianRepository) CreateRelation(ctx context.Context, relation *entities.GuardianRelation) error {
	query := `INSERT INTO guardian_relations (id, guardian_id, student_id, relationship_type, is_active, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		relation.ID, relation.GuardianID, relation.StudentID, relation.RelationshipType,
		relation.IsActive, relation.CreatedAt, relation.UpdatedAt, relation.CreatedBy,
	)
	return err
}

func (r *postgresGuardianRepository) FindRelationByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error) {
	query := `SELECT id, guardian_id, student_id, relationship_type, is_active, created_at, updated_at, created_by
		FROM guardian_relations WHERE id = $1`
	relation := &entities.GuardianRelation{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&relation.ID, &relation.GuardianID, &relation.StudentID, &relation.RelationshipType,
		&relation.IsActive, &relation.CreatedAt, &relation.UpdatedAt, &relation.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return relation, err
}

func (r *postgresGuardianRepository) FindRelationsByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error) {
	query := `SELECT id, guardian_id, student_id, relationship_type, is_active, created_at, updated_at, created_by
		FROM guardian_relations WHERE guardian_id = $1 AND is_active = true`
	return r.scanRelations(ctx, query, guardianID)
}

func (r *postgresGuardianRepository) FindRelationsByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error) {
	query := `SELECT id, guardian_id, student_id, relationship_type, is_active, created_at, updated_at, created_by
		FROM guardian_relations WHERE student_id = $1 AND is_active = true`
	return r.scanRelations(ctx, query, studentID)
}

func (r *postgresGuardianRepository) UpdateRelation(ctx context.Context, relation *entities.GuardianRelation) error {
	query := `UPDATE guardian_relations SET relationship_type = $1, is_active = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, relation.RelationshipType, relation.IsActive, relation.UpdatedAt, relation.ID)
	return err
}

func (r *postgresGuardianRepository) DeleteRelation(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE guardian_relations SET is_active = false, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *postgresGuardianRepository) scanRelations(ctx context.Context, query string, id uuid.UUID) ([]*entities.GuardianRelation, error) {
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var relations []*entities.GuardianRelation
	for rows.Next() {
		relation := &entities.GuardianRelation{}
		err := rows.Scan(&relation.ID, &relation.GuardianID, &relation.StudentID, &relation.RelationshipType,
			&relation.IsActive, &relation.CreatedAt, &relation.UpdatedAt, &relation.CreatedBy)
		if err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}
	return relations, rows.Err()
}

// Alias methods for compatibility with services
func (r *postgresGuardianRepository) Create(ctx context.Context, relation *entities.GuardianRelation) error {
	return r.CreateRelation(ctx, relation)
}

func (r *postgresGuardianRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error) {
	return r.FindRelationByID(ctx, id)
}

func (r *postgresGuardianRepository) FindByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error) {
	return r.FindRelationsByGuardian(ctx, guardianID)
}

func (r *postgresGuardianRepository) FindByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error) {
	return r.FindRelationsByStudent(ctx, studentID)
}

func (r *postgresGuardianRepository) Update(ctx context.Context, relation *entities.GuardianRelation) error {
	return r.UpdateRelation(ctx, relation)
}

func (r *postgresGuardianRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteRelation(ctx, id)
}

func (r *postgresGuardianRepository) ExistsActiveRelation(ctx context.Context, guardianID, studentID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM guardian_relations WHERE guardian_id = $1 AND student_id = $2 AND is_active = true)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, guardianID, studentID).Scan(&exists)
	return exists, err
}
