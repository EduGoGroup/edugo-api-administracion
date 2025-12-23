package container

import (
	"database/sql"
	"log"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	authHandler "github.com/EduGoGroup/edugo-api-administracion/internal/auth/handler"
	authService "github.com/EduGoGroup/edugo-api-administracion/internal/auth/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/config"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/factory"
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-api-administracion/internal/shared/crypto"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// Container es el contenedor de dependencias de la aplicación
// Implementa el patrón Dependency Injection
type Container struct {
	// Infrastructure
	DB         *sql.DB
	Logger     logger.Logger
	JWTManager *auth.JWTManager

	// Auth (centralizado)
	PasswordHasher    *crypto.PasswordHasher
	InternalJWTManager *crypto.JWTManager
	TokenService      *authService.TokenService
	AuthService       authService.AuthService
	AuthHandler       *authHandler.AuthHandler
	VerifyHandler     *authHandler.VerifyHandler

	// Repositories
	UserRepository           repository.UserRepository
	SchoolRepository         repository.SchoolRepository
	AcademicUnitRepository   repository.AcademicUnitRepository
	UnitMembershipRepository repository.UnitMembershipRepository
	UnitRepository           repository.UnitRepository
	SubjectRepository        repository.SubjectRepository
	MaterialRepository       repository.MaterialRepository
	StatsRepository          repository.StatsRepository
	GuardianRepository       repository.GuardianRepository

	// Services
	UserService           service.UserService
	SchoolService         service.SchoolService
	AcademicUnitService   service.AcademicUnitService
	UnitMembershipService service.UnitMembershipService
	UnitService           service.UnitService
	SubjectService        service.SubjectService
	MaterialService       service.MaterialService
	StatsService          service.StatsService
	GuardianService       service.GuardianService

	// Handlers
	UserHandler           *handler.UserHandler
	SchoolHandler         *handler.SchoolHandler
	AcademicUnitHandler   *handler.AcademicUnitHandler
	UnitMembershipHandler *handler.UnitMembershipHandler
	UnitHandler           *handler.UnitHandler
	SubjectHandler        *handler.SubjectHandler
	MaterialHandler       *handler.MaterialHandler
	StatsHandler          *handler.StatsHandler
	GuardianHandler       *handler.GuardianHandler
}

// NewContainer crea un nuevo contenedor e inicializa todas las dependencias
func NewContainer(db *sql.DB, logger logger.Logger, jwtSecret string, cfg *config.Config) *Container {
	c := &Container{
		DB:         db,
		Logger:     logger,
		JWTManager: auth.NewJWTManager(jwtSecret, "edugo-central"),
	}

	// ==================== AUTH (Centralizado) ====================
	// Password Hasher (costo 12 para producción)
	c.PasswordHasher = crypto.NewPasswordHasher(12)

	// JWT Manager interno para tokens (usando el crypto package local)
	jwtConfig := crypto.JWTConfig{
		Secret:               jwtSecret,
		Issuer:               "edugo-central",
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
	}
	internalJWTManager, err := crypto.NewJWTManager(jwtConfig)
	if err != nil {
		log.Fatalf("❌ Error creando JWTManager: %v", err)
	}
	c.InternalJWTManager = internalJWTManager

	// Token Service
	tokenConfig := authService.TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   false, // Por ahora sin cache (se habilitará con Redis)
		BlacklistCheck: false, // Por ahora sin blacklist (se habilitará con Redis)
	}
	c.TokenService = authService.NewTokenService(internalJWTManager, nil, tokenConfig)

	// ==================== REPOSITORY FACTORY ====================
	// Decidir entre Mock o PostgreSQL según configuración
	var repositoryFactory factory.RepositoryFactory

	if cfg.Database.UseMockRepositories {
		logger.Info("usando mock repositories", "mock_enabled", true, "postgres_required", false)
		repositoryFactory = factory.NewMockRepositoryFactory()
	} else {
		logger.Info("usando postgresql repositories", "mock_enabled", false, "postgres_required", true)
		repositoryFactory = factory.NewPostgresRepositoryFactory(db)
	}

	// Inicializar repositories usando el factory
	c.SchoolRepository = repositoryFactory.CreateSchoolRepository()
	c.UserRepository = repositoryFactory.CreateUserRepository()
	c.AcademicUnitRepository = repositoryFactory.CreateAcademicUnitRepository()
	c.UnitMembershipRepository = repositoryFactory.CreateUnitMembershipRepository()
	c.UnitRepository = repositoryFactory.CreateUnitRepository()
	c.SubjectRepository = repositoryFactory.CreateSubjectRepository()
	c.MaterialRepository = repositoryFactory.CreateMaterialRepository()
	c.StatsRepository = repositoryFactory.CreateStatsRepository()
	c.GuardianRepository = repositoryFactory.CreateGuardianRepository()

	// Auth Service (usa UserRepository y TokenService)
	c.AuthService = authService.NewAuthService(
		c.UnitMembershipRepository,
		c.UserRepository,
		c.TokenService,
		c.PasswordHasher,
		logger,
	)

	// Auth Handler
	c.AuthHandler = authHandler.NewAuthHandler(c.AuthService)

	// Verify Handler (para /v1/auth/verify)
	c.VerifyHandler = authHandler.NewVerifyHandler(
		c.TokenService,
		[]string{"127.0.0.1/32", "::1/128", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
		map[string]string{"api-mobile": "internal-mobile-key", "api-worker": "internal-worker-key"},
	)

	// Inicializar services (capa de aplicación)
	c.UserService = service.NewUserService(
		c.UserRepository,
		logger,
	)
	c.SchoolService = service.NewSchoolService(
		c.SchoolRepository,
		logger,
		cfg.Defaults.School,
	)
	c.AcademicUnitService = service.NewAcademicUnitService(
		c.AcademicUnitRepository,
		c.SchoolRepository,
		logger,
	)
	c.UnitMembershipService = service.NewUnitMembershipService(
		c.UnitMembershipRepository,
		c.AcademicUnitRepository,
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
	c.AcademicUnitHandler = handler.NewAcademicUnitHandler(
		c.AcademicUnitService,
		logger,
	)
	c.UnitMembershipHandler = handler.NewUnitMembershipHandler(
		c.UnitMembershipService,
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
