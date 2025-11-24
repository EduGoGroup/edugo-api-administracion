// Package cache proporciona utilidades de cache usando Redis
package cache

// RedisClient wrappea las operaciones de cache con Redis
// Proporciona métodos de alto nivel para cache de tokens y validaciones
// Será implementado en FASE 2 del Sprint 1
type RedisClient interface {
	// Get obtiene un valor del cache
	// Get(ctx context.Context, key string) (string, error)

	// Set almacena un valor en el cache con TTL
	// Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete elimina un valor del cache
	// Delete(ctx context.Context, key string) error

	// Exists verifica si una key existe
	// Exists(ctx context.Context, key string) (bool, error)

	// Incr incrementa un contador (para rate limiting)
	// Incr(ctx context.Context, key string) (int64, error)

	// Expire establece TTL en una key existente
	// Expire(ctx context.Context, key string, ttl time.Duration) error
}
