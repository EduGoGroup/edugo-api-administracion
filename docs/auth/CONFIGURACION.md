# Configuración del Sistema de Autenticación Centralizada

## Variables de Entorno

### JWT (Requeridas)

| Variable | Descripción | Ejemplo | Validación |
|----------|-------------|---------|------------|
| `AUTH_JWT_SECRET` | Clave secreta para firmar tokens | `mi-clave-super-secreta-de-32-chars` | Mínimo 32 caracteres |
| `AUTH_JWT_ISSUER` | Identificador del emisor | `edugo-central` | Debe ser exactamente `edugo-central` |
| `AUTH_JWT_ACCESS_DURATION` | Duración del access token | `15m` | Formato Go duration |
| `AUTH_JWT_REFRESH_DURATION` | Duración del refresh token | `168h` | Formato Go duration (7 días) |

### Rate Limiting

| Variable | Descripción | Default |
|----------|-------------|---------|
| `AUTH_RATE_LIMIT_INTERNAL_MAX` | Requests/min servicios internos | `1000` |
| `AUTH_RATE_LIMIT_INTERNAL_WINDOW` | Ventana para internos | `1m` |
| `AUTH_RATE_LIMIT_EXTERNAL_MAX` | Requests/min clientes externos | `60` |
| `AUTH_RATE_LIMIT_EXTERNAL_WINDOW` | Ventana para externos | `1m` |

### Servicios Internos

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `AUTH_INTERNAL_API_KEYS` | Lista de API Keys válidas | `key1,key2,key3` |
| `AUTH_INTERNAL_IP_RANGES` | Rangos CIDR internos | `10.0.0.0/8,192.168.0.0/16` |

### Cache (Redis)

| Variable | Descripción | Default |
|----------|-------------|---------|
| `REDIS_HOST` | Host de Redis | `localhost` |
| `REDIS_PORT` | Puerto de Redis | `6379` |
| `REDIS_PASSWORD` | Contraseña | `` |
| `REDIS_DB` | Base de datos | `0` |
| `AUTH_CACHE_TTL` | TTL del caché de tokens | `60s` |
| `AUTH_CACHE_ENABLED` | Habilitar caché | `true` |
| `AUTH_BLACKLIST_CHECK` | Verificar blacklist | `true` |

### Password (Bcrypt)

| Variable | Descripción | Default |
|----------|-------------|---------|
| `AUTH_PASSWORD_MIN_LENGTH` | Longitud mínima | `8` |
| `AUTH_PASSWORD_BCRYPT_COST` | Costo de bcrypt | `12` |

---

## Archivo de Configuración YAML

También se puede configurar vía `configs/config.yaml`:

```yaml
auth:
  jwt:
    secret: "${AUTH_JWT_SECRET}"
    issuer: "edugo-central"
    access_duration: "15m"
    refresh_duration: "168h"
  
  password:
    min_length: 8
    bcrypt_cost: 12
  
  rate_limit:
    internal:
      max_requests: 1000
      window: "1m"
    external:
      max_requests: 60
      window: "1m"
  
  internal_services:
    api_keys:
      - "api-mobile-key"
      - "worker-key"
    ip_ranges:
      - "10.0.0.0/8"
      - "172.16.0.0/12"
  
  cache:
    enabled: true
    ttl: "60s"
    blacklist_check: true

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

---

## Configuración por Ambiente

### Desarrollo (`configs/config-dev.yaml`)

```yaml
auth:
  jwt:
    access_duration: "1h"      # Mayor duración para debug
    refresh_duration: "24h"
  rate_limit:
    external:
      max_requests: 1000       # Sin límite restrictivo
  cache:
    enabled: false             # Deshabilitado para ver cambios inmediatos
```

### Testing (`configs/config-test.yaml`)

```yaml
auth:
  jwt:
    secret: "test-secret-key-at-least-32-characters-long"
    issuer: "edugo-central"
    access_duration: "15m"
  cache:
    enabled: false
    blacklist_check: false
```

### Producción (`configs/config-prod.yaml`)

```yaml
auth:
  jwt:
    access_duration: "15m"
    refresh_duration: "168h"
  rate_limit:
    internal:
      max_requests: 1000
    external:
      max_requests: 60
  cache:
    enabled: true
    ttl: "60s"
    blacklist_check: true
```

---

## Validación de Configuración

El sistema valida la configuración al iniciar:

### Script de Validación

```bash
./scripts/validate-env.sh
```

### Validaciones Automáticas

1. **JWT Secret**: Mínimo 32 caracteres
2. **JWT Issuer**: Debe ser `edugo-central`
3. **Rate Limits**: Valores positivos
4. **Password Config**: min_length ≥ 6, bcrypt_cost entre 10-14

### Errores Comunes

| Error | Causa | Solución |
|-------|-------|----------|
| `JWT secret debe tener al menos 32 caracteres` | Secret muy corto | Usar clave más larga |
| `JWT issuer es requerido` | Issuer vacío | Configurar `edugo-central` |
| `invalid rate limit config` | Valores negativos | Usar valores positivos |

---

## Integración con Otros Servicios

### api-mobile

```env
# .env de api-mobile
AUTH_SERVICE_URL=http://api-admin:8080
AUTH_API_KEY=api-mobile-secret-key
```

### worker

```env
# .env de worker
AUTH_SERVICE_URL=http://api-admin:8080
AUTH_API_KEY=worker-secret-key
```

---

## Rotación de Secretos

### Proceso de Rotación de JWT Secret

1. Generar nuevo secret
2. Configurar ambos secrets temporalmente (grace period)
3. Esperar expiración de tokens antiguos
4. Remover secret antiguo

### Rotación de API Keys

1. Generar nueva API Key
2. Agregar a lista de keys válidas
3. Actualizar servicio cliente
4. Remover key antigua después de confirmar

---

## Monitoreo

### Métricas Recomendadas

- `auth_verify_requests_total` - Total de verificaciones
- `auth_verify_valid_total` - Tokens válidos
- `auth_verify_invalid_total` - Tokens inválidos
- `auth_rate_limit_hits_total` - Veces rate limited
- `auth_cache_hits_total` - Cache hits
- `auth_cache_misses_total` - Cache misses

### Logs

```
level=info msg="Token verified" user_id=123 valid=true duration=2ms
level=warn msg="Rate limit exceeded" client_ip=203.0.113.1 
level=error msg="Token validation failed" error="issuer inválido"
```
