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
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	log.Println("🔄 EduGo API Administración iniciando...")

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
	c := container.NewContainer(resources.PostgreSQL, resources.Logger)
	defer c.Close()

	resources.Logger.Info("✅ API Administración iniciada", "port", cfg.Server.Port)

	// 4. Configurar Gin
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "edugo-api-admin"})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rutas v1
	v1 := r.Group("/v1")
	v1.Use(AdminAuthRequired())
	{
		// ==================== SCHOOLS ====================
		schools := v1.Group("/schools")
		{
			schools.POST("", c.SchoolHandler.CreateSchool)
			schools.GET("", c.SchoolHandler.ListSchools)
			schools.GET("/:id", c.SchoolHandler.GetSchool)
			schools.GET("/code/:code", c.SchoolHandler.GetSchoolByCode)
			schools.PUT("/:id", c.SchoolHandler.UpdateSchool)
			schools.DELETE("/:id", c.SchoolHandler.DeleteSchool)

			// Units nested under school
			schools.POST("/:schoolId/units", c.AcademicUnitHandler.CreateUnit)
			schools.GET("/:schoolId/units", c.AcademicUnitHandler.ListUnitsBySchool)
			schools.GET("/:schoolId/units/tree", c.AcademicUnitHandler.GetUnitTree)
			schools.GET("/:schoolId/units/by-type", c.AcademicUnitHandler.ListUnitsByType)
		}

		// ==================== ACADEMIC UNITS ====================
		units := v1.Group("/units")
		{
			units.GET("/:id", c.AcademicUnitHandler.GetUnit)
			units.PUT("/:id", c.AcademicUnitHandler.UpdateUnit)
			units.DELETE("/:id", c.AcademicUnitHandler.DeleteUnit)
			units.POST("/:id/restore", c.AcademicUnitHandler.RestoreUnit)
			units.GET("/:id/hierarchy-path", c.AcademicUnitHandler.GetHierarchyPath)

			// Memberships nested under unit
			units.GET("/:unitId/memberships", c.UnitMembershipHandler.ListMembershipsByUnit)
			units.GET("/:unitId/memberships/by-role", c.UnitMembershipHandler.ListMembershipsByRole)
		}

		// ==================== MEMBERSHIPS ====================
		memberships := v1.Group("/memberships")
		{
			memberships.POST("", c.UnitMembershipHandler.CreateMembership)
			memberships.GET("/:id", c.UnitMembershipHandler.GetMembership)
			memberships.PUT("/:id", c.UnitMembershipHandler.UpdateMembership)
			memberships.DELETE("/:id", c.UnitMembershipHandler.DeleteMembership)
			memberships.POST("/:id/expire", c.UnitMembershipHandler.ExpireMembership)
		}

		// ==================== USERS ====================
		users := v1.Group("/users")
		{
			users.GET("/:userId/memberships", c.UnitMembershipHandler.ListMembershipsByUser)
			// Legacy routes (mantener por compatibilidad)
			users.POST("", CreateUser)
			users.PATCH("/:id", UpdateUser)
			users.DELETE("/:id", DeleteUser)
		}

		// ==================== LEGACY ROUTES ====================
		// Subjects
		v1.POST("/subjects", CreateSubject)

		// Materials
		v1.DELETE("/materials/:id", DeleteMaterial)

		// Stats
		v1.GET("/stats/global", GetGlobalStats)
	}

	// 5. Servidor HTTP con graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: r,
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

// AdminAuthRequired es un middleware que valida autenticación de admin
func AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implementar validación real con JWT
		// Por ahora, solo dejamos pasar todas las peticiones
		c.Next()
	}
}

// ==================== LEGACY HANDLERS (TODO: Migrar a handlers reales) ====================

// CreateUser godoc
// @Summary Crear usuario
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 201 {object} CreateUserResponse
// @Router /users [post]
func CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"user_id": "mock-uuid"})
}

// UpdateUser godoc
// @Summary Actualizar usuario
// @Tags Users
// @Produce json
// @Param id path string true "ID del usuario"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Router /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado"})
}

// DeleteUser godoc
// @Summary Eliminar usuario
// @Tags Users
// @Produce json
// @Param id path string true "ID del usuario"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
}

// CreateSubject godoc
// @Summary Crear materia
// @Tags Subjects
// @Produce json
// @Security BearerAuth
// @Success 201 {object} CreateSubjectResponse
// @Router /subjects [post]
func CreateSubject(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"subject_id": "mock-uuid"})
}

// DeleteMaterial godoc
// @Summary Eliminar material
// @Tags Materials
// @Produce json
// @Param id path string true "ID del material"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Router /materials/{id} [delete]
func DeleteMaterial(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Material eliminado, limpieza en proceso"})
}

// GetGlobalStats godoc
// @Summary Estadísticas globales
// @Tags Stats
// @Produce json
// @Security BearerAuth
// @Success 200 {object} GlobalStatsResponse
// @Router /stats/global [get]
func GetGlobalStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_users":      1250,
		"total_materials":  450,
		"active_users_30d": 980,
	})
}

// ==================== RESPONSE TYPES ====================

type SuccessResponse struct {
	Message string `json:"message"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

type CreateSubjectResponse struct {
	SubjectID string `json:"subject_id"`
}

type GlobalStatsResponse struct {
	TotalUsers     int `json:"total_users"`
	TotalMaterials int `json:"total_materials"`
	ActiveUsers30d int `json:"active_users_30d"`
}
