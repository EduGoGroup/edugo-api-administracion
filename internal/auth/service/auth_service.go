// Package service contiene la lógica de negocio de autenticación
package service

// AuthService contiene la lógica de negocio para autenticación
// Responsabilidades:
// - Validar credenciales de usuario
// - Coordinar creación de tokens
// - Manejar flujo de login/logout/refresh
// Será implementado en FASE 2 del Sprint 1
type AuthService interface {
	// Login valida credenciales y retorna tokens
	// Login(ctx context.Context, email, password string) (*dto.TokenPair, error)

	// Logout invalida los tokens del usuario
	// Logout(ctx context.Context, userID string) error

	// RefreshToken genera nuevos tokens usando el refresh token
	// RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenPair, error)
}
