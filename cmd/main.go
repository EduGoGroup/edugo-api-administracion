package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/edugo/api-administracion/internal/config"
	_ "github.com/edugo/api-administracion/docs" // Swagger docs generados
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
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Error loading configuration: %v", err)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	log.Printf("🌍 Environment: %s", env)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "edugo-api-admin"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1")
	v1.Use(AdminAuthRequired())
	{
		// Users CRUD
		v1.POST("/users", CreateUser)       // POST /v1/users
		v1.PATCH("/users/:id", UpdateUser)  // PATCH /v1/users/:id
		v1.DELETE("/users/:id", DeleteUser) // DELETE /v1/users/:id

		// Units (Jerarquía Académica)
		v1.POST("/schools", CreateSchool)               // POST /v1/schools
		v1.POST("/units", CreateUnit)                   // POST /v1/units
		v1.PATCH("/units/:id", UpdateUnit)              // PATCH /v1/units/:id
		v1.POST("/units/:id/members", AssignMembership) // POST /v1/units/:id/members

		// Subjects
		v1.POST("/subjects", CreateSubject) // POST /v1/subjects

		// Materials Admin
		v1.DELETE("/materials/:id", DeleteMaterial) // DELETE /v1/materials/:id

		// Stats
		v1.GET("/stats/global", GetGlobalStats) // GET /v1/stats/global
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🔧 API Administración running on http://localhost:%d", cfg.Server.Port)
	log.Printf("📚 Swagger: http://localhost:%d/swagger/index.html", cfg.Server.Port)

	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}

// Middleware Admin
func AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Validar JWT y verificar role = admin
		c.Set("admin_id", "admin-mock")
		c.Next()
	}
}

// Request/Response types para Swagger

// CreateUserRequest representa la petición para crear usuario
type CreateUserRequest struct {
	Email    string `json:"email" example:"usuario@example.com"`
	Password string `json:"password" example:"password123"`
	Name     string `json:"name" example:"Juan Pérez"`
	Role     string `json:"role" example:"teacher"`
	SchoolID string `json:"school_id" example:"school-uuid-123"`
} // @name CreateUserRequest

// CreateUserResponse representa la respuesta al crear usuario
type CreateUserResponse struct {
	UserID  string `json:"user_id" example:"user-uuid-123"`
	Message string `json:"message" example:"Usuario creado"`
} // @name CreateUserResponse

// SuccessResponse representa una respuesta genérica de éxito
type SuccessResponse struct {
	Message string `json:"message" example:"Operación exitosa"`
} // @name SuccessResponse

// CreateSchoolResponse representa la respuesta al crear escuela
type CreateSchoolResponse struct {
	SchoolID string `json:"school_id" example:"school-uuid-123"`
} // @name CreateSchoolResponse

// CreateUnitResponse representa la respuesta al crear unidad
type CreateUnitResponse struct {
	UnitID string `json:"unit_id" example:"unit-uuid-123"`
} // @name CreateUnitResponse

// CreateSubjectResponse representa la respuesta al crear materia
type CreateSubjectResponse struct {
	SubjectID string `json:"subject_id" example:"subject-uuid-123"`
} // @name CreateSubjectResponse

// GlobalStatsResponse representa las estadísticas globales
type GlobalStatsResponse struct {
	TotalUsers     int `json:"total_users" example:"1250"`
	TotalMaterials int `json:"total_materials" example:"450"`
	ActiveUsers30d int `json:"active_users_30d" example:"980"`
} // @name GlobalStatsResponse

// Mock handlers con Swagger annotations

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
