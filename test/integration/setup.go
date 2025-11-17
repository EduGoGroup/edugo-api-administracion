//go:build integration

package integration

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
)

// getMigrationScripts retorna paths a migraciones de infrastructure
func getMigrationScripts() []string {
	// Obtener path al directorio raíz del proyecto
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	// Intentar primero desde GOPATH (módulos descargados)
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}

	// Path a migraciones en pkg/mod (go modules cache)
	modPath := filepath.Join(gopath, "pkg", "mod", "github.com", "!edu!go!group", "edugo-infrastructure", "postgres@v0.7.1", "migrations")

	// Si no existe en mod, intentar desde vendor
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		modPath = filepath.Join(projectRoot, "vendor", "github.com", "EduGoGroup", "edugo-infrastructure", "postgres", "migrations")
	}

	// Migraciones necesarias para api-administracion (001-004 son las básicas)
	migrations := []string{
		"001_create_users.up.sql",
		"002_create_schools.up.sql",
		"003_create_academic_units.up.sql",
		"004_create_memberships.up.sql",
	}

	var fullPaths []string
	for _, migration := range migrations {
		fullPaths = append(fullPaths, filepath.Join(modPath, migration))
	}

	return fullPaths
}

// setupTestDB crea una instancia de PostgreSQL para tests usando shared/testing
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	// Obtener migraciones de infrastructure
	migrationScripts := getMigrationScripts()

	// Log para debug
	t.Logf("Usando %d migraciones desde infrastructure v0.7.1", len(migrationScripts))
	for i, script := range migrationScripts {
		t.Logf("  [%d] %s", i+1, filepath.Base(script))
	}

	// Configurar PostgreSQL CON migraciones de infrastructure
	cfg := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Database:    "edugo_test",
			Username:    "edugo_user",
			Password:    "edugo_pass",
			InitScripts: migrationScripts,
		}).
		Build()

	manager, err := containers.GetManager(t, cfg)
	if err != nil {
		t.Fatalf("Failed to get manager: %v", err)
	}

	pg := manager.PostgreSQL()
	if pg == nil {
		t.Fatal("Failed to get PostgreSQL container")
	}

	connString, err := pg.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get Postgres connection string: %v", err)
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		t.Fatalf("Failed to connect to Postgres: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("Failed to ping Postgres: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

// setupTestMongoDB crea una instancia de MongoDB para tests usando shared/testing
func setupTestMongoDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()

	// Configurar MongoDB
	cfg := containers.NewConfig().
		WithMongoDB(&containers.MongoConfig{
			Database: "edugo_test",
			Username: "edugo_admin",
			Password: "edugo_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, cfg)
	if err != nil {
		t.Fatalf("Failed to get manager: %v", err)
	}

	mongoDB := manager.MongoDB()
	if mongoDB == nil {
		t.Fatal("Failed to get MongoDB container")
	}

	db := mongoDB.Database()

	cleanup := func() {
		// Limpiar colecciones al terminar
		mongoDB.DropAllCollections(ctx)
	}

	return db, cleanup
}

// setupTestDBWithMigrations crea una BD con las migraciones aplicadas
// Ahora usa migraciones centralizadas de infrastructure v0.7.1
func setupTestDBWithMigrations(t *testing.T) (*sql.DB, func()) {
	return setupTestDB(t)
}
