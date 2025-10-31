package repository

import (
	"context"
	"database/sql"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/types/enum"
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
func (r *postgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (
			id, email, first_name, last_name, role, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID().String(),
		user.Email().String(),
		user.FirstName(),
		user.LastName(),
		user.Role().String(),
		user.IsActive(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	return err
}

// FindByID busca un usuario por ID
func (r *postgresUserRepository) FindByID(
	ctx context.Context,
	id valueobject.UserID,
) (*entity.User, error) {
	query := `
		SELECT id, email, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var (
		idStr      string
		email      string
		firstName  string
		lastName   string
		roleStr    string
		isActive   bool
		createdAt  sql.NullTime
		updatedAt  sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &email, &firstName, &lastName, &roleStr, &isActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.scanToEntity(idStr, email, firstName, lastName, roleStr, isActive, createdAt, updatedAt)
}

// FindByEmail busca un usuario por email
func (r *postgresUserRepository) FindByEmail(
	ctx context.Context,
	email valueobject.Email,
) (*entity.User, error) {
	query := `
		SELECT id, email, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var (
		idStr      string
		emailStr   string
		firstName  string
		lastName   string
		roleStr    string
		isActive   bool
		createdAt  sql.NullTime
		updatedAt  sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&idStr, &emailStr, &firstName, &lastName, &roleStr, &isActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.scanToEntity(idStr, emailStr, firstName, lastName, roleStr, isActive, createdAt, updatedAt)
}

// Update actualiza un usuario
func (r *postgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, role = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName(),
		user.LastName(),
		user.Role().String(),
		user.IsActive(),
		user.UpdatedAt(),
		user.ID().String(),
	)

	return err
}

// Delete elimina un usuario (soft delete)
func (r *postgresUserRepository) Delete(ctx context.Context, id valueobject.UserID) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

// List lista usuarios con filtros
func (r *postgresUserRepository) List(
	ctx context.Context,
	filters repository.ListFilters,
) ([]*entity.User, error) {
	query := `
		SELECT id, email, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE 1=1
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
	defer rows.Close()

	return r.scanRows(rows)
}

// ExistsByEmail verifica si existe un usuario con ese email
func (r *postgresUserRepository) ExistsByEmail(
	ctx context.Context,
	email valueobject.Email,
) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(&exists)
	return exists, err
}

// Helper methods

func (r *postgresUserRepository) scanRows(rows *sql.Rows) ([]*entity.User, error) {
	var users []*entity.User

	for rows.Next() {
		var (
			idStr      string
			email      string
			firstName  string
			lastName   string
			roleStr    string
			isActive   bool
			createdAt  sql.NullTime
			updatedAt  sql.NullTime
		)

		err := rows.Scan(&idStr, &email, &firstName, &lastName, &roleStr, &isActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		user, err := r.scanToEntity(idStr, email, firstName, lastName, roleStr, isActive, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *postgresUserRepository) scanToEntity(
	idStr string,
	emailStr string,
	firstName string,
	lastName string,
	roleStr string,
	isActive bool,
	createdAt sql.NullTime,
	updatedAt sql.NullTime,
) (*entity.User, error) {
	// Parsear UserID
	userID, err := valueobject.UserIDFromString(idStr)
	if err != nil {
		return nil, err
	}

	// Parsear Email
	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	// Parsear Role
	role := enum.SystemRole(roleStr)

	// Reconstruir entidad
	return entity.ReconstructUser(
		userID,
		email,
		firstName,
		lastName,
		role,
		isActive,
		createdAt.Time,
		updatedAt.Time,
	), nil
}
