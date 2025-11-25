package config

import (
	"fmt"
	"strings"
)

// Validate valida que la configuración tenga los campos obligatorios y valores válidos
func Validate(cfg *Config) error {
	var validationErrors []string

	// ============================================
	// Validar Database
	// ============================================
	if cfg.Database.Postgres.Password == "" {
		validationErrors = append(validationErrors, "POSTGRES_PASSWORD is required")
	}
	if cfg.Database.MongoDB.URI == "" {
		validationErrors = append(validationErrors, "MONGODB_URI is required")
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
	if cfg.Database.Postgres.MaxConnections <= 0 {
		validationErrors = append(validationErrors, "database.postgres.max_connections must be positive")
	}

	// ============================================
	// Validar Server
	// ============================================
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		validationErrors = append(validationErrors, "server.port must be between 1 and 65535")
	}

	// ============================================
	// Validar Auth JWT
	// ============================================
	if cfg.Auth.JWT.Secret == "" {
		validationErrors = append(validationErrors, "AUTH_JWT_SECRET is required")
	} else if len(cfg.Auth.JWT.Secret) < 32 {
		validationErrors = append(validationErrors, "AUTH_JWT_SECRET must be at least 32 characters")
	}

	if cfg.Auth.JWT.Issuer == "" {
		validationErrors = append(validationErrors, "AUTH_JWT_ISSUER is required")
	}

	if cfg.Auth.JWT.AccessTokenDuration <= 0 {
		validationErrors = append(validationErrors, "auth.jwt.access_token_duration must be positive")
	}

	if cfg.Auth.JWT.RefreshTokenDuration <= 0 {
		validationErrors = append(validationErrors, "auth.jwt.refresh_token_duration must be positive")
	}

	// ============================================
	// Validar Auth Rate Limiting
	// ============================================
	if cfg.Auth.RateLimit.Login.MaxAttempts <= 0 {
		validationErrors = append(validationErrors, "auth.rate_limit.login.max_attempts must be positive")
	}

	if cfg.Auth.RateLimit.Login.Window <= 0 {
		validationErrors = append(validationErrors, "auth.rate_limit.login.window must be positive")
	}

	// ============================================
	// Validar Auth Password
	// ============================================
	if cfg.Auth.Password.MinLength < 6 {
		validationErrors = append(validationErrors, "auth.password.min_length must be at least 6")
	}

	if cfg.Auth.Password.BcryptCost < 4 || cfg.Auth.Password.BcryptCost > 31 {
		validationErrors = append(validationErrors, "auth.password.bcrypt_cost must be between 4 and 31")
	}

	// ============================================
	// Validar Redis (opcional pero recomendado para cache)
	// ============================================
	if cfg.Redis.Port <= 0 || cfg.Redis.Port > 65535 {
		validationErrors = append(validationErrors, "redis.port must be between 1 and 65535")
	}

	// ============================================
	// Retornar errores si existen
	// ============================================
	if len(validationErrors) > 0 {
		errorMsg := "Configuration validation failed:\n  - " +
			strings.Join(validationErrors, "\n  - ") +
			"\n\nPlease check your .env file or environment variables.\n" +
			"For local development, copy .env.example to .env and fill in the values.\n" +
			"Run: ./scripts/validate-env.sh to check your configuration."
		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
