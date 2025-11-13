//go:build integration

package integration

import (
	"context"
	"testing"
)

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
