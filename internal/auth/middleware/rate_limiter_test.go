package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/dto"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRateLimiter(config RateLimitConfig) *RateLimiter {
	return NewRateLimiter(config)
}

func TestNewRateLimiter_DefaultValues(t *testing.T) {
	// Arrange & Act
	rl := NewRateLimiter(RateLimitConfig{})
	defer rl.Stop()

	// Assert - Verificar valores por defecto
	stats := rl.GetStats()
	assert.Equal(t, 1000, stats["internal_limit"])
	assert.Equal(t, 60, stats["external_limit"])
	assert.Equal(t, "1m0s", stats["internal_window"])
	assert.Equal(t, "1m0s", stats["external_window"])
}

func TestNewRateLimiter_CustomValues(t *testing.T) {
	// Arrange & Act
	rl := NewRateLimiter(RateLimitConfig{
		InternalMaxRequests: 500,
		InternalWindow:      2 * time.Minute,
		ExternalMaxRequests: 30,
		ExternalWindow:      30 * time.Second,
	})
	defer rl.Stop()

	// Assert
	stats := rl.GetStats()
	assert.Equal(t, 500, stats["internal_limit"])
	assert.Equal(t, 30, stats["external_limit"])
	assert.Equal(t, "2m0s", stats["internal_window"])
	assert.Equal(t, "30s", stats["external_window"])
}

func TestNewRateLimiter_ParseCIDR(t *testing.T) {
	testCases := []struct {
		name        string
		ips         []string
		expectCount int
	}{
		{"CIDR format", []string{"10.0.0.0/8"}, 1},
		{"Single IPv4", []string{"192.168.1.1"}, 1},
		{"Single IPv6", []string{"::1"}, 1},
		{"Mixed valid", []string{"10.0.0.0/8", "192.168.1.1"}, 2},
		{"Invalid IP", []string{"invalid"}, 0},
		{"Empty", []string{}, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rl := NewRateLimiter(RateLimitConfig{
				InternalIPs: tc.ips,
			})
			defer rl.Stop()

			stats := rl.GetStats()
			assert.Equal(t, tc.expectCount, stats["configured_ip_ranges"])
		})
	}
}

func TestNewRateLimiter_APIKeys(t *testing.T) {
	// Arrange & Act
	rl := NewRateLimiter(RateLimitConfig{
		InternalAPIKeys: []string{"key1", "key2", "key3"},
	})
	defer rl.Stop()

	// Assert
	stats := rl.GetStats()
	assert.Equal(t, 3, stats["configured_api_keys"])
}

func TestRateLimiter_Middleware_AllowRequest(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		ExternalMaxRequests: 10,
		ExternalWindow:      time.Minute,
	})
	defer rl.Stop()

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Limit"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Reset"))
}

func TestRateLimiter_Middleware_BlockWhenExceeded(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		ExternalMaxRequests: 3,
		ExternalWindow:      time.Minute,
	})
	defer rl.Stop()

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - Hacer 4 requests (límite es 3)
	for i := 0; i < 4; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if i < 3 {
			// Primeros 3 deben pasar
			assert.Equal(t, http.StatusOK, rec.Code, "Request %d should pass", i+1)
		} else {
			// El 4to debe ser bloqueado
			assert.Equal(t, http.StatusTooManyRequests, rec.Code)

			var response dto.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "RATE_LIMIT", response.Code)
			assert.NotEmpty(t, rec.Header().Get("Retry-After"))
		}
	}
}

func TestRateLimiter_Middleware_InternalHigherLimit(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		InternalMaxRequests: 100,
		InternalWindow:      time.Minute,
		ExternalMaxRequests: 5,
		ExternalWindow:      time.Minute,
		InternalAPIKeys:     []string{"internal-api-key"},
	})
	defer rl.Stop()

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - Hacer 10 requests con API Key interna
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Service-API-Key", "internal-api-key")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Todos deben pasar (límite interno es 100)
		assert.Equal(t, http.StatusOK, rec.Code, "Internal request %d should pass", i+1)
	}
}

func TestRateLimiter_Middleware_InternalByIP(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		InternalMaxRequests: 100,
		ExternalMaxRequests: 5,
		InternalIPs:         []string{"10.0.0.0/8"},
	})
	defer rl.Stop()

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - Request desde IP interna
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.5:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert - Debe usar límite interno (100)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "100", rec.Header().Get("X-RateLimit-Limit"))
}

func TestRateLimiter_Middleware_ExternalByIP(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		InternalMaxRequests: 100,
		ExternalMaxRequests: 60,
		InternalIPs:         []string{"10.0.0.0/8"},
	})
	defer rl.Stop()

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - Request desde IP externa
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "203.0.113.50:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert - Debe usar límite externo (60)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "60", rec.Header().Get("X-RateLimit-Limit"))
}

func TestRateLimiter_Middleware_RemainingDecrements(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		ExternalMaxRequests: 10,
		ExternalWindow:      time.Minute,
	})
	defer rl.Stop()

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act & Assert - Verificar que remaining decrementa
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		expectedRemaining := 10 - (i + 1)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, 
			string(rune('0'+expectedRemaining)), 
			rec.Header().Get("X-RateLimit-Remaining"),
			"Remaining should be %d after request %d", expectedRemaining, i+1,
		)
	}
}

func TestRateLimiter_GetIdentifier(t *testing.T) {
	rl := setupRateLimiter(RateLimitConfig{})
	defer rl.Stop()

	testCases := []struct {
		name           string
		apiKey         string
		clientIP       string
		expectContains string
	}{
		{"With API Key", "test-key", "192.168.1.1", "apikey:"},
		{"Without API Key", "", "192.168.1.1", "ip:"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			var identifier string
			router.GET("/test", func(c *gin.Context) {
				identifier = rl.getIdentifier(c)
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tc.apiKey != "" {
				req.Header.Set("X-Service-API-Key", tc.apiKey)
			}
			req.RemoteAddr = tc.clientIP + ":12345"
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Contains(t, identifier, tc.expectContains)
		})
	}
}

func TestRateLimiter_IsInternalService(t *testing.T) {
	rl := setupRateLimiter(RateLimitConfig{
		InternalAPIKeys: []string{"valid-key"},
		InternalIPs:     []string{"10.0.0.0/8"},
	})
	defer rl.Stop()

	testCases := []struct {
		name     string
		apiKey   string
		clientIP string
		expected bool
	}{
		{"Valid API Key", "valid-key", "203.0.113.1", true},
		{"Invalid API Key", "invalid-key", "203.0.113.1", false},
		{"Internal IP", "", "10.0.0.5", true},
		{"External IP", "", "203.0.113.1", false},
		{"Internal IP with invalid key", "invalid-key", "10.0.0.5", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			var isInternal bool
			router.GET("/test", func(c *gin.Context) {
				isInternal = rl.isInternalService(c)
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tc.apiKey != "" {
				req.Header.Set("X-Service-API-Key", tc.apiKey)
			}
			req.RemoteAddr = tc.clientIP + ":12345"
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expected, isInternal, "Case: %s", tc.name)
		})
	}
}

func TestRateLimiter_Cleanup(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		ExternalMaxRequests: 10,
		ExternalWindow:      100 * time.Millisecond, // Ventana corta para test
	})
	defer rl.Stop()

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - Hacer un request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Verificar que hay una entrada
	stats := rl.GetStats()
	assert.Equal(t, 1, stats["active_entries"])

	// Esperar a que expire la ventana
	time.Sleep(150 * time.Millisecond)

	// Ejecutar cleanup manualmente
	rl.cleanup()

	// Verificar que se eliminó la entrada
	stats = rl.GetStats()
	assert.Equal(t, 0, stats["active_entries"])
}

func TestRateLimiter_GetStats(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		InternalMaxRequests: 500,
		InternalWindow:      2 * time.Minute,
		ExternalMaxRequests: 30,
		ExternalWindow:      30 * time.Second,
		InternalAPIKeys:     []string{"key1", "key2"},
		InternalIPs:         []string{"10.0.0.0/8", "192.168.0.0/16"},
	})
	defer rl.Stop()

	// Act
	stats := rl.GetStats()

	// Assert
	assert.Equal(t, 500, stats["internal_limit"])
	assert.Equal(t, 30, stats["external_limit"])
	assert.Equal(t, "2m0s", stats["internal_window"])
	assert.Equal(t, "30s", stats["external_window"])
	assert.Equal(t, 2, stats["configured_api_keys"])
	assert.Equal(t, 2, stats["configured_ip_ranges"])
	assert.Equal(t, 0, stats["active_entries"])
}

func TestRateLimiter_WindowReset(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		ExternalMaxRequests: 2,
		ExternalWindow:      100 * time.Millisecond,
	})
	defer rl.Stop()

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - Usar todo el límite
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Tercer request debe fallar
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)

	// Esperar a que se reinicie la ventana
	time.Sleep(150 * time.Millisecond)

	// Ahora debe permitir de nuevo
	req, _ = http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRateLimiter_DifferentClients(t *testing.T) {
	// Arrange
	rl := setupRateLimiter(RateLimitConfig{
		ExternalMaxRequests: 2,
		ExternalWindow:      time.Minute,
	})
	defer rl.Stop()

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - Cliente 1 usa su límite
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Cliente 1 bloqueado
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)

	// Cliente 2 aún puede hacer requests
	req, _ = http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRateLimiter_Stop(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(RateLimitConfig{})

	// Act & Assert - No debe hacer panic
	rl.Stop()
}
