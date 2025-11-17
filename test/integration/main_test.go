//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	_ "github.com/lib/pq"
)

// TestMain se ejecuta UNA VEZ antes de todos los tests
// Crea los containers y los destruye al final
func TestMain(m *testing.M) {
	var code int

	// Setup: crear containers
	if err := setupSharedContainers(); err != nil {
		log.Fatalf("Failed to setup containers: %v", err)
	}

	// Ejecutar todos los tests
	code = m.Run()

	// Cleanup: destruir containers
	cleanupSharedContainers()

	os.Exit(code)
}

// setupSharedContainers crea los containers de PostgreSQL y MongoDB una sola vez
func setupSharedContainers() error {
	ctx := context.Background()

	// Configurar PostgreSQL y MongoDB
	cfg := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Database: "edugo_test",
			Username: "edugo_user",
			Password: "edugo_pass",
		}).
		WithMongoDB(&containers.MongoConfig{
			Database: "edugo_test",
			Username: "edugo_admin",
			Password: "edugo_pass",
		}).
		Build()

	// Crear manager sin *testing.T (nil es válido según el código de GetManager)
	manager, err := containers.GetManager(nil, cfg)
	if err != nil {
		return err
	}
	sharedManager = manager

	// Configurar PostgreSQL
	pg := manager.PostgreSQL()
	if pg == nil {
		return fmt.Errorf("failed to get PostgreSQL container from manager")
	}

	connString, err := pg.ConnectionString(ctx)
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	sharedDB = db

	// Ejecutar migraciones básicas
	if err := runMigrations(ctx, db); err != nil {
		db.Close()
		return err
	}

	// Configurar MongoDB
	mongoDB := manager.MongoDB()
	if mongoDB != nil {
		sharedMongoDB = mongoDB.Database()
	}

	log.Println("✅ Containers compartidos inicializados correctamente")
	return nil
}

// runMigrations ejecuta las migraciones SQL necesarias para los tests
func runMigrations(ctx context.Context, db *sql.DB) error {
	// SQL para crear las tablas básicas necesarias para los tests
	migrations := []string{
		// Tabla schools
		`CREATE TABLE IF NOT EXISTS schools (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			code VARCHAR(50) NOT NULL UNIQUE,
			address TEXT,
			email VARCHAR(255),
			phone VARCHAR(50),
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)`,

		// Tabla academic_units
		`CREATE TABLE IF NOT EXISTS academic_units (
			id UUID PRIMARY KEY,
			parent_unit_id UUID REFERENCES academic_units(id),
			school_id UUID NOT NULL REFERENCES schools(id),
			unit_type VARCHAR(50) NOT NULL,
			display_name VARCHAR(255) NOT NULL,
			code VARCHAR(100),
			description TEXT,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)`,

		// Tabla users (simplificada para tests)
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Tabla unit_memberships
		`CREATE TABLE IF NOT EXISTS unit_memberships (
			id UUID PRIMARY KEY,
			unit_id UUID NOT NULL REFERENCES academic_units(id),
			user_id UUID NOT NULL,
			role VARCHAR(50) NOT NULL,
			valid_from TIMESTAMP NOT NULL,
			valid_until TIMESTAMP,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(unit_id, user_id, role)
		)`,

		// Tablas adicionales que pueden ser referenciadas
		`CREATE TABLE IF NOT EXISTS guardians (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS materials (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS subjects (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
	}

	for _, migration := range migrations {
		if _, err := db.ExecContext(ctx, migration); err != nil {
			return err
		}
	}

	log.Println("✅ Migraciones ejecutadas correctamente")
	return nil
}

// cleanupSharedContainers destruye los containers al finalizar todos los tests
func cleanupSharedContainers() {
	if sharedDB != nil {
		sharedDB.Close()
	}

	log.Println("✅ Containers compartidos destruidos correctamente")
}
