package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/EduGoGroup/edugo-api-administracion/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		corsConfig     *config.CORSConfig
		requestOrigin  string
		requestMethod  string
		expectedOrigin string
		expectedStatus int
		shouldHaveCORS bool
	}{
		{
			name: "Origen permitido - localhost:3000",
			corsConfig: &config.CORSConfig{
				AllowedOrigins: "http://localhost:3000,http://localhost:5173",
				AllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
				AllowedHeaders: "Content-Type,Authorization",
			},
			requestOrigin:  "http://localhost:3000",
			requestMethod:  "GET",
			expectedOrigin: "http://localhost:3000",
			expectedStatus: 200,
			shouldHaveCORS: true,
		},
		{
			name: "Origen permitido - localhost:5173",
			corsConfig: &config.CORSConfig{
				AllowedOrigins: "http://localhost:3000,http://localhost:5173",
				AllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
				AllowedHeaders: "Content-Type,Authorization",
			},
			requestOrigin:  "http://localhost:5173",
			requestMethod:  "GET",
			expectedOrigin: "http://localhost:5173",
			expectedStatus: 200,
			shouldHaveCORS: true,
		},
		{
			name: "Origen no permitido",
			corsConfig: &config.CORSConfig{
				AllowedOrigins: "http://localhost:3000,http://localhost:5173",
				AllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
				AllowedHeaders: "Content-Type,Authorization",
			},
			requestOrigin:  "http://evil.com",
			requestMethod:  "GET",
			expectedOrigin: "",
			expectedStatus: 200,
			shouldHaveCORS: false,
		},
		{
			name: "Wildcard - todos los orígenes permitidos",
			corsConfig: &config.CORSConfig{
				AllowedOrigins: "*",
				AllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
				AllowedHeaders: "Content-Type,Authorization",
			},
			requestOrigin:  "http://any-origin.com",
			requestMethod:  "GET",
			expectedOrigin: "http://any-origin.com",
			expectedStatus: 200,
			shouldHaveCORS: true,
		},
		{
			name: "Preflight request (OPTIONS) - origen permitido",
			corsConfig: &config.CORSConfig{
				AllowedOrigins: "http://localhost:3000",
				AllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
				AllowedHeaders: "Content-Type,Authorization",
			},
			requestOrigin:  "http://localhost:3000",
			requestMethod:  "OPTIONS",
			expectedOrigin: "http://localhost:3000",
			expectedStatus: 204,
			shouldHaveCORS: true,
		},
		{
			name: "Preflight request (OPTIONS) - origen no permitido",
			corsConfig: &config.CORSConfig{
				AllowedOrigins: "http://localhost:3000",
				AllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
				AllowedHeaders: "Content-Type,Authorization",
			},
			requestOrigin:  "http://evil.com",
			requestMethod:  "OPTIONS",
			expectedOrigin: "",
			expectedStatus: 204,
			shouldHaveCORS: false,
		},
		{
			name: "Sin origen en request",
			corsConfig: &config.CORSConfig{
				AllowedOrigins: "http://localhost:3000",
				AllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
				AllowedHeaders: "Content-Type,Authorization",
			},
			requestOrigin:  "",
			requestMethod:  "GET",
			expectedOrigin: "",
			expectedStatus: 200,
			shouldHaveCORS: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()
			router.Use(CORSMiddleware(tt.corsConfig))
			router.GET("/test", func(c *gin.Context) {
				c.Status(200)
			})
			router.POST("/test", func(c *gin.Context) {
				c.Status(200)
			})

			// Request
			req := httptest.NewRequest(tt.requestMethod, "/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.shouldHaveCORS {
				assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, tt.corsConfig.AllowedMethods, w.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, tt.corsConfig.AllowedHeaders, w.Header().Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
				assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}

func TestParseCSV(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "CSV simple",
			input:    "a,b,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "CSV con espacios",
			input:    "a, b , c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "CSV vacío",
			input:    "",
			expected: []string{},
		},
		{
			name:     "CSV con URLs",
			input:    "http://localhost:3000,http://localhost:5173",
			expected: []string{"http://localhost:3000", "http://localhost:5173"},
		},
		{
			name:     "CSV con elementos vacíos",
			input:    "a,,b,  ,c",
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCSV(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		expected       bool
	}{
		{
			name:           "Origen en la lista",
			origin:         "http://localhost:3000",
			allowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
			expected:       true,
		},
		{
			name:           "Origen no en la lista",
			origin:         "http://evil.com",
			allowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
			expected:       false,
		},
		{
			name:           "Wildcard permite todo",
			origin:         "http://any-origin.com",
			allowedOrigins: []string{"*"},
			expected:       true,
		},
		{
			name:           "Origen vacío",
			origin:         "",
			allowedOrigins: []string{"http://localhost:3000"},
			expected:       false,
		},
		{
			name:           "Lista vacía",
			origin:         "http://localhost:3000",
			allowedOrigins: []string{},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOriginAllowed(tt.origin, tt.allowedOrigins)
			assert.Equal(t, tt.expected, result)
		})
	}
}
