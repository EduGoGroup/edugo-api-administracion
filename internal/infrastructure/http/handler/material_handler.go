package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/edugo/api-administracion/internal/application/service"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

// MaterialHandler maneja las peticiones HTTP relacionadas con materiales
type MaterialHandler struct {
	materialService service.MaterialService
	logger          logger.Logger
}

func NewMaterialHandler(materialService service.MaterialService, logger logger.Logger) *MaterialHandler {
	return &MaterialHandler{
		materialService: materialService,
		logger:          logger,
	}
}

// DeleteMaterial godoc
// @Summary Delete a material
// @Description Soft delete a material (mark as deleted)
// @Tags materials
// @Produce json
// @Param id path string true "Material ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse "Material not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials/{id} [delete]
// @Security BearerAuth
func (h *MaterialHandler) DeleteMaterial(c *gin.Context) {
	id := c.Param("id")

	err := h.materialService.DeleteMaterial(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("delete material failed", "error", appErr.Message, "id", id)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("material deleted successfully", "material_id", id)
	c.Status(http.StatusNoContent)
}
