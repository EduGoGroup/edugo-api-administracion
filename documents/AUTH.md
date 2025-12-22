# 🔐 Sistema de Autenticación

> Autenticación centralizada para el ecosistema EduGo

## 🎯 Visión General

**EduGo API Administración** actúa como el **servicio central de autenticación** para todo el ecosistema EduGo. Esto significa:

- Un único punto de login/logout
- Tokens JWT válidos en todos los servicios
- Verificación centralizada de tokens
- Gestión unificada de sesiones

```
┌─────────────────────────────────────────────────────────────┐
│                     ECOSISTEMA EDUGO                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │  API Mobile  │    │   Workers    │    │   Web Admin  │  │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘  │
│         │                   │                   │           │
│         └───────────────────┼───────────────────┘           │
│                             │                                │
│                             ▼                                │
│              ┌──────────────────────────┐                   │
│              │   API ADMINISTRACIÓN     │                   │
│              │   (Auth Centralizado)    │                   │
│              │                          │                   │
│              │  /v1/auth/login          │                   │
│              │  /v1/auth/refresh        │                   │
│              │  /v1/auth/verify         │                   │
│              │  /v1/auth/logout         │                   │
│              └──────────────────────────┘                   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔑 JWT (JSON Web Token)

### Estructura del Token

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiI1NTBlODQwMC1lMjliLTQxZDQtYTcxNi00NDY2NTU0NDAwMDAiLCJlbWFpbCI6ImFkbWluQGVkdWdvLmNvbSIsInJvbGUiOiJzdXBlcl9hZG1pbiIsImlzcyI6ImVkdWdvLWNlbnRyYWwiLCJleHAiOjE3MzM0OTU4MDB9.
signature
```

### Claims del Token

| Claim | Tipo | Descripción |
|-------|------|-------------|
| `sub` | string | User ID (UUID) |
| `email` | string | Email del usuario |
| `role` | string | Rol del sistema |
| `iss` | string | Issuer: `edugo-central` |
| `exp` | int64 | Timestamp de expiración |
| `iat` | int64 | Timestamp de creación |
| `jti` | string | JWT ID único (para blacklist) |

### Configuración JWT

```yaml
auth:
  jwt:
    issuer: "edugo-central"          # Issuer unificado
    algorithm: "HS256"               # Algoritmo de firma
    access_token_duration: 15m       # Duración access token
    refresh_token_duration: 168h     # 7 días refresh token
```

---

## 🔄 Flujo de Autenticación

### 1. Login

```
┌────────┐                    ┌─────────────────┐              ┌────────┐
│ Client │                    │ API Admin       │              │   DB   │
└───┬────┘                    └────────┬────────┘              └───┬────┘
    │                                  │                           │
    │  POST /v1/auth/login             │                           │
    │  {email, password}               │                           │
    │─────────────────────────────────▶│                           │
    │                                  │                           │
    │                                  │  Find user by email       │
    │                                  │──────────────────────────▶│
    │                                  │◀──────────────────────────│
    │                                  │                           │
    │                                  │  Verify password (bcrypt) │
    │                                  │                           │
    │                                  │  Generate tokens          │
    │                                  │                           │
    │  {access_token, refresh_token,   │                           │
    │   user_info}                     │                           │
    │◀─────────────────────────────────│                           │
    │                                  │                           │
```

### 2. Request Autenticado

```
┌────────┐                    ┌─────────────────┐              ┌─────────────┐
│ Client │                    │ API Admin       │              │  Service    │
└───┬────┘                    └────────┬────────┘              └──────┬──────┘
    │                                  │                              │
    │  GET /v1/schools                 │                              │
    │  Authorization: Bearer {token}   │                              │
    │─────────────────────────────────▶│                              │
    │                                  │                              │
    │                                  │  JWT Middleware:             │
    │                                  │  - Validate signature        │
    │                                  │  - Check expiration          │
    │                                  │  - Extract claims            │
    │                                  │                              │
    │                                  │  Call handler                │
    │                                  │─────────────────────────────▶│
    │                                  │◀─────────────────────────────│
    │                                  │                              │
    │  Response                        │                              │
    │◀─────────────────────────────────│                              │
```

### 3. Refresh Token

```
┌────────┐                    ┌─────────────────┐
│ Client │                    │ API Admin       │
└───┬────┘                    └────────┬────────┘
    │                                  │
    │  POST /v1/auth/refresh           │
    │  {refresh_token}                 │
    │─────────────────────────────────▶│
    │                                  │
    │                                  │  Validate refresh token
    │                                  │  Check user is active
    │                                  │  Generate NEW access token
    │                                  │  (refresh token NO cambia)
    │                                  │
    │  {access_token, expires_in}      │
    │◀─────────────────────────────────│
```

### 4. Verificación (Otros Servicios)

```
┌────────────┐              ┌─────────────────┐              ┌─────────────────┐
│ API Mobile │              │ API Admin       │              │     Client      │
└─────┬──────┘              └────────┬────────┘              └────────┬────────┘
      │                              │                                │
      │                              │  GET /mobile/resource          │
      │                              │  Authorization: Bearer {token} │
      │◀─────────────────────────────────────────────────────────────│
      │                              │                                │
      │  POST /v1/auth/verify        │                                │
      │  {token}                     │                                │
      │─────────────────────────────▶│                                │
      │                              │                                │
      │  {valid: true, user_id, ...} │                                │
      │◀─────────────────────────────│                                │
      │                              │                                │
      │  Response to client          │                                │
      │──────────────────────────────────────────────────────────────▶│
```

---

## 🛡️ Seguridad

### Password Hashing

```go
// Bcrypt con cost 12 (producción)
hasher := crypto.NewPasswordHasher(12)

// Hash
hash, _ := hasher.Hash("password")
// $2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X4J...

// Verify
valid := hasher.Compare("password", hash)
```

### Validación de Passwords

```yaml
auth:
  password:
    min_length: 8              # Mínimo 8 caracteres
    require_uppercase: true    # Al menos 1 mayúscula
    require_lowercase: true    # Al menos 1 minúscula
    require_number: true       # Al menos 1 número
    require_special: false     # Caracteres especiales (opcional)
    bcrypt_cost: 10            # Cost factor
```

### Rate Limiting

```yaml
auth:
  rate_limit:
    login:
      max_attempts: 5          # Máximo 5 intentos
      window: 15m              # En ventana de 15 minutos
      block_duration: 1h       # Bloqueo de 1 hora

    internal_services:
      max_requests: 1000       # Para /verify
      window: 1m

    external_clients:
      max_requests: 60         # Para clientes externos
      window: 1m
```

---

## 🔌 Integración de Servicios

### API Keys para Servicios Internos

Los servicios internos (api-mobile, workers) pueden verificar tokens usando API Keys:

```http
POST /v1/auth/verify
X-Service-API-Key: internal-mobile-key
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Configuración:**
```env
# .env
AUTH_INTERNAL_SERVICES_API_KEYS=api-mobile:mobile-secret-key,worker:worker-secret-key
AUTH_INTERNAL_SERVICES_IP_RANGES=127.0.0.1/32,10.0.0.0/8,172.16.0.0/12
```

### Verificación Bulk

Para servicios que necesitan validar múltiples tokens:

```http
POST /v1/auth/verify-bulk
X-Service-API-Key: internal-mobile-key
Content-Type: application/json

{
  "tokens": [
    "token1...",
    "token2...",
    "token3..."
  ]
}
```

**Límite:** Máximo 100 tokens por request.

---

## 🔧 Configuración Completa

### Variables de Entorno

```env
# JWT Secret (mínimo 32 caracteres)
AUTH_JWT_SECRET=your-production-secret-minimum-32-characters-long

# Issuer unificado
AUTH_JWT_ISSUER=edugo-central

# Duración de tokens
AUTH_JWT_ACCESS_TOKEN_DURATION=15m
AUTH_JWT_REFRESH_TOKEN_DURATION=168h

# Rate limiting
AUTH_RATE_LIMIT_LOGIN_ATTEMPTS=5
AUTH_RATE_LIMIT_LOGIN_WINDOW=15m
AUTH_RATE_LIMIT_LOGIN_BLOCK=1h

# Servicios internos
AUTH_INTERNAL_SERVICES_API_KEYS=api-mobile:key1,worker:key2
AUTH_INTERNAL_SERVICES_IP_RANGES=127.0.0.1/32,10.0.0.0/8

# Cache
AUTH_CACHE_TOKEN_VALIDATION_TTL=60s
AUTH_CACHE_USER_INFO_TTL=300s
```

### Archivo YAML

```yaml
auth:
  jwt:
    issuer: "edugo-central"
    access_token_duration: 15m
    refresh_token_duration: 168h
    algorithm: "HS256"

  password:
    min_length: 8
    require_uppercase: true
    require_lowercase: true
    require_number: true
    bcrypt_cost: 10

  rate_limit:
    login:
      max_attempts: 5
      window: 15m
      block_duration: 1h

  cache:
    token_validation:
      enabled: true
      ttl: 60s
      max_size: 10000
```

---

## 📝 Ejemplos de Código

### Go Client

```go
package main

import (
    "github.com/EduGoGroup/edugo-shared/auth"
)

func main() {
    client := auth.NewAuthClient("https://api-admin.edugo.com")

    // Login
    tokens, _ := client.Login("user@edugo.com", "password")

    // Usar token
    req, _ := http.NewRequest("GET", "https://api-admin.edugo.com/v1/schools", nil)
    req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

    // Refresh
    newTokens, _ := client.Refresh(tokens.RefreshToken)
}
```

### JavaScript/TypeScript

```typescript
const API_URL = 'https://api-admin.edugo.com';

// Login
const loginResponse = await fetch(`${API_URL}/v1/auth/login`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email: 'user@edugo.com', password: 'password' })
});
const { access_token, refresh_token } = await loginResponse.json();

// Authenticated request
const response = await fetch(`${API_URL}/v1/schools`, {
  headers: { 'Authorization': `Bearer ${access_token}` }
});

// Refresh token
const refreshResponse = await fetch(`${API_URL}/v1/auth/refresh`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ refresh_token })
});
```

### cURL

```bash
# Login
curl -X POST https://api-admin.edugo.com/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@edugo.com","password":"SecurePass123"}'

# Authenticated request
curl https://api-admin.edugo.com/v1/schools \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."

# Refresh
curl -X POST https://api-admin.edugo.com/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGciOiJIUzI1NiIs..."}'

# Verify (servicio interno)
curl -X POST https://api-admin.edugo.com/v1/auth/verify \
  -H "Content-Type: application/json" \
  -H "X-Service-API-Key: internal-mobile-key" \
  -d '{"token":"eyJhbGciOiJIUzI1NiIs..."}'
```

---

## 🚨 Errores Comunes

| Código | Error | Causa | Solución |
|--------|-------|-------|----------|
| 401 | `INVALID_CREDENTIALS` | Email/password incorrectos | Verificar credenciales |
| 401 | `TOKEN_EXPIRED` | Access token expirado | Usar refresh token |
| 401 | `INVALID_REFRESH_TOKEN` | Refresh token inválido | Re-login |
| 403 | `USER_INACTIVE` | Usuario desactivado | Contactar admin |
| 429 | `RATE_LIMIT` | Demasiados intentos | Esperar `window` time |

---

## 🔒 Blacklist de Tokens

Cuando un usuario hace logout, su token se agrega a una blacklist:

```go
// Al hacer logout
tokenService.Blacklist(ctx, token, expirationTime)

// Al verificar
if tokenService.IsBlacklisted(ctx, token) {
    return ErrTokenRevoked
}
```

**Implementación actual:** En memoria (Redis próximamente)
