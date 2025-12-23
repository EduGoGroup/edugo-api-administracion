# ⚙️ Configuración y Setup

> Guía completa para configurar y ejecutar EduGo API Administración

## 📋 Requisitos Previos

### Software Requerido

| Software | Versión Mínima | Propósito |
|----------|----------------|-----------|
| **Go** | 1.21+ | Lenguaje de programación |
| **Docker** | 20.10+ | Contenedores (tests, BD local) |
| **PostgreSQL** | 15+ | Base de datos principal |
| **MongoDB** | 7.0+ | Logs y eventos |
| **Git** | 2.30+ | Control de versiones |
| **Make** | 3.81+ | Automatización de comandos |

### Acceso a Repositorios Privados

Este proyecto depende de paquetes privados de EduGo:

```bash
# Configurar GOPRIVATE
export GOPRIVATE=github.com/EduGoGroup/*

# Configurar git para usar SSH o token
git config --global url."git@github.com:".insteadOf "https://github.com/"
# o con token
git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
```

---

## 🚀 Instalación Paso a Paso

### 1. Clonar Repositorio

```bash
git clone git@github.com:EduGoGroup/edugo-api-administracion.git
cd edugo-api-administracion
```

### 2. Instalar Dependencias

```bash
go mod download
```

### 3. Configurar Variables de Entorno

```bash
# Copiar archivo de ejemplo
cp .env.example .env

# Editar con tus valores
vim .env
```

### 4. Levantar Servicios de Infraestructura

**Opción A: Docker Compose (recomendado para desarrollo)**
```bash
# Levantar PostgreSQL y MongoDB
docker-compose up -d postgres mongodb

# Verificar que estén corriendo
docker-compose ps
```

**Opción B: Servicios locales**
```bash
# PostgreSQL
brew services start postgresql@15

# MongoDB
brew services start mongodb-community@7.0
```

### 5. Ejecutar la Aplicación

```bash
# Opción 1: Con Make
make run

# Opción 2: Directo con Go
go run cmd/main.go

# Opción 3: Con hot-reload (air)
air
```

### 6. Verificar Funcionamiento

```bash
# Health check
curl http://localhost:8081/health

# Swagger UI
open http://localhost:8081/swagger/index.html
```

---

## 🌍 Variables de Entorno

### Completas

```env
# ============================================
# AMBIENTE
# ============================================
APP_ENV=local                    # local | dev | qa | prod

# ============================================
# SERVIDOR
# ============================================
SERVER_PORT=8081                 # Puerto HTTP
SERVER_HOST=0.0.0.0              # Host de escucha

# ============================================
# POSTGRESQL
# ============================================
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DATABASE=edugo
POSTGRES_USER=edugo_user
POSTGRES_PASSWORD=your-secure-password
POSTGRES_MAX_CONNECTIONS=25
POSTGRES_SSL_MODE=disable        # disable | require | verify-full

# ============================================
# MONGODB
# ============================================
MONGODB_URI=mongodb://localhost:27017/edugo
MONGODB_DATABASE=edugo

# ============================================
# AUTENTICACIÓN JWT
# ============================================
AUTH_JWT_SECRET=your-production-secret-minimum-32-characters-long
AUTH_JWT_ISSUER=edugo-central
AUTH_JWT_ACCESS_TOKEN_DURATION=15m
AUTH_JWT_REFRESH_TOKEN_DURATION=168h

# ============================================
# RATE LIMITING
# ============================================
AUTH_RATE_LIMIT_LOGIN_ATTEMPTS=5
AUTH_RATE_LIMIT_LOGIN_WINDOW=15m
AUTH_RATE_LIMIT_LOGIN_BLOCK=1h

# ============================================
# SERVICIOS INTERNOS
# ============================================
AUTH_INTERNAL_SERVICES_API_KEYS=api-mobile:key1,worker:key2
AUTH_INTERNAL_SERVICES_IP_RANGES=127.0.0.1/32,10.0.0.0/8

# ============================================
# CACHE
# ============================================
AUTH_CACHE_TOKEN_VALIDATION_TTL=60s
AUTH_CACHE_USER_INFO_TTL=300s

# ============================================
# REDIS (opcional)
# ============================================
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# ============================================
# LOGGING
# ============================================
LOGGING_LEVEL=info               # debug | info | warn | error
LOGGING_FORMAT=json              # json | text

# ============================================
# DESARROLLO
# ============================================
USE_MOCK_REPOSITORIES=false      # true para tests sin DB
```

### Por Ambiente

| Variable | Local | Dev | QA | Prod |
|----------|-------|-----|----|----- |
| `APP_ENV` | local | dev | qa | prod |
| `POSTGRES_SSL_MODE` | disable | require | require | verify-full |
| `LOGGING_LEVEL` | debug | info | info | warn |
| `USE_MOCK_REPOSITORIES` | false | false | false | false |

---

## 📁 Archivos de Configuración

### Estructura

```
config/
├── config.yaml          # Configuración base (todos los ambientes)
├── config-local.yaml    # Override para desarrollo local
├── config-dev.yaml      # Override para desarrollo
├── config-qa.yaml       # Override para QA
├── config-prod.yaml     # Override para producción
└── config-test.yaml     # Override para tests
```

### Prioridad de Configuración

1. Variables de entorno (máxima prioridad)
2. `config-{APP_ENV}.yaml`
3. `config.yaml` (base)

### Ejemplo config-local.yaml

```yaml
server:
  port: 8081
  host: "localhost"

database:
  use_mock_repositories: false
  postgres:
    host: "localhost"
    port: 5432
    database: "edugo"
    user: "edugo_user"
    ssl_mode: "disable"

logging:
  level: "debug"
  format: "text"

auth:
  jwt:
    # Para desarrollo local, usar un secret simple
    secret: "local-development-secret-change-in-production-min-32-chars"
```

---

## 🐳 Docker

### Dockerfile

El proyecto incluye un Dockerfile multi-stage optimizado:

```dockerfile
# Build stage
FROM golang:alpine AS builder
# ... build ...

# Final stage
FROM alpine:latest
EXPOSE 8081
CMD ["./main"]
```

### Docker Compose

```yaml
version: '3.8'

services:
  api-administracion:
    build: .
    container_name: edugo-api-admin
    ports:
      - "8081:8081"
    environment:
      APP_ENV: ${APP_ENV:-local}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    networks:
      - edugo-network

networks:
  edugo-network:
    external: true
```

### Comandos Docker

```bash
# Build imagen
docker build -t edugo-api-admin:latest .

# Build con token de GitHub (para dependencias privadas)
docker build --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t edugo-api-admin:latest .

# Ejecutar contenedor
docker run -p 8081:8081 --env-file .env edugo-api-admin:latest

# Con docker-compose
docker-compose up -d
docker-compose logs -f api-administracion
docker-compose down
```

---

## 🛠️ Comandos Make

```bash
# Ver todos los comandos disponibles
make help

# ============================================
# BUILD
# ============================================
make build              # Compilar binario
make build-linux        # Compilar para Linux

# ============================================
# RUN
# ============================================
make run                # Ejecutar aplicación
make dev                # Ejecutar con hot-reload (air)

# ============================================
# TESTS
# ============================================
make test               # Tests unitarios
make test-integration   # Tests de integración (requiere Docker)
make test-all           # Todos los tests
make coverage           # Generar reporte de cobertura

# ============================================
# QUALITY
# ============================================
make lint               # Linting con golangci-lint
make fmt                # Formatear código
make vet                # Go vet

# ============================================
# SWAGGER
# ============================================
make swagger            # Generar documentación Swagger
make swagger-validate   # Validar spec Swagger

# ============================================
# DATABASE
# ============================================
make db-migrate         # Ejecutar migraciones
make db-seed            # Seed de datos iniciales

# ============================================
# DOCKER
# ============================================
make docker-build       # Build imagen Docker
make docker-push        # Push a registry
make docker-run         # Run contenedor local
```

---

## 🧪 Testing

### Tests Unitarios

```bash
# Ejecutar todos
go test ./internal/... -v

# Con cobertura
go test ./internal/... -cover -coverprofile=coverage.out

# Ver reporte HTML
go tool cover -html=coverage.out
```

### Tests de Integración

Requieren Docker para levantar contenedores de prueba:

```bash
# Ejecutar tests de integración
go test ./test/integration/... -v -tags=integration

# Con testcontainers (levanta PostgreSQL y MongoDB automáticamente)
make test-integration
```

### Mocks

Para tests sin base de datos real:

```bash
# Activar mocks en .env
USE_MOCK_REPOSITORIES=true

# O por variable de entorno
USE_MOCK_REPOSITORIES=true go test ./...
```

---

## 📊 Servicios Externos Requeridos

### 1. PostgreSQL

**Desarrollo local:**
```bash
# macOS
brew install postgresql@15
brew services start postgresql@15

# Crear base de datos
createdb edugo
createuser edugo_user -P
psql -c "GRANT ALL PRIVILEGES ON DATABASE edugo TO edugo_user;"
```

**Docker:**
```bash
docker run -d \
  --name postgres-edugo \
  -e POSTGRES_USER=edugo_user \
  -e POSTGRES_PASSWORD=edugo_pass \
  -e POSTGRES_DB=edugo \
  -p 5432:5432 \
  postgres:15-alpine
```

### 2. MongoDB

**Desarrollo local:**
```bash
# macOS
brew tap mongodb/brew
brew install mongodb-community@7.0
brew services start mongodb-community@7.0
```

**Docker:**
```bash
docker run -d \
  --name mongodb-edugo \
  -e MONGO_INITDB_ROOT_USERNAME=edugo_admin \
  -e MONGO_INITDB_ROOT_PASSWORD=edugo_pass \
  -p 27017:27017 \
  mongo:7.0
```

### 3. Redis (Opcional - para cache)

```bash
docker run -d \
  --name redis-edugo \
  -p 6379:6379 \
  redis:7-alpine
```

---

## 🔍 Troubleshooting

### Error: "cannot find package github.com/EduGoGroup/..."

```bash
# Configurar acceso a repos privados
export GOPRIVATE=github.com/EduGoGroup/*
git config --global url."git@github.com:".insteadOf "https://github.com/"
go mod download
```

### Error: "connection refused" a PostgreSQL

```bash
# Verificar que PostgreSQL esté corriendo
docker ps | grep postgres
# o
brew services list | grep postgresql

# Verificar credenciales en .env
echo $POSTGRES_HOST $POSTGRES_PORT $POSTGRES_USER
```

### Error: "JWT_SECRET no está configurado"

```bash
# Asegurar que AUTH_JWT_SECRET tenga al menos 32 caracteres
export AUTH_JWT_SECRET="your-production-secret-minimum-32-characters-long"
```

### Puerto 8081 ocupado

```bash
# Encontrar proceso
lsof -i :8081

# Matar proceso
kill -9 <PID>

# O cambiar puerto en .env
SERVER_PORT=8082
```

---

## 📈 Monitoreo y Logs

### Formato de Logs

**JSON (producción):**
```json
{"level":"info","msg":"Servidor escuchando","port":8081,"time":"2025-12-06T10:30:00Z"}
```

**Text (desarrollo):**
```
INFO[2025-12-06T10:30:00Z] Servidor escuchando port=8081
```

### Health Check

```bash
# Endpoint de salud
curl http://localhost:8081/health

# Response esperado
{"status":"healthy","service":"edugo-api-admin"}
```

### Métricas (próximamente)

```
GET /metrics  # Prometheus metrics
```
