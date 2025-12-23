package service

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/shared/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTokenCache implementa TokenCache para tests
type mockTokenCache struct {
	cache     map[string]*dto.VerifyTokenResponse
	blacklist map[string]bool
}

func newMockTokenCache() *mockTokenCache {
	return &mockTokenCache{
		cache:     make(map[string]*dto.VerifyTokenResponse),
		blacklist: make(map[string]bool),
	}
}

func (m *mockTokenCache) Get(_ context.Context, key string) (*dto.VerifyTokenResponse, bool) {
	val, ok := m.cache[key]
	return val, ok
}

func (m *mockTokenCache) Set(_ context.Context, key string, value *dto.VerifyTokenResponse, _ time.Duration) error {
	m.cache[key] = value
	return nil
}

func (m *mockTokenCache) Delete(_ context.Context, key string) error {
	delete(m.cache, key)
	return nil
}

func (m *mockTokenCache) IsBlacklisted(_ context.Context, tokenID string) bool {
	return m.blacklist[tokenID]
}

func (m *mockTokenCache) Blacklist(_ context.Context, tokenID string, _ time.Duration) error {
	m.blacklist[tokenID] = true
	return nil
}

// createTestJWTManager crea un JWTManager para tests
func createTestJWTManager(t *testing.T) *crypto.JWTManager {
	t.Helper()
	manager, err := crypto.NewJWTManager(crypto.JWTConfig{
		Secret:              "test-secret-key-at-least-32-characters-long",
		Issuer:              "edugo-central",
		AccessTokenDuration: 15 * time.Minute,
		RefreshTokenDuration:     7 * 24 * time.Hour,
	})
	require.NoError(t, err)
	return manager
}

func TestTokenService_VerifyToken_Valid(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	}

	service := NewTokenService(jwtManager, cache, config)

	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "")
	require.NoError(t, err)

	// Act
	result, err := service.VerifyToken(context.Background(), token)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, "user-123", result.UserID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "admin", result.Role)
	assert.NotNil(t, result.ExpiresAt)
	assert.Empty(t, result.Error)
}

func TestTokenService_VerifyToken_Invalid(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   false,
		BlacklistCheck: false,
	}

	service := NewTokenService(jwtManager, cache, config)

	// Act
	result, err := service.VerifyToken(context.Background(), "invalid-token")

	// Assert
	require.NoError(t, err) // El servicio no retorna error, solo marca como inválido
	assert.False(t, result.Valid)
	assert.Empty(t, result.UserID)
	assert.NotEmpty(t, result.Error)
}

func TestTokenService_VerifyToken_Blacklisted(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	}

	service := NewTokenService(jwtManager, cache, config)

	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "")
	require.NoError(t, err)

	// Agregar el token al blacklist
	tokenID, err := jwtManager.GetTokenID(token)
	require.NoError(t, err)
	err = cache.Blacklist(context.Background(), tokenID, time.Hour)
	require.NoError(t, err)

	// Act
	result, err := service.VerifyToken(context.Background(), token)

	// Assert
	require.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Contains(t, result.Error, "revocado")
}

func TestTokenService_VerifyToken_CachedResult(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	}

	service := NewTokenService(jwtManager, cache, config)

	token, _, err := jwtManager.GenerateAccessToken("user-456", "cached@example.com", "user", "")
	require.NoError(t, err)

	// Primera verificación - debería cachear
	result1, err := service.VerifyToken(context.Background(), token)
	require.NoError(t, err)
	assert.True(t, result1.Valid)

	// Segunda verificación - debería venir de caché
	result2, err := service.VerifyToken(context.Background(), token)
	require.NoError(t, err)
	assert.True(t, result2.Valid)
	assert.Equal(t, result1.UserID, result2.UserID)
}

func TestTokenService_VerifyTokenBulk(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	}

	service := NewTokenService(jwtManager, cache, config)

	token1, _, err := jwtManager.GenerateAccessToken("user-1", "user1@example.com", "admin", "")
	require.NoError(t, err)

	token2, _, err := jwtManager.GenerateAccessToken("user-2", "user2@example.com", "user", "")
	require.NoError(t, err)

	tokens := []string{token1, token2, "invalid-token"}

	// Act
	result, err := service.VerifyTokenBulk(context.Background(), tokens)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result.Results, 3)

	// Contar válidos e inválidos
	validCount := 0
	invalidCount := 0
	for _, r := range result.Results {
		if r.Valid {
			validCount++
		} else {
			invalidCount++
		}
	}
	assert.Equal(t, 2, validCount)
	assert.Equal(t, 1, invalidCount)
}

func TestTokenService_VerifyTokenBulk_Empty(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   false,
		BlacklistCheck: false,
	}

	service := NewTokenService(jwtManager, cache, config)

	// Act
	result, err := service.VerifyTokenBulk(context.Background(), []string{})

	// Assert
	require.NoError(t, err)
	assert.Len(t, result.Results, 0)
}

func TestTokenService_RevokeToken(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	}

	service := NewTokenService(jwtManager, cache, config)

	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "")
	require.NoError(t, err)

	// Verificar que el token es válido antes de revocar
	result1, err := service.VerifyToken(context.Background(), token)
	require.NoError(t, err)
	assert.True(t, result1.Valid)

	// Act - Revocar el token
	err = service.RevokeToken(context.Background(), token)
	require.NoError(t, err)

	// Assert - Verificar que ahora está revocado
	result2, err := service.VerifyToken(context.Background(), token)
	require.NoError(t, err)
	assert.False(t, result2.Valid)
	assert.Contains(t, result2.Error, "revocado")
}

func TestTokenService_RevokeToken_Invalid(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	}

	service := NewTokenService(jwtManager, cache, config)

	// Act - Intentar revocar token inválido
	err := service.RevokeToken(context.Background(), "invalid-token")

	// Assert
	require.Error(t, err)
}

func TestTokenService_GenerateTokenPair(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   false,
		BlacklistCheck: false,
	}

	service := NewTokenService(jwtManager, cache, config)

	// Act
	result, err := service.GenerateTokenPair("user-123", "test@example.com", "admin", "")

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, "Bearer", result.TokenType)
	assert.Greater(t, result.ExpiresIn, int64(0))
}

func TestNewTokenService_DefaultTTL(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	cache := newMockTokenCache()
	config := TokenServiceConfig{
		CacheTTL:       0, // Sin TTL configurado
		CacheEnabled:   true,
		BlacklistCheck: true,
	}

	// Act
	service := NewTokenService(jwtManager, cache, config)

	// Assert - El servicio debe usar TTL por defecto (60s)
	assert.NotNil(t, service)

	// Verificar que funciona correctamente
	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "")
	require.NoError(t, err)

	result, err := service.VerifyToken(context.Background(), token)
	require.NoError(t, err)
	assert.True(t, result.Valid)
}

func TestTokenService_VerifyToken_NilCache(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	}

	service := NewTokenService(jwtManager, nil, config)

	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "")
	require.NoError(t, err)

	// Act
	result, err := service.VerifyToken(context.Background(), token)

	// Assert - Debe funcionar sin cache
	require.NoError(t, err)
	assert.True(t, result.Valid)
}

func TestTokenService_RevokeToken_NilCache(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	config := TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	}

	service := NewTokenService(jwtManager, nil, config)

	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "")
	require.NoError(t, err)

	// Act - Sin cache, revocar no hace nada pero no debe fallar
	err = service.RevokeToken(context.Background(), token)

	// Assert
	require.NoError(t, err)
}

func TestTokenService_HashToken(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	config := TokenServiceConfig{}
	service := NewTokenService(jwtManager, nil, config)

	// Act
	hash1 := service.hashToken("token1")
	hash2 := service.hashToken("token1")
	hash3 := service.hashToken("token2")

	// Assert
	assert.Equal(t, hash1, hash2)       // Mismo token = mismo hash
	assert.NotEqual(t, hash1, hash3)    // Diferente token = diferente hash
	assert.Contains(t, hash1, "auth:token:") // Prefijo correcto
}

func TestTokenService_TruncateToken(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	config := TokenServiceConfig{}
	service := NewTokenService(jwtManager, nil, config)

	// Act
	short := service.truncateToken("short")
	long := service.truncateToken("this-is-a-very-long-token-string-here")

	// Assert
	assert.Equal(t, "short", short)              // Token corto no se trunca
	assert.Contains(t, long, "...")              // Token largo se trunca
	assert.LessOrEqual(t, len(long), 23)         // 10 + ... + 10
}

func TestTokenService_CalculateCacheTTL(t *testing.T) {
	// Arrange
	jwtManager := createTestJWTManager(t)
	config := TokenServiceConfig{
		CacheTTL: 60 * time.Second,
	}
	service := NewTokenService(jwtManager, nil, config)

	// Act & Assert
	// Token que expira en 2 minutos - TTL debe ser el configurado (60s)
	ttl1 := service.calculateCacheTTL(time.Now().Add(2 * time.Minute))
	assert.Equal(t, 60*time.Second, ttl1)

	// Token que expira en 30 segundos - TTL debe ser 30s (menor que configurado)
	ttl2 := service.calculateCacheTTL(time.Now().Add(30 * time.Second))
	assert.LessOrEqual(t, ttl2, 30*time.Second)
}
