package handler

import (
	"net/http"

	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// ErrorResponse representa una respuesta de error HTTP estándar
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// handleError maneja errores de forma centralizada en los handlers
// Reduce código duplicado y estandariza las respuestas de error
func handleError(c *gin.Context, log logger.Logger, err error, operation string) {
	if appErr, ok := errors.GetAppError(err); ok {
		log.Error(operation+" failed", "error", appErr.Message, "code", appErr.Code)
		c.JSON(appErr.StatusCode, ErrorResponse{
			Error: appErr.Message,
			Code:  string(appErr.Code),
		})
		return
	}

	log.Error("unexpected error", "error", err, "operation", operation)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "internal server error",
		Code:  "INTERNAL_ERROR",
	})
}
