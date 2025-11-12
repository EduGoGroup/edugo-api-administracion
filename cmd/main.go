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
	defer func() { if err := cleanup(); err != nil { resources.Logger.Error("Error durante cleanup", "error", err) } }()

	resources.Logger.Info("✅ API Administración iniciada", "port", cfg.Server.Port)

	// 3. Configurar Gin
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
		// Users CRUD
		v1.POST("/users", CreateUser)
		v1.PATCH("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)

		// Units (Jerarquía Académica)
		v1.POST("/schools", CreateSchool)
		v1.POST("/units", CreateUnit)
		v1.PATCH("/units/:id", UpdateUnit)
		v1.POST("/units/:id/members", AssignMembership)

		// Subjects
		v1.POST("/subjects", CreateSubject)

		// Materials Admin
		v1.DELETE("/materials/:id", DeleteMaterial)

		// Stats
		v1.GET("/stats/global", GetGlobalStats)
	}

	// 4. Servidor con graceful shutdown
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	resources.Logger.Info("🔧 API Administración running", 
		"addr", addr,
		"swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.Server.Port))

	// Iniciar servidor en goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			resources.Logger.Error("Server error", "error", err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	resources.Logger.Info("🛑 Apagando servidor...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		resources.Logger.Error("Server shutdown error", "error", err.Error())
	}

	resources.Logger.Info("✅ Servidor apagado correctamente")
}

// Middleware y Handlers (mantener por ahora como mock)

func AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("admin_id", "admin-mock")
		c.Next()
	}
}

// Handlers mock con Swagger annotations

type CreateUserRequest struct {
	Email    string `json:"email" example:"usuario@example.com"`
	Password string `json:"password" example:"password123"`
	Name     string `json:"name" example:"Juan Pérez"`
	Role     string `json:"role" example:"teacher"`
	SchoolID string `json:"school_id" example:"school-uuid-123"`
}

type CreateUserResponse struct {
	UserID  string `json:"user_id" example:"user-uuid-123"`
	Message string `json:"message" example:"Usuario creado"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"Operación exitosa"`
}

type CreateSchoolResponse struct {
	SchoolID string `json:"school_id" example:"school-uuid-123"`
}

type CreateUnitResponse struct {
	UnitID string `json:"unit_id" example:"unit-uuid-123"`
}

type CreateSubjectResponse struct {
	SubjectID string `json:"subject_id" example:"subject-uuid-123"`
}

type GlobalStatsResponse struct {
	TotalUsers     int `json:"total_users" example:"1250"`
	TotalMaterials int `json:"total_materials" example:"450"`
	ActiveUsers30d int `json:"active_users_30d" example:"980"`
}

// CreateUser godoc
// @Summary Crear usuario
// @Description Crear nuevo usuario con rol y perfil
// @Tags Users
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "Datos del usuario"
// @Security BearerAuth
// @Success 201 {object} CreateUserResponse
// @Router /users [post]
func CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"user_id": "mock-uuid", "message": "Usuario creado"})
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

// CreateSchool godoc
// @Summary Crear escuela
// @Tags Schools
// @Produce json
// @Security BearerAuth
// @Success 201 {object} CreateSchoolResponse
// @Router /schools [post]
func CreateSchool(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"school_id": "mock-uuid"})
}

// CreateUnit godoc
// @Summary Crear unidad académica
// @Tags Units
// @Produce json
// @Security BearerAuth
// @Success 201 {object} CreateUnitResponse
// @Router /units [post]
func CreateUnit(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"unit_id": "mock-uuid"})
}

// UpdateUnit godoc
// @Summary Actualizar unidad
// @Tags Units
// @Produce json
// @Param id path string true "ID de la unidad"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Router /units/{id} [patch]
func UpdateUnit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Unidad actualizada"})
}

// AssignMembership godoc
// @Summary Asignar membresía
// @Tags Units
// @Produce json
// @Param id path string true "ID de la unidad"
// @Security BearerAuth
// @Success 201 {object} SuccessResponse
// @Router /units/{id}/members [post]
func AssignMembership(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Membresía asignada"})
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
