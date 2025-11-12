package bootstrap

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-api-administracion/internal/bootstrap/adapter"
	"github.com/EduGoGroup/edugo-api-administracion/internal/config"
	sharedBootstrap "github.com/EduGoGroup/edugo-shared/bootstrap"
	"github.com/EduGoGroup/edugo-shared/lifecycle"
	gormLogger "gorm.io/gorm/logger"
)

func bridgeToSharedBootstrap(ctx context.Context, cfg *config.Config) (*Resources, func() error, error) {
	// 1. Logger GORM
	gormLogLevel := gormLogger.Silent
	if cfg.Logging.Level == "debug" {
		gormLogLevel = gormLogger.Info
	}
	gormLog := gormLogger.Default.LogMode(gormLogLevel)

	// 2. Factories
	sharedFactories := &sharedBootstrap.Factories{
		Logger:     sharedBootstrap.NewDefaultLoggerFactory(),
		PostgreSQL: sharedBootstrap.NewDefaultPostgreSQLFactory(gormLog),
	}

	wrapper := newCustomFactoriesWrapper(sharedFactories)
	customFactories := createCustomFactories(wrapper)

	lifecycleManager := lifecycle.NewManager(nil)

	// 3. Config
	bootstrapConfig := struct {
		Environment string
		PostgreSQL  sharedBootstrap.PostgreSQLConfig
	}{
		Environment: cfg.Environment,
		PostgreSQL: sharedBootstrap.PostgreSQLConfig{
			Host:     cfg.Database.Postgres.Host,
			Port:     cfg.Database.Postgres.Port,
			User:     cfg.Database.Postgres.User,
			Password: cfg.Database.Postgres.Password,
			Database: cfg.Database.Postgres.Database,
			SSLMode:  cfg.Database.Postgres.SSLMode,
		},
	}

	// 4. Bootstrap
	_, err := sharedBootstrap.Bootstrap(ctx, bootstrapConfig, customFactories, lifecycleManager,
		sharedBootstrap.WithRequiredResources("logger", "postgresql"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to bootstrap: %w", err)
	}

	// 5. Logger adapter
	if wrapper.logrusLogger == nil {
		return nil, nil, fmt.Errorf("logger not initialized")
	}
	loggerAdapter := adapter.NewLoggerAdapter(wrapper.logrusLogger)
	lifecycleWithLogger := lifecycle.NewManager(loggerAdapter)

	// 6. Resources
	resources := &Resources{
		Logger:     loggerAdapter,
		PostgreSQL: wrapper.sqlDB,
		JWTSecret:  "", // api-admin no usa JWT por ahora
	}

	cleanup := func() error {
		resources.Logger.Info("starting api-admin cleanup")
		err := lifecycleWithLogger.Cleanup()
		if err != nil {
			resources.Logger.Error("cleanup errors", "error", err.Error())
			return err
		}
		resources.Logger.Info("api-admin cleanup completed")
		return nil
	}

	return resources, cleanup, nil
}
