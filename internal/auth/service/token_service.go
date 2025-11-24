// Package service contiene la lógica de negocio de autenticación
package service

// TokenService gestiona la creación y validación de JWT
// Responsabilidades:
// - Generar access tokens y refresh tokens
// - Validar y parsear tokens
// - Verificar tokens para servicios internos
// Será implementado en FASE 2 del Sprint 1
type TokenService interface {
	// GenerateTokenPair genera un par de access/refresh tokens
	// GenerateTokenPair(ctx context.Context, user *entity.User) (*dto.TokenPair, error)

	// ValidateAccessToken valida un access token y retorna los claims
	// ValidateAccessToken(ctx context.Context, token string) (*dto.TokenClaims, error)

	// ValidateRefreshToken valida un refresh token
	// ValidateRefreshToken(ctx context.Context, token string) (*dto.RefreshTokenClaims, error)

	// VerifyTokenForService valida un token para servicios internos (endpoint /v1/auth/verify)
	// VerifyTokenForService(ctx context.Context, token string) (*dto.VerifyResponse, error)
}
