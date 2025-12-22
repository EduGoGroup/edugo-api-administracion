package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
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
		handleError(c, h.logger, err, "create guardian relation")
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
		handleError(c, h.logger, err, "get guardian relation")
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
		handleError(c, h.logger, err, "get guardian relations")
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
		handleError(c, h.logger, err, "get student guardians")
		return
	}

	c.JSON(http.StatusOK, relations)
}
