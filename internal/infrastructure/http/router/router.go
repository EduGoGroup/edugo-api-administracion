package router

import (
	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/config"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/middleware"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// Config configuración para el router
type Config struct {
	SchoolRepo     repository.SchoolRepository
	UnitRepo       repository.AcademicUnitRepository
	Logger         logger.Logger
	SchoolDefaults config.SchoolDefaults
	CORSConfig     *config.CORSConfig
}

// SetupRouter configura todas las rutas de la API
func SetupRouter(cfg *Config) *gin.Engine {
	router := gin.Default()

	// Middleware global
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware(cfg.CORSConfig))
	router.Use(middleware.ErrorHandler(cfg.Logger))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Inicializar servicios
		schoolService := service.NewSchoolService(cfg.SchoolRepo, cfg.Logger, cfg.SchoolDefaults)
		academicUnitService := service.NewAcademicUnitService(cfg.UnitRepo, cfg.SchoolRepo, cfg.Logger)

		// Handlers
		schoolHandler := handler.NewSchoolHandler(schoolService, cfg.Logger)
		unitHandler := handler.NewAcademicUnitHandler(academicUnitService, cfg.Logger)

		// School routes
		schools := v1.Group("/schools")
		{
			schools.POST("", schoolHandler.CreateSchool)
			schools.GET("", schoolHandler.ListSchools)
			schools.GET("/:id", schoolHandler.GetSchool)
			schools.GET("/code/:code", schoolHandler.GetSchoolByCode)
			schools.PUT("/:id", schoolHandler.UpdateSchool)
			schools.DELETE("/:id", schoolHandler.DeleteSchool)

			// School-scoped unit routes
			schools.POST("/:schoolId/units", unitHandler.CreateUnit)
			schools.GET("/:schoolId/units", unitHandler.ListUnitsBySchool)
			schools.GET("/:schoolId/units/tree", unitHandler.GetUnitTree) // ltree endpoint!
			schools.GET("/:schoolId/units/by-type", unitHandler.ListUnitsByType)
		}

		// Academic Unit routes (not scoped to school)
		units := v1.Group("/units")
		{
			units.GET("/:id", unitHandler.GetUnit)
			units.PUT("/:id", unitHandler.UpdateUnit)
			units.DELETE("/:id", unitHandler.DeleteUnit)
			units.POST("/:id/restore", unitHandler.RestoreUnit)
			units.GET("/:id/hierarchy-path", unitHandler.GetHierarchyPath) // ltree endpoint!
		}
	}

	return router
}
