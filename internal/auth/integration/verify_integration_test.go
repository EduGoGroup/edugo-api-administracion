// Package integration contiene tests de integración para auth
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/handler"
	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/middleware"
	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/shared/crypto"
)

// mockTokenCache para tests de integración
type mockTokenCache struct {
	mu        sync.RWMutex
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
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.cache[key]
	return val, ok
}

func (m *mockTokenCache) Set(_ context.Context, key string, value *dto.VerifyTokenResponse, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache[key] = value
	return nil
}

func (m *mockTokenCache) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cache, key)
	return nil
}

func (m *mockTokenCache) IsBlacklisted(_ context.Context, tokenID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.blacklist[tokenID]
}

func (m *mockTokenCache) Blacklist(_ context.Context, tokenID string, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blacklist[tokenID] = true
	return nil
}

// setupIntegrationServer crea un servidor completo para tests de integración
func setupIntegrationServer(t *testing.T) (*gin.Engine, *crypto.JWTManager, *middleware.RateLimiter) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// Crear JWTManager
	jwtManager, err := crypto.NewJWTManager(crypto.JWTConfig{
		Secret:               "integration-test-secret-at-least-32-chars",
		Issuer:               "edugo-central",
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
	})
	require.NoError(t, err)

	// Crear TokenService con cache mock
	cache := newMockTokenCache()
	tokenService := service.NewTokenService(jwtManager, cache, service.TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	})

	// Crear RateLimiter
	rateLimiter := middleware.NewRateLimiter(middleware.RateLimitConfig{
		InternalMaxRequests: 1000,
		InternalWindow:      time.Minute,
		ExternalMaxRequests: 60,
		ExternalWindow:      time.Minute,
		InternalAPIKeys:     []string{"internal-api-key"},
		InternalIPs:         []string{"10.0.0.0/8"},
	})

	// Crear VerifyHandler
	verifyHandler := handler.NewVerifyHandler(
		tokenService,
		[]string{"10.0.0.0/8"},
		map[string]string{
			"api-mobile": "internal-api-key",
		},
	)

	// Configurar router
	router := gin.New()
	router.Use(gin.Recovery())

	// Grupo con rate limiting
	v1 := router.Group("/v1")
	v1.Use(rateLimiter.Middleware())
	verifyHandler.RegisterRoutes(v1)

	return router, jwtManager, rateLimiter
}

func TestIntegration_VerifyToken_FullFlow(t *testing.T) {
	// Arrange
	router, jwtManager, rateLimiter := setupIntegrationServer(t)
	defer rateLimiter.Stop()

	// Generar token
	token, _, err := jwtManager.GenerateAccessToken("user-integration", "integration@test.com", "admin")
	require.NoError(t, err)

	// Act
	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: token})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "203.0.113.1:12345" // IP externa
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.VerifyTokenResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Valid)
	assert.Equal(t, "user-integration", response.UserID)
	assert.Equal(t, "integration@test.com", response.Email)
	assert.Equal(t, "admin", response.Role)

	// Verificar headers de rate limit
	assert.Equal(t, "60", rec.Header().Get("X-RateLimit-Limit"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Reset"))
	assert.NotEmpty(t, rec.Header().Get("X-Response-Time"))
}

func TestIntegration_VerifyToken_InternalService(t *testing.T) {
	// Arrange
	router, jwtManager, rateLimiter := setupIntegrationServer(t)
	defer rateLimiter.Stop()

	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@test.com", "user")
	require.NoError(t, err)

	// Act - Request con API Key interna
	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: token})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-API-Key", "internal-api-key")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	// Debe usar límite interno (1000)
	assert.Equal(t, "1000", rec.Header().Get("X-RateLimit-Limit"))
}

func TestIntegration_VerifyToken_RateLimitEnforced(t *testing.T) {
	// Arrange - Rate limit bajo para probar
	gin.SetMode(gin.TestMode)

	jwtManager, _ := crypto.NewJWTManager(crypto.JWTConfig{
		Secret:               "integration-test-secret-at-least-32-chars",
		Issuer:               "edugo-central",
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
	})

	cache := newMockTokenCache()
	tokenService := service.NewTokenService(jwtManager, cache, service.TokenServiceConfig{})

	rateLimiter := middleware.NewRateLimiter(middleware.RateLimitConfig{
		ExternalMaxRequests: 3, // Límite muy bajo
		ExternalWindow:      time.Minute,
	})
	defer rateLimiter.Stop()

	verifyHandler := handler.NewVerifyHandler(tokenService, nil, nil)

	router := gin.New()
	v1 := router.Group("/v1")
	v1.Use(rateLimiter.Middleware())
	verifyHandler.RegisterRoutes(v1)

	token, _, _ := jwtManager.GenerateAccessToken("user-123", "test@test.com", "user")

	// Act - Hacer 4 requests (límite es 3)
	for i := 0; i < 4; i++ {
		body, _ := json.Marshal(dto.VerifyTokenRequest{Token: token})
		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "203.0.113.1:12345"
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if i < 3 {
			assert.Equal(t, http.StatusOK, rec.Code, "Request %d should pass", i+1)
		} else {
			// El 4to debe ser rate limited
			assert.Equal(t, http.StatusTooManyRequests, rec.Code)

			var errorResponse dto.ErrorResponse
			_ = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
			assert.Equal(t, "RATE_LIMIT", errorResponse.Code)
		}
	}
}

func TestIntegration_VerifyTokenBulk_FullFlow(t *testing.T) {
	// Arrange
	router, jwtManager, rateLimiter := setupIntegrationServer(t)
	defer rateLimiter.Stop()

	// Generar múltiples tokens
	token1, _, _ := jwtManager.GenerateAccessToken("user-1", "user1@test.com", "admin")
	token2, _, _ := jwtManager.GenerateAccessToken("user-2", "user2@test.com", "user")

	// Act
	body, _ := json.Marshal(dto.VerifyTokenBulkRequest{
		Tokens: []string{token1, token2, "invalid-token"},
	})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify-bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-API-Key", "internal-api-key") // Requerido para bulk
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.VerifyTokenBulkResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Results, 3)

	// Contar válidos e inválidos
	validCount := 0
	for _, r := range response.Results {
		if r.Valid {
			validCount++
		}
	}
	assert.Equal(t, 2, validCount)
}

func TestIntegration_VerifyTokenBulk_RequiresAPIKey(t *testing.T) {
	// Arrange
	router, jwtManager, rateLimiter := setupIntegrationServer(t)
	defer rateLimiter.Stop()

	token, _, _ := jwtManager.GenerateAccessToken("user-1", "user1@test.com", "admin")

	// Act - Sin API Key
	body, _ := json.Marshal(dto.VerifyTokenBulkRequest{
		Tokens: []string{token},
	})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify-bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "203.0.113.1:12345" // IP externa
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var errorResponse dto.ErrorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.Equal(t, "API_KEY_REQUIRED", errorResponse.Code)
}

func TestIntegration_WrongIssuer(t *testing.T) {
	// Arrange - Dos JWTManagers con diferentes issuers
	gin.SetMode(gin.TestMode)

	// Manager que genera tokens
	wrongIssuerManager, _ := crypto.NewJWTManager(crypto.JWTConfig{
		Secret:              "integration-test-secret-at-least-32-chars",
		Issuer:              "wrong-issuer",
		AccessTokenDuration: 15 * time.Minute,
	})

	// Manager que valida (servidor)
	correctIssuerManager, _ := crypto.NewJWTManager(crypto.JWTConfig{
		Secret:              "integration-test-secret-at-least-32-chars",
		Issuer:              "edugo-central",
		AccessTokenDuration: 15 * time.Minute,
	})

	cache := newMockTokenCache()
	tokenService := service.NewTokenService(correctIssuerManager, cache, service.TokenServiceConfig{})
	verifyHandler := handler.NewVerifyHandler(tokenService, nil, nil)

	router := gin.New()
	v1 := router.Group("/v1")
	verifyHandler.RegisterRoutes(v1)

	// Generar token con issuer incorrecto
	token, _, _ := wrongIssuerManager.GenerateAccessToken("user-123", "test@test.com", "user")

	// Act
	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: token})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.VerifyTokenResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.False(t, response.Valid)
	assert.Contains(t, response.Error, "issuer")
}

func TestIntegration_ConcurrentRequests(t *testing.T) {
	// Arrange
	router, jwtManager, rateLimiter := setupIntegrationServer(t)
	defer rateLimiter.Stop()

	token, _, _ := jwtManager.GenerateAccessToken("user-concurrent", "concurrent@test.com", "user")

	// Act - Hacer múltiples requests concurrentes
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			body, _ := json.Marshal(dto.VerifyTokenRequest{Token: token})
			req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Service-API-Key", "internal-api-key")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			// Todos deben pasar (límite interno alto)
			assert.Equal(t, http.StatusOK, rec.Code)
			done <- true
		}()
	}

	// Esperar a que terminen todos
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestIntegration_InvalidToken_Response(t *testing.T) {
	// Arrange
	router, _, rateLimiter := setupIntegrationServer(t)
	defer rateLimiter.Stop()

	// Act - Token inválido
	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: "completely-invalid-token"})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code) // Siempre 200

	var response dto.VerifyTokenResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response.Valid)
	assert.NotEmpty(t, response.Error)
	assert.Empty(t, response.UserID)
}

func TestIntegration_TokenCaching(t *testing.T) {
	// Arrange
	router, jwtManager, rateLimiter := setupIntegrationServer(t)
	defer rateLimiter.Stop()

	token, _, _ := jwtManager.GenerateAccessToken("user-cache", "cache@test.com", "user")

	// Act - Primera llamada
	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: token})
	req1, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-Service-API-Key", "internal-api-key")
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)

	// Segunda llamada (debería venir de cache)
	body, _ = json.Marshal(dto.VerifyTokenRequest{Token: token})
	req2, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Service-API-Key", "internal-api-key")
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	// Assert - Ambas deben tener el mismo resultado
	var response1, response2 dto.VerifyTokenResponse
	_ = json.Unmarshal(rec1.Body.Bytes(), &response1)
	_ = json.Unmarshal(rec2.Body.Bytes(), &response2)

	assert.True(t, response1.Valid)
	assert.True(t, response2.Valid)
	assert.Equal(t, response1.UserID, response2.UserID)
}
