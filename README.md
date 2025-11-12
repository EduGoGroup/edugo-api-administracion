# EduGo API Administración

API REST para operaciones administrativas y CRUD en la plataforma EduGo.

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

# Logging
LOGGING_LEVEL=info
LOGGING_FORMAT=json
```

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

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/v1/users` | Crear usuario |
| PATCH | `/v1/users/:id` | Actualizar usuario |
| DELETE | `/v1/users/:id` | Eliminar usuario |
| POST | `/v1/schools` | Crear escuela |
| POST | `/v1/units` | Crear unidad académica |
| PATCH | `/v1/units/:id` | Actualizar unidad |
| POST | `/v1/units/:id/members` | Asignar membresía |
| POST | `/v1/subjects` | Crear materia |
| DELETE | `/v1/materials/:id` | Eliminar material |
| GET | `/v1/stats/global` | Estadísticas globales |
| GET | `/health` | Health check |

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
