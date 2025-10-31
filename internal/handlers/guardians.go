package handlers

import (
	"net/http"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/models/request"
	"github.com/EduGoGroup/edugo-api-administracion/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateGuardianRelation godoc
// @Summary Crear vínculo tutor-estudiante (Post-MVP)
// @Description Crea una relación entre un tutor (guardian) y un estudiante
// @Tags Guardians
// @Accept json
// @Produce json
// @Param body body request.CreateGuardianRelationRequest true "Datos de la relación"
// @Success 201 {object} response.GuardianRelationResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse "Relación ya existe"
// @Security BearerAuth
// @Router /guardian-relations [post]
func CreateGuardianRelation(c *gin.Context) {
	var req request.CreateGuardianRelationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// TODO: Validar que guardian_id tenga rol 'guardian' en PostgreSQL
	// TODO: Validar que student_id tenga rol 'student' en PostgreSQL
	// TODO: Verificar que no exista duplicado (UNIQUE constraint)
	// TODO: INSERT INTO guardian_student_relation (guardian_id, student_id, relationship_type)
	// TODO: Manejar errores de constraint violations

	// Mock response
	mockResponse := response.GuardianRelationResponse{
		ID:               uuid.New().String(),
		GuardianID:       req.GuardianID,
		StudentID:        req.StudentID,
		RelationshipType: req.RelationshipType,
		CreatedAt:        time.Now(),
	}

	c.JSON(http.StatusCreated, mockResponse)
}
