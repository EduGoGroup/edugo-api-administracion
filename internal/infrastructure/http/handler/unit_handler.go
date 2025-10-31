package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/logger"
)

// UnitHandler maneja las peticiones HTTP relacionadas con unidades
type UnitHandler struct {
	unitService service.UnitService
	logger      logger.Logger
}

func NewUnitHandler(unitService service.UnitService, logger logger.Logger) *UnitHandler {
	return &UnitHandler{
		unitService: unitService,
		logger:      logger,
	}
}

// CreateUnit godoc
// @Summary Create a new unit
// @Description Creates a new organizational unit (department, grade, group, etc.)
// @Tags units
// @Accept json
// @Produce json
// @Param request body dto.CreateUnitRequest true "Unit data"
// @Success 201 {object} dto.UnitResponse
// @Failure 400 {object} ErrorResponse
// @Router /v1/units [post]
// @Security BearerAuth
func (h *UnitHandler) CreateUnit(c *gin.Context) {
	var req dto.CreateUnitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	unit, err := h.unitService.CreateUnit(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("create unit failed", "error", appErr.Message, "code", appErr.Code)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("unit created", "unit_id", unit.ID, "name", unit.Name)
	c.JSON(http.StatusCreated, unit)
}

// UpdateUnit godoc
// @Summary Update unit
// @Description Update unit information
// @Tags units
// @Accept json
// @Produce json
// @Param id path string true "Unit ID"
// @Param request body dto.UpdateUnitRequest true "Update data"
// @Success 200 {object} dto.UnitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/units/{id} [patch]
// @Security BearerAuth
func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUnitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	unit, err := h.unitService.UpdateUnit(c.Request.Context(), id, req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("unit updated", "unit_id", unit.ID)
	c.JSON(http.StatusOK, unit)
}

// AssignMember godoc
// @Summary Assign member to unit
// @Description Assign a user (teacher or student) to a unit
// @Tags units
// @Accept json
// @Produce json
// @Param id path string true "Unit ID"
// @Param request body dto.AssignMemberRequest true "Member data"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /v1/units/{id}/members [post]
// @Security BearerAuth
func (h *UnitHandler) AssignMember(c *gin.Context) {
	id := c.Param("id")
	var req dto.AssignMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	err := h.unitService.AssignMember(c.Request.Context(), id, req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("assign member failed", "error", appErr.Message, "unit_id", id)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("member assigned", "unit_id", id, "user_id", req.UserID)
	c.Status(http.StatusNoContent)
}
