package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	httpdto "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"
	"github.com/EduGoGroup/edugo-shared/logger"
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
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	unit, err := h.unitService.CreateUnit(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
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
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	unit, err := h.unitService.UpdateUnit(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
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
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /v1/units/{id}/members [post]
// @Security BearerAuth
func (h *UnitHandler) AssignMember(c *gin.Context) {
	// TODO: Migrar esta funcionalidad a usar MembershipService
	// Esta funcionalidad debería usar el servicio de membresías en lugar del servicio de unidades
	h.logger.Warn("AssignMember endpoint deprecated - use membership service instead")
	c.JSON(http.StatusNotImplemented, httpdto.ErrorResponse{
		Error: "This endpoint is being migrated. Please use /v1/memberships endpoint instead",
		Code:  "NOT_IMPLEMENTED",
	})
}
