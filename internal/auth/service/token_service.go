// Package service contiene la lógica de negocio de autenticación
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/shared/crypto"
)

// Errores del servicio
var (
	ErrTokenNotFound    = errors.New("token no encontrado")
	ErrTokenBlacklisted = errors.New("token en blacklist")
	ErrCacheUnavailable = errors.New("cache no disponible")
)

// TokenCache define la interfaz para cache de tokens
type TokenCache interface {
	Get(ctx context.Context, key string) (*dto.VerifyTokenResponse, bool)
	Set(ctx context.Context, key string, value *dto.VerifyTokenResponse, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	IsBlacklisted(ctx context.Context, tokenID string) bool
	Blacklist(ctx context.Context, tokenID string, ttl time.Duration) error
}

// TokenServiceConfig configuración del servicio
type TokenServiceConfig struct {
	CacheTTL       time.Duration
	CacheEnabled   bool
	BlacklistCheck bool
}

// TokenService gestiona operaciones de tokens
type TokenService struct {
	jwtManager *crypto.JWTManager
	cache      TokenCache
	config     TokenServiceConfig
}

// NewTokenService crea una nueva instancia
func NewTokenService(
	jwtManager *crypto.JWTManager,
	cache TokenCache,
	config TokenServiceConfig,
) *TokenService {
	if config.CacheTTL == 0 {
		config.CacheTTL = 60 * time.Second
	}

	return &TokenService{
		jwtManager: jwtManager,
		cache:      cache,
		config:     config,
	}
}

// VerifyToken verifica un token y retorna información del usuario
func (s *TokenService) VerifyToken(ctx context.Context, token string) (*dto.VerifyTokenResponse, error) {
	// 1. Generar hash del token para cache key
	cacheKey := s.hashToken(token)

	// 2. Verificar cache
	if s.config.CacheEnabled && s.cache != nil {
		if cached, found := s.cache.Get(ctx, cacheKey); found {
			return cached, nil
		}
	}

	// 3. Validar JWT
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		response := &dto.VerifyTokenResponse{
			Valid: false,
			Error: err.Error(),
		}
		return response, nil // No retornar error, retornar response con valid=false
	}

	// 4. Verificar blacklist
	if s.config.BlacklistCheck && s.cache != nil {
		if s.cache.IsBlacklisted(ctx, claims.ID) {
			response := &dto.VerifyTokenResponse{
				Valid: false,
				Error: "token revocado",
			}
			return response, nil
		}
	}

	// 5. Construir response
	expiresAt := claims.ExpiresAt.Time
	response := &dto.VerifyTokenResponse{
		Valid:     true,
		UserID:    claims.UserID,
		Email:     claims.Email,
		Role:      claims.Role,
		ExpiresAt: &expiresAt,
	}

	// 6. Guardar en cache
	if s.config.CacheEnabled && s.cache != nil {
		// TTL del cache debe ser menor que el tiempo restante del token
		ttl := s.calculateCacheTTL(expiresAt)
		if ttl > 0 {
			_ = s.cache.Set(ctx, cacheKey, response, ttl)
		}
	}

	return response, nil
}

// VerifyTokenBulk verifica múltiples tokens
func (s *TokenService) VerifyTokenBulk(ctx context.Context, tokens []string) (*dto.VerifyTokenBulkResponse, error) {
	results := make(map[string]*dto.VerifyTokenResponse, len(tokens))

	for _, token := range tokens {
		response, err := s.VerifyToken(ctx, token)
		if err != nil {
			results[s.truncateToken(token)] = &dto.VerifyTokenResponse{
				Valid: false,
				Error: err.Error(),
			}
			continue
		}
		results[s.truncateToken(token)] = response
	}

	return &dto.VerifyTokenBulkResponse{Results: results}, nil
}

// RevokeToken agrega un token a la blacklist
func (s *TokenService) RevokeToken(ctx context.Context, token string) error {
	// Extraer token ID
	tokenID, err := s.jwtManager.GetTokenID(token)
	if err != nil {
		return fmt.Errorf("error extrayendo token ID: %w", err)
	}

	// Obtener tiempo de expiración para calcular TTL del blacklist
	expiresAt, err := s.jwtManager.GetExpirationTime(token)
	if err != nil {
		// Si no podemos obtener expiración, usar TTL por defecto
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token ya expirado, no necesita blacklist
		return nil
	}

	// Agregar a blacklist
	if s.cache != nil {
		if err := s.cache.Blacklist(ctx, tokenID, ttl); err != nil {
			return fmt.Errorf("error agregando a blacklist: %w", err)
		}
	}

	// Invalidar cache
	cacheKey := s.hashToken(token)
	if s.cache != nil {
		_ = s.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// GenerateTokenPair genera un par de tokens (access + refresh) para login
func (s *TokenService) GenerateTokenPair(userID, email, role string) (*dto.LoginResponse, error) {
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(userID, email, role)
	if err != nil {
		return nil, fmt.Errorf("error generando access token: %w", err)
	}

	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("error generando refresh token: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Until(expiresAt).Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// GenerateAccessToken genera solo un nuevo access token (para refresh)
func (s *TokenService) GenerateAccessToken(userID, email, role string) (*dto.RefreshResponse, error) {
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(userID, email, role)
	if err != nil {
		return nil, fmt.Errorf("error generando access token: %w", err)
	}

	return &dto.RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(time.Until(expiresAt).Seconds()),
		TokenType:   "Bearer",
	}, nil
}

// Helper functions

func (s *TokenService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return "auth:token:" + hex.EncodeToString(hash[:])
}

func (s *TokenService) truncateToken(token string) string {
	if len(token) > 20 {
		return token[:10] + "..." + token[len(token)-10:]
	}
	return token
}

func (s *TokenService) calculateCacheTTL(expiresAt time.Time) time.Duration {
	timeRemaining := time.Until(expiresAt)

	// Si queda menos tiempo que el TTL configurado, usar el tiempo restante
	if timeRemaining < s.config.CacheTTL {
		return timeRemaining
	}

	return s.config.CacheTTL
}
