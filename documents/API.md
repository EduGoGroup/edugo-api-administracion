# 🔌 API Reference

> Documentación completa de endpoints, requests y responses

## 📍 Base URL

| Ambiente | URL |
|----------|-----|
| **Local** | `http://localhost:8081` |
| **Development** | `https://api-admin-dev.edugo.com` |
| **QA** | `https://api-admin-qa.edugo.com` |
| **Production** | `https://api-admin.edugo.com` |

---

## 🔐 Autenticación

Todos los endpoints `/v1/*` (excepto auth) requieren autenticación JWT.

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Endpoints públicos (sin auth):**
- `GET /health`
- `GET /swagger/*`
- `POST /v1/auth/login`
- `POST /v1/auth/refresh`
- `POST /v1/auth/verify`

---

## 📚 Índice de Endpoints

| Módulo | Endpoints |
|--------|-----------|
| [Health](#health) | Health check |
| [Auth](#auth) | Login, Logout, Refresh, Verify |
| [Schools](#schools) | CRUD de escuelas |
| [Academic Units](#academic-units) | Gestión de unidades jerárquicas |
| [Memberships](#memberships) | Asignación de usuarios a unidades |
| [Users](#users) | Membresías de usuarios |

---

## ❤️ Health

### GET /health

Health check del servicio.

**Response 200:**
```json
{
  "status": "healthy",
  "service": "edugo-api-admin"
}
```

---

## 🔐 Auth

### POST /v1/auth/login

Autenticar usuario y obtener tokens.

**Request:**
```json
{
  "email": "admin@edugo.com",
  "password": "SecurePass123"
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@edugo.com",
    "first_name": "Admin",
    "last_name": "User",
    "role": "super_admin"
  }
}
```

**Response 401:**
```json
{
  "error": "unauthorized",
  "message": "Credenciales inválidas",
  "code": "INVALID_CREDENTIALS"
}
```

---

### POST /v1/auth/refresh

Renovar access token con refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

---

### POST /v1/auth/logout

Invalidar token actual.

**Headers:**
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response 200:**
```json
{
  "message": "Logout exitoso"
}
```

---

### POST /v1/auth/verify

Verificar validez de un token (para otros servicios).

**Request:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response 200:**
```json
{
  "valid": true,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "admin@edugo.com",
  "role": "super_admin",
  "expires_at": "2025-12-06T15:30:00Z"
}
```

**Response 200 (token inválido):**
```json
{
  "valid": false,
  "error": "token_expired"
}
```

---

### POST /v1/auth/verify-bulk

Verificar múltiples tokens (solo servicios internos con API Key).

**Headers:**
```http
X-Service-API-Key: internal-mobile-key
```

**Request:**
```json
{
  "tokens": [
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  ]
}
```

**Response 200:**
```json
{
  "results": [
    {
      "token": "eyJ...",
      "valid": true,
      "user_id": "550e8400-e29b-41d4-a716-446655440000"
    },
    {
      "token": "eyJ...",
      "valid": false,
      "error": "token_expired"
    }
  ],
  "total": 2,
  "valid_count": 1,
  "invalid_count": 1
}
```

---

## 🏫 Schools

### POST /v1/schools

Crear nueva escuela.

**Headers:**
```http
Authorization: Bearer {token}
```

**Request:**
```json
{
  "name": "Colegio San Martín",
  "code": "san-martin",
  "address": "Av. Principal 123",
  "contact_email": "contacto@sanmartin.edu",
  "contact_phone": "+54 11 4567-8900",
  "metadata": {
    "foundation_year": 1985,
    "type": "private"
  }
}
```

**Response 201:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "Colegio San Martín",
  "code": "san-martin",
  "address": "Av. Principal 123",
  "contact_email": "contacto@sanmartin.edu",
  "contact_phone": "+54 11 4567-8900",
  "metadata": {
    "foundation_year": 1985,
    "type": "private"
  },
  "created_at": "2025-12-06T10:30:00Z",
  "updated_at": "2025-12-06T10:30:00Z"
}
```

---

### GET /v1/schools

Listar escuelas.

**Query Parameters:**
| Param | Tipo | Default | Descripción |
|-------|------|---------|-------------|
| `limit` | int | 20 | Máximo de resultados |
| `offset` | int | 0 | Offset para paginación |

**Response 200:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Colegio San Martín",
    "code": "san-martin",
    "address": "Av. Principal 123",
    "created_at": "2025-12-06T10:30:00Z",
    "updated_at": "2025-12-06T10:30:00Z"
  }
]
```

---

### GET /v1/schools/:id

Obtener escuela por ID.

**Response 200:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "Colegio San Martín",
  "code": "san-martin",
  "address": "Av. Principal 123",
  "contact_email": "contacto@sanmartin.edu",
  "contact_phone": "+54 11 4567-8900",
  "metadata": {...},
  "created_at": "2025-12-06T10:30:00Z",
  "updated_at": "2025-12-06T10:30:00Z"
}
```

**Response 404:**
```json
{
  "error": "not_found",
  "message": "Escuela no encontrada"
}
```

---

### GET /v1/schools/code/:code

Obtener escuela por código único.

**Response 200:** (igual que GET /v1/schools/:id)

---

### PUT /v1/schools/:id

Actualizar escuela.

**Request:**
```json
{
  "name": "Colegio San Martín Actualizado",
  "address": "Nueva Av. 456"
}
```

**Response 200:** (escuela actualizada)

---

### DELETE /v1/schools/:id

Eliminar escuela (soft delete).

**Response 204:** No Content

---

## 🏛️ Academic Units

### POST /v1/schools/:id/units

Crear unidad académica dentro de una escuela.

**Request:**
```json
{
  "parent_unit_id": null,
  "type": "grade",
  "display_name": "1° Primaria",
  "code": "1-primaria",
  "description": "Primer grado de primaria",
  "metadata": {
    "year": 2025,
    "capacity": 30
  }
}
```

**Tipos de unidad válidos:**
- `school` - Nivel escuela
- `grade` - Grado
- `section` - Sección
- `club` - Club
- `department` - Departamento

**Response 201:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "school_id": "550e8400-e29b-41d4-a716-446655440001",
  "parent_unit_id": null,
  "type": "grade",
  "display_name": "1° Primaria",
  "code": "1-primaria",
  "description": "Primer grado de primaria",
  "metadata": {...},
  "created_at": "2025-12-06T10:35:00Z",
  "updated_at": "2025-12-06T10:35:00Z"
}
```

---

### GET /v1/schools/:id/units

Listar unidades de una escuela.

**Response 200:**
```json
[
  {
    "id": "...",
    "school_id": "...",
    "parent_unit_id": null,
    "type": "grade",
    "display_name": "1° Primaria",
    "code": "1-primaria"
  },
  {
    "id": "...",
    "school_id": "...",
    "parent_unit_id": "...",
    "type": "section",
    "display_name": "Sección A",
    "code": "1-primaria-a"
  }
]
```

---

### GET /v1/schools/:id/units/tree

Obtener árbol jerárquico de unidades.

**Response 200:**
```json
[
  {
    "id": "...",
    "type": "grade",
    "display_name": "1° Primaria",
    "code": "1-primaria",
    "depth": 1,
    "children": [
      {
        "id": "...",
        "type": "section",
        "display_name": "Sección A",
        "code": "1-primaria-a",
        "depth": 2,
        "children": []
      },
      {
        "id": "...",
        "type": "section",
        "display_name": "Sección B",
        "code": "1-primaria-b",
        "depth": 2,
        "children": []
      }
    ]
  }
]
```

---

### GET /v1/schools/:id/units/by-type

Filtrar unidades por tipo.

**Query Parameters:**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `type` | string | grade, section, club, department |

**Response 200:** (lista de unidades del tipo especificado)

---

### GET /v1/units/:id

Obtener unidad por ID.

---

### PUT /v1/units/:id

Actualizar unidad.

**Request:**
```json
{
  "display_name": "Primer Grado",
  "description": "Descripción actualizada"
}
```

---

### DELETE /v1/units/:id

Eliminar unidad (soft delete).

---

### POST /v1/units/:id/restore

Restaurar unidad eliminada.

---

### GET /v1/units/:id/hierarchy-path

Obtener path jerárquico desde la raíz hasta la unidad.

**Response 200:**
```json
[
  {
    "id": "...",
    "type": "school",
    "display_name": "Colegio San Martín",
    "depth": 0
  },
  {
    "id": "...",
    "type": "grade",
    "display_name": "1° Primaria",
    "depth": 1
  },
  {
    "id": "...",
    "type": "section",
    "display_name": "Sección A",
    "depth": 2
  }
]
```

---

## 👥 Memberships

### POST /v1/memberships

Crear membresía (asignar usuario a unidad).

**Request:**
```json
{
  "unit_id": "550e8400-e29b-41d4-a716-446655440002",
  "user_id": "550e8400-e29b-41d4-a716-446655440003",
  "role": "teacher",
  "valid_from": "2025-03-01T00:00:00Z",
  "valid_until": "2025-12-31T23:59:59Z"
}
```

**Roles válidos:**
- `director`
- `coordinator`
- `teacher`
- `assistant`
- `student`
- `observer`

**Response 201:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440004",
  "unit_id": "550e8400-e29b-41d4-a716-446655440002",
  "user_id": "550e8400-e29b-41d4-a716-446655440003",
  "role": "teacher",
  "enrolled_at": "2025-03-01T00:00:00Z",
  "is_active": true,
  "created_at": "2025-12-06T10:40:00Z",
  "updated_at": "2025-12-06T10:40:00Z"
}
```

---

### GET /v1/memberships

Listar membresías de una unidad.

**Query Parameters:**
| Param | Tipo | Required | Descripción |
|-------|------|----------|-------------|
| `unit_id` | UUID | Sí | ID de la unidad |

---

### GET /v1/memberships/by-role

Listar membresías por rol.

**Query Parameters:**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `unit_id` | UUID | ID de la unidad |
| `role` | string | Rol a filtrar |

---

### GET /v1/memberships/:id

Obtener membresía por ID.

---

### PUT /v1/memberships/:id

Actualizar membresía.

**Request:**
```json
{
  "role": "coordinator",
  "valid_until": "2026-06-30T23:59:59Z"
}
```

---

### DELETE /v1/memberships/:id

Eliminar membresía.

---

### POST /v1/memberships/:id/expire

Marcar membresía como expirada.

---

## 👤 Users

### GET /v1/users/:userId/memberships

Listar todas las membresías de un usuario.

**Response 200:**
```json
[
  {
    "id": "...",
    "unit_id": "...",
    "user_id": "...",
    "role": "teacher",
    "is_active": true,
    "enrolled_at": "2025-03-01T00:00:00Z"
  }
]
```

---

## 🛑 Códigos de Error

| HTTP Code | Error Code | Descripción |
|-----------|------------|-------------|
| 400 | `INVALID_REQUEST` | Request malformado o validación fallida |
| 401 | `UNAUTHORIZED` | Token faltante o inválido |
| 401 | `INVALID_CREDENTIALS` | Credenciales de login incorrectas |
| 401 | `TOKEN_EXPIRED` | Token JWT expirado |
| 403 | `FORBIDDEN` | Sin permisos para la operación |
| 403 | `USER_INACTIVE` | Usuario desactivado |
| 404 | `NOT_FOUND` | Recurso no encontrado |
| 409 | `CONFLICT` | Conflicto (ej: código duplicado) |
| 429 | `RATE_LIMIT` | Demasiadas requests |
| 500 | `INTERNAL_ERROR` | Error interno del servidor |

**Formato de Error:**
```json
{
  "error": "not_found",
  "message": "Escuela no encontrada",
  "code": "NOT_FOUND"
}
```

---

## 📝 Ejemplos Completos por Lenguaje

### cURL

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@edugo.com","password":"SecurePass123"}' \
  | jq -r '.access_token')

# Crear escuela
curl -X POST http://localhost:8081/v1/schools \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Colegio San Martín",
    "code": "san-martin",
    "address": "Av. Principal 123",
    "contact_email": "info@sanmartin.edu"
  }'

# Listar escuelas
curl http://localhost:8081/v1/schools \
  -H "Authorization: Bearer $TOKEN"

# Crear unidad académica
curl -X POST http://localhost:8081/v1/schools/{school_id}/units \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "grade",
    "display_name": "1° Primaria",
    "code": "1-primaria"
  }'

# Obtener árbol jerárquico
curl http://localhost:8081/v1/schools/{school_id}/units/tree \
  -H "Authorization: Bearer $TOKEN"
```

### JavaScript/TypeScript

```typescript
const API_URL = 'http://localhost:8081';

class EduGoClient {
  private accessToken: string = '';
  private refreshToken: string = '';

  async login(email: string, password: string): Promise<void> {
    const response = await fetch(`${API_URL}/v1/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
    
    if (!response.ok) {
      throw new Error('Login failed');
    }
    
    const data = await response.json();
    this.accessToken = data.access_token;
    this.refreshToken = data.refresh_token;
  }

  private async authFetch(url: string, options: RequestInit = {}): Promise<Response> {
    const response = await fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json'
      }
    });

    // Auto-refresh en 401
    if (response.status === 401) {
      await this.refresh();
      return this.authFetch(url, options);
    }

    return response;
  }

  async refresh(): Promise<void> {
    const response = await fetch(`${API_URL}/v1/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: this.refreshToken })
    });
    
    const data = await response.json();
    this.accessToken = data.access_token;
  }

  // Schools
  async createSchool(data: CreateSchoolRequest): Promise<School> {
    const response = await this.authFetch(`${API_URL}/v1/schools`, {
      method: 'POST',
      body: JSON.stringify(data)
    });
    return response.json();
  }

  async listSchools(): Promise<School[]> {
    const response = await this.authFetch(`${API_URL}/v1/schools`);
    return response.json();
  }

  async getSchool(id: string): Promise<School> {
    const response = await this.authFetch(`${API_URL}/v1/schools/${id}`);
    return response.json();
  }

  // Academic Units
  async createUnit(schoolId: string, data: CreateUnitRequest): Promise<AcademicUnit> {
    const response = await this.authFetch(`${API_URL}/v1/schools/${schoolId}/units`, {
      method: 'POST',
      body: JSON.stringify(data)
    });
    return response.json();
  }

  async getUnitTree(schoolId: string): Promise<UnitTreeNode[]> {
    const response = await this.authFetch(`${API_URL}/v1/schools/${schoolId}/units/tree`);
    return response.json();
  }

  // Memberships
  async createMembership(data: CreateMembershipRequest): Promise<Membership> {
    const response = await this.authFetch(`${API_URL}/v1/memberships`, {
      method: 'POST',
      body: JSON.stringify(data)
    });
    return response.json();
  }
}

// Tipos
interface CreateSchoolRequest {
  name: string;
  code: string;
  address?: string;
  contact_email?: string;
  contact_phone?: string;
  metadata?: Record<string, any>;
}

interface School {
  id: string;
  name: string;
  code: string;
  address: string;
  contact_email: string;
  contact_phone: string;
  metadata: Record<string, any>;
  created_at: string;
  updated_at: string;
}

interface CreateUnitRequest {
  parent_unit_id?: string;
  type: 'school' | 'grade' | 'section' | 'club' | 'department';
  display_name: string;
  code?: string;
  description?: string;
  metadata?: Record<string, any>;
}

interface AcademicUnit {
  id: string;
  school_id: string;
  parent_unit_id?: string;
  type: string;
  display_name: string;
  code: string;
  description: string;
  metadata: Record<string, any>;
  created_at: string;
  updated_at: string;
}

interface UnitTreeNode {
  id: string;
  type: string;
  display_name: string;
  code: string;
  depth: number;
  children: UnitTreeNode[];
}

interface CreateMembershipRequest {
  unit_id: string;
  user_id: string;
  role: string;
  valid_from?: string;
  valid_until?: string;
}

interface Membership {
  id: string;
  unit_id: string;
  user_id: string;
  role: string;
  enrolled_at: string;
  withdrawn_at?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Go

```go
package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type EduGoClient struct {
    baseURL      string
    accessToken  string
    refreshToken string
    httpClient   *http.Client
}

func NewEduGoClient(baseURL string) *EduGoClient {
    return &EduGoClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
    }
}

func (c *EduGoClient) Login(email, password string) error {
    payload := map[string]string{"email": email, "password": password}
    body, _ := json.Marshal(payload)
    
    resp, err := c.httpClient.Post(
        c.baseURL+"/v1/auth/login",
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var result struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    
    c.accessToken = result.AccessToken
    c.refreshToken = result.RefreshToken
    return nil
}

func (c *EduGoClient) doRequest(method, path string, body interface{}) (*http.Response, error) {
    var bodyReader *bytes.Buffer
    if body != nil {
        jsonBody, _ := json.Marshal(body)
        bodyReader = bytes.NewBuffer(jsonBody)
    } else {
        bodyReader = bytes.NewBuffer(nil)
    }
    
    req, _ := http.NewRequest(method, c.baseURL+path, bodyReader)
    req.Header.Set("Authorization", "Bearer "+c.accessToken)
    req.Header.Set("Content-Type", "application/json")
    
    return c.httpClient.Do(req)
}

func (c *EduGoClient) ListSchools() ([]School, error) {
    resp, err := c.doRequest("GET", "/v1/schools", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var schools []School
    json.NewDecoder(resp.Body).Decode(&schools)
    return schools, nil
}

func (c *EduGoClient) CreateSchool(req CreateSchoolRequest) (*School, error) {
    resp, err := c.doRequest("POST", "/v1/schools", req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var school School
    json.NewDecoder(resp.Body).Decode(&school)
    return &school, nil
}

// Types
type CreateSchoolRequest struct {
    Name         string                 `json:"name"`
    Code         string                 `json:"code"`
    Address      string                 `json:"address,omitempty"`
    ContactEmail string                 `json:"contact_email,omitempty"`
    ContactPhone string                 `json:"contact_phone,omitempty"`
    Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type School struct {
    ID           string                 `json:"id"`
    Name         string                 `json:"name"`
    Code         string                 `json:"code"`
    Address      string                 `json:"address"`
    ContactEmail string                 `json:"contact_email"`
    ContactPhone string                 `json:"contact_phone"`
    Metadata     map[string]interface{} `json:"metadata"`
    CreatedAt    string                 `json:"created_at"`
    UpdatedAt    string                 `json:"updated_at"`
}
```

---

## 📖 Swagger

Documentación interactiva disponible en:
```
http://localhost:8081/swagger/index.html
```

Para regenerar:
```bash
make swagger
# o
swag init -g cmd/main.go -o docs
```

---

## 📊 Rate Limiting

| Endpoint | Límite | Ventana |
|----------|--------|--------|
| `/v1/auth/login` | 5 intentos | 15 min |
| `/v1/auth/verify` (interno) | 1000 req | 1 min |
| `/v1/auth/verify` (externo) | 60 req | 1 min |
| Otros endpoints | Sin límite | - |

**Headers de Rate Limit:**
```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 55
X-RateLimit-Reset: 1733495800
```

---

## 🔄 Paginación

Para endpoints que retornan listas:

**Query Parameters:**
```
?limit=20&offset=0
```

| Parámetro | Default | Máx | Descripción |
|-----------|---------|-----|------------|
| `limit` | 20 | 100 | Items por página |
| `offset` | 0 | - | Items a saltar |

**Response Headers (futuro):**
```
X-Total-Count: 150
X-Page: 1
X-Page-Size: 20
```
