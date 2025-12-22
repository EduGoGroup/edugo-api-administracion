# 🟢 FASE 5: Refactorización de Error Handling

**Prioridad**: Media  
**Estimación**: 8 horas  
**Rama**: `refactor/fase-5-error-handling`

---

## Preparación Git

```bash
git checkout dev
git pull origin dev
git checkout -b refactor/fase-5-error-handling
```

---

## Problema Actual

Cada handler repite el mismo patrón de manejo de errores (~15 líneas por método):

```go
if err != nil {
    if appErr, ok := errors.GetAppError(err); ok {
        h.logger.Error("operation failed", "error", appErr.Message, "code", appErr.Code)
        c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
        return
    }
    h.logger.Error("unexpected error", "error", err)
    c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
    return
}
```

**Impacto**: ~300 líneas de código duplicado entre todos los handlers.

---

## 5.1 Crear Middleware de Error Handling

### Crear Archivo
```
internal/infrastructure/http/middleware/error_handler.go
```

### Código
```go
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/EduGoGroup/edugo-api-administracion/internal/application/errors"
    httpdto "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/dto"
    "github.com/EduGoGroup/edugo-shared/logger"
)

// ErrorHandler middleware que procesa errores de forma centralizada
func ErrorHandler(log logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // Procesar errores si existen
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            handleError(c, log, err)
        }
    }
}

func handleError(c *gin.Context, log logger.Logger, err error) {
    // Intentar obtener AppError
    if appErr, ok := errors.GetAppError(err); ok {
        log.Error("request failed",
            "path", c.Request.URL.Path,
            "method", c.Request.Method,
            "error", appErr.Message,
            "code", appErr.Code,
            "status", appErr.StatusCode,
        )
        c.JSON(appErr.StatusCode, httpdto.ErrorResponse{
            Error:   appErr.Message,
            Code:    string(appErr.Code),
        })
        return
    }

    // Error genérico
    log.Error("unexpected error",
        "path", c.Request.URL.Path,
        "method", c.Request.Method,
        "error", err.Error(),
    )
    c.JSON(http.StatusInternalServerError, httpdto.ErrorResponse{
        Error: "internal server error",
        Code:  "INTERNAL_ERROR",
    })
}
```

### Registrar Middleware
En el router principal:
```go
router.Use(middleware.ErrorHandler(logger))
```

### Esfuerzo
3 horas

---

## 5.2 Refactorizar Handlers

### Antes (Código Actual)
```go
func (h *SchoolHandler) GetSchool(c *gin.Context) {
    id := c.Param("id")
    uid, err := uuid.Parse(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid school ID", Code: "INVALID_ID"})
        return
    }

    school, err := h.schoolService.GetSchool(c.Request.Context(), uid)
    if err != nil {
        if appErr, ok := errors.GetAppError(err); ok {
            h.logger.Error("get school failed", "error", appErr.Message, "code", appErr.Code)
            c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
            return
        }
        h.logger.Error("unexpected error", "error", err)
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
        return
    }

    c.JSON(http.StatusOK, school)
}
```

### Después (Código Refactorizado)
```go
func (h *SchoolHandler) GetSchool(c *gin.Context) {
    id := c.Param("id")
    uid, err := uuid.Parse(id)
    if err != nil {
        _ = c.Error(errors.NewValidationError("invalid school ID"))
        return
    }

    school, err := h.schoolService.GetSchool(c.Request.Context(), uid)
    if err != nil {
        _ = c.Error(err)
        return
    }

    c.JSON(http.StatusOK, school)
}
```

### Archivos a Refactorizar
- `internal/infrastructure/http/handler/school_handler.go`
- `internal/infrastructure/http/handler/academic_unit_handler.go`
- `internal/infrastructure/http/handler/unit_membership_handler.go`
- `internal/infrastructure/http/handler/user_handler.go`
- `internal/infrastructure/http/handler/subject_handler.go`
- `internal/infrastructure/http/handler/material_handler.go`
- `internal/infrastructure/http/handler/stats_handler.go`
- `internal/infrastructure/http/handler/guardian_handler.go`
- `internal/auth/handler/auth_handler.go`
- `internal/auth/handler/verify_handler.go`

### Esfuerzo
4 horas

---

## 5.3 Tests para el Middleware

### Crear Archivo
```
internal/infrastructure/http/middleware/error_handler_test.go
```

### Código
```go
package middleware_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/EduGoGroup/edugo-api-administracion/internal/application/errors"
    "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/http/middleware"
)

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
}

type mockLogger struct{}
func (m mockLogger) Error(msg string, args ...interface{}) {}
func (m mockLogger) Info(msg string, args ...interface{}) {}
func (m mockLogger) Warn(msg string, args ...interface{}) {}
func (m mockLogger) Debug(msg string, args ...interface{}) {}
```

### Esfuerzo
1 hora

---

## Documentación a Actualizar

Al completar esta fase, actualizar:

- `documents/improvements/REFACTORING.md` - Eliminar sección 1 (error handling)
- `documents/ARCHITECTURE.md` - Agregar documentación sobre el middleware de errores en la sección de Infrastructure Layer

---

## Finalización

```bash
git add .
git commit -m "refactor: centralizar error handling en middleware"
git push origin refactor/fase-5-error-handling
```

### Crear PR a dev con:
- Título: `refactor: centralizar error handling en middleware`
- Descripción: Fase 5 del plan de mejoras - Error Handling centralizado

---

## Checklist

- [ ] `error_handler.go` creado en middleware
- [ ] Middleware registrado en router
- [ ] `school_handler.go` refactorizado
- [ ] `academic_unit_handler.go` refactorizado
- [ ] `unit_membership_handler.go` refactorizado
- [ ] `user_handler.go` refactorizado
- [ ] Otros handlers refactorizados
- [ ] Auth handlers refactorizados
- [ ] Tests del middleware creados
- [ ] Todos los tests pasan
- [ ] Documentación actualizada
- [ ] PR creado a dev
