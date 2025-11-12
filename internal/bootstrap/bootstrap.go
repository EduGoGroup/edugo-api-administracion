package bootstrap

import (
	"context"
	"database/sql"

	"github.com/EduGoGroup/edugo-api-administracion/internal/config"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// Resources contiene los recursos inicializados
type Resources struct {
	Logger     logger.Logger
	PostgreSQL *sql.DB
	JWTSecret  string
}

// Initialize inicializa la infraestructura
func Initialize(ctx context.Context, cfg *config.Config) (*Resources, func() error, error) {
	return bridgeToSharedBootstrap(ctx, cfg)
}
