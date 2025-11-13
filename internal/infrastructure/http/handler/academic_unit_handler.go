package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// AcademicUnitHandler maneja las peticiones HTTP relacionadas con unidades académicas
type AcademicUnitHandler struct {
	unitService service.AcademicUnitService
	logger      logger.Logger
}

// NewAcademicUnitHandler crea un nuevo AcademicUnitHandler
func NewAcademicUnitHandler(unitService service.AcademicUnitService, logger logger.Logger) *AcademicUnitHandler {
	return &AcademicUnitHandler{
		unitService: unitService,
		logger:      logger,
	}
}

// CreateUnit godoc
// @Summary Create a new academic unit
// @Description Creates a new academic unit (grade, section, club, department) within a school
// @Tags academic-units
// @Accept json
// @Produce json
// @Param schoolId path string true "School ID"
// @Param request body dto.CreateAcademicUnitRequest true "Academic unit data"
// @Success 201 {object} dto.AcademicUnitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/schools/{schoolId}/units [post]
// @Security BearerAuth
func (h *AcademicUnitHandler) CreateUnit(c *gin.Context) {
	schoolID := c.Param("schoolId")

	if schoolID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}

	var req dto.CreateAcademicUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err, "school_id", schoolID)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	unit, err := h.unitService.CreateUnit(c.Request.Context(), schoolID, req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("create unit failed", "error", appErr.Message, "code", appErr.Code, "school_id", schoolID)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "school_id", schoolID)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("unit created", "unit_id", unit.ID, "school_id", schoolID, "type", unit.Type)
	c.JSON(http.StatusCreated, unit)
}

// GetUnit godoc
// @Summary Get an academic unit by ID
// @Description Retrieves an academic unit by its ID
// @Tags academic-units
// @Produce json
// @Param id path string true "Unit ID"
// @Success 200 {object} dto.AcademicUnitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/units/{id} [get]
// @Security BearerAuth
func (h *AcademicUnitHandler) GetUnit(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unit ID is required", Code: "INVALID_REQUEST"})
		return
	}

	unit, err := h.unitService.GetUnit(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("get unit failed", "error", appErr.Message, "code", appErr.Code, "unit_id", id)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "unit_id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, unit)
}

// GetUnitTree godoc
// @Summary Get the hierarchical tree of units for a school
// @Description Retrieves the complete hierarchy tree of academic units for a school
// @Tags academic-units
// @Produce json
// @Param schoolId path string true "School ID"
// @Success 200 {array} dto.UnitTreeNode
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/schools/{schoolId}/units/tree [get]
// @Security BearerAuth
func (h *AcademicUnitHandler) GetUnitTree(c *gin.Context) {
	schoolID := c.Param("schoolId")

	if schoolID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}

	tree, err := h.unitService.GetUnitTree(c.Request.Context(), schoolID)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("get unit tree failed", "error", appErr.Message, "code", appErr.Code, "school_id", schoolID)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "school_id", schoolID)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// ListUnitsBySchool godoc
// @Summary List all units for a school
// @Description Retrieves all academic units for a specific school
// @Tags academic-units
// @Produce json
// @Param schoolId path string true "School ID"
// @Param includeDeleted query bool false "Include soft-deleted units"
// @Success 200 {array} dto.AcademicUnitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/schools/{schoolId}/units [get]
// @Security BearerAuth
func (h *AcademicUnitHandler) ListUnitsBySchool(c *gin.Context) {
	schoolID := c.Param("schoolId")

	if schoolID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}

	includeDeleted := c.DefaultQuery("includeDeleted", "false") == "true"

	units, err := h.unitService.ListUnitsBySchool(c.Request.Context(), schoolID, includeDeleted)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("list units failed", "error", appErr.Message, "code", appErr.Code, "school_id", schoolID)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "school_id", schoolID)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, units)
}

// ListUnitsByType godoc
// @Summary List units by type
// @Description Retrieves academic units filtered by type (grade, section, club, department)
// @Tags academic-units
// @Produce json
// @Param schoolId path string true "School ID"
// @Param type query string true "Unit type (grade, section, club, department)"
// @Success 200 {array} dto.AcademicUnitResponse
// @Failure 400 {object} ErrorResponse
// @Router /v1/schools/{schoolId}/units/by-type [get]
// @Security BearerAuth
func (h *AcademicUnitHandler) ListUnitsByType(c *gin.Context) {
	schoolID := c.Param("schoolId")
	unitType := c.Query("type")

	if schoolID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}

	if unitType == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unit type is required", Code: "INVALID_REQUEST"})
		return
	}

	units, err := h.unitService.ListUnitsByType(c.Request.Context(), schoolID, unitType)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("list units by type failed", "error", appErr.Message, "code", appErr.Code, "school_id", schoolID, "type", unitType)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "school_id", schoolID, "type", unitType)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, units)
}

// UpdateUnit godoc
// @Summary Update an academic unit
// @Description Updates an existing academic unit
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "Unit ID"
// @Param request body dto.UpdateAcademicUnitRequest true "Updated unit data"
// @Success 200 {object} dto.AcademicUnitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/units/{id} [put]
// @Security BearerAuth
func (h *AcademicUnitHandler) UpdateUnit(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unit ID is required", Code: "INVALID_REQUEST"})
		return
	}

	var req dto.UpdateAcademicUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err, "unit_id", id)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	unit, err := h.unitService.UpdateUnit(c.Request.Context(), id, req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("update unit failed", "error", appErr.Message, "code", appErr.Code, "unit_id", id)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "unit_id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("unit updated", "unit_id", id)
	c.JSON(http.StatusOK, unit)
}

// DeleteUnit godoc
// @Summary Delete an academic unit
// @Description Soft deletes an academic unit
// @Tags academic-units
// @Produce json
// @Param id path string true "Unit ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/units/{id} [delete]
// @Security BearerAuth
func (h *AcademicUnitHandler) DeleteUnit(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unit ID is required", Code: "INVALID_REQUEST"})
		return
	}

	err := h.unitService.DeleteUnit(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("delete unit failed", "error", appErr.Message, "code", appErr.Code, "unit_id", id)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "unit_id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("unit deleted", "unit_id", id)
	c.Status(http.StatusNoContent)
}

// RestoreUnit godoc
// @Summary Restore a soft-deleted academic unit
// @Description Restores a previously soft-deleted academic unit
// @Tags academic-units
// @Produce json
// @Param id path string true "Unit ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/units/{id}/restore [post]
// @Security BearerAuth
func (h *AcademicUnitHandler) RestoreUnit(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unit ID is required", Code: "INVALID_REQUEST"})
		return
	}

	err := h.unitService.RestoreUnit(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("restore unit failed", "error", appErr.Message, "code", appErr.Code, "unit_id", id)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "unit_id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("unit restored", "unit_id", id)
	c.Status(http.StatusNoContent)
}

// GetHierarchyPath godoc
// @Summary Get the hierarchy path for a unit
// @Description Retrieves the complete hierarchy path from root to the specified unit
// @Tags academic-units
// @Produce json
// @Param id path string true "Unit ID"
// @Success 200 {array} dto.AcademicUnitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/units/{id}/hierarchy-path [get]
// @Security BearerAuth
func (h *AcademicUnitHandler) GetHierarchyPath(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unit ID is required", Code: "INVALID_REQUEST"})
		return
	}

	path, err := h.unitService.GetHierarchyPath(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("get hierarchy path failed", "error", appErr.Message, "code", appErr.Code, "unit_id", id)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "unit_id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, path)
}
