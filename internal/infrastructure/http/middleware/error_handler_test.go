package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/middleware"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// mockLogger implementa logger.Logger para tests
type mockLogger struct{}

func (m mockLogger) Error(msg string, args ...interface{})      {}
func (m mockLogger) Info(msg string, args ...interface{})       {}
func (m mockLogger) Warn(msg string, args ...interface{})       {}
func (m mockLogger) Debug(msg string, args ...interface{})      {}
func (m mockLogger) Fatal(msg string, args ...interface{})      {}
func (m mockLogger) With(fields ...interface{}) logger.Logger   { return m }
func (m mockLogger) Sync() error                                { return nil }

func TestErrorHandler_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler(mockLogger{}))

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.NewValidationError("invalid input"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid input")
	assert.Contains(t, w.Body.String(), "VALIDATION_ERROR")
}

func TestErrorHandler_NotFoundError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler(mockLogger{}))

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.NewNotFoundError("school"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "school not found")
}

func TestErrorHandler_ConflictError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler(mockLogger{}))

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.NewConflictError("school code already exists"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "school code already exists")
}

func TestErrorHandler_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler(mockLogger{}))

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(fmt.Errorf("database connection failed"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal server error")
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

func TestErrorHandler_NoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler(mockLogger{}))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestErrorHandler_MultipleErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler(mockLogger{}))

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.NewValidationError("first error"))
		_ = c.Error(errors.NewNotFoundError("second error"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// El middleware procesa el último error
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "second error not found")
}
