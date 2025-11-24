// Package middleware contiene middlewares para autenticación
package middleware

// AuthMiddleware valida tokens JWT en requests entrantes
// Extrae el token del header Authorization y valida con TokenService
// Será implementado en FASE 2 del Sprint 1
type AuthMiddleware struct {
	// tokenService service.TokenService
}

// NewAuthMiddleware crea una nueva instancia de AuthMiddleware
// func NewAuthMiddleware(tokenService service.TokenService) *AuthMiddleware {
// 	return &AuthMiddleware{tokenService: tokenService}
// }
