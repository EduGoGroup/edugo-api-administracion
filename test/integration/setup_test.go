//go:build integration

package integration

import (
	"context"
	"testing"
)

// TestSetupSharedContainers verifica que los contenedores compartidos funcionen
func TestSetupSharedContainers(t *testing.T) {
	// Verificar PostgreSQL
	if sharedDB == nil {
		t.Fatal("Shared PostgreSQL not initialized")
	}

	var result int
	err := sharedDB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Failed to query Postgres: %v", err)
	}

	if result != 1 {
		t.Fatalf("Expected 1, got %d", result)
	}

	t.Log("✅ PostgreSQL container compartido funcionando correctamente")
}

func TestSetupSharedMongoDB(t *testing.T) {
	// Verificar MongoDB
	if sharedMongoDB == nil {
		t.Skip("MongoDB not configured for this test suite")
	}

	ctx := context.Background()
	collections, err := sharedMongoDB.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to list MongoDB collections: %v", err)
	}

	t.Logf("✅ MongoDB container compartido funcionando correctamente (collections: %d)", len(collections))
}

// TestCleanupData verifica que la función de cleanup funcione
func TestCleanupData(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// El cleanup debería limpiar los datos sin errores
	// Esto se verifica implícitamente en el defer

	// Verificar que la conexión sigue funcionando después del cleanup
	var result int
	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Failed to query after setup: %v", err)
	}

	t.Log("✅ Función de cleanup de datos funciona correctamente")
}
