package handler

import (
	"net/http"

	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

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

// handleErrorWithContext maneja errores incluyendo contexto adicional para logging
func handleErrorWithContext(c *gin.Context, log logger.Logger, err error, operation string, context map[string]interface{}) {
	args := []interface{}{"operation", operation}
	for k, v := range context {
		args = append(args, k, v)
	}

	if appErr, ok := errors.GetAppError(err); ok {
		args = append(args, "error", appErr.Message, "code", appErr.Code)
		log.Error(operation+" failed", args...)
		c.JSON(appErr.StatusCode, ErrorResponse{
			Error: appErr.Message,
			Code:  string(appErr.Code),
		})
		return
	}

	args = append(args, "error", err)
	log.Error("unexpected error", args...)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "internal server error",
		Code:  "INTERNAL_ERROR",
	})
}
