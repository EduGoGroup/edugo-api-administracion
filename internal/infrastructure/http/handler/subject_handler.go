package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/edugo/api-administracion/internal/application/dto"
	"github.com/edugo/api-administracion/internal/application/service"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	subject, err := h.subjectService.CreateSubject(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("create subject failed", "error", appErr.Message)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	subject, err := h.subjectService.UpdateSubject(c.Request.Context(), id, req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("subject updated", "subject_id", subject.ID)
	c.JSON(http.StatusOK, subject)
}
