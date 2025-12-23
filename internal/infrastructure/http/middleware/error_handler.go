package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	httpdto "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// ErrorHandler middleware que procesa errores de forma centralizada
func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Procesar errores si existen
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, log, err)
		}
	}
}

func handleError(c *gin.Context, log logger.Logger, err error) {
	// Intentar obtener AppError
	if appErr, ok := errors.GetAppError(err); ok {
		log.Error("request failed",
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"error", appErr.Message,
			"code", appErr.Code,
			"status", appErr.StatusCode,
		)
		c.JSON(appErr.StatusCode, httpdto.ErrorResponse{
			Error: appErr.Message,
			Code:  string(appErr.Code),
		})
		return
	}

	// Error genérico
	log.Error("unexpected error",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"error", err.Error(),
	)
	c.JSON(http.StatusInternalServerError, httpdto.ErrorResponse{
		Error: "internal server error",
		Code:  "INTERNAL_ERROR",
	})
}
