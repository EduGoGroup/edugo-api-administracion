package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type postgresSubjectRepository struct {
	db *sql.DB
}

func NewPostgresSubjectRepository(db *sql.DB) repository.SubjectRepository {
	return &postgresSubjectRepository{db: db}
}

func (r *postgresSubjectRepository) Create(ctx context.Context, subject *entities.Subject) error {
	query := `
		INSERT INTO subjects (id, name, description, metadata, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		subject.ID, subject.Name, subject.Description, subject.Metadata,
		subject.IsActive, subject.CreatedAt, subject.UpdatedAt,
	)
	return err
}

func (r *postgresSubjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Subject, error) {
	query := `
		SELECT id, name, description, metadata, is_active, created_at, updated_at
		FROM subjects WHERE id = $1
	`
	subject := &entities.Subject{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subject.ID, &subject.Name, &subject.Description, &subject.Metadata,
		&subject.IsActive, &subject.CreatedAt, &subject.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return subject, err
}

func (r *postgresSubjectRepository) Update(ctx context.Context, subject *entities.Subject) error {
	query := `
		UPDATE subjects SET name = $1, description = $2, metadata = $3,
		       is_active = $4, updated_at = $5 WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		subject.Name, subject.Description, subject.Metadata,
		subject.IsActive, subject.UpdatedAt, subject.ID,
	)
	return err
}

func (r *postgresSubjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE subjects SET is_active = false, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *postgresSubjectRepository) List(ctx context.Context) ([]*entities.Subject, error) {
	query := `
		SELECT id, name, description, metadata, is_active, created_at, updated_at
		FROM subjects WHERE is_active = true ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var subjects []*entities.Subject
	for rows.Next() {
		subject := &entities.Subject{}
		err := rows.Scan(
			&subject.ID, &subject.Name, &subject.Description, &subject.Metadata,
			&subject.IsActive, &subject.CreatedAt, &subject.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, rows.Err()
}

// FindBySchoolID lista materias activas filtradas por school_id
// Nota: La entidad Subject actualmente no tiene campo school_id en la base de datos
// Este método retorna todas las materias activas por ahora hasta que se agregue la columna
func (r *postgresSubjectRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID) ([]*entities.Subject, error) {
	// Por ahora retornamos todas las materias activas ya que no existe school_id en la tabla
	query := `
		SELECT id, name, description, metadata, is_active, created_at, updated_at
		FROM subjects WHERE is_active = true ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var subjects []*entities.Subject
	for rows.Next() {
		subject := &entities.Subject{}
		err := rows.Scan(
			&subject.ID, &subject.Name, &subject.Description, &subject.Metadata,
			&subject.IsActive, &subject.CreatedAt, &subject.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, rows.Err()
}

func (r *postgresSubjectRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM subjects WHERE name = $1 AND is_active = true)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	return exists, err
}
