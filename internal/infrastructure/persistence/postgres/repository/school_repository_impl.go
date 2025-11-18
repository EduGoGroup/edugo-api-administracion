package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

type postgresSchoolRepository struct {
	db *sql.DB
}

func NewPostgresSchoolRepository(db *sql.DB) repository.SchoolRepository {
	return &postgresSchoolRepository{db: db}
}

func (r *postgresSchoolRepository) Create(ctx context.Context, school *entity.School) error {
	query := `
		INSERT INTO schools (id, name, code, address, email, phone, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	var contactEmail *string
	if school.ContactEmail() != nil {
		email := school.ContactEmail().String()
		contactEmail = &email
	}

	// Serializar metadata (siempre enviar al menos {} para JSONB)
	metadata := school.Metadata()
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return errors.NewDatabaseError("marshal metadata", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		school.ID().String(),
		school.Name(),
		school.Code(),
		school.Address(),
		contactEmail,
		school.ContactPhone(),
		metadataJSON,
		school.CreatedAt(),
		school.UpdatedAt(),
	)

	return err
}

func (r *postgresSchoolRepository) FindByID(ctx context.Context, id valueobject.SchoolID) (*entity.School, error) {
	query := `
		SELECT id, name, code, address, email, phone, metadata, created_at, updated_at
		FROM schools
		WHERE id = $1 AND deleted_at IS NULL
	`

	var (
		idStr        string
		name         string
		code         string
		address      string
		contactEmail sql.NullString
		contactPhone string
		metadataJSON []byte
		createdAt    sql.NullTime
		updatedAt    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &name, &code, &address, &contactEmail, &contactPhone, &metadataJSON, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("school")
	}
	if err != nil {
		return nil, err
	}

	return r.scanSchool(idStr, name, code, address, contactEmail, contactPhone, metadataJSON, createdAt, updatedAt)
}

func (r *postgresSchoolRepository) FindByCode(ctx context.Context, code string) (*entity.School, error) {
	query := `
		SELECT id, name, code, address, email, phone, metadata, created_at, updated_at
		FROM schools
		WHERE code = $1
	`

	var (
		idStr        string
		name         string
		codeStr      string
		address      string
		contactEmail sql.NullString
		contactPhone string
		metadataJSON []byte
		createdAt    sql.NullTime
		updatedAt    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&idStr, &name, &codeStr, &address, &contactEmail, &contactPhone, &metadataJSON, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.scanSchool(idStr, name, codeStr, address, contactEmail, contactPhone, metadataJSON, createdAt, updatedAt)
}

func (r *postgresSchoolRepository) FindByName(ctx context.Context, name string) (*entity.School, error) {
	query := `
		SELECT id, name, code, address, email, phone, metadata, created_at, updated_at
		FROM schools
		WHERE name = $1
	`

	var (
		idStr        string
		nameStr      string
		code         string
		address      string
		contactEmail sql.NullString
		contactPhone string
		metadataJSON []byte
		createdAt    sql.NullTime
		updatedAt    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&idStr, &nameStr, &code, &address, &contactEmail, &contactPhone, &metadataJSON, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.scanSchool(idStr, nameStr, code, address, contactEmail, contactPhone, metadataJSON, createdAt, updatedAt)
}

func (r *postgresSchoolRepository) Update(ctx context.Context, school *entity.School) error {
	query := `
		UPDATE schools
		SET name = $1, address = $2, email = $3, phone = $4, metadata = $5, updated_at = $6
		WHERE id = $7
	`

	var contactEmail *string
	if school.ContactEmail() != nil {
		email := school.ContactEmail().String()
		contactEmail = &email
	}

	// Asegurar que metadata nunca sea nil para JSONB
	metadata := school.Metadata()
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return errors.NewDatabaseError("marshal metadata", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		school.Name(),
		school.Address(),
		contactEmail,
		school.ContactPhone(),
		metadataJSON,
		school.UpdatedAt(),
		school.ID().String(),
	)

	return err
}

func (r *postgresSchoolRepository) Delete(ctx context.Context, id valueobject.SchoolID) error {
	query := `DELETE FROM schools WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *postgresSchoolRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entity.School, error) {
	query := `
		SELECT id, name, code, address, email, phone, metadata, created_at, updated_at
		FROM schools
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schools []*entity.School
	for rows.Next() {
		var (
			idStr        string
			name         string
			code         string
			address      string
			contactEmail sql.NullString
			contactPhone string
			metadataJSON []byte
			createdAt    sql.NullTime
			updatedAt    sql.NullTime
		)

		if err := rows.Scan(&idStr, &name, &code, &address, &contactEmail, &contactPhone, &metadataJSON, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		school, err := r.scanSchool(idStr, name, code, address, contactEmail, contactPhone, metadataJSON, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}

		schools = append(schools, school)
	}

	return schools, rows.Err()
}

func (r *postgresSchoolRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM schools WHERE name = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	return exists, err
}

func (r *postgresSchoolRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM schools WHERE code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	return exists, err
}

// Helper para escanear una escuela
func (r *postgresSchoolRepository) scanSchool(
	idStr, name, code, address string,
	contactEmail sql.NullString,
	contactPhone string,
	metadataJSON []byte,
	createdAt, updatedAt sql.NullTime,
) (*entity.School, error) {
	schoolID, err := valueobject.SchoolIDFromString(idStr)
	if err != nil {
		return nil, errors.NewDatabaseError("parse school ID", err)
	}

	var email *valueobject.Email
	if contactEmail.Valid && contactEmail.String != "" {
		e, err := valueobject.NewEmail(contactEmail.String)
		if err != nil {
			return nil, errors.NewDatabaseError("parse email", err)
		}
		email = &e
	}

	var metadata map[string]interface{}
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			return nil, errors.NewDatabaseError("unmarshal metadata", err)
		}
	}

	return entity.ReconstructSchool(
		schoolID,
		name,
		code,
		address,
		email,
		contactPhone,
		metadata,
		createdAt.Time,
		updatedAt.Time,
	), nil
}
