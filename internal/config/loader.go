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
	v.SetDefault("logging.level", "info")

	// Ambiente
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	// Config files
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")

	// Base
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading base config: %w", err)
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
	v.BindEnv("database.postgres.password", "POSTGRES_PASSWORD")
	v.BindEnv("database.mongodb.uri", "MONGODB_URI")

	// Unmarshal
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
