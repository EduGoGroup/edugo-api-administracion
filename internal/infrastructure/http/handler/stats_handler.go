package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
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
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, stats)
}
