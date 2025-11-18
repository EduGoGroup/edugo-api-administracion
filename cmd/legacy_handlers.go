package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==================== LEGACY HANDLERS ====================
//
// DEPRECATED: Estos endpoints están deprecated y serán removidos en v0.6.0
// NO implementan lógica real, solo retornan datos mock para compatibilidad.
//
// Si necesitas estas funcionalidades, deberás:
// 1. Implementar los handlers reales en internal/interface/http/handler/
// 2. Crear los services correspondientes en internal/application/service/
// 3. Actualizar la documentación Swagger

// CreateUser godoc
// @Summary [DEPRECATED] Crear usuario
// @Description DEPRECATED: Este endpoint será removido en v0.6.0. Retorna datos mock.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 410 {object} DeprecatedResponse
// @Router /users [post]
// @deprecated
func CreateUser(c *gin.Context) {
	c.Header("X-Deprecated", "true")
	c.Header("X-Deprecated-Remove-Version", "v0.6.0")
	c.JSON(http.StatusGone, DeprecatedResponse{
		Error:           "Este endpoint está deprecated y será removido en v0.6.0",
		DeprecatedSince: "v0.5.0",
		RemoveInVersion: "v0.6.0",
		MockResponse:    gin.H{"user_id": "mock-uuid"},
	})
}

// UpdateUser godoc
// @Summary [DEPRECATED] Actualizar usuario
// @Description DEPRECATED: Este endpoint será removido en v0.6.0. Retorna datos mock.
// @Tags Users
// @Produce json
// @Param id path string true "ID del usuario"
// @Security BearerAuth
// @Success 410 {object} DeprecatedResponse
// @Router /users/{id} [patch]
// @deprecated
func UpdateUser(c *gin.Context) {
	c.Header("X-Deprecated", "true")
	c.Header("X-Deprecated-Remove-Version", "v0.6.0")
	c.JSON(http.StatusGone, DeprecatedResponse{
		Error:           "Este endpoint está deprecated y será removido en v0.6.0",
		DeprecatedSince: "v0.5.0",
		RemoveInVersion: "v0.6.0",
		MockResponse:    gin.H{"message": "Usuario actualizado"},
	})
}

// DeleteUser godoc
// @Summary [DEPRECATED] Eliminar usuario
// @Description DEPRECATED: Este endpoint será removido en v0.6.0. Retorna datos mock.
// @Tags Users
// @Produce json
// @Param id path string true "ID del usuario"
// @Security BearerAuth
// @Success 410 {object} DeprecatedResponse
// @Router /users/{id} [delete]
// @deprecated
func DeleteUser(c *gin.Context) {
	c.Header("X-Deprecated", "true")
	c.Header("X-Deprecated-Remove-Version", "v0.6.0")
	c.JSON(http.StatusGone, DeprecatedResponse{
		Error:           "Este endpoint está deprecated y será removido en v0.6.0",
		DeprecatedSince: "v0.5.0",
		RemoveInVersion: "v0.6.0",
		MockResponse:    gin.H{"message": "Usuario eliminado"},
	})
}

// CreateSubject godoc
// @Summary [DEPRECATED] Crear materia
// @Description DEPRECATED: Este endpoint será removido en v0.6.0. Retorna datos mock.
// @Tags Subjects
// @Produce json
// @Security BearerAuth
// @Success 410 {object} DeprecatedResponse
// @Router /subjects [post]
// @deprecated
func CreateSubject(c *gin.Context) {
	c.Header("X-Deprecated", "true")
	c.Header("X-Deprecated-Remove-Version", "v0.6.0")
	c.JSON(http.StatusGone, DeprecatedResponse{
		Error:           "Este endpoint está deprecated y será removido en v0.6.0",
		DeprecatedSince: "v0.5.0",
		RemoveInVersion: "v0.6.0",
		MockResponse:    gin.H{"subject_id": "mock-uuid"},
	})
}

// DeleteMaterial godoc
// @Summary [DEPRECATED] Eliminar material
// @Description DEPRECATED: Este endpoint será removido en v0.6.0. Retorna datos mock.
// @Tags Materials
// @Produce json
// @Param id path string true "ID del material"
// @Security BearerAuth
// @Success 410 {object} DeprecatedResponse
// @Router /materials/{id} [delete]
// @deprecated
func DeleteMaterial(c *gin.Context) {
	c.Header("X-Deprecated", "true")
	c.Header("X-Deprecated-Remove-Version", "v0.6.0")
	c.JSON(http.StatusGone, DeprecatedResponse{
		Error:           "Este endpoint está deprecated y será removido en v0.6.0",
		DeprecatedSince: "v0.5.0",
		RemoveInVersion: "v0.6.0",
		MockResponse:    gin.H{"message": "Material eliminado, limpieza en proceso"},
	})
}

// GetGlobalStats godoc
// @Summary [DEPRECATED] Estadísticas globales
// @Description DEPRECATED: Este endpoint será removido en v0.6.0. Retorna datos mock.
// @Tags Stats
// @Produce json
// @Security BearerAuth
// @Success 410 {object} DeprecatedResponse
// @Router /stats/global [get]
// @deprecated
func GetGlobalStats(c *gin.Context) {
	c.Header("X-Deprecated", "true")
	c.Header("X-Deprecated-Remove-Version", "v0.6.0")
	c.JSON(http.StatusGone, DeprecatedResponse{
		Error:           "Este endpoint está deprecated y será removido en v0.6.0",
		DeprecatedSince: "v0.5.0",
		RemoveInVersion: "v0.6.0",
		MockResponse: gin.H{
			"total_users":      1250,
			"total_materials":  450,
			"active_users_30d": 980,
		},
	})
}

// ==================== DEPRECATED RESPONSE TYPES ====================

// DeprecatedResponse es la respuesta estándar para endpoints deprecated
type DeprecatedResponse struct {
	Error           string      `json:"error" example:"Este endpoint está deprecated y será removido en v0.6.0"`
	DeprecatedSince string      `json:"deprecated_since" example:"v0.5.0"`
	RemoveInVersion string      `json:"remove_in_version" example:"v0.6.0"`
	MockResponse    interface{} `json:"mock_response,omitempty"`
}

// DEPRECATED: Legacy response types - mantener solo para Swagger
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
