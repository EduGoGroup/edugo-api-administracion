package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/logger"
)

// StatsHandler maneja las peticiones HTTP relacionadas con estadísticas
type StatsHandler struct {
	statsService service.StatsService
	logger       logger.Logger
}

func NewStatsHandler(statsService service.StatsService, logger logger.Logger) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
		logger:       logger,
	}
}

// GetGlobalStats godoc
// @Summary Get global statistics
// @Description Get system-wide statistics (users, schools, subjects, etc.)
// @Tags stats
// @Produce json
// @Success 200 {object} dto.GlobalStatsResponse
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/stats/global [get]
// @Security BearerAuth
func (h *StatsHandler) GetGlobalStats(c *gin.Context) {
	stats, err := h.statsService.GetGlobalStats(c.Request.Context())
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("get global stats failed", "error", appErr.Message)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
