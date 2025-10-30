//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestContainers struct {
	Postgres *postgres.PostgresContainer
	MongoDB  *mongodb.MongoDBContainer
}

func SetupContainers(t *testing.T) (*TestContainers, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase("edugo"),
		postgres.WithUsername("edugo_user"),
		postgres.WithPassword("edugo_pass"),
		postgres.WithInitScripts("../../../scripts/postgresql/01_schema.sql"),
	)
	if err != nil {
		t.Fatalf("Failed to start Postgres: %v", err)
	}

	mongoContainer, err := mongodb.Run(ctx, "mongo:7.0",
		mongodb.WithUsername("edugo_admin"),
		mongodb.WithPassword("edugo_pass"),
	)
	if err != nil {
		pgContainer.Terminate(ctx)
		t.Fatalf("Failed to start MongoDB: %v", err)
	}

	cleanup := func() {
		pgContainer.Terminate(ctx)
		mongoContainer.Terminate(ctx)
	}

	return &TestContainers{Postgres: pgContainer, MongoDB: mongoContainer}, cleanup
}
