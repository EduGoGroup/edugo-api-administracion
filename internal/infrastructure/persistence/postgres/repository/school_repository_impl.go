package repository

import (
	"context"
	"database/sql"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
)

type postgresSchoolRepository struct {
	db *sql.DB
}

func NewPostgresSchoolRepository(db *sql.DB) repository.SchoolRepository {
	return &postgresSchoolRepository{db: db}
}

func (r *postgresSchoolRepository) Create(ctx context.Context, school *entity.School) error {
	query := `
		INSERT INTO schools (id, name, address, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		school.ID().String(),
		school.Name(),
		school.Address(),
		school.IsActive(),
		school.CreatedAt(),
		school.UpdatedAt(),
	)

	return err
}

func (r *postgresSchoolRepository) FindByID(ctx context.Context, id valueobject.SchoolID) (*entity.School, error) {
	query := `
		SELECT id, name, address, is_active, created_at, updated_at
		FROM schools
		WHERE id = $1
	`

	var (
		idStr     string
		name      string
		address   string
		isActive  bool
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&idStr, &name, &address, &isActive, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	schoolID, _ := valueobject.SchoolIDFromString(idStr)
	return entity.ReconstructSchool(schoolID, name, address, isActive, createdAt.Time, updatedAt.Time), nil
}

func (r *postgresSchoolRepository) FindByName(ctx context.Context, name string) (*entity.School, error) {
	query := `
		SELECT id, name, address, is_active, created_at, updated_at
		FROM schools
		WHERE name = $1
	`

	var (
		idStr     string
		nameStr   string
		address   string
		isActive  bool
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, name).Scan(&idStr, &nameStr, &address, &isActive, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	schoolID, _ := valueobject.SchoolIDFromString(idStr)
	return entity.ReconstructSchool(schoolID, nameStr, address, isActive, createdAt.Time, updatedAt.Time), nil
}

func (r *postgresSchoolRepository) Update(ctx context.Context, school *entity.School) error {
	query := `
		UPDATE schools
		SET name = $1, address = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		school.Name(),
		school.Address(),
		school.IsActive(),
		school.UpdatedAt(),
		school.ID().String(),
	)

	return err
}

func (r *postgresSchoolRepository) Delete(ctx context.Context, id valueobject.SchoolID) error {
	query := `UPDATE schools SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *postgresSchoolRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entity.School, error) {
	query := `
		SELECT id, name, address, is_active, created_at, updated_at
		FROM schools
		WHERE 1=1
	`

	args := []interface{}{}
	if filters.IsActive != nil {
		query += ` AND is_active = $1`
		args = append(args, *filters.IsActive)
	}

	query += ` ORDER BY created_at DESC`

	if filters.Limit > 0 {
		query += ` LIMIT $2`
		args = append(args, filters.Limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schools []*entity.School
	for rows.Next() {
		var (
			idStr     string
			name      string
			address   string
			isActive  bool
			createdAt sql.NullTime
			updatedAt sql.NullTime
		)

		if err := rows.Scan(&idStr, &name, &address, &isActive, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		schoolID, _ := valueobject.SchoolIDFromString(idStr)
		schools = append(schools, entity.ReconstructSchool(schoolID, name, address, isActive, createdAt.Time, updatedAt.Time))
	}

	return schools, rows.Err()
}

func (r *postgresSchoolRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM schools WHERE name = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	return exists, err
}
