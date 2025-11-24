// Package ratelimit proporciona utilidades de rate limiting
package ratelimit

// Limiter implementa rate limiting genérico usando Redis
// Soporta sliding window y token bucket algorithms
// Será implementado en FASE 2 del Sprint 1
type Limiter interface {
	// Allow verifica si una request está permitida
	// Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)

	// GetRemaining retorna requests restantes en la ventana actual
	// GetRemaining(ctx context.Context, key string) (int, error)

	// Reset reinicia el contador para una key
	// Reset(ctx context.Context, key string) error
}

// Config contiene la configuración de rate limiting
type Config struct {
	MaxRequests int    // Máximo de requests permitidas
	Window      string // Ventana de tiempo (ej: "1m", "1h")
}
