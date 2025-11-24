// Package dto contiene los Data Transfer Objects para autenticación
package dto

// LoginRequest representa la solicitud de login
// Será usado en FASE 2 del Sprint 1
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// TokenPair representa el par de tokens retornado en login/refresh
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // "Bearer"
	ExpiresIn    int64  `json:"expires_in"` // segundos hasta expiración
}

// RefreshRequest representa la solicitud de refresh token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// VerifyRequest representa la solicitud de verificación de token (servicios internos)
type VerifyRequest struct {
	Token string `json:"token" validate:"required"`
}

// VerifyResponse representa la respuesta de verificación de token
type VerifyResponse struct {
	Valid     bool       `json:"valid"`
	UserID    string     `json:"user_id,omitempty"`
	Email     string     `json:"email,omitempty"`
	Roles     []string   `json:"roles,omitempty"`
	ExpiresAt int64      `json:"expires_at,omitempty"`
	Error     string     `json:"error,omitempty"`
}

// TokenClaims representa los claims del access token
type TokenClaims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	TokenType string   `json:"token_type"` // "access" o "refresh"
}

// RefreshTokenClaims representa los claims del refresh token
type RefreshTokenClaims struct {
	UserID    string `json:"user_id"`
	TokenID   string `json:"token_id"` // Para revocación
	TokenType string `json:"token_type"`
}
