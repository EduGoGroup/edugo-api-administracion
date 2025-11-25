# EduGo API Administración

API REST para operaciones administrativas y CRUD en la plataforma EduGo.

## 🔐 Autenticación

**⚠️ IMPORTANTE:** Todos los endpoints `/v1/*` requieren autenticación JWT.

📖 **[Ver Guía Completa de Autenticación](docs/AUTH_GUIDE.md)** - Ejemplos en JavaScript, Kotlin, Swift, Python, Go, Java

**Quick Start:**
```bash
# Incluir header en todas las requests
curl -H "Authorization: Bearer {tu-token-jwt}" \
  https://api-admin.edugo.com/v1/schools
```

**Ecosistema Unificado:** Esta API usa el **mismo mecanismo de autenticación** que `edugo-api-mobile`.  
Un token funciona en ambas APIs. [Ver más](docs/AUTH_GUIDE.md#ecosistema-unificado)

---

## Descripción

Esta API maneja:
- Gestión de usuarios (crear, editar, eliminar)
- Gestión de jerarquía académica (escuelas, unidades)
- Gestión de materias
- Moderación de contenidos
- Estadísticas globales del sistema

## Arquitectura

Este proyecto utiliza **Clean Architecture** siguiendo las mejores prácticas de Go:

```
├── cmd/                          # Punto de entrada de la aplicación
├── internal/
│   ├── application/              # Capa de aplicación
│   │   ├── dto/                  # Data Transfer Objects
│   │   └── service/              # Servicios de aplicación
│   ├── domain/                   # Capa de dominio
│   │   ├── entity/               # Entidades de dominio
│   │   ├── repository/           # Interfaces de repositorios
│   │   └── valueobject/          # Value Objects
│   ├── infrastructure/           # Capa de infraestructura
│   │   ├── http/handler/         # Handlers HTTP (Gin)
│   │   └── persistence/postgres/ # Implementaciones de repositorios
│   ├── bootstrap/                # Inicialización de infraestructura
│   ├── config/                   # Configuración
│   └── container/                # Inyección de dependencias
└── test/
    ├── integration/              # Tests de integración con testcontainers
    └── unit/                     # Tests unitarios
```

## Tecnología

- **Go 1.21+** + Gin + Swagger
- **PostgreSQL 15** (base de datos relacional)
- **MongoDB 7.0** (logs y eventos)
- **shared/bootstrap** (componentes compartidos de EduGo)
- **Testcontainers** (tests de integración)
- Puerto: `8081`

## Instalación

### Requisitos

- Go 1.21+
- Docker (para tests de integración)
- PostgreSQL 15+ (para desarrollo local)

### Setup

```bash
# Instalar dependencias
go mod download

# Generar documentación Swagger
swag init -g cmd/main.go -o docs

# Ejecutar la aplicación
go run cmd/main.go
```

### Variables de Entorno

Crear archivo `.env` en la raíz del proyecto:

```env
# Ambiente
APP_ENV=local

# Server
SERVER_PORT=8081
SERVER_HOST=0.0.0.0

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DATABASE=edugo
POSTGRES_USER=edugo_user
POSTGRES_PASSWORD=your_password_here
POSTGRES_MAX_CONNECTIONS=25
POSTGRES_SSL_MODE=disable

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=edugo

# Autenticación JWT
AUTH_JWT_SECRET=local-development-secret-change-in-production-min-32-chars

# Logging
LOGGING_LEVEL=info
LOGGING_FORMAT=json
```

**Nota:** Para local, el secret ya está configurado en `config/config-local.yaml`.  
Para dev/qa/prod, la variable `AUTH_JWT_SECRET` es **OBLIGATORIA**.

## Comandos Disponibles

```bash
# Compilar
make build

# Ejecutar
make run

# Tests unitarios
make test

# Tests de integración (requiere Docker)
make test-integration

# Linting
make lint

# Coverage
make coverage

# Generar Swagger
make swagger
```

## Endpoints

🔐 **Todos los endpoints `/v1/*` requieren autenticación JWT.** [Ver guía](docs/AUTH_GUIDE.md)

| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| POST | `/v1/schools` | Crear escuela | 🔐 JWT |
| GET | `/v1/schools` | Listar escuelas | 🔐 JWT |
| GET | `/v1/schools/:id` | Obtener escuela | 🔐 JWT |
| PUT | `/v1/schools/:id` | Actualizar escuela | 🔐 JWT |
| DELETE | `/v1/schools/:id` | Eliminar escuela | 🔐 JWT |
| POST | `/v1/schools/:id/units` | Crear unidad académica | 🔐 JWT |
| GET | `/v1/units/:id` | Obtener unidad | 🔐 JWT |
| PUT | `/v1/units/:id` | Actualizar unidad | 🔐 JWT |
| DELETE | `/v1/units/:id` | Eliminar unidad | 🔐 JWT |
| POST | `/v1/memberships` | Crear membresía | 🔐 JWT |
| GET | `/v1/memberships` | Listar membresías | 🔐 JWT |
| GET | `/health` | Health check | ❌ Público |
| GET | `/swagger/*` | Documentación | ❌ Público |

**Ejemplo con autenticación:**
```bash
curl -H "Authorization: Bearer {token}" \
  https://api-admin.edugo.com/v1/schools
```

📖 **[Ver ejemplos completos en todos los lenguajes](docs/AUTH_GUIDE.md#ejemplos-por-lenguaje)**

## Swagger

Documentación interactiva disponible en:  
`http://localhost:8081/swagger/index.html`

## Testing

### Tests Unitarios

```bash
go test ./internal/... -v
```

### Tests de Integración

Los tests de integración usan **testcontainers** para levantar PostgreSQL y MongoDB automáticamente:

```bash
go test ./test/integration/... -v -tags=integration
```

## Estado del Proyecto

### ✅ Completado

- Arquitectura Clean Architecture implementada
- Bootstrap con shared/bootstrap v0.1.0
- Configuración modular con validación
- Tests de integración con testcontainers
- Health check endpoint
- Documentación Swagger

### 🚧 En Desarrollo

- Jerarquía académica (FASE 2-7)
- Tests unitarios completos
- CI/CD pipelines

### 📋 Pendiente

- Validación de rol admin en middleware
- Auditoría completa en `audit_log`
- Métricas y observabilidad

## Contribuir

1. Crear rama desde `dev`: `git checkout -b feature/mi-feature`
2. Implementar cambios siguiendo Clean Architecture
3. Agregar tests (unitarios + integración)
4. Ejecutar linting: `make lint`
5. Crear PR hacia `dev`

## Licencia

Privado - EduGo © 2025

---

## 🔑 Autenticación Centralizada (Nuevo)

A partir de la versión 1.1.0, `api-administracion` actúa como el **servicio central de autenticación** para todo el ecosistema EduGo.

### Endpoint de Verificación

```bash
POST /v1/auth/verify
```

Permite a otros servicios (api-mobile, worker) validar tokens JWT de manera centralizada.

### Documentación

| Documento | Descripción |
|-----------|-------------|
| [API Verify Endpoint](docs/auth/API-VERIFY-ENDPOINT.md) | Documentación completa del endpoint |
| [Configuración](docs/auth/CONFIGURACION.md) | Variables de entorno y configuración |
| [Guía de Integración](docs/auth/GUIA-INTEGRACION.md) | Cómo integrar otros servicios |

### Características

- ✅ Verificación individual y bulk de tokens
- ✅ Rate limiting diferenciado (interno/externo)
- ✅ Cache de resultados con Redis
- ✅ Blacklist para tokens revocados
- ✅ Identificación de servicios internos por API Key o IP
- ✅ Issuer unificado: `edugo-central`

### Quick Start para Servicios

```go
// En tu servicio (api-mobile, worker, etc.)
client := auth.NewAuthClient()

result, err := client.VerifyToken(ctx, "eyJhbG...")
if result.Valid {
    fmt.Printf("Usuario: %s, Rol: %s\n", result.UserID, result.Role)
}
```

### Configuración Mínima

```env
# .env de api-administracion
AUTH_JWT_SECRET=tu-clave-secreta-de-al-menos-32-caracteres
AUTH_JWT_ISSUER=edugo-central
```

