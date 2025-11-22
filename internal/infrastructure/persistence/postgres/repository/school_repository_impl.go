package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/google/uuid"
)

type postgresSchoolRepository struct {
	db *sql.DB
}

func NewPostgresSchoolRepository(db *sql.DB) repository.SchoolRepository {
	return &postgresSchoolRepository{db: db}
}

func (r *postgresSchoolRepository) Create(ctx context.Context, school *entities.School) error {
	query := `
		INSERT INTO schools (
			id, name, code, address, city, country, phone, email, metadata,
			is_active, subscription_tier, max_teachers, max_students,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.ExecContext(ctx, query,
		school.ID,
		school.Name,
		school.Code,
		school.Address,
		school.City,
		school.Country,
		school.Phone,
		school.Email,
		school.Metadata,
		school.IsActive,
		school.SubscriptionTier,
		school.MaxTeachers,
		school.MaxStudents,
		school.CreatedAt,
		school.UpdatedAt,
	)

	return err
}

func (r *postgresSchoolRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error) {
	query := `
		SELECT id, name, code, address, city, country, phone, email, metadata,
		       is_active, subscription_tier, max_teachers, max_students,
		       created_at, updated_at, deleted_at
		FROM schools
		WHERE id = $1 AND deleted_at IS NULL
	`

	school := &entities.School{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&school.ID,
		&school.Name,
		&school.Code,
		&school.Address,
		&school.City,
		&school.Country,
		&school.Phone,
		&school.Email,
		&school.Metadata,
		&school.IsActive,
		&school.SubscriptionTier,
		&school.MaxTeachers,
		&school.MaxStudents,
		&school.CreatedAt,
		&school.UpdatedAt,
		&school.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("school")
	}
	if err != nil {
		return nil, err
	}

	return school, nil
}

func (r *postgresSchoolRepository) FindByCode(ctx context.Context, code string) (*entities.School, error) {
	query := `
		SELECT id, name, code, address, city, country, phone, email, metadata,
		       is_active, subscription_tier, max_teachers, max_students,
		       created_at, updated_at, deleted_at
		FROM schools
		WHERE code = $1 AND deleted_at IS NULL
	`

	school := &entities.School{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&school.ID,
		&school.Name,
		&school.Code,
		&school.Address,
		&school.City,
		&school.Country,
		&school.Phone,
		&school.Email,
		&school.Metadata,
		&school.IsActive,
		&school.SubscriptionTier,
		&school.MaxTeachers,
		&school.MaxStudents,
		&school.CreatedAt,
		&school.UpdatedAt,
		&school.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("school")
	}
	if err != nil {
		return nil, err
	}

	return school, nil
}

func (r *postgresSchoolRepository) FindByName(ctx context.Context, name string) (*entities.School, error) {
	query := `
		SELECT id, name, code, address, city, country, phone, email, metadata,
		       is_active, subscription_tier, max_teachers, max_students,
		       created_at, updated_at, deleted_at
		FROM schools
		WHERE name = $1 AND deleted_at IS NULL
	`

	school := &entities.School{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&school.ID,
		&school.Name,
		&school.Code,
		&school.Address,
		&school.City,
		&school.Country,
		&school.Phone,
		&school.Email,
		&school.Metadata,
		&school.IsActive,
		&school.SubscriptionTier,
		&school.MaxTeachers,
		&school.MaxStudents,
		&school.CreatedAt,
		&school.UpdatedAt,
		&school.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("school")
	}
	if err != nil {
		return nil, err
	}

	return school, nil
}

func (r *postgresSchoolRepository) Update(ctx context.Context, school *entities.School) error {
	query := `
		UPDATE schools
		SET name = $1, address = $2, city = $3, country = $4, phone = $5,
		    email = $6, metadata = $7, is_active = $8, subscription_tier = $9,
		    max_teachers = $10, max_students = $11, updated_at = $12
		WHERE id = $13 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		school.Name,
		school.Address,
		school.City,
		school.Country,
		school.Phone,
		school.Email,
		school.Metadata,
		school.IsActive,
		school.SubscriptionTier,
		school.MaxTeachers,
		school.MaxStudents,
		school.UpdatedAt,
		school.ID,
	)

	return err
}

func (r *postgresSchoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE schools
		SET deleted_at = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	return err
}

func (r *postgresSchoolRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entities.School, error) {
	query := `
		SELECT id, name, code, address, city, country, phone, email, metadata,
		       is_active, subscription_tier, max_teachers, max_students,
		       created_at, updated_at, deleted_at
		FROM schools
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 1

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

	var schools []*entities.School
	for rows.Next() {
		school := &entities.School{}
		err := rows.Scan(
			&school.ID,
			&school.Name,
			&school.Code,
			&school.Address,
			&school.City,
			&school.Country,
			&school.Phone,
			&school.Email,
			&school.Metadata,
			&school.IsActive,
			&school.SubscriptionTier,
			&school.MaxTeachers,
			&school.MaxStudents,
			&school.CreatedAt,
			&school.UpdatedAt,
			&school.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		schools = append(schools, school)
	}

	return schools, rows.Err()
}

func (r *postgresSchoolRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM schools WHERE name = $1 AND deleted_at IS NULL)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	return exists, err
}

func (r *postgresSchoolRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM schools WHERE code = $1 AND deleted_at IS NULL)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	return exists, err
}
