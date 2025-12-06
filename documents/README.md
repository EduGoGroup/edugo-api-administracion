# 📚 Documentación - EduGo API Administración

> Visión 360° del proyecto EduGo API de Administración

## 🗂️ Índice de Documentos

| Documento | Descripción |
|-----------|-------------|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Arquitectura del sistema, patrones de diseño, capas |
| [DATABASE.md](./DATABASE.md) | Modelo de datos, entidades, relaciones, diagramas ER |
| [API.md](./API.md) | Documentación completa de endpoints, request/response |
| [AUTH.md](./AUTH.md) | Sistema de autenticación centralizada, JWT, flujos |
| [SETUP.md](./SETUP.md) | Configuración, variables de entorno, servicios requeridos |
| [FLOWS.md](./FLOWS.md) | Diagramas de procesos, flujos de negocio |
| [GLOSSARY.md](./GLOSSARY.md) | Glosario de términos y conceptos del dominio |
| [improvements/](./improvements/) | Código deprecado, mejoras pendientes, deuda técnica |

---

## 🎯 Resumen Ejecutivo

**EduGo API Administración** es el servicio central de administración del ecosistema EduGo. Gestiona:

- 🏫 **Escuelas** - CRUD completo de instituciones educativas
- 🏛️ **Unidades Académicas** - Jerarquía de grados, secciones, departamentos
- 👥 **Membresías** - Asignación de usuarios a unidades con roles
- 🔐 **Autenticación Centralizada** - Servicio de auth para todo el ecosistema
- 📊 **Estadísticas** - Métricas globales del sistema

---

## 🏗️ Stack Tecnológico

| Componente | Tecnología | Versión |
|------------|------------|---------|
| **Lenguaje** | Go | 1.21+ |
| **Framework HTTP** | Gin | 1.11 |
| **Base de Datos Principal** | PostgreSQL | 15+ |
| **Base de Datos Secundaria** | MongoDB | 7.0 (logs/eventos) |
| **ORM** | GORM | 1.31 |
| **Autenticación** | JWT (HS256) | - |
| **Documentación API** | Swagger/OpenAPI | 3.0 |
| **Contenedores** | Docker | - |
| **Testing** | Testcontainers | 0.40 |

---

## 📐 Arquitectura en Alto Nivel

```
┌─────────────────────────────────────────────────────────────┐
│                        CLIENTES                              │
│  (Web Admin, API Mobile, Workers, Servicios Externos)       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    API ADMINISTRACIÓN                        │
│                      Puerto: 8081                            │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Auth    │  │  Schools │  │  Units   │  │ Members  │   │
│  │ Handler  │  │ Handler  │  │ Handler  │  │ Handler  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
├─────────────────────────────────────────────────────────────┤
│                    APPLICATION SERVICES                      │
├─────────────────────────────────────────────────────────────┤
│                    DOMAIN / REPOSITORIES                     │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  PostgreSQL  │    │   MongoDB    │    │    Redis     │
│   (Datos)    │    │   (Logs)     │    │   (Cache)    │
└──────────────┘    └──────────────┘    └──────────────┘
```

---

## 🚀 Quick Start

```bash
# 1. Clonar repositorio
git clone https://github.com/EduGoGroup/edugo-api-administracion.git

# 2. Instalar dependencias
go mod download

# 3. Configurar variables de entorno
cp .env.example .env
# Editar .env con tus valores

# 4. Ejecutar
make run
# o
go run cmd/main.go

# 5. Acceder a Swagger
open http://localhost:8081/swagger/index.html
```

---

## 🔗 Dependencias del Ecosistema

Este proyecto depende de paquetes internos de EduGo:

| Paquete | Propósito |
|---------|-----------|
| `edugo-infrastructure/postgres` | Entidades compartidas, conexión DB |
| `edugo-shared/auth` | JWT Manager compartido |
| `edugo-shared/bootstrap` | Inicialización de infraestructura |
| `edugo-shared/common` | Tipos, enums, validadores |
| `edugo-shared/logger` | Logger estructurado |
| `edugo-shared/middleware/gin` | Middlewares HTTP |
| `edugo-shared/testing` | Utilidades para tests |

---

## 📁 Estructura del Proyecto

```
edugo-api-administracion/
├── cmd/                          # Punto de entrada
│   └── main.go                   # Inicialización y router
├── config/                       # Archivos de configuración
│   ├── config.yaml               # Config base
│   └── config-{env}.yaml         # Override por ambiente
├── internal/                     # Código interno
│   ├── application/              # Capa de aplicación
│   │   ├── dto/                  # Data Transfer Objects
│   │   └── service/              # Servicios de negocio
│   ├── auth/                     # Módulo de autenticación
│   │   ├── handler/              # Handlers HTTP auth
│   │   ├── service/              # Servicios auth
│   │   └── dto/                  # DTOs auth
│   ├── bootstrap/                # Inicialización
│   ├── config/                   # Carga de configuración
│   ├── container/                # Dependency Injection
│   ├── domain/                   # Capa de dominio
│   │   └── repository/           # Interfaces de repositorios
│   └── infrastructure/           # Capa de infraestructura
│       ├── http/handler/         # Handlers HTTP
│       └── persistence/          # Implementaciones DB
├── docs/                         # Swagger generado
├── documents/                    # Esta documentación
├── test/                         # Tests
└── postman/                      # Colecciones Postman
```

---

## 📊 Estado Actual del Proyecto

### ✅ Funcionalidades Implementadas

| Módulo | Estado | Cobertura Tests | Notas |
|--------|--------|-----------------|-------|
| **Autenticación** | ✅ Completo | ~85% | Login, logout, refresh, verify |
| **Escuelas** | ✅ Completo | ~80% | CRUD completo |
| **Unidades Académicas** | ✅ Completo | ~75% | Jerarquía, árbol, CRUD |
| **Membresías** | ✅ Completo | ~70% | Asignación usuarios-unidades |
| **Verificación Tokens** | ✅ Completo | ~90% | Para servicios internos |

### 🚧 En Desarrollo

| Funcionalidad | Prioridad | Sprint |
|---------------|-----------|--------|
| Validación de permisos por rol | Alta | Sprint 5 |
| Cache con Redis | Media | Sprint 5 |
| Auditoría completa | Media | Sprint 6 |
| Métricas Prometheus | Baja | Sprint 6 |

### 📝 Deuda Técnica

Ver carpeta [improvements/](./improvements/) para:
- Código deprecado a eliminar
- Refactorizaciones pendientes
- Malas prácticas identificadas

---

## 🔄 Versionado

| Versión | Fecha | Cambios Principales |
|---------|-------|--------------------|
| v1.1.0 | 2025-12 | Auth centralizada, verify endpoint |
| v1.0.0 | 2025-11 | Release inicial, CRUD básico |
| v0.5.0 | 2025-10 | Membresías, unidades académicas |
| v0.1.0 | 2025-09 | Bootstrap, estructura inicial |

---

## 🏢 Contacto

- **Equipo**: EduGo Development Team
- **Repositorio**: github.com/EduGoGroup/edugo-api-administracion
- **Licencia**: Privado - EduGo 2025

---

## 📚 Referencias Adicionales

- **Swagger UI**: `http://localhost:8081/swagger/index.html`
- **Postman Collection**: `postman/edugo-api-admin.json`
- **Changelog**: `CHANGELOG.md` en raíz del proyecto
- **Archivos Históricos**: `archivado-documentos/`
