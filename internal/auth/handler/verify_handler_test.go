package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/service"
	"github.com/EduGoGroup/edugo-api-administracion/internal/shared/crypto"
)

// mockTokenCache para tests
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

// setupTestHandler crea un handler para tests
func setupTestHandler(t *testing.T) (*VerifyHandler, *crypto.JWTManager) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	jwtManager, err := crypto.NewJWTManager(crypto.JWTConfig{
		Secret:               "test-secret-key-at-least-32-characters-long",
		Issuer:               "edugo-central",
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
	})
	require.NoError(t, err)

	cache := newMockTokenCache()
	tokenService := service.NewTokenService(jwtManager, cache, service.TokenServiceConfig{
		CacheTTL:       60 * time.Second,
		CacheEnabled:   true,
		BlacklistCheck: true,
	})

	handler := NewVerifyHandler(
		tokenService,
		[]string{"10.0.0.0/8", "192.168.1.0/24"},
		map[string]string{
			"api-mobile": "test-api-key-mobile",
			"worker":     "test-api-key-worker",
		},
	)

	return handler, jwtManager
}

func TestVerifyHandler_VerifyToken_Success(t *testing.T) {
	// Arrange
	handler, jwtManager := setupTestHandler(t)
	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "")
	require.NoError(t, err)

	router := gin.New()
	router.POST("/v1/auth/verify", handler.VerifyToken)

	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: token})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.VerifyTokenResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Valid)
	assert.Equal(t, "user-123", response.UserID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "admin", response.Role)
}

func TestVerifyHandler_VerifyToken_WithBearerPrefix(t *testing.T) {
	// Arrange
	handler, jwtManager := setupTestHandler(t)
	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "")
	require.NoError(t, err)

	router := gin.New()
	router.POST("/v1/auth/verify", handler.VerifyToken)

	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: "Bearer " + token})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.VerifyTokenResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Valid)
}

func TestVerifyHandler_VerifyToken_InvalidToken(t *testing.T) {
	// Arrange
	handler, _ := setupTestHandler(t)

	router := gin.New()
	router.POST("/v1/auth/verify", handler.VerifyToken)

	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: "invalid-token"})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code) // Siempre 200, valid=false

	var response dto.VerifyTokenResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Valid)
	assert.NotEmpty(t, response.Error)
}

func TestVerifyHandler_VerifyToken_EmptyToken(t *testing.T) {
	// Arrange
	handler, _ := setupTestHandler(t)

	router := gin.New()
	router.POST("/v1/auth/verify", handler.VerifyToken)

	// Token vacío falla binding validation (required)
	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: ""})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert - Gin binding valida "required" y retorna INVALID_REQUEST
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "INVALID_REQUEST", response.Code)
}

func TestVerifyHandler_VerifyToken_OnlySpaces(t *testing.T) {
	// Arrange
	handler, _ := setupTestHandler(t)

	router := gin.New()
	router.POST("/v1/auth/verify", handler.VerifyToken)

	// Token con solo espacios pasa binding pero falla después de trim
	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: "   "})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "EMPTY_TOKEN", response.Code)
}

func TestVerifyHandler_VerifyToken_MissingBody(t *testing.T) {
	// Arrange
	handler, _ := setupTestHandler(t)

	router := gin.New()
	router.POST("/v1/auth/verify", handler.VerifyToken)

	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyHandler_VerifyToken_ResponseTimeHeader(t *testing.T) {
	// Arrange
	handler, jwtManager := setupTestHandler(t)
	token, _, err := jwtManager.GenerateAccessToken("user-123", "test@example.com", "admin", "")
	require.NoError(t, err)

	router := gin.New()
	router.POST("/v1/auth/verify", handler.VerifyToken)

	body, _ := json.Marshal(dto.VerifyTokenRequest{Token: token})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.NotEmpty(t, rec.Header().Get("X-Response-Time"))
}

func TestVerifyHandler_VerifyTokenBulk_Success(t *testing.T) {
	// Arrange
	handler, jwtManager := setupTestHandler(t)
	token1, _, _ := jwtManager.GenerateAccessToken("user-1", "user1@example.com", "admin", "")
	token2, _, _ := jwtManager.GenerateAccessToken("user-2", "user2@example.com", "user", "")

	router := gin.New()
	router.POST("/v1/auth/verify-bulk", handler.VerifyTokenBulk)

	body, _ := json.Marshal(dto.VerifyTokenBulkRequest{Tokens: []string{token1, token2}})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify-bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-API-Key", "test-api-key-mobile")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.VerifyTokenBulkResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response.Results, 2)
}

func TestVerifyHandler_VerifyTokenBulk_Unauthorized(t *testing.T) {
	// Arrange
	handler, jwtManager := setupTestHandler(t)
	token, _, _ := jwtManager.GenerateAccessToken("user-1", "user1@example.com", "admin", "")

	router := gin.New()
	router.POST("/v1/auth/verify-bulk", handler.VerifyTokenBulk)

	body, _ := json.Marshal(dto.VerifyTokenBulkRequest{Tokens: []string{token}})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify-bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// No API Key
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "API_KEY_REQUIRED", response.Code)
}

func TestVerifyHandler_VerifyTokenBulk_EmptyTokens(t *testing.T) {
	// Arrange
	handler, _ := setupTestHandler(t)

	router := gin.New()
	router.POST("/v1/auth/verify-bulk", handler.VerifyTokenBulk)

	// Tokens vacío falla binding validation (min=1)
	body, _ := json.Marshal(dto.VerifyTokenBulkRequest{Tokens: []string{}})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify-bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-API-Key", "test-api-key-mobile")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert - Gin binding valida min=1 primero
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "INVALID_REQUEST", response.Code)
}

func TestVerifyHandler_VerifyTokenBulk_TooManyTokens(t *testing.T) {
	// Arrange
	handler, _ := setupTestHandler(t)

	router := gin.New()
	router.POST("/v1/auth/verify-bulk", handler.VerifyTokenBulk)

	// Crear 101 tokens - excede max=100 del binding
	tokens := make([]string, 101)
	for i := range tokens {
		tokens[i] = "token"
	}

	body, _ := json.Marshal(dto.VerifyTokenBulkRequest{Tokens: tokens})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify-bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-API-Key", "test-api-key-mobile")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert - Gin binding valida max=100 primero
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "INVALID_REQUEST", response.Code)
}

func TestVerifyHandler_IsInternalService_APIKey(t *testing.T) {
	// Arrange
	handler, _ := setupTestHandler(t)

	testCases := []struct {
		name     string
		apiKey   string
		expected bool
	}{
		{"valid_mobile_key", "test-api-key-mobile", true},
		{"valid_worker_key", "test-api-key-worker", true},
		{"invalid_key", "invalid-key", false},
		{"empty_key", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				result := handler.IsInternalService(c)
				c.JSON(http.StatusOK, gin.H{"internal": result})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tc.apiKey != "" {
				req.Header.Set("X-Service-API-Key", tc.apiKey)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			var response map[string]bool
			_ = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expected, response["internal"])
		})
	}
}

func TestVerifyHandler_RegisterRoutes(t *testing.T) {
	// Arrange
	handler, _ := setupTestHandler(t)

	router := gin.New()
	group := router.Group("/v1")
	handler.RegisterRoutes(group)

	// Act & Assert - Verificar que las rutas están registradas
	routes := router.Routes()

	var foundVerify, foundBulk bool
	for _, route := range routes {
		if route.Path == "/v1/auth/verify" && route.Method == "POST" {
			foundVerify = true
		}
		if route.Path == "/v1/auth/verify-bulk" && route.Method == "POST" {
			foundBulk = true
		}
	}

	assert.True(t, foundVerify, "Route /v1/auth/verify should be registered")
	assert.True(t, foundBulk, "Route /v1/auth/verify-bulk should be registered")
}

func TestNewVerifyHandler_ParseCIDR(t *testing.T) {
	// Arrange
	jwtManager, _ := crypto.NewJWTManager(crypto.JWTConfig{
		Secret:               "test-secret-key-at-least-32-characters-long",
		Issuer:               "edugo-central",
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
	})

	tokenService := service.NewTokenService(jwtManager, nil, service.TokenServiceConfig{})

	// Test con diferentes formatos de IP
	testCases := []struct {
		name   string
		ranges []string
	}{
		{"CIDR format", []string{"10.0.0.0/8"}},
		{"Single IP", []string{"192.168.1.1"}},
		{"Mixed", []string{"10.0.0.0/8", "192.168.1.1"}},
		{"IPv6 CIDR", []string{"::1/128"}},
		{"Invalid", []string{"invalid-ip"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			handler := NewVerifyHandler(tokenService, tc.ranges, nil)

			// Assert
			assert.NotNil(t, handler)
		})
	}
}

func TestVerifyHandler_VerifyTokenBulk_MixedResults(t *testing.T) {
	// Arrange
	handler, jwtManager := setupTestHandler(t)
	validToken, _, _ := jwtManager.GenerateAccessToken("user-1", "user1@example.com", "admin", "")
	invalidToken := "invalid-token-here"

	router := gin.New()
	router.POST("/v1/auth/verify-bulk", handler.VerifyTokenBulk)

	body, _ := json.Marshal(dto.VerifyTokenBulkRequest{Tokens: []string{validToken, invalidToken}})
	req, _ := http.NewRequest(http.MethodPost, "/v1/auth/verify-bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-API-Key", "test-api-key-mobile")
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.VerifyTokenBulkResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response.Results, 2)

	// Verificar que hay resultados válidos e inválidos
	validCount := 0
	invalidCount := 0
	for _, r := range response.Results {
		if r.Valid {
			validCount++
		} else {
			invalidCount++
		}
	}
	assert.Equal(t, 1, validCount)
	assert.Equal(t, 1, invalidCount)
}
