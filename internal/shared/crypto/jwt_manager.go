// Package crypto proporciona utilidades criptográficas
package crypto

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Errores de JWT
var (
	ErrInvalidToken     = errors.New("token inválido")
	ErrTokenExpired     = errors.New("token expirado")
	ErrInvalidIssuer    = errors.New("issuer inválido")
	ErrInvalidSignature = errors.New("firma inválida")
	ErrTokenRevoked     = errors.New("token revocado")
	ErrMalformedToken   = errors.New("token malformado")
)

// JWTConfig contiene la configuración para JWT
type JWTConfig struct {
	Secret               string
	Issuer               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// Claims representa los claims personalizados del JWT
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	SchoolID string `json:"school_id,omitempty"` // Escuela principal del usuario (vacío para super_admin)
	jwt.RegisteredClaims
}

// JWTManager gestiona operaciones JWT
type JWTManager struct {
	config JWTConfig
}

// NewJWTManager crea una nueva instancia de JWTManager
func NewJWTManager(config JWTConfig) (*JWTManager, error) {
	if len(config.Secret) < 32 {
		return nil, fmt.Errorf("JWT secret debe tener al menos 32 caracteres, tiene %d", len(config.Secret))
	}
	if config.Issuer == "" {
		return nil, errors.New("JWT issuer es requerido")
	}
	if config.AccessTokenDuration == 0 {
		config.AccessTokenDuration = 15 * time.Minute
	}
	if config.RefreshTokenDuration == 0 {
		config.RefreshTokenDuration = 7 * 24 * time.Hour
	}

	return &JWTManager{config: config}, nil
}

// GenerateAccessToken genera un nuevo access token
func (m *JWTManager) GenerateAccessToken(userID, email, role, schoolID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.config.AccessTokenDuration)

	claims := Claims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		SchoolID: schoolID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    m.config.Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.config.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error firmando token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken genera un nuevo refresh token
func (m *JWTManager) GenerateRefreshToken(userID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.config.RefreshTokenDuration)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    m.config.Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.config.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error firmando refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateToken valida un token y retorna los claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar algoritmo
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(m.config.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrMalformedToken
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrInvalidSignature
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Verificar issuer
	if claims.Issuer != m.config.Issuer {
		return nil, fmt.Errorf("%w: esperado '%s', recibido '%s'",
			ErrInvalidIssuer, m.config.Issuer, claims.Issuer)
	}

	return claims, nil
}

// GetTokenID extrae el JTI (token ID) de un token sin validar completamente
// Útil para operaciones de blacklist
func (m *JWTManager) GetTokenID(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return "", fmt.Errorf("error parseando token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", ErrInvalidToken
	}

	return claims.ID, nil
}

// GetExpirationTime retorna el tiempo de expiración de un token
func (m *JWTManager) GetExpirationTime(tokenString string) (time.Time, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return time.Time{}, err
	}

	if claims != nil && claims.ExpiresAt != nil {
		return claims.ExpiresAt.Time, nil
	}

	return time.Time{}, ErrInvalidToken
}

// GetConfig retorna la configuración del JWTManager (solo lectura)
func (m *JWTManager) GetConfig() JWTConfig {
	return m.config
}
