package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/EduGoGroup/edugo-api-administracion/docs"
	"github.com/EduGoGroup/edugo-api-administracion/internal/bootstrap"
	"github.com/EduGoGroup/edugo-api-administracion/internal/config"
	"github.com/EduGoGroup/edugo-api-administracion/internal/container"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	// Version is the application version, injected at build time via ldflags
	Version = "dev"
	// BuildTime is the build timestamp, injected at build time via ldflags
	BuildTime = "unknown"
)

// @title EduGo API Administración
// @version 1.0
// @description API para operaciones CRUD y administrativas en EduGo
// @host localhost:8081
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	log.Printf("🔄 EduGo API Administración iniciando... (Version: %s, Build: %s)", Version, BuildTime)

	ctx := context.Background()

	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Error loading configuration: %v", err)
	}

	// 2. Inicializar infraestructura con shared/bootstrap
	resources, cleanup, err := bootstrap.Initialize(ctx, cfg)
	if err != nil {
		log.Fatalf("❌ Error inicializando infraestructura: %v", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			resources.Logger.Error("Error durante cleanup", "error", err)
		}
	}()

	// 3. Crear container de dependencias
	jwtSecret := cfg.Auth.JWT.Secret
	if jwtSecret == "" {
		log.Fatalf("❌ JWT_SECRET no está configurado")
	}
	c := container.NewContainer(resources.PostgreSQL, resources.Logger, jwtSecret, cfg)
	defer func() { _ = c.Close() }()

	resources.Logger.Info("✅ API Administración iniciada", "port", cfg.Server.Port)

	// 4. Configurar Gin
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "edugo-api-admin"})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ==================== RUTAS PÚBLICAS (sin autenticación) ====================
	v1Public := r.Group("/v1")
	{
		// Auth endpoints (públicos)
		c.AuthHandler.RegisterRoutes(v1Public)

		// Verify endpoint (para otros servicios)
		c.VerifyHandler.RegisterRoutes(v1Public)
	}

	// ==================== RUTAS PROTEGIDAS (requieren JWT) ====================
	v1 := r.Group("/v1")
	// Middleware de autenticación JWT (todas las rutas requieren token válido)
	v1.Use(ginmiddleware.JWTAuthMiddleware(c.JWTManager))
	{
		// ==================== SCHOOLS ====================
		schools := v1.Group("/schools")
		{
			schools.POST("", c.SchoolHandler.CreateSchool)
			schools.GET("", c.SchoolHandler.ListSchools)
			schools.GET("/code/:code", c.SchoolHandler.GetSchoolByCode)

			// Academic Units nested under school (usando :id como parámetro)
			schools.POST("/:id/units", c.AcademicUnitHandler.CreateUnit)
			schools.GET("/:id/units", c.AcademicUnitHandler.ListUnitsBySchool)
			schools.GET("/:id/units/tree", c.AcademicUnitHandler.GetUnitTree)
			schools.GET("/:id/units/by-type", c.AcademicUnitHandler.ListUnitsByType)

			// School CRUD (mismo parámetro :id)
			schools.GET("/:id", c.SchoolHandler.GetSchool)
			schools.PUT("/:id", c.SchoolHandler.UpdateSchool)
			schools.DELETE("/:id", c.SchoolHandler.DeleteSchool)
		}

		// ==================== ACADEMIC UNITS ====================
		units := v1.Group("/units")
		{
			units.GET("/:id", c.AcademicUnitHandler.GetUnit)
			units.PUT("/:id", c.AcademicUnitHandler.UpdateUnit)
			units.DELETE("/:id", c.AcademicUnitHandler.DeleteUnit)
			units.POST("/:id/restore", c.AcademicUnitHandler.RestoreUnit)
			units.GET("/:id/hierarchy-path", c.AcademicUnitHandler.GetHierarchyPath)
		}

		// ==================== MEMBERSHIPS ====================
		memberships := v1.Group("/memberships")
		{
			memberships.POST("", c.UnitMembershipHandler.CreateMembership)
			memberships.GET("", c.UnitMembershipHandler.ListMembershipsByUnit)         // Usa query param unit_id
			memberships.GET("/by-role", c.UnitMembershipHandler.ListMembershipsByRole) // Usa query params
			memberships.GET("/:id", c.UnitMembershipHandler.GetMembership)
			memberships.PUT("/:id", c.UnitMembershipHandler.UpdateMembership)
			memberships.DELETE("/:id", c.UnitMembershipHandler.DeleteMembership)
			memberships.POST("/:id/expire", c.UnitMembershipHandler.ExpireMembership)
		}

		// ==================== USERS ====================
		users := v1.Group("/users")
		{
			users.GET("/:userId/memberships", c.UnitMembershipHandler.ListMembershipsByUser)
		}

		// ==================== SUBJECTS ====================
		subjects := v1.Group("/subjects")
		{
			subjects.POST("", c.SubjectHandler.CreateSubject)
			subjects.GET("", c.SubjectHandler.ListSubjects)
			subjects.GET("/:id", c.SubjectHandler.GetSubject)
			subjects.PATCH("/:id", c.SubjectHandler.UpdateSubject)
			subjects.DELETE("/:id", c.SubjectHandler.DeleteSubject)
		}

		// ==================== GUARDIAN RELATIONS ====================
		guardianRelations := v1.Group("/guardian-relations")
		{
			guardianRelations.POST("", c.GuardianHandler.CreateGuardianRelation)
			guardianRelations.GET("/:id", c.GuardianHandler.GetGuardianRelation)
			guardianRelations.PUT("/:id", c.GuardianHandler.UpdateGuardianRelation)
			guardianRelations.DELETE("/:id", c.GuardianHandler.DeleteGuardianRelation)
		}

		// Guardian relations by guardian or student
		guardians := v1.Group("/guardians")
		{
			guardians.GET("/:guardian_id/relations", c.GuardianHandler.GetGuardianRelations)
		}

		students := v1.Group("/students")
		{
			students.GET("/:student_id/guardians", c.GuardianHandler.GetStudentGuardians)
		}
	}

	// 5. Servidor HTTP con graceful shutdown
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server
	go func() {
		resources.Logger.Info("🚀 Servidor escuchando", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			resources.Logger.Error("Error en servidor HTTP", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	resources.Logger.Info("🛑 Apagando servidor...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		resources.Logger.Error("Error en shutdown", "error", err)
	}

	resources.Logger.Info("✅ Servidor detenido correctamente")
}
