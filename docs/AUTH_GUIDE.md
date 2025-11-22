# 🔐 Guía de Autenticación JWT - EduGo API Administración

**Última actualización:** 18 de Noviembre, 2025  
**Versión API:** v0.5.0+  
**Estándar del Ecosistema:** ✅ Unificado con api-mobile

---

## 📋 Tabla de Contenidos

1. [Introducción](#introducción)
2. [Autenticación Requerida](#autenticación-requerida)
3. [Obtener un Token JWT](#obtener-un-token-jwt)
4. [Usar el Token en Requests](#usar-el-token-en-requests)
5. [Ejemplos por Lenguaje](#ejemplos-por-lenguaje)
6. [Manejo de Errores](#manejo-de-errores)
7. [Claims Disponibles](#claims-disponibles)
8. [Tokens Expirados](#tokens-expirados)
9. [Buenas Prácticas](#buenas-prácticas)

---

## 🎯 Introducción

A partir de la **versión 0.5.0**, todos los endpoints de la API de Administración requieren autenticación JWT (JSON Web Token).

### ¿Por qué JWT?

- ✅ **Seguridad:** Tokens firmados que no pueden ser falsificados
- ✅ **Stateless:** El servidor no necesita almacenar sesiones
- ✅ **Claims:** Información del usuario incluida en el token
- ✅ **Estándar:** Mismo mecanismo en todo el ecosistema EduGo

### Ecosistema Unificado

Esta API usa **exactamente el mismo mecanismo** que `edugo-api-mobile`:

| Aspecto | Valor |
|---------|-------|
| **Variable de entorno** | `AUTH_JWT_SECRET` |
| **Header** | `Authorization: Bearer {token}` |
| **Claims** | `user_id`, `email`, `role` |
| **Status sin auth** | 401 Unauthorized |
| **Dependencia** | `edugo-shared/auth@v0.7.0` |

**Beneficio:** Un cliente (app móvil, web, etc.) puede usar la **misma lógica** para ambas APIs.

---

## 🔒 Autenticación Requerida

### Endpoints Protegidos

**TODOS los endpoints bajo `/v1/*` requieren JWT:**

```
✅ Requieren JWT:
  POST   /v1/schools
  GET    /v1/schools
  GET    /v1/schools/:id
  PUT    /v1/schools/:id
  DELETE /v1/schools/:id
  POST   /v1/schools/:id/units
  GET    /v1/units/:id
  PUT    /v1/units/:id
  DELETE /v1/units/:id
  POST   /v1/memberships
  GET    /v1/memberships
  ... (todos los endpoints /v1/*)
```

### Endpoints Públicos (Sin JWT)

```
❌ NO requieren JWT:
  GET /health        - Health check
  GET /swagger/*     - Documentación Swagger
```

---

## 🎫 Obtener un Token JWT

### Opción 1: Servicio de Autenticación Centralizado (Recomendado)

Si tu ecosistema EduGo tiene un servicio de autenticación:

```bash
# Login
curl -X POST https://auth.edugo.com/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@edugo.com",
    "password": "your-password"
  }'

# Respuesta
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2025-11-19T18:00:00Z",
  "user": {
    "id": "user-uuid",
    "email": "admin@edugo.com",
    "role": "admin"
  }
}
```

### Opción 2: Generar Token para Testing (Solo Desarrollo)

Para desarrollo local, puedes generar un token de prueba:

```go
package main

import (
    "fmt"
    "time"
    "github.com/EduGoGroup/edugo-shared/auth"
)

func main() {
    jwtManager := auth.NewJWTManager(
        "local-development-secret-change-in-production-min-32-chars",
        "edugo-admin",
    )
    
    claims := &auth.Claims{
        UserID: "test-user-id",
        Email:  "test@edugo.com",
        Role:   "admin",
    }
    
    token, err := jwtManager.GenerateToken(claims, 24*time.Hour)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Token:", token)
}
```

---

## 🚀 Usar el Token en Requests

### Header Requerido

```
Authorization: Bearer {tu-token-jwt}
```

**Formato:**
- Palabra clave: `Bearer` (con B mayúscula)
- Espacio
- Token JWT completo

### Ejemplo cURL

```bash
curl -X GET https://api-admin.edugo.com/v1/schools \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Ejemplo con Variables

```bash
# Guardar token en variable
export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Usar en requests
curl -X GET https://api-admin.edugo.com/v1/schools \
  -H "Authorization: Bearer $JWT_TOKEN"
```

---

## 💻 Ejemplos por Lenguaje

### JavaScript / TypeScript

```typescript
// Usando fetch (nativo)
const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...';

const response = await fetch('https://api-admin.edugo.com/v1/schools', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});

const schools = await response.json();
```

```typescript
// Usando axios
import axios from 'axios';

const client = axios.create({
  baseURL: 'https://api-admin.edugo.com',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

// Todas las requests usan el token automáticamente
const schools = await client.get('/v1/schools');
const school = await client.get('/v1/schools/123');
```

### Kotlin (Android)

```kotlin
// Usando Retrofit
interface EduGoAdminAPI {
    @GET("/v1/schools")
    suspend fun getSchools(
        @Header("Authorization") authHeader: String
    ): List<School>
}

// Uso
val api = retrofit.create(EduGoAdminAPI::class.java)
val token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
val schools = api.getSchools("Bearer $token")
```

```kotlin
// Usando OkHttp Interceptor (Recomendado)
class AuthInterceptor(private val tokenProvider: () -> String) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request().newBuilder()
            .addHeader("Authorization", "Bearer ${tokenProvider()}")
            .build()
        return chain.proceed(request)
    }
}

val client = OkHttpClient.Builder()
    .addInterceptor(AuthInterceptor { jwtToken })
    .build()

val retrofit = Retrofit.Builder()
    .baseUrl("https://api-admin.edugo.com")
    .client(client)
    .build()

// Todas las requests automáticamente incluyen el token
val schools = api.getSchools() // ✅ Token agregado automáticamente
```

### Swift (iOS)

```swift
// Usando URLSession
let token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
let url = URL(string: "https://api-admin.edugo.com/v1/schools")!

var request = URLRequest(url: url)
request.httpMethod = "GET"
request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

let (data, response) = try await URLSession.shared.data(for: request)
let schools = try JSONDecoder().decode([School].self, from: data)
```

```swift
// Usando Alamofire con Interceptor
class AuthInterceptor: RequestInterceptor {
    private let tokenProvider: () -> String
    
    init(tokenProvider: @escaping () -> String) {
        self.tokenProvider = tokenProvider
    }
    
    func adapt(_ urlRequest: URLRequest, for session: Session, completion: @escaping (Result<URLRequest, Error>) -> Void) {
        var request = urlRequest
        request.setValue("Bearer \(tokenProvider())", forHTTPHeaderField: "Authorization")
        completion(.success(request))
    }
}

let session = Session(interceptor: AuthInterceptor { jwtToken })

// Todas las requests automáticamente incluyen el token
let schools: [School] = try await session.request("https://api-admin.edugo.com/v1/schools")
    .serializingDecodable([School].self)
    .value
```

### Python

```python
import requests

token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Opción 1: Header por request
response = requests.get(
    "https://api-admin.edugo.com/v1/schools",
    headers={"Authorization": f"Bearer {token}"}
)
schools = response.json()

# Opción 2: Session con header permanente (Recomendado)
session = requests.Session()
session.headers.update({"Authorization": f"Bearer {token}"})

# Todas las requests usan el token automáticamente
schools = session.get("https://api-admin.edugo.com/v1/schools").json()
school = session.get("https://api-admin.edugo.com/v1/schools/123").json()
```

### Go

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    
    req, _ := http.NewRequest("GET", "https://api-admin.edugo.com/v1/schools", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Java

```java
// Usando OkHttp
OkHttpClient client = new OkHttpClient.Builder()
    .addInterceptor(chain -> {
        Request original = chain.request();
        Request request = original.newBuilder()
            .header("Authorization", "Bearer " + jwtToken)
            .build();
        return chain.proceed(request);
    })
    .build();

Request request = new Request.Builder()
    .url("https://api-admin.edugo.com/v1/schools")
    .build();

Response response = client.newCall(request).execute();
String json = response.body().string();
```

---

## ❌ Manejo de Errores

### Error 401: Authorization Required

**Request sin header:**
```bash
curl https://api-admin.edugo.com/v1/schools
```

**Respuesta:**
```json
{
  "error": "authorization required",
  "code": "UNAUTHORIZED"
}
```

**Status:** `401 Unauthorized`

---

### Error 401: Invalid Token

**Request con token inválido o expirado:**
```bash
curl https://api-admin.edugo.com/v1/schools \
  -H "Authorization: Bearer token-invalido"
```

**Respuesta:**
```json
{
  "error": "invalid or expired token",
  "code": "UNAUTHORIZED"
}
```

**Status:** `401 Unauthorized`

---

### Manejo en Cliente

```typescript
async function callAPI(endpoint: string, token: string) {
  const response = await fetch(`https://api-admin.edugo.com${endpoint}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  if (response.status === 401) {
    // Token inválido o expirado
    const error = await response.json();
    
    if (error.error === 'invalid or expired token') {
      // Renovar token
      const newToken = await refreshToken();
      // Reintentar request con nuevo token
      return callAPI(endpoint, newToken);
    } else {
      // No hay token, redirigir a login
      redirectToLogin();
    }
  }
  
  return response.json();
}
```

---

## 📦 Claims Disponibles

Después de validación exitosa, el middleware inyecta estos claims en el contexto:

| Claim | Tipo | Descripción | Ejemplo |
|-------|------|-------------|---------|
| `user_id` | string | ID único del usuario | `"550e8400-e29b-41d4-a716-446655440000"` |
| `email` | string | Email del usuario | `"admin@edugo.com"` |
| `role` | string | Rol del usuario | `"admin"`, `"teacher"`, `"student"` |

### Acceder a Claims en Handlers (Backend)

Si estás desarrollando handlers en Go:

```go
func (h *SchoolHandler) CreateSchool(c *gin.Context) {
    // Obtener información del usuario autenticado
    userID := c.GetString("user_id")
    email := c.GetString("email")
    role := c.GetString("role")
    
    h.logger.Info("creating school", 
        "user_id", userID, 
        "role", role,
    )
    
    // Usar para validaciones de permisos
    if role != "admin" {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "only admins can create schools",
        })
        return
    }
    
    // ... resto de la lógica
}
```

---

## ⏰ Tokens Expirados

### Detección

Los tokens JWT tienen un tiempo de expiración definido al momento de generación (usualmente 24 horas).

**Respuesta cuando token expira:**
```json
{
  "error": "invalid or expired token",
  "code": "UNAUTHORIZED"
}
```

### Solución: Refresh Token

```typescript
class EduGoClient {
  private accessToken: string;
  private refreshToken: string;
  
  async callAPI(endpoint: string) {
    try {
      return await this.request(endpoint, this.accessToken);
    } catch (error) {
      if (error.status === 401) {
        // Token expirado, renovar
        this.accessToken = await this.refresh(this.refreshToken);
        // Reintentar
        return await this.request(endpoint, this.accessToken);
      }
      throw error;
    }
  }
  
  async refresh(refreshToken: string): Promise<string> {
    const response = await fetch('https://auth.edugo.com/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
    });
    const data = await response.json();
    return data.access_token;
  }
}
```

---

## ✅ Buenas Prácticas

### 1. Usar Interceptors/Middleware en Cliente

**❌ Malo (repetitivo):**
```javascript
// Agregar header manualmente en cada request
fetch('/v1/schools', { headers: { 'Authorization': 'Bearer ...' }});
fetch('/v1/units', { headers: { 'Authorization': 'Bearer ...' }});
fetch('/v1/memberships', { headers: { 'Authorization': 'Bearer ...' }});
```

**✅ Bueno (centralizado):**
```javascript
// Configurar una vez
const client = axios.create({
  baseURL: 'https://api-admin.edugo.com',
  headers: { 'Authorization': `Bearer ${token}` }
});

// Usar sin repetir header
client.get('/v1/schools');
client.get('/v1/units');
client.get('/v1/memberships');
```

### 2. Almacenar Token de Forma Segura

**Android:**
```kotlin
// ✅ Usar EncryptedSharedPreferences
val sharedPreferences = EncryptedSharedPreferences.create(
    "edugo_secure_prefs",
    MasterKey.DEFAULT_MASTER_KEY_ALIAS,
    context,
    EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
    EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
)

sharedPreferences.edit()
    .putString("jwt_token", token)
    .apply()
```

**iOS:**
```swift
// ✅ Usar Keychain
let keychain = KeychainSwift()
keychain.set(token, forKey: "jwt_token")
```

**Web:**
```javascript
// ✅ Usar httpOnly cookies (si el backend lo soporta)
// o localStorage con precauciones
localStorage.setItem('jwt_token', token);
```

### 3. Validar Token Antes de Usar

```typescript
function isTokenExpired(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    const exp = payload.exp * 1000; // Convertir a ms
    return Date.now() >= exp;
  } catch {
    return true;
  }
}

// Antes de cada request
if (isTokenExpired(token)) {
  token = await refreshToken();
}
```

### 4. Implementar Retry con Renovación

```typescript
async function fetchWithAuth(url: string, options = {}) {
  const maxRetries = 1;
  let attempt = 0;
  
  while (attempt <= maxRetries) {
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          ...options.headers,
          'Authorization': `Bearer ${getToken()}`
        }
      });
      
      if (response.status === 401 && attempt < maxRetries) {
        // Renovar token
        await refreshToken();
        attempt++;
        continue;
      }
      
      return response;
    } catch (error) {
      throw error;
    }
  }
}
```

---

## 🔄 SDK Unificado para Ecosistema EduGo

### Ejemplo de Cliente Completo

```typescript
class EduGoEcosystemClient {
  private token: string;
  private refreshToken: string;
  
  constructor(token: string, refreshToken: string) {
    this.token = token;
    this.refreshToken = refreshToken;
  }
  
  // ✅ Mismo método funciona para AMBAS APIs
  private async request(service: 'mobile' | 'admin', endpoint: string, options = {}) {
    const baseURLs = {
      mobile: 'https://api-mobile.edugo.com',
      admin: 'https://api-admin.edugo.com'
    };
    
    const response = await fetch(`${baseURLs[service]}${endpoint}`, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${this.token}`, // ✅ Mismo header
        'Content-Type': 'application/json'
      }
    });
    
    if (response.status === 401) {
      // Token expirado, renovar
      this.token = await this.refresh();
      // Reintentar
      return this.request(service, endpoint, options);
    }
    
    return response.json();
  }
  
  // API Mobile
  async getMaterials() {
    return this.request('mobile', '/v1/materials');
  }
  
  async getAssessments() {
    return this.request('mobile', '/v1/assessments');
  }
  
  // API Admin
  async getSchools() {
    return this.request('admin', '/v1/schools');
  }
  
  async getAcademicUnits() {
    return this.request('admin', '/v1/units');
  }
  
  // Refresh token (mismo para ambas APIs)
  private async refresh(): Promise<string> {
    const response = await fetch('https://auth.edugo.com/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: this.refreshToken })
    });
    const data = await response.json();
    return data.access_token;
  }
}

// Uso unificado
const client = new EduGoEcosystemClient(accessToken, refreshToken);

// ✅ Misma lógica de auth para ambas APIs
const materials = await client.getMaterials();   // api-mobile
const schools = await client.getSchools();       // api-admin
```

---

## 🧪 Testing

### Generar Token para Tests

```bash
# Usar el mismo secret que en config-local.yaml
JWT_SECRET="local-development-secret-change-in-production-min-32-chars"

# Generar token con herramienta (ejemplo: jwt.io)
# O usar script de Go (ver Opción 2 arriba)
```

### Postman

1. **Variables de colección:**
   - `baseURL`: `https://api-admin.edugo.com`
   - `jwt_token`: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`

2. **Pre-request Script (para renovación automática):**
```javascript
// Verificar si token está expirado
const token = pm.collectionVariables.get("jwt_token");
if (isExpired(token)) {
    // Renovar token
    pm.sendRequest({
        url: 'https://auth.edugo.com/refresh',
        method: 'POST',
        body: {
            refresh_token: pm.collectionVariables.get("refresh_token")
        }
    }, (err, res) => {
        pm.collectionVariables.set("jwt_token", res.json().access_token);
    });
}
```

3. **Authorization Tab:**
   - Type: `Bearer Token`
   - Token: `{{jwt_token}}`

---

## 🌐 Ecosistema EduGo - Consistencia Total

### Mismo Token, Múltiples APIs

```typescript
// ✅ UN SOLO TOKEN funciona en TODAS las APIs del ecosistema
const token = await login('admin@edugo.com', 'password');

// Usar en api-mobile
await fetch('https://api-mobile.edugo.com/v1/materials', {
  headers: { 'Authorization': `Bearer ${token}` }
});

// ✅ Usar en api-admin (MISMO token, MISMO header)
await fetch('https://api-admin.edugo.com/v1/schools', {
  headers: { 'Authorization': `Bearer ${token}` }
});
```

### Configuración Unificada

| Aspecto | api-mobile | api-admin | Consistencia |
|---------|-----------|-----------|--------------|
| **Config path** | `auth.jwt.secret` | `auth.jwt.secret` | ✅ Idéntico |
| **Variable ENV** | `AUTH_JWT_SECRET` | `AUTH_JWT_SECRET` | ✅ Idéntico |
| **Header** | `Authorization: Bearer` | `Authorization: Bearer` | ✅ Idéntico |
| **Claims** | user_id, email, role | user_id, email, role | ✅ Idéntico |
| **Status 401** | Unauthorized | Unauthorized | ✅ Idéntico |
| **Dependencia** | shared/auth@v0.7.0 | shared/auth@v0.7.0 | ✅ Idéntico |

**Beneficio:** Aprende una vez, funciona en todo el ecosistema.

---

## 🔍 Debugging

### Verificar Token

Puedes decodificar tu token en https://jwt.io para ver su contenido:

```
Header:
{
  "alg": "HS256",
  "typ": "JWT"
}

Payload:
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "admin@edugo.com",
  "role": "admin",
  "iss": "edugo-admin",
  "exp": 1700409600
}
```

### Logs del Servidor

El servidor loguea información de autenticación:

```
✅ Auth exitoso:
  DEBUG auth successful user_id=550e8400... role=admin

❌ Auth fallido:
  WARN missing authorization header
  WARN invalid token error=token is expired
```

---

## 📞 Soporte

### Problemas Comunes

**1. "authorization required"**
- ✅ Verifica que incluyes el header `Authorization`
- ✅ Verifica el formato: `Bearer {token}` (con espacio)

**2. "invalid or expired token"**
- ✅ Verifica que el token no haya expirado
- ✅ Verifica que el token fue generado con el mismo secret
- ✅ Verifica que el token no está corrupto

**3. Token funciona en api-mobile pero no en api-admin**
- ⚠️ **NO debería pasar** - Ambas APIs usan el mismo estándar
- Si pasa, reporta un bug

---

## 🔗 Referencias

- **Especificación JWT:** https://jwt.io/introduction
- **edugo-shared/auth:** Módulo compartido de autenticación
- **Middleware:** `edugo-shared/middleware/gin`
- **Swagger:** `/swagger/index.html` (ver ejemplos con auth)

---

## 📝 Changelog

### v0.5.0 (2025-11-18)
- ✅ Implementación inicial de JWT
- ✅ Alineación con api-mobile (mismo estándar)
- ✅ Variable ENV: `AUTH_JWT_SECRET`
- ✅ Configuración: `auth.jwt.secret`

---

## ⚙️ Configuración del Servidor

### Variables de Entorno Requeridas

```bash
# Desarrollo Local
# (Ya configurado en config-local.yaml, no requiere ENV)

# Development/QA/Production
export AUTH_JWT_SECRET="your-super-secret-key-minimum-32-characters-long"
```

### Generar Secret Seguro

```bash
# Opción 1: OpenSSL
openssl rand -base64 32

# Opción 2: Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"

# Opción 3: Python
python -c "import secrets; print(secrets.token_urlsafe(32))"
```

---

**Última actualización:** 18 de Noviembre, 2025  
**Mantenido por:** Equipo EduGo Backend  
**Versión del documento:** 1.0.0
