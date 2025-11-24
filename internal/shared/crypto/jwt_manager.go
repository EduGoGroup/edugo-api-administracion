// Package crypto proporciona utilidades criptográficas
package crypto

// JWTManager centraliza operaciones de JWT
// Usa la configuración centralizada (issuer: edugo-central)
// Será implementado en FASE 2 del Sprint 1
type JWTManager interface {
	// GenerateAccessToken genera un access token
	// GenerateAccessToken(claims map[string]interface{}) (string, error)

	// GenerateRefreshToken genera un refresh token
	// GenerateRefreshToken(userID string) (string, error)

	// ParseToken parsea y valida un token
	// ParseToken(tokenString string) (*jwt.Token, error)

	// GetClaims extrae los claims de un token validado
	// GetClaims(token *jwt.Token) (map[string]interface{}, error)
}
