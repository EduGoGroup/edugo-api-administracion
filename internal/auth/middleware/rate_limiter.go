// Package middleware contiene middlewares para autenticación
package middleware

// RateLimiter implementa rate limiting diferenciado para autenticación
// - Servicios internos: 1000 req/min (api-mobile, worker)
// - Clientes externos: 60 req/min
// Será implementado en FASE 2 del Sprint 1
type RateLimiter struct {
	// cache cache.RedisClient
	// config *config.RateLimitConfig
}

// NewRateLimiter crea una nueva instancia de RateLimiter
// func NewRateLimiter(cache cache.RedisClient, config *config.RateLimitConfig) *RateLimiter {
// 	return &RateLimiter{cache: cache, config: config}
// }
