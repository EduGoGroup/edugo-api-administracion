// Package middleware contiene middlewares para autenticación
package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/dto"
)

// RateLimitConfig configuración del rate limiter
type RateLimitConfig struct {
	// Límites para servicios internos (api-mobile, worker)
	InternalMaxRequests int
	InternalWindow      time.Duration

	// Límites para clientes externos
	ExternalMaxRequests int
	ExternalWindow      time.Duration

	// Identificación de servicios internos
	InternalAPIKeys []string
	InternalIPs     []string
}

// rateLimitEntry representa una entrada en el rate limiter
type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

// RateLimiter implementa rate limiting en memoria
// Para producción, usar Redis para estado compartido
type RateLimiter struct {
	config       RateLimitConfig
	entries      map[string]*rateLimitEntry
	mutex        sync.RWMutex
	internalNets []*net.IPNet
	apiKeys      map[string]bool
	stopCleanup  chan struct{}
}

// NewRateLimiter crea un nuevo rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	// Valores por defecto
	if config.InternalMaxRequests == 0 {
		config.InternalMaxRequests = 1000
	}
	if config.InternalWindow == 0 {
		config.InternalWindow = time.Minute
	}
	if config.ExternalMaxRequests == 0 {
		config.ExternalMaxRequests = 60
	}
	if config.ExternalWindow == 0 {
		config.ExternalWindow = time.Minute
	}

	rl := &RateLimiter{
		config:      config,
		entries:     make(map[string]*rateLimitEntry),
		apiKeys:     make(map[string]bool),
		stopCleanup: make(chan struct{}),
	}

	// Parsear rangos CIDR
	for _, cidr := range config.InternalIPs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			// Si no es CIDR, intentar como IP simple
			ip := net.ParseIP(cidr)
			if ip != nil {
				if ip.To4() != nil {
					_, ipNet, _ = net.ParseCIDR(cidr + "/32")
				} else {
					_, ipNet, _ = net.ParseCIDR(cidr + "/128")
				}
			}
		}
		if ipNet != nil {
			rl.internalNets = append(rl.internalNets, ipNet)
		}
	}

	for _, key := range config.InternalAPIKeys {
		rl.apiKeys[key] = true
	}

	// Iniciar limpieza periódica
	go rl.cleanupLoop()

	return rl
}

// Middleware retorna el middleware de Gin
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := rl.getIdentifier(c)
		isInternal := rl.isInternalService(c)

		allowed, remaining, resetTime := rl.checkLimit(identifier, isInternal)

		// Agregar headers de rate limit
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.getLimit(isInternal)))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

		if !allowed {
			retryAfter := time.Until(resetTime).Seconds()
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter))

			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Error:   "rate_limit_exceeded",
				Message: "Demasiadas solicitudes. Intente de nuevo más tarde.",
				Code:    "RATE_LIMIT",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkLimit verifica si el request está permitido
func (rl *RateLimiter) checkLimit(identifier string, isInternal bool) (allowed bool, remaining int, resetTime time.Time) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	maxRequests := rl.config.ExternalMaxRequests
	window := rl.config.ExternalWindow

	if isInternal {
		maxRequests = rl.config.InternalMaxRequests
		window = rl.config.InternalWindow
	}

	now := time.Now()
	entry, exists := rl.entries[identifier]

	// Si no existe o la ventana expiró, crear nueva entrada
	if !exists || now.After(entry.resetTime) {
		rl.entries[identifier] = &rateLimitEntry{
			count:     1,
			resetTime: now.Add(window),
		}
		return true, maxRequests - 1, now.Add(window)
	}

	// Verificar límite
	if entry.count >= maxRequests {
		return false, 0, entry.resetTime
	}

	// Incrementar contador
	entry.count++
	return true, maxRequests - entry.count, entry.resetTime
}

// getLimit retorna el límite aplicable
func (rl *RateLimiter) getLimit(isInternal bool) int {
	if isInternal {
		return rl.config.InternalMaxRequests
	}
	return rl.config.ExternalMaxRequests
}

// getIdentifier obtiene el identificador único para rate limiting
func (rl *RateLimiter) getIdentifier(c *gin.Context) string {
	// Prioridad: API Key > IP
	apiKey := c.GetHeader("X-Service-API-Key")
	if apiKey != "" {
		return "apikey:" + apiKey
	}

	return "ip:" + c.ClientIP()
}

// isInternalService verifica si es un servicio interno
func (rl *RateLimiter) isInternalService(c *gin.Context) bool {
	// Verificar API Key
	apiKey := c.GetHeader("X-Service-API-Key")
	if apiKey != "" && rl.apiKeys[apiKey] {
		return true
	}

	// Verificar IP
	clientIP := net.ParseIP(c.ClientIP())
	if clientIP != nil {
		for _, ipNet := range rl.internalNets {
			if ipNet.Contains(clientIP) {
				return true
			}
		}
	}

	return false
}

// cleanupLoop limpia entradas expiradas periódicamente
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopCleanup:
			return
		}
	}
}

// cleanup elimina entradas expiradas
func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for key, entry := range rl.entries {
		if now.After(entry.resetTime) {
			delete(rl.entries, key)
		}
	}
}

// Stop detiene el rate limiter (para tests)
func (rl *RateLimiter) Stop() {
	close(rl.stopCleanup)
}

// GetStats retorna estadísticas del rate limiter
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	return map[string]interface{}{
		"active_entries":         len(rl.entries),
		"internal_limit":         rl.config.InternalMaxRequests,
		"external_limit":         rl.config.ExternalMaxRequests,
		"internal_window":        rl.config.InternalWindow.String(),
		"external_window":        rl.config.ExternalWindow.String(),
		"configured_api_keys":    len(rl.apiKeys),
		"configured_ip_ranges":   len(rl.internalNets),
	}
}
