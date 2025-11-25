# Guía de Integración - Autenticación Centralizada

Esta guía explica cómo integrar los servicios del ecosistema EduGo con el sistema de autenticación centralizada de `api-administracion`.

---

## Arquitectura

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   api-mobile    │     │     worker      │     │  otros servicios │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         │    POST /v1/auth/verify                       │
         └───────────────────────┼───────────────────────┘
                                 │
                                 ▼
                    ┌────────────────────────┐
                    │   api-administracion   │
                    │  (Auth Centralizado)   │
                    └────────────────────────┘
                                 │
                                 ▼
                    ┌────────────────────────┐
                    │    Redis (Cache)       │
                    └────────────────────────┘
```

---

## Paso 1: Configuración del Cliente

### Variables de Entorno

Agregar al `.env` del servicio cliente:

```env
# URL del servicio de autenticación
AUTH_SERVICE_URL=http://api-admin:8080

# API Key para identificación como servicio interno
AUTH_API_KEY=tu-api-key-secreta

# Timeout para requests de auth
AUTH_TIMEOUT=5s
```

### Solicitar API Key

Contactar al equipo de infraestructura para obtener una API Key única para el servicio. La API Key permite:

- Rate limiting extendido (1000 req/min vs 60 req/min)
- Acceso al endpoint de verificación bulk
- Identificación en logs y métricas

---

## Paso 2: Implementar Cliente de Verificación

### Go

```go
package auth

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

// AuthClient cliente para verificación de tokens
type AuthClient struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}

// VerifyResponse respuesta de verificación
type VerifyResponse struct {
    Valid     bool       `json:"valid"`
    UserID    string     `json:"user_id,omitempty"`
    Email     string     `json:"email,omitempty"`
    Role      string     `json:"role,omitempty"`
    ExpiresAt *time.Time `json:"expires_at,omitempty"`
    Error     string     `json:"error,omitempty"`
}

// NewAuthClient crea un nuevo cliente
func NewAuthClient() *AuthClient {
    timeout, _ := time.ParseDuration(os.Getenv("AUTH_TIMEOUT"))
    if timeout == 0 {
        timeout = 5 * time.Second
    }

    return &AuthClient{
        baseURL: os.Getenv("AUTH_SERVICE_URL"),
        apiKey:  os.Getenv("AUTH_API_KEY"),
        httpClient: &http.Client{
            Timeout: timeout,
        },
    }
}

// VerifyToken verifica un token JWT
func (c *AuthClient) VerifyToken(ctx context.Context, token string) (*VerifyResponse, error) {
    payload := map[string]string{"token": token}
    body, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("error marshaling request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST",
        c.baseURL+"/v1/auth/verify",
        bytes.NewBuffer(body))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    if c.apiKey != "" {
        req.Header.Set("X-Service-API-Key", c.apiKey)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusTooManyRequests {
        return nil, fmt.Errorf("rate limit exceeded")
    }

    var result VerifyResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("error decoding response: %w", err)
    }

    return &result, nil
}

// VerifyTokenBulk verifica múltiples tokens
func (c *AuthClient) VerifyTokenBulk(ctx context.Context, tokens []string) (map[string]*VerifyResponse, error) {
    if c.apiKey == "" {
        return nil, fmt.Errorf("API key required for bulk verification")
    }

    payload := map[string][]string{"tokens": tokens}
    body, _ := json.Marshal(payload)

    req, _ := http.NewRequestWithContext(ctx, "POST",
        c.baseURL+"/v1/auth/verify-bulk",
        bytes.NewBuffer(body))

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Service-API-Key", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Results map[string]*VerifyResponse `json:"results"`
    }
    json.NewDecoder(resp.Body).Decode(&result)

    return result.Results, nil
}
```

### Middleware de Autenticación

```go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "tu-servicio/internal/auth"
)

func AuthMiddleware(authClient *auth.AuthClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extraer token del header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header required",
            })
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")

        // Verificar token con servicio centralizado
        result, err := authClient.VerifyToken(c.Request.Context(), token)
        if err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "error": "Auth service unavailable",
            })
            c.Abort()
            return
        }

        if !result.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": result.Error,
            })
            c.Abort()
            return
        }

        // Guardar información del usuario en context
        c.Set("user_id", result.UserID)
        c.Set("user_email", result.Email)
        c.Set("user_role", result.Role)

        c.Next()
    }
}
```

---

## Paso 3: Manejo de Errores

### Retry con Backoff Exponencial

```go
func (c *AuthClient) VerifyTokenWithRetry(ctx context.Context, token string) (*VerifyResponse, error) {
    var lastErr error
    
    for attempt := 0; attempt < 3; attempt++ {
        result, err := c.VerifyToken(ctx, token)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // No reintentar si es rate limit
        if strings.Contains(err.Error(), "rate limit") {
            break
        }
        
        // Backoff exponencial
        delay := time.Duration(1<<attempt) * 100 * time.Millisecond
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(delay):
        }
    }
    
    return nil, lastErr
}
```

### Circuit Breaker (Recomendado)

```go
import "github.com/sony/gobreaker"

func NewAuthClientWithCircuitBreaker() *AuthClient {
    client := NewAuthClient()
    
    cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        "auth-service",
        MaxRequests: 3,
        Interval:    10 * time.Second,
        Timeout:     30 * time.Second,
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            return counts.ConsecutiveFailures > 5
        },
    })
    
    client.circuitBreaker = cb
    return client
}
```

---

## Paso 4: Cache Local (Opcional)

Para reducir latencia y carga en el servicio central:

```go
import (
    "sync"
    "time"
)

type CachedAuthClient struct {
    *AuthClient
    cache    map[string]*cachedResult
    mutex    sync.RWMutex
    cacheTTL time.Duration
}

type cachedResult struct {
    response  *VerifyResponse
    expiresAt time.Time
}

func (c *CachedAuthClient) VerifyToken(ctx context.Context, token string) (*VerifyResponse, error) {
    // Verificar cache
    c.mutex.RLock()
    if cached, ok := c.cache[token]; ok && time.Now().Before(cached.expiresAt) {
        c.mutex.RUnlock()
        return cached.response, nil
    }
    c.mutex.RUnlock()

    // Llamar al servicio
    result, err := c.AuthClient.VerifyToken(ctx, token)
    if err != nil {
        return nil, err
    }

    // Guardar en cache solo si es válido
    if result.Valid {
        c.mutex.Lock()
        c.cache[token] = &cachedResult{
            response:  result,
            expiresAt: time.Now().Add(c.cacheTTL),
        }
        c.mutex.Unlock()
    }

    return result, nil
}
```

**Nota**: El cache local debe tener un TTL menor al del token para evitar usar tokens revocados.

---

## Paso 5: Testing

### Mock del Cliente

```go
type MockAuthClient struct {
    VerifyFunc func(token string) (*VerifyResponse, error)
}

func (m *MockAuthClient) VerifyToken(ctx context.Context, token string) (*VerifyResponse, error) {
    if m.VerifyFunc != nil {
        return m.VerifyFunc(token)
    }
    return &VerifyResponse{Valid: true, UserID: "test-user"}, nil
}
```

### Test de Integración

```go
func TestAuthMiddleware(t *testing.T) {
    mockClient := &MockAuthClient{
        VerifyFunc: func(token string) (*VerifyResponse, error) {
            if token == "valid-token" {
                return &VerifyResponse{
                    Valid:  true,
                    UserID: "user-123",
                    Email:  "test@test.com",
                    Role:   "admin",
                }, nil
            }
            return &VerifyResponse{Valid: false, Error: "invalid token"}, nil
        },
    }

    router := gin.New()
    router.Use(AuthMiddleware(mockClient))
    router.GET("/protected", func(c *gin.Context) {
        c.JSON(200, gin.H{"user": c.GetString("user_id")})
    })

    // Test con token válido
    req, _ := http.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer valid-token")
    rec := httptest.NewRecorder()
    router.ServeHTTP(rec, req)

    assert.Equal(t, 200, rec.Code)
}
```

---

## Checklist de Integración

- [ ] Configurar variables de entorno
- [ ] Solicitar API Key al equipo de infraestructura
- [ ] Implementar cliente de verificación
- [ ] Agregar middleware de autenticación
- [ ] Implementar manejo de errores con retry
- [ ] Considerar cache local si hay alta carga
- [ ] Agregar tests unitarios y de integración
- [ ] Configurar circuit breaker para producción
- [ ] Monitorear métricas de latencia y errores

---

## Troubleshooting

| Problema | Causa Posible | Solución |
|----------|---------------|----------|
| 401 Unauthorized | Token inválido o expirado | Verificar token, renovar si expirado |
| 429 Too Many Requests | Rate limit excedido | Implementar backoff, solicitar mayor límite |
| 503 Service Unavailable | Servicio auth caído | Circuit breaker, retry con backoff |
| Latencia alta | Sin cache local | Implementar cache con TTL corto |
| Token siempre inválido | Issuer incorrecto | Verificar tokens generados con issuer `edugo-central` |
