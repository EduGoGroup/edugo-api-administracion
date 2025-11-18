//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	globalDB      *sql.DB
	globalCleanup func()
)

// TestMain se ejecuta UNA VEZ antes de todos los tests
func TestMain(m *testing.M) {
	// Setup: Levantar 1 container PostgreSQL para TODOS los tests
	db, cleanup := setupTestDBWithMigrations(&testing.T{})
	globalDB = db
	globalCleanup = cleanup

	// Ejecutar todos los tests
	code := m.Run()

	// Cleanup: Destruir container al final
	if globalCleanup != nil {
		globalCleanup()
	}

	os.Exit(code)
}

// getTestDB retorna la conexión global y una función de cleanup por test
func getTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Cleanup por test: Truncar tablas para aislar tests
	cleanup := func() {
		ctx := context.Background()

		// Truncar en orden inverso a las FK
		tables := []string{
			"memberships",
			"academic_units",
			"schools",
			"users",
		}

		for _, table := range tables {
			_, err := globalDB.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
			if err != nil {
				t.Logf("Warning: Error truncating %s: %v", table, err)
			}
		}
	}

	return globalDB, cleanup
}
