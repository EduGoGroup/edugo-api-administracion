package middleware

import (
	"strings"

	"github.com/EduGoGroup/edugo-api-administracion/internal/config"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura CORS basado en variables de entorno
func CORSMiddleware(cfg *config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Validar si el origen está permitido
		allowedOrigins := parseCSV(cfg.AllowedOrigins)
		if isOriginAllowed(origin, allowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		
		// Configurar métodos y headers permitidos
		c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.AllowedMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.AllowedHeaders)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 horas
		
		// Manejar preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// parseCSV convierte un string CSV en slice de strings
func parseCSV(csv string) []string {
	if csv == "" {
		return []string{}
	}
	
	parts := strings.Split(csv, ",")
	result := make([]string, 0, len(parts))
	
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}

// isOriginAllowed verifica si un origen está en la lista de permitidos
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}
	
	for _, allowed := range allowedOrigins {
		// Soporte para wildcard
		if allowed == "*" {
			return true
		}
		
		// Comparación exacta
		if allowed == origin {
			return true
		}
	}
	
	return false
}
