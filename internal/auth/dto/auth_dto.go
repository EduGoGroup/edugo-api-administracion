// Package dto contiene los Data Transfer Objects para autenticación
package dto

import "time"

// ===============================================
// REQUEST DTOs
// ===============================================

// VerifyTokenRequest representa el request para verificar un token
type VerifyTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyTokenBulkRequest representa el request para verificar múltiples tokens
type VerifyTokenBulkRequest struct {
	Tokens []string `json:"tokens" binding:"required,min=1,max=100"`
}

// LoginRequest representa el request de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RefreshTokenRequest representa el request para refrescar token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ===============================================
// RESPONSE DTOs
// ===============================================

// VerifyTokenResponse representa la respuesta de verificación de token
type VerifyTokenResponse struct {
	Valid     bool       `json:"valid"`
	UserID    string     `json:"user_id,omitempty"`
	Email     string     `json:"email,omitempty"`
	Role      string     `json:"role,omitempty"`
	SchoolID  string     `json:"school_id,omitempty"` // Escuela principal del usuario
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Error     string     `json:"error,omitempty"`
}

// VerifyTokenBulkResponse representa la respuesta de verificación bulk
type VerifyTokenBulkResponse struct {
	Results map[string]*VerifyTokenResponse `json:"results"`
}

// LoginResponse representa la respuesta de login exitoso
// Compatible con api-mobile (edugo-api-mobile/internal/application/dto/auth_dto.go)
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	User         *UserInfo `json:"user"`
}

// RefreshResponse representa la respuesta de refresh token
// Compatible con api-mobile (solo access_token, expires_in, token_type)
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// UserInfo representa información básica del usuario
// Compatible con api-mobile (mismo contrato JSON)
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	SchoolID  string `json:"school_id,omitempty"` // Escuela principal del usuario
}

// ErrorResponse representa una respuesta de error estándar
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// ===============================================
// INTERNAL DTOs (para comunicación entre servicios)
// ===============================================

// TokenClaims representa los claims extraídos de un JWT
type TokenClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	SchoolID  string    `json:"school_id,omitempty"` // Escuela principal del usuario
	TokenID   string    `json:"jti"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	Issuer    string    `json:"iss"`
}
