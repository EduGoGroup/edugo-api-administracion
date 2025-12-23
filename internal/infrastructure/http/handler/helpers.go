package handler

import (
	"net/http"

	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// handleError maneja errores de forma centralizada en los handlers
// Reduce código duplicado y estandariza las respuestas de error
func handleError(c *gin.Context, log logger.Logger, err error, operation string) {
	if appErr, ok := errors.GetAppError(err); ok {
		log.Error(operation+" failed", "error", appErr.Message, "code", appErr.Code)
		c.JSON(appErr.StatusCode, dto.ErrorResponse{
			Error: appErr.Message,
			Code:  string(appErr.Code),
		})
		return
	}

	log.Error("unexpected error", "error", err, "operation", operation)
	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Error: "internal server error",
		Code:  "INTERNAL_ERROR",
	})
}
