package container

import (
	"database/sql"

	"github.com/edugo/api-administracion/internal/application/service"
	"github.com/edugo/api-administracion/internal/domain/repository"
	"github.com/edugo/api-administracion/internal/infrastructure/http/handler"
	postgresRepo "github.com/edugo/api-administracion/internal/infrastructure/persistence/postgres/repository"
	"github.com/edugo/shared/pkg/logger"
)

// Container es el contenedor de dependencias de la aplicación
// Implementa el patrón Dependency Injection
type Container struct {
	// Infrastructure
	DB     *sql.DB
	Logger logger.Logger

	// Repositories
	UserRepository     repository.UserRepository
	SchoolRepository   repository.SchoolRepository
	UnitRepository     repository.UnitRepository
	SubjectRepository  repository.SubjectRepository
	MaterialRepository repository.MaterialRepository
	StatsRepository    repository.StatsRepository
	GuardianRepository repository.GuardianRepository

	// Services
	UserService     service.UserService
	SchoolService   service.SchoolService
	UnitService     service.UnitService
	SubjectService  service.SubjectService
	MaterialService service.MaterialService
	StatsService    service.StatsService
	GuardianService service.GuardianService

	// Handlers
	UserHandler     *handler.UserHandler
	SchoolHandler   *handler.SchoolHandler
	UnitHandler     *handler.UnitHandler
	SubjectHandler  *handler.SubjectHandler
	MaterialHandler *handler.MaterialHandler
	StatsHandler    *handler.StatsHandler
	GuardianHandler *handler.GuardianHandler
}

// NewContainer crea un nuevo contenedor e inicializa todas las dependencias
func NewContainer(db *sql.DB, logger logger.Logger) *Container {
	c := &Container{
		DB:     db,
		Logger: logger,
	}

	// Inicializar repositories (capa de infraestructura)
	c.UserRepository = postgresRepo.NewPostgresUserRepository(db)
	c.SchoolRepository = postgresRepo.NewPostgresSchoolRepository(db)
	c.UnitRepository = postgresRepo.NewPostgresUnitRepository(db)
	c.SubjectRepository = postgresRepo.NewPostgresSubjectRepository(db)
	c.MaterialRepository = postgresRepo.NewPostgresMaterialRepository(db)
	c.StatsRepository = postgresRepo.NewPostgresStatsRepository(db)
	c.GuardianRepository = postgresRepo.NewPostgresGuardianRepository(db)

	// Inicializar services (capa de aplicación)
	c.UserService = service.NewUserService(
		c.UserRepository,
		logger,
	)
	c.SchoolService = service.NewSchoolService(
		c.SchoolRepository,
		logger,
	)
	c.UnitService = service.NewUnitService(
		c.UnitRepository,
		logger,
	)
	c.SubjectService = service.NewSubjectService(
		c.SubjectRepository,
		logger,
	)
	c.MaterialService = service.NewMaterialService(
		c.MaterialRepository,
		logger,
	)
	c.StatsService = service.NewStatsService(
		c.StatsRepository,
		logger,
	)
	c.GuardianService = service.NewGuardianService(
		c.GuardianRepository,
		logger,
	)

	// Inicializar handlers (capa de infraestructura HTTP)
	c.UserHandler = handler.NewUserHandler(
		c.UserService,
		logger,
	)
	c.SchoolHandler = handler.NewSchoolHandler(
		c.SchoolService,
		logger,
	)
	c.UnitHandler = handler.NewUnitHandler(
		c.UnitService,
		logger,
	)
	c.SubjectHandler = handler.NewSubjectHandler(
		c.SubjectService,
		logger,
	)
	c.MaterialHandler = handler.NewMaterialHandler(
		c.MaterialService,
		logger,
	)
	c.StatsHandler = handler.NewStatsHandler(
		c.StatsService,
		logger,
	)
	c.GuardianHandler = handler.NewGuardianHandler(
		c.GuardianService,
		logger,
	)

	return c
}

// Close cierra los recursos del contenedor
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
