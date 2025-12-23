package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	httpdto "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"
	"github.com/EduGoGroup/edugo-shared/logger"
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
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	school, err := h.schoolService.CreateSchool(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.logger.Info("school created", "school_id", school.ID, "name", school.Name)
	c.JSON(http.StatusCreated, school)
}

// GetSchool godoc
// @Summary Get a school by ID
// @Description Retrieves a school by its ID
// @Tags schools
// @Produce json
// @Param id path string true "School ID"
// @Success 200 {object} dto.SchoolResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/schools/{id} [get]
// @Security BearerAuth
func (h *SchoolHandler) GetSchool(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}

	school, err := h.schoolService.GetSchool(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, school)
}

// GetSchoolByCode godoc
// @Summary Get a school by code
// @Description Retrieves a school by its unique code
// @Tags schools
// @Produce json
// @Param code path string true "School Code"
// @Success 200 {object} dto.SchoolResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/schools/code/{code} [get]
// @Security BearerAuth
func (h *SchoolHandler) GetSchoolByCode(c *gin.Context) {
	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "school code is required", Code: "INVALID_REQUEST"})
		return
	}

	school, err := h.schoolService.GetSchoolByCode(c.Request.Context(), code)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, school)
}

// ListSchools godoc
// @Summary List all schools
// @Description Retrieves a list of all schools in the system
// @Tags schools
// @Produce json
// @Success 200 {array} dto.SchoolResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/schools [get]
// @Security BearerAuth
func (h *SchoolHandler) ListSchools(c *gin.Context) {
	schools, err := h.schoolService.ListSchools(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, schools)
}

// UpdateSchool godoc
// @Summary Update a school
// @Description Updates an existing school
// @Tags schools
// @Accept json
// @Produce json
// @Param id path string true "School ID"
// @Param request body dto.UpdateSchoolRequest true "Updated school data"
// @Success 200 {object} dto.SchoolResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/schools/{id} [put]
// @Security BearerAuth
func (h *SchoolHandler) UpdateSchool(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}

	var req dto.UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err, "school_id", id)
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	school, err := h.schoolService.UpdateSchool(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.logger.Info("school updated", "school_id", id)
	c.JSON(http.StatusOK, school)
}

// DeleteSchool godoc
// @Summary Delete a school
// @Description Soft deletes a school from the system
// @Tags schools
// @Produce json
// @Param id path string true "School ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/schools/{id} [delete]
// @Security BearerAuth
func (h *SchoolHandler) DeleteSchool(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}

	err := h.schoolService.DeleteSchool(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.logger.Info("school deleted", "school_id", id)
	c.Status(http.StatusNoContent)
}
