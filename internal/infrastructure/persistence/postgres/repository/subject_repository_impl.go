package repository

import (
	"context"
	"database/sql"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
)

type postgresSubjectRepository struct {
	db *sql.DB
}

func NewPostgresSubjectRepository(db *sql.DB) repository.SubjectRepository {
	return &postgresSubjectRepository{db: db}
}

func (r *postgresSubjectRepository) Create(ctx context.Context, subject *entity.Subject) error {
	query := `
		INSERT INTO subjects (id, name, description, metadata, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		subject.ID().String(),
		subject.Name(),
		subject.Description(),
		subject.Metadata(),
		subject.IsActive(),
		subject.CreatedAt(),
		subject.UpdatedAt(),
	)

	return err
}

func (r *postgresSubjectRepository) FindByID(ctx context.Context, id valueobject.SubjectID) (*entity.Subject, error) {
	query := `
		SELECT id, name, description, metadata, is_active, created_at, updated_at
		FROM subjects
		WHERE id = $1
	`

	var (
		idStr       string
		name        string
		description sql.NullString
		metadata    sql.NullString
		isActive    bool
		createdAt   sql.NullTime
		updatedAt   sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &name, &description, &metadata, &isActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	subjectID, _ := valueobject.SubjectIDFromString(idStr)
	return entity.ReconstructSubject(
		subjectID,
		name,
		description.String,
		metadata.String,
		isActive,
		createdAt.Time,
		updatedAt.Time,
	), nil
}

func (r *postgresSubjectRepository) Update(ctx context.Context, subject *entity.Subject) error {
	query := `
		UPDATE subjects
		SET name = $1, description = $2, metadata = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		subject.Name(),
		subject.Description(),
		subject.Metadata(),
		subject.UpdatedAt(),
		subject.ID().String(),
	)

	return err
}

func (r *postgresSubjectRepository) Delete(ctx context.Context, id valueobject.SubjectID) error {
	query := `UPDATE subjects SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}
