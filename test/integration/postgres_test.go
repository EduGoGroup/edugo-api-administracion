//go:build integration

package integration

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgresConnection(t *testing.T) {
	containers, cleanup := SetupContainers(t)
	defer cleanup()

	connStr, err := containers.Postgres.ConnectionString(context.Background(), "sslmode=disable")
	assert.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	assert.NoError(t, err)
	t.Log("✅ PostgreSQL connection successful")
}
