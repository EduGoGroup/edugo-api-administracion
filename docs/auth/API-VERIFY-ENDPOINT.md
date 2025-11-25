# API de Verificación de Tokens

## Descripción General

El endpoint `/v1/auth/verify` permite a los servicios del ecosistema EduGo validar tokens JWT de manera centralizada. Este es el componente principal del sistema de autenticación centralizada.

---

## Endpoints

### POST /v1/auth/verify

Verifica la validez de un token JWT individual.

#### Request

```http
POST /v1/auth/verify
Content-Type: application/json
X-Service-API-Key: <api-key-opcional>
```

**Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| token | string | Sí | Token JWT a verificar. Puede incluir prefijo "Bearer " |

#### Response (200 OK)

**Token Válido:**
```json
{
  "valid": true,
  "user_id": "user-123",
  "email": "usuario@ejemplo.com",
  "role": "admin",
  "expires_at": "2025-01-24T15:30:00Z"
}
```

**Token Inválido:**
```json
{
  "valid": false,
  "error": "token expirado"
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| valid | boolean | Indica si el token es válido |
| user_id | string | ID del usuario (solo si válido) |
| email | string | Email del usuario (solo si válido) |
| role | string | Rol del usuario (solo si válido) |
| expires_at | string | Fecha de expiración ISO 8601 (solo si válido) |
| error | string | Mensaje de error (solo si inválido) |

#### Response (400 Bad Request)

```json
{
  "error": "bad_request",
  "message": "Token es requerido",
  "code": "INVALID_REQUEST"
}
```

#### Response (429 Too Many Requests)

```json
{
  "error": "rate_limit_exceeded",
  "message": "Demasiadas solicitudes. Intente de nuevo más tarde.",
  "code": "RATE_LIMIT"
}
```

---

### POST /v1/auth/verify-bulk

Verifica múltiples tokens JWT en una sola llamada. **Requiere API Key de servicio interno.**

#### Request

```http
POST /v1/auth/verify-bulk
Content-Type: application/json
X-Service-API-Key: <api-key-requerida>
```

**Body:**
```json
{
  "tokens": [
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  ]
}
```

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| tokens | array | Sí | Lista de tokens (máximo 100) |

#### Response (200 OK)

```json
{
  "results": {
    "eyJhbGci...xyz": {
      "valid": true,
      "user_id": "user-123",
      "email": "usuario1@ejemplo.com",
      "role": "admin",
      "expires_at": "2025-01-24T15:30:00Z"
    },
    "eyJhbGci...abc": {
      "valid": false,
      "error": "token expirado"
    }
  }
}
```

#### Response (401 Unauthorized)

```json
{
  "error": "unauthorized",
  "message": "API Key requerida para verificación en lote",
  "code": "API_KEY_REQUIRED"
}
```

---

## Headers de Rate Limiting

Todas las respuestas incluyen headers de rate limiting:

| Header | Descripción |
|--------|-------------|
| X-RateLimit-Limit | Límite máximo de requests por ventana |
| X-RateLimit-Remaining | Requests restantes en la ventana actual |
| X-RateLimit-Reset | Timestamp Unix cuando se reinicia el contador |
| Retry-After | Segundos a esperar (solo cuando rate limited) |
| X-Response-Time | Tiempo de procesamiento del request |

### Límites por Tipo de Cliente

| Tipo | Límite | Ventana |
|------|--------|---------|
| Servicios Internos | 1000 requests | 1 minuto |
| Clientes Externos | 60 requests | 1 minuto |

**Identificación de Servicios Internos:**
- Header `X-Service-API-Key` con API Key válida
- IP de origen en rango CIDR configurado (ej: `10.0.0.0/8`)

---

## Códigos de Error

| Código | HTTP Status | Descripción |
|--------|-------------|-------------|
| INVALID_REQUEST | 400 | Request malformado o token vacío |
| EMPTY_TOKEN | 400 | Token solo contiene espacios |
| EMPTY_TOKENS | 400 | Lista de tokens vacía |
| TOO_MANY_TOKENS | 400 | Más de 100 tokens en bulk |
| API_KEY_REQUIRED | 401 | Bulk requiere API Key |
| RATE_LIMIT | 429 | Rate limit excedido |
| VERIFICATION_ERROR | 500 | Error interno de verificación |

---

## Ejemplos de Uso

### cURL - Verificar Token Individual

```bash
curl -X POST https://api-admin.edugo.com/v1/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}'
```

### cURL - Verificación Bulk (Servicio Interno)

```bash
curl -X POST https://api-admin.edugo.com/v1/auth/verify-bulk \
  -H "Content-Type: application/json" \
  -H "X-Service-API-Key: mi-api-key-secreta" \
  -d '{
    "tokens": [
      "token1...",
      "token2..."
    ]
  }'
```

### Go - Cliente de Verificación

```go
func VerifyToken(token string) (*VerifyResponse, error) {
    payload := map[string]string{"token": token}
    body, _ := json.Marshal(payload)
    
    req, _ := http.NewRequest("POST", 
        "https://api-admin.edugo.com/v1/auth/verify", 
        bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Service-API-Key", os.Getenv("AUTH_API_KEY"))
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result VerifyResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

---

## Notas de Seguridad

1. **HTTPS Obligatorio**: Todos los requests deben usar HTTPS en producción
2. **API Keys**: Rotarlas periódicamente y no exponerlas en código cliente
3. **Rate Limiting**: Implementar backoff exponencial ante 429
4. **Tokens Sensibles**: No loguear tokens completos, usar truncado
5. **Issuer Único**: Solo se aceptan tokens con issuer `edugo-central`
