package middleware

import (
	"strings"

	"github.com/EduGoGroup/edugo-api-administracion/internal/config"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura CORS basado en variables de entorno
func CORSMiddleware(cfg *config.CORSConfig) gin.HandlerFunc {
	// Parsear orígenes permitidos una sola vez (Issue #5 - Performance)
	allowedOrigins := parseCSV(cfg.AllowedOrigins)
	
	// Detectar si hay wildcard en la configuración (Issue #2 - Seguridad)
	hasWildcard := false
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			hasWildcard = true
			break
		}
	}
	
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Validar y configurar origen permitido (Issue #2 - Wildcard + Credentials)
		if isOriginAllowed(origin, allowedOrigins) {
			if hasWildcard {
				// Cuando se usa wildcard, no se deben permitir credenciales
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				// Lista explícita de orígenes: permitir credenciales para el origen concreto
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}
		
		// Manejar preflight requests (Issue #1 - Headers solo en OPTIONS)
		if c.Request.Method == "OPTIONS" {
			// Configurar métodos y headers permitidos solo para preflight
			c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.AllowedMethods)
			c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.AllowedHeaders)
			c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 horas
			
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
