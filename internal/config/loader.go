package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("server.port", 8081)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("database.postgres.max_connections", 25)
	v.SetDefault("database.postgres.ssl_mode", "disable")
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

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

	// ENV vars
	v.AutomaticEnv()
	v.SetEnvPrefix("EDUGO_ADMIN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Secrets
	_ = v.BindEnv("database.postgres.password", "POSTGRES_PASSWORD")
	_ = v.BindEnv("database.mongodb.uri", "MONGODB_URI")

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
