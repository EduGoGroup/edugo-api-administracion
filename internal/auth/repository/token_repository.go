// Package repository define las interfaces de persistencia para autenticación
package repository

// TokenRepository define las operaciones de persistencia para refresh tokens
// Almacena tokens en Redis para revocación rápida
// Será implementado en FASE 2 del Sprint 1
type TokenRepository interface {
	// StoreRefreshToken almacena un refresh token
	// StoreRefreshToken(ctx context.Context, userID string, token string, expiration time.Duration) error

	// GetRefreshToken obtiene un refresh token almacenado
	// GetRefreshToken(ctx context.Context, userID string) (string, error)

	// DeleteRefreshToken elimina un refresh token (logout)
	// DeleteRefreshToken(ctx context.Context, userID string) error

	// IsTokenRevoked verifica si un token fue revocado
	// IsTokenRevoked(ctx context.Context, tokenID string) (bool, error)

	// RevokeToken marca un token como revocado
	// RevokeToken(ctx context.Context, tokenID string, expiration time.Duration) error
}
