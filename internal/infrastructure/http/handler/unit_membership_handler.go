package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	httpdto "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// UnitMembershipHandler maneja las peticiones HTTP relacionadas con membresías de unidades
type UnitMembershipHandler struct {
	membershipService service.UnitMembershipService
	logger            logger.Logger
}

// NewUnitMembershipHandler crea un nuevo UnitMembershipHandler
func NewUnitMembershipHandler(membershipService service.UnitMembershipService, logger logger.Logger) *UnitMembershipHandler {
	return &UnitMembershipHandler{
		membershipService: membershipService,
		logger:            logger,
	}
}

// CreateMembership godoc
// @Summary Create a new membership
// @Description Assigns a user to an academic unit with a specific role
// @Tags memberships
// @Accept json
// @Produce json
// @Param request body dto.CreateMembershipRequest true "Membership data"
// @Success 201 {object} dto.MembershipResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /v1/memberships [post]
// @Security BearerAuth
func (h *UnitMembershipHandler) CreateMembership(c *gin.Context) {
	var req dto.CreateMembershipRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	membership, err := h.membershipService.CreateMembership(c.Request.Context(), req)
	if err != nil {
		handleError(c, h.logger, err, "create membership")
		return
	}

	h.logger.Info("membership created", "membership_id", membership.ID, "unit_id", membership.UnitID, "user_id", membership.UserID)
	c.JSON(http.StatusCreated, membership)
}

// GetMembership godoc
// @Summary Get a membership by ID
// @Description Retrieves a membership by its ID
// @Tags memberships
// @Produce json
// @Param id path string true "Membership ID"
// @Success 200 {object} dto.MembershipResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/memberships/{id} [get]
// @Security BearerAuth
func (h *UnitMembershipHandler) GetMembership(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "membership ID is required", Code: "INVALID_REQUEST"})
		return
	}

	membership, err := h.membershipService.GetMembership(c.Request.Context(), id)
	if err != nil {
		handleError(c, h.logger, err, "get membership")
		return
	}

	c.JSON(http.StatusOK, membership)
}

// ListMembershipsByUnit godoc
// @Summary List memberships for a unit
// @Description Retrieves all memberships for a specific academic unit
// @Tags memberships
// @Produce json
// @Param unitId path string true "Unit ID"
// @Param activeOnly query bool false "Show only active memberships"
// @Success 200 {array} dto.MembershipResponse
// @Failure 400 {object} ErrorResponse
// @Router /v1/units/{unitId}/memberships [get]
// @Security BearerAuth
func (h *UnitMembershipHandler) ListMembershipsByUnit(c *gin.Context) {
	unitID := c.Param("unitId")

	if unitID == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "unit ID is required", Code: "INVALID_REQUEST"})
		return
	}

	activeOnly := c.DefaultQuery("activeOnly", "true") == "true"

	memberships, err := h.membershipService.ListMembershipsByUnit(c.Request.Context(), unitID, activeOnly)
	if err != nil {
		handleError(c, h.logger, err, "list memberships by unit")
		return
	}

	c.JSON(http.StatusOK, memberships)
}

// ListMembershipsByUser godoc
// @Summary List memberships for a user
// @Description Retrieves all memberships for a specific user across all units
// @Tags memberships
// @Produce json
// @Param userId path string true "User ID"
// @Param activeOnly query bool false "Show only active memberships"
// @Success 200 {array} dto.MembershipResponse
// @Failure 400 {object} ErrorResponse
// @Router /v1/users/{userId}/memberships [get]
// @Security BearerAuth
func (h *UnitMembershipHandler) ListMembershipsByUser(c *gin.Context) {
	userID := c.Param("userId")

	if userID == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "user ID is required", Code: "INVALID_REQUEST"})
		return
	}

	activeOnly := c.DefaultQuery("activeOnly", "true") == "true"

	memberships, err := h.membershipService.ListMembershipsByUser(c.Request.Context(), userID, activeOnly)
	if err != nil {
		handleError(c, h.logger, err, "list memberships by user")
		return
	}

	c.JSON(http.StatusOK, memberships)
}

// ListMembershipsByRole godoc
// @Summary List memberships by role
// @Description Retrieves memberships filtered by role within a unit
// @Tags memberships
// @Produce json
// @Param unitId path string true "Unit ID"
// @Param role query string true "Role (owner, teacher, assistant, student, guardian)"
// @Param activeOnly query bool false "Show only active memberships"
// @Success 200 {array} dto.MembershipResponse
// @Failure 400 {object} ErrorResponse
// @Router /v1/units/{unitId}/memberships/by-role [get]
// @Security BearerAuth
func (h *UnitMembershipHandler) ListMembershipsByRole(c *gin.Context) {
	unitID := c.Param("unitId")
	role := c.Query("role")

	if unitID == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "unit ID is required", Code: "INVALID_REQUEST"})
		return
	}

	if role == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "role is required", Code: "INVALID_REQUEST"})
		return
	}

	activeOnly := c.DefaultQuery("activeOnly", "true") == "true"

	memberships, err := h.membershipService.ListMembershipsByRole(c.Request.Context(), unitID, role, activeOnly)
	if err != nil {
		handleError(c, h.logger, err, "list memberships by role")
		return
	}

	c.JSON(http.StatusOK, memberships)
}

// UpdateMembership godoc
// @Summary Update a membership
// @Description Updates an existing membership (e.g., change role or validity dates)
// @Tags memberships
// @Accept json
// @Produce json
// @Param id path string true "Membership ID"
// @Param request body dto.UpdateMembershipRequest true "Updated membership data"
// @Success 200 {object} dto.MembershipResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/memberships/{id} [put]
// @Security BearerAuth
func (h *UnitMembershipHandler) UpdateMembership(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "membership ID is required", Code: "INVALID_REQUEST"})
		return
	}

	var req dto.UpdateMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err, "membership_id", id)
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	membership, err := h.membershipService.UpdateMembership(c.Request.Context(), id, req)
	if err != nil {
		handleError(c, h.logger, err, "update membership")
		return
	}

	h.logger.Info("membership updated", "membership_id", id)
	c.JSON(http.StatusOK, membership)
}

// ExpireMembership godoc
// @Summary Expire a membership
// @Description Sets a membership as expired (sets valid_until to now)
// @Tags memberships
// @Produce json
// @Param id path string true "Membership ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/memberships/{id}/expire [post]
// @Security BearerAuth
func (h *UnitMembershipHandler) ExpireMembership(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "membership ID is required", Code: "INVALID_REQUEST"})
		return
	}

	err := h.membershipService.ExpireMembership(c.Request.Context(), id)
	if err != nil {
		handleError(c, h.logger, err, "expire membership")
		return
	}

	h.logger.Info("membership expired", "membership_id", id)
	c.Status(http.StatusNoContent)
}

// DeleteMembership godoc
// @Summary Delete a membership
// @Description Permanently deletes a membership from the system
// @Tags memberships
// @Produce json
// @Param id path string true "Membership ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/memberships/{id} [delete]
// @Security BearerAuth
func (h *UnitMembershipHandler) DeleteMembership(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "membership ID is required", Code: "INVALID_REQUEST"})
		return
	}

	err := h.membershipService.DeleteMembership(c.Request.Context(), id)
	if err != nil {
		handleError(c, h.logger, err, "delete membership")
		return
	}

	h.logger.Info("membership deleted", "membership_id", id)
	c.Status(http.StatusNoContent)
}
