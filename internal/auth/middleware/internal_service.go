// Package middleware contiene middlewares para autenticación
package middleware

// InternalServiceMiddleware valida que las requests vengan de servicios internos autorizados
// Verifica API Key y rango IP
// Será implementado en FASE 2 del Sprint 1
type InternalServiceMiddleware struct {
	// config *config.InternalServicesConfig
}

// NewInternalServiceMiddleware crea una nueva instancia
// func NewInternalServiceMiddleware(config *config.InternalServicesConfig) *InternalServiceMiddleware {
// 	return &InternalServiceMiddleware{config: config}
// }
