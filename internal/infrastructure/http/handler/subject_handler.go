package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	httpdto "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// SubjectHandler maneja las peticiones HTTP relacionadas con materias
type SubjectHandler struct {
	subjectService service.SubjectService
	logger         logger.Logger
}

func NewSubjectHandler(subjectService service.SubjectService, logger logger.Logger) *SubjectHandler {
	return &SubjectHandler{
		subjectService: subjectService,
		logger:         logger,
	}
}

// CreateSubject godoc
// @Summary Create a new subject
// @Tags subjects
// @Accept json
// @Produce json
// @Param request body dto.CreateSubjectRequest true "Subject data"
// @Success 201 {object} dto.SubjectResponse
// @Router /v1/subjects [post]
// @Security BearerAuth
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var req dto.CreateSubjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	subject, err := h.subjectService.CreateSubject(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.logger.Info("subject created", "subject_id", subject.ID)
	c.JSON(http.StatusCreated, subject)
}

// UpdateSubject godoc
// @Summary Update subject
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID"
// @Param request body dto.UpdateSubjectRequest true "Update data"
// @Success 200 {object} dto.SubjectResponse
// @Router /v1/subjects/{id} [patch]
// @Security BearerAuth
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateSubjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httpdto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	subject, err := h.subjectService.UpdateSubject(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.logger.Info("subject updated", "subject_id", subject.ID)
	c.JSON(http.StatusOK, subject)
}

// GetSubject godoc
// @Summary Get subject by ID
// @Description Retrieves a subject by its unique identifier
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID" format(uuid)
// @Success 200 {object} dto.SubjectResponse
// @Failure 404 {object} httpdto.ErrorResponse
// @Router /v1/subjects/{id} [get]
// @Security BearerAuth
func (h *SubjectHandler) GetSubject(c *gin.Context) {
	id := c.Param("id")

	subject, err := h.subjectService.GetSubject(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, subject)
}

// ListSubjects godoc
// @Summary List subjects
// @Description Lists all subjects, optionally filtered by school
// @Tags subjects
// @Accept json
// @Produce json
// @Param school_id query string false "Filter by school ID" format(uuid)
// @Success 200 {array} dto.SubjectResponse
// @Router /v1/subjects [get]
// @Security BearerAuth
func (h *SubjectHandler) ListSubjects(c *gin.Context) {
	schoolID := c.Query("school_id")

	subjects, err := h.subjectService.ListSubjects(c.Request.Context(), schoolID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// DeleteSubject godoc
// @Summary Delete subject
// @Description Soft deletes a subject
// @Tags subjects
// @Param id path string true "Subject ID" format(uuid)
// @Success 204 "No Content"
// @Failure 404 {object} httpdto.ErrorResponse
// @Router /v1/subjects/{id} [delete]
// @Security BearerAuth
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	id := c.Param("id")

	err := h.subjectService.DeleteSubject(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.logger.Info("subject deleted", "subject_id", id)
	c.Status(http.StatusNoContent)
}
