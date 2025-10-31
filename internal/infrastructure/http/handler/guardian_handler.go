package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/logger"
)

// GuardianHandler maneja las peticiones HTTP relacionadas con guardians
type GuardianHandler struct {
	guardianService service.GuardianService
	logger          logger.Logger
}

// NewGuardianHandler crea un nuevo GuardianHandler
func NewGuardianHandler(
	guardianService service.GuardianService,
	logger logger.Logger,
) *GuardianHandler {
	return &GuardianHandler{
		guardianService: guardianService,
		logger:          logger,
	}
}

// CreateGuardianRelation godoc
// @Summary Create guardian-student relation
// @Description Creates a new relationship between a guardian and a student
// @Tags guardians
// @Accept json
// @Produce json
// @Param request body dto.CreateGuardianRelationRequest true "Guardian relation data"
// @Success 201 {object} dto.GuardianRelationResponse
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 409 {object} ErrorResponse "Relation already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/guardian-relations [post]
// @Security BearerAuth
func (h *GuardianHandler) CreateGuardianRelation(c *gin.Context) {
	var req dto.CreateGuardianRelationRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Obtener admin_id del contexto (agregado por middleware de autenticación)
	adminID, exists := c.Get("admin_id")
	if !exists {
		adminID = "system" // fallback para desarrollo
	}

	// Llamar al servicio
	relation, err := h.guardianService.CreateGuardianRelation(
		c.Request.Context(),
		req,
		adminID.(string),
	)

	if err != nil {
		// Manejar errores usando shared/errors
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("create guardian relation failed",
				"error", appErr.Message,
				"code", appErr.Code,
				"guardian_id", req.GuardianID,
				"student_id", req.StudentID,
			)

			c.JSON(appErr.StatusCode, ErrorResponse{
				Error: appErr.Message,
				Code:  string(appErr.Code),
			})
			return
		}

		// Error no manejado
		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	// Log de éxito
	h.logger.Info("guardian relation created successfully",
		"relation_id", relation.ID,
		"guardian_id", relation.GuardianID,
		"student_id", relation.StudentID,
		"relationship_type", relation.RelationshipType,
		"created_by", adminID,
	)

	c.JSON(http.StatusCreated, relation)
}

// GetGuardianRelation godoc
// @Summary Get guardian relation by ID
// @Description Get details of a specific guardian-student relation
// @Tags guardians
// @Produce json
// @Param id path string true "Relation ID"
// @Success 200 {object} dto.GuardianRelationResponse
// @Failure 404 {object} ErrorResponse "Relation not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/guardian-relations/{id} [get]
// @Security BearerAuth
func (h *GuardianHandler) GetGuardianRelation(c *gin.Context) {
	id := c.Param("id")

	relation, err := h.guardianService.GetGuardianRelation(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Warn("get guardian relation failed",
				"error", appErr.Message,
				"code", appErr.Code,
				"id", id,
			)

			c.JSON(appErr.StatusCode, ErrorResponse{
				Error: appErr.Message,
				Code:  string(appErr.Code),
			})
			return
		}

		h.logger.Error("unexpected error", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, relation)
}

// GetGuardianRelations godoc
// @Summary Get all relations for a guardian
// @Description Get all student relations for a specific guardian
// @Tags guardians
// @Produce json
// @Param guardian_id path string true "Guardian ID"
// @Success 200 {array} dto.GuardianRelationResponse
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/guardians/{guardian_id}/relations [get]
// @Security BearerAuth
func (h *GuardianHandler) GetGuardianRelations(c *gin.Context) {
	guardianID := c.Param("guardian_id")

	relations, err := h.guardianService.GetGuardianRelations(c.Request.Context(), guardianID)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("get guardian relations failed",
				"error", appErr.Message,
				"code", appErr.Code,
				"guardian_id", guardianID,
			)

			c.JSON(appErr.StatusCode, ErrorResponse{
				Error: appErr.Message,
				Code:  string(appErr.Code),
			})
			return
		}

		h.logger.Error("unexpected error", "error", err, "guardian_id", guardianID)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, relations)
}

// GetStudentGuardians godoc
// @Summary Get all guardians for a student
// @Description Get all guardian relations for a specific student
// @Tags guardians
// @Produce json
// @Param student_id path string true "Student ID"
// @Success 200 {array} dto.GuardianRelationResponse
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/students/{student_id}/guardians [get]
// @Security BearerAuth
func (h *GuardianHandler) GetStudentGuardians(c *gin.Context) {
	studentID := c.Param("student_id")

	relations, err := h.guardianService.GetStudentGuardians(c.Request.Context(), studentID)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("get student guardians failed",
				"error", appErr.Message,
				"code", appErr.Code,
				"student_id", studentID,
			)

			c.JSON(appErr.StatusCode, ErrorResponse{
				Error: appErr.Message,
				Code:  string(appErr.Code),
			})
			return
		}

		h.logger.Error("unexpected error", "error", err, "student_id", studentID)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, relations)
}

// ErrorResponse representa una respuesta de error HTTP
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
