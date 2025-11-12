package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	sharedBootstrap "github.com/EduGoGroup/edugo-shared/bootstrap"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type customFactoriesWrapper struct {
	factories    *sharedBootstrap.Factories
	logrusLogger *logrus.Logger
	gormDB       *gorm.DB
	sqlDB        *sql.DB
}

func newCustomFactoriesWrapper(factories *sharedBootstrap.Factories) *customFactoriesWrapper {
	return &customFactoriesWrapper{factories: factories}
}

func createCustomFactories(wrapper *customFactoriesWrapper) *sharedBootstrap.Factories {
	return &sharedBootstrap.Factories{
		Logger:     &customLoggerFactory{wrapper: wrapper},
		PostgreSQL: &customPostgreSQLFactory{wrapper: wrapper},
	}
}

type customLoggerFactory struct {
	wrapper *customFactoriesWrapper
}

func (f *customLoggerFactory) CreateLogger(ctx context.Context, level, version string) (*logrus.Logger, error) {
	logger, err := f.wrapper.factories.Logger.CreateLogger(ctx, level, version)
	if err != nil {
		return nil, err
	}
	f.wrapper.logrusLogger = logger
	return logger, nil
}

type customPostgreSQLFactory struct {
	wrapper *customFactoriesWrapper
}

func (f *customPostgreSQLFactory) CreateConnection(ctx context.Context, cfg sharedBootstrap.PostgreSQLConfig) (*gorm.DB, error) {
	db, err := f.wrapper.factories.PostgreSQL.CreateConnection(ctx, cfg)
	if err != nil {
		return nil, err
	}
	f.wrapper.gormDB = db
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	f.wrapper.sqlDB = sqlDB
	return db, nil
}

func (f *customPostgreSQLFactory) CreateRawConnection(ctx context.Context, cfg sharedBootstrap.PostgreSQLConfig) (*sql.DB, error) {
	db, err := f.wrapper.factories.PostgreSQL.CreateRawConnection(ctx, cfg)
	if err != nil {
		return nil, err
	}
	f.wrapper.sqlDB = db
	return db, nil
}

func (f *customPostgreSQLFactory) Ping(ctx context.Context, db *gorm.DB) error {
	return f.wrapper.factories.PostgreSQL.Ping(ctx, db)
}

func (f *customPostgreSQLFactory) Close(db *gorm.DB) error {
	return f.wrapper.factories.PostgreSQL.Close(db)
}
