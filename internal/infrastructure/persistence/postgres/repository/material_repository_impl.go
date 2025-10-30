package repository

import (
	"context"
	"database/sql"

	"github.com/edugo/api-administracion/internal/domain/repository"
	"github.com/edugo/api-administracion/internal/domain/valueobject"
)

type postgresMaterialRepository struct {
	db *sql.DB
}

func NewPostgresMaterialRepository(db *sql.DB) repository.MaterialRepository {
	return &postgresMaterialRepository{db: db}
}

func (r *postgresMaterialRepository) Delete(ctx context.Context, id valueobject.MaterialID) error {
	query := `
		UPDATE materials
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *postgresMaterialRepository) Exists(ctx context.Context, id valueobject.MaterialID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM materials WHERE id = $1 AND is_deleted = false)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&exists)
	return exists, err
}
