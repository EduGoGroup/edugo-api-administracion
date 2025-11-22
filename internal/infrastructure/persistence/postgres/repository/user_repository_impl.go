package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// postgresUserRepository implementa repository.UserRepository para PostgreSQL
type postgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository crea un nuevo repository de PostgreSQL
func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{db: db}
}

// Create crea un nuevo usuario
func (r *postgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, first_name, last_name, role, is_active,
			email_verified, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
		user.EmailVerified,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// FindByID busca un usuario por ID
func (r *postgresUserRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*entities.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active,
		       email_verified, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByEmail busca un usuario por email
func (r *postgresUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active,
		       email_verified, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update actualiza un usuario
func (r *postgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, role = $3, is_active = $4,
		    email_verified = $5, updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
		user.EmailVerified,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// Delete elimina un usuario (soft delete)
func (r *postgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET deleted_at = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	return err
}

// List lista usuarios con filtros
func (r *postgresUserRepository) List(
	ctx context.Context,
	filters repository.ListFilters,
) ([]*entities.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active,
		       email_verified, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 1

	if filters.Role != nil {
		query += ` AND role = $` + string(rune(argCount+'0'))
		args = append(args, *filters.Role)
		argCount++
	}

	if filters.IsActive != nil {
		query += ` AND is_active = $` + string(rune(argCount+'0'))
		args = append(args, *filters.IsActive)
		argCount++
	}

	query += ` ORDER BY created_at DESC`

	if filters.Limit > 0 {
		query += ` LIMIT $` + string(rune(argCount+'0'))
		args = append(args, filters.Limit)
		argCount++
	}

	if filters.Offset > 0 {
		query += ` OFFSET $` + string(rune(argCount+'0'))
		args = append(args, filters.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return r.scanRows(rows)
}

// ExistsByEmail verifica si existe un usuario con ese email
func (r *postgresUserRepository) ExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

// Helper methods

func (r *postgresUserRepository) scanRows(rows *sql.Rows) ([]*entities.User, error) {
	var users []*entities.User

	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.IsActive,
			&user.EmailVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, rows.Err()
}
