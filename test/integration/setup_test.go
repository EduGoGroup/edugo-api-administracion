//go:build integration

package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupTestDB crea una instancia de PostgreSQL para tests usando testcontainers
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	containers, cleanup := SetupContainers(t)

	connString, err := containers.Postgres.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		cleanup()
		t.Fatalf("Failed to get Postgres connection string: %v", err)
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to connect to Postgres: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		cleanup()
		t.Fatalf("Failed to ping Postgres: %v", err)
	}

	return db, func() {
		db.Close()
		cleanup()
	}
}

// setupTestMongoDB crea una instancia de MongoDB para tests usando testcontainers
func setupTestMongoDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	containers, cleanup := SetupContainers(t)

	connString, err := containers.MongoDB.ConnectionString(ctx)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to get MongoDB connection string: %v", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		cleanup()
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		cleanup()
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database("edugo_test")

	return db, func() {
		client.Disconnect(ctx)
		cleanup()
	}
}

// TestSetupContainers verifica que los contenedores se inicialicen correctamente
func TestSetupContainers(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Verificar que la conexión funciona
	var result int
	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Failed to query Postgres: %v", err)
	}

	if result != 1 {
		t.Fatalf("Expected 1, got %d", result)
	}

	t.Log("✅ PostgreSQL testcontainer inicializado correctamente")
}

func TestSetupMongoDB(t *testing.T) {
	mongoDB, cleanup := setupTestMongoDB(t)
	defer cleanup()

	// Verificar que la conexión funciona
	ctx := context.Background()
	collections, err := mongoDB.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to list MongoDB collections: %v", err)
	}

	t.Logf("✅ MongoDB testcontainer inicializado correctamente (collections: %d)", len(collections))
}

// execSQLFile ejecuta un archivo SQL completo en la base de datos
func execSQLFile(t *testing.T, db *sql.DB, filepath string) {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read SQL file %s: %v", filepath, err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		t.Fatalf("Failed to execute SQL file %s: %v", filepath, err)
	}
}

// setupTestDBWithMigrations crea una BD con las migraciones aplicadas
func setupTestDBWithMigrations(t *testing.T) (*sql.DB, func()) {
	db, cleanup := setupTestDB(t)

	// Ejecutar migraciones
	execSQLFile(t, db, "../../scripts/postgresql/01_academic_hierarchy.sql")

	return db, cleanup
}
