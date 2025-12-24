package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	v := viper.New()

	// Defaults - Server
	v.SetDefault("server.port", 8081)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")

	// Defaults - Database
	v.SetDefault("database.postgres.max_connections", 25)
	v.SetDefault("database.postgres.ssl_mode", "disable")

	// Defaults - Logging
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Defaults - Auth JWT
	v.SetDefault("auth.jwt.issuer", "edugo-central")
	v.SetDefault("auth.jwt.access_token_duration", "15m")
	v.SetDefault("auth.jwt.refresh_token_duration", "168h")
	v.SetDefault("auth.jwt.algorithm", "HS256")

	// Defaults - Auth Password
	v.SetDefault("auth.password.min_length", 8)
	v.SetDefault("auth.password.require_uppercase", true)
	v.SetDefault("auth.password.require_lowercase", true)
	v.SetDefault("auth.password.require_number", true)
	v.SetDefault("auth.password.require_special", false)
	v.SetDefault("auth.password.bcrypt_cost", 10)

	// Defaults - Rate Limiting
	v.SetDefault("auth.rate_limit.login.max_attempts", 5)
	v.SetDefault("auth.rate_limit.login.window", "15m")
	v.SetDefault("auth.rate_limit.login.block_duration", "1h")
	v.SetDefault("auth.rate_limit.internal_services.max_requests", 1000)
	v.SetDefault("auth.rate_limit.internal_services.window", "1m")
	v.SetDefault("auth.rate_limit.external_clients.max_requests", 60)
	v.SetDefault("auth.rate_limit.external_clients.window", "1m")

	// Defaults - Cache
	v.SetDefault("auth.cache.token_validation.enabled", true)
	v.SetDefault("auth.cache.token_validation.ttl", "60s")
	v.SetDefault("auth.cache.token_validation.max_size", 10000)
	v.SetDefault("auth.cache.user_info.enabled", true)
	v.SetDefault("auth.cache.user_info.ttl", "300s")
	v.SetDefault("auth.cache.user_info.max_size", 1000)

	// Defaults - Redis
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	// Defaults - CORS
	v.SetDefault("cors.allowed_origins", "http://localhost:3000,http://localhost:5173")
	v.SetDefault("cors.allowed_methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
	v.SetDefault("cors.allowed_headers", "Content-Type,Authorization,X-Requested-With")

	// Ambiente
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	// Config files
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")

	// Base (opcional en Docker)
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		// En Docker, el archivo puede no existir (se usa solo env vars)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading base config: %w", err)
		}
		// Archivo no encontrado es OK, continuamos con defaults + env vars
	}

	// Merge environment
	v.SetConfigName(fmt.Sprintf("config-%s", env))
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error merging %s config: %w", env, err)
		}
	}

	// ENV vars automáticos
	v.AutomaticEnv()
	v.SetEnvPrefix("EDUGO_ADMIN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bindings explícitos para variables de entorno sensibles
	// Database
	_ = v.BindEnv("database.postgres.password", "POSTGRES_PASSWORD")
	_ = v.BindEnv("database.mongodb.uri", "MONGODB_URI")

	// Mock repositories - binding explícito sin prefijo para compatibilidad con debug.json
	_ = v.BindEnv("database.use_mock_repositories", "USE_MOCK_REPOSITORIES")

	// Auth JWT
	_ = v.BindEnv("auth.jwt.secret", "AUTH_JWT_SECRET")
	_ = v.BindEnv("auth.jwt.issuer", "AUTH_JWT_ISSUER")
	_ = v.BindEnv("auth.jwt.access_token_duration", "AUTH_JWT_ACCESS_TOKEN_DURATION")
	_ = v.BindEnv("auth.jwt.refresh_token_duration", "AUTH_JWT_REFRESH_TOKEN_DURATION")

	// Rate Limiting
	_ = v.BindEnv("auth.rate_limit.login.max_attempts", "AUTH_RATE_LIMIT_LOGIN_ATTEMPTS")
	_ = v.BindEnv("auth.rate_limit.login.window", "AUTH_RATE_LIMIT_LOGIN_WINDOW")
	_ = v.BindEnv("auth.rate_limit.login.block_duration", "AUTH_RATE_LIMIT_LOGIN_BLOCK")

	// Internal Services
	_ = v.BindEnv("auth.internal_services.api_keys", "AUTH_INTERNAL_SERVICES_API_KEYS")
	_ = v.BindEnv("auth.internal_services.ip_ranges", "AUTH_INTERNAL_SERVICES_IP_RANGES")

	// Cache
	_ = v.BindEnv("auth.cache.token_validation.ttl", "AUTH_CACHE_TOKEN_VALIDATION_TTL")
	_ = v.BindEnv("auth.cache.user_info.ttl", "AUTH_CACHE_USER_INFO_TTL")

	// Redis
	_ = v.BindEnv("redis.host", "REDIS_HOST")
	_ = v.BindEnv("redis.port", "REDIS_PORT")
	_ = v.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = v.BindEnv("redis.db", "REDIS_DB")

	// CORS
	_ = v.BindEnv("cors.allowed_origins", "ALLOWED_ORIGINS")
	_ = v.BindEnv("cors.allowed_methods", "ALLOWED_METHODS")
	_ = v.BindEnv("cors.allowed_headers", "ALLOWED_HEADERS")

	// Unmarshal
	var cfg Config
	cfg.Environment = env
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate usando función separada
	if err := Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
