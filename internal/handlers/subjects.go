package handlers

import (
	"net/http"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/models/request"
	"github.com/EduGoGroup/edugo-api-administracion/internal/models/response"
	"github.com/gin-gonic/gin"
)

// UpdateSubject godoc
// @Summary Actualizar materia (Post-MVP)
// @Description Actualiza información de una materia del catálogo
// @Tags Subjects
// @Accept json
// @Produce json
// @Param id path string true "ID de la materia"
// @Param body body request.UpdateSubjectRequest true "Datos a actualizar"
// @Success 200 {object} response.SubjectResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse "Solo administradores"
// @Failure 404 {object} response.ErrorResponse "Materia no encontrada"
// @Security BearerAuth
// @Router /subjects/{id} [patch]
func UpdateSubject(c *gin.Context) {
	subjectID := c.Param("id")
	var req request.UpdateSubjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// TODO: Verificar que usuario tenga rol 'admin'
	// TODO: Verificar que materia existe en subject table
	// TODO: UPDATE subject SET ... WHERE id = $1
	// TODO: Actualizar solo campos presentes en request (COALESCE)
	// TODO: Actualizar updated_at = CURRENT_TIMESTAMP

	// Mock response
	now := time.Now()
	mockResponse := response.SubjectResponse{
		ID:          subjectID,
		Name:        "Matemáticas Avanzadas",
		Description: "Curso avanzado de matemáticas para nivel secundario",
		Metadata: map[string]interface{}{
			"level": "advanced",
			"hours": 120,
		},
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: &now,
	}

	c.JSON(http.StatusOK, mockResponse)
}
