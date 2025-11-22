package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/google/uuid"
)

type postgresMaterialRepository struct {
	db *sql.DB
}

func NewPostgresMaterialRepository(db *sql.DB) repository.MaterialRepository {
	return &postgresMaterialRepository{db: db}
}

func (r *postgresMaterialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE materials
		SET is_deleted = true, deleted_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *postgresMaterialRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM materials WHERE id = $1 AND is_deleted = false)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}
