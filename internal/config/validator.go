package config

import (
	"fmt"
	"strings"
)

// Validate valida que la configuración tenga los campos obligatorios y valores válidos
func Validate(cfg *Config) error {
	var validationErrors []string

	// Validar secretos requeridos
	if cfg.Database.Postgres.Password == "" {
		validationErrors = append(validationErrors, "DATABASE_POSTGRES_PASSWORD is required")
	}
	if cfg.Database.MongoDB.URI == "" {
		validationErrors = append(validationErrors, "DATABASE_MONGODB_URI is required")
	}

	// Validar valores públicos
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		validationErrors = append(validationErrors, "server.port must be between 1 and 65535")
	}
	if cfg.Database.Postgres.MaxConnections <= 0 {
		validationErrors = append(validationErrors, "database.postgres.max_connections must be positive")
	}
	if cfg.Database.Postgres.Host == "" {
		validationErrors = append(validationErrors, "database.postgres.host is required")
	}
	if cfg.Database.Postgres.Database == "" {
		validationErrors = append(validationErrors, "database.postgres.database is required")
	}
	if cfg.Database.Postgres.User == "" {
		validationErrors = append(validationErrors, "database.postgres.user is required")
	}

	// Si hay errores, retornar un error compuesto con mensaje claro
	if len(validationErrors) > 0 {
		errorMsg := "Configuration validation failed:\n  - " +
			strings.Join(validationErrors, "\n  - ") +
			"\n\nPlease check your .env file or environment variables.\nFor local development, copy .env.example to .env and fill in the values."
		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
