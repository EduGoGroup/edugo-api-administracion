package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/logger"
)

// SchoolHandler maneja las peticiones HTTP relacionadas con escuelas
type SchoolHandler struct {
	schoolService service.SchoolService
	logger        logger.Logger
}

// NewSchoolHandler crea un nuevo SchoolHandler
func NewSchoolHandler(schoolService service.SchoolService, logger logger.Logger) *SchoolHandler {
	return &SchoolHandler{
		schoolService: schoolService,
		logger:        logger,
	}
}

// CreateSchool godoc
// @Summary Create a new school
// @Description Creates a new school in the system
// @Tags schools
// @Accept json
// @Produce json
// @Param request body dto.CreateSchoolRequest true "School data"
// @Success 201 {object} dto.SchoolResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /v1/schools [post]
// @Security BearerAuth
func (h *SchoolHandler) CreateSchool(c *gin.Context) {
	var req dto.CreateSchoolRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	school, err := h.schoolService.CreateSchool(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("create school failed", "error", appErr.Message, "code", appErr.Code)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("school created", "school_id", school.ID, "name", school.Name)
	c.JSON(http.StatusCreated, school)
}
