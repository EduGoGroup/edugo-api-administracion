//go:build integration

package integration

import (
	"context"
	"database/sql"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// Containers compartidos para todos los tests
	// Se inicializan en TestMain (main_test.go)
	sharedDB      *sql.DB
	sharedMongoDB *mongo.Database
	sharedManager *containers.Manager
)

// setupTestDB retorna la conexión compartida de PostgreSQL
// y una función de cleanup que limpia los DATOS (no el container)
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	if sharedDB == nil {
		t.Fatal("Shared DB not initialized. TestMain should have set it up.")
	}

	cleanup := func() {
		cleanupTestData(t)
	}

	return sharedDB, cleanup
}

// setupTestMongoDB retorna la conexión compartida de MongoDB
// y una función de cleanup que limpia los DATOS (no el container)
func setupTestMongoDB(t *testing.T) (*mongo.Database, func()) {
	t.Helper()

	if sharedMongoDB == nil {
		t.Fatal("Shared MongoDB not initialized. TestMain should have set it up.")
	}

	cleanup := func() {
		cleanupMongoData(t)
	}

	return sharedMongoDB, cleanup
}

// setupTestDBWithMigrations retorna la BD compartida con migraciones ya aplicadas
func setupTestDBWithMigrations(t *testing.T) (*sql.DB, func()) {
	return setupTestDB(t)
}

// cleanupTestData limpia todas las tablas de PostgreSQL
// Esto es MUCHO más rápido que destruir y recrear el container
func cleanupTestData(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	// Orden correcto para evitar violaciones de foreign key
	tables := []string{
		"unit_memberships",
		"academic_units",
		"schools",
		"users",
		"guardians",
		"materials",
		"subjects",
	}

	for _, table := range tables {
		// TRUNCATE CASCADE es más rápido que DELETE
		_, err := sharedDB.ExecContext(ctx, "TRUNCATE TABLE "+table+" CASCADE")
		if err != nil {
			// Si la tabla no existe, ignorar el error
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}
}

// cleanupMongoData limpia todas las colecciones de MongoDB
func cleanupMongoData(t *testing.T) {
	t.Helper()

	if sharedMongoDB == nil {
		return
	}

	ctx := context.Background()

	// Listar y limpiar todas las colecciones
	collections, err := sharedMongoDB.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		t.Logf("Warning: Failed to list MongoDB collections: %v", err)
		return
	}

	for _, collection := range collections {
		err := sharedMongoDB.Collection(collection).Drop(ctx)
		if err != nil {
			t.Logf("Warning: Failed to drop collection %s: %v", collection, err)
		}
	}
}
