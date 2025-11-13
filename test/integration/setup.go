//go:build integration

package integration

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/EduGoGroup/edugo-shared/testing/containers"
	"go.mongodb.org/mongo-driver/mongo"
)

// setupTestDB crea una instancia de PostgreSQL para tests usando shared/testing
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	
	// Configurar PostgreSQL con scripts de inicialización
	cfg := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Database:    "edugo_test",
			Username:    "edugo_user",
			Password:    "edugo_pass",
			InitScripts: []string{"../../scripts/postgresql/01_academic_hierarchy.sql"},
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
// Alias para setupTestDB ya que InitScripts aplica las migraciones automáticamente
func setupTestDBWithMigrations(t *testing.T) (*sql.DB, func()) {
	return setupTestDB(t)
}
