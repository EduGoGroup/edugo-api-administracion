package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
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
		_ = c.Error(err)
		return
	}

	h.logger.Info("material deleted successfully", "material_id", id)
	c.Status(http.StatusNoContent)
}
