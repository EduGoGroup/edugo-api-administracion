//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
)

// getMigrationScripts retorna paths a migraciones de infrastructure
// Retorna error si no encuentra las migraciones para facilitar debugging
func getMigrationScripts(t *testing.T) ([]string, error) {
	// Obtener path al directorio raíz del proyecto (solo para vendor fallback)
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	// Usar go/build para obtener GOPATH de forma cross-platform (Windows compatible)
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	// IMPORTANTE: Versión hardcodeada v0.9.0
	// Si actualizas edugo-infrastructure en go.mod, actualiza esta versión aquí también
	const infrastructureVersion = "v0.9.0"

	// Path a migraciones en pkg/mod (go modules cache)
	modPath := filepath.Join(gopath, "pkg", "mod", "github.com", "!edu!go!group", "edugo-infrastructure", "postgres@"+infrastructureVersion, "migrations")

	// Si no existe en mod, intentar desde vendor (fallback)
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		t.Logf("Migraciones no encontradas en pkg/mod, intentando vendor...")
		modPath = filepath.Join(projectRoot, "vendor", "github.com", "EduGoGroup", "edugo-infrastructure", "postgres", "migrations")
	} else {
		t.Logf("Usando migraciones desde: %s", modPath)
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
		fullPath := filepath.Join(modPath, migration)

		// Validar que cada archivo existe antes de agregarlo
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("migración no encontrada: %s (buscado en: %s)", migration, modPath)
		}

		fullPaths = append(fullPaths, fullPath)
	}

	// Agregar migración 013 ltree desde el proyecto local (Sprint-03)
	localMigration := filepath.Join(projectRoot, "docs", "isolated", "04-Implementation", "Sprint-00-Integrar-Infrastructure", "migrations", "013_add_ltree_to_academic_units.up.sql")
	if _, err := os.Stat(localMigration); err == nil {
		fullPaths = append(fullPaths, localMigration)
		t.Logf("✅ Agregada migración local ltree: 013_add_ltree_to_academic_units.up.sql")
	} else {
		t.Logf("⚠️  Migración ltree no encontrada en: %s", localMigration)
	}

	// Verificar que se encontraron todas las migraciones
	if len(fullPaths) == 0 {
		return nil, fmt.Errorf("no se encontraron migraciones en %s", modPath)
	}

	return fullPaths, nil
}

// setupTestDB crea una instancia de PostgreSQL para tests usando shared/testing
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	// Obtener migraciones de infrastructure con validación
	migrationScripts, err := getMigrationScripts(t)
	if err != nil {
		t.Fatalf("Error obteniendo migraciones de infrastructure: %v", err)
	}

	// Log para debug
	t.Logf("✅ Usando %d migraciones desde infrastructure", len(migrationScripts))
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
