# 🏗️ Arquitectura del Sistema

> Documentación técnica de la arquitectura de EduGo API Administración

## 📐 Patrón Arquitectónico: Clean Architecture

Este proyecto implementa **Clean Architecture** (también conocida como Onion Architecture o Hexagonal Architecture), que garantiza:

- **Independencia de frameworks** - La lógica de negocio no depende de Gin, GORM, etc.
- **Testabilidad** - Cada capa puede testearse de forma aislada
- **Independencia de UI** - Los handlers HTTP pueden reemplazarse
- **Independencia de DB** - PostgreSQL puede cambiarse sin afectar el dominio

---

## 🧅 Capas de la Aplicación

```
┌────────────────────────────────────────────────────────────────────┐
│                                                                     │
│                     INFRASTRUCTURE LAYER                            │
│         (HTTP Handlers, PostgreSQL Repos, External Services)        │
│                                                                     │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│                     APPLICATION LAYER                               │
│              (Services, DTOs, Use Cases, Orchestration)             │
│                                                                     │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│                       DOMAIN LAYER                                  │
│         (Entities, Repository Interfaces, Business Rules)           │
│                                                                     │
└────────────────────────────────────────────────────────────────────┘

         Regla de Dependencia: Las capas internas NO conocen las externas
                    Domain ← Application ← Infrastructure
```

---

## 📁 Estructura Detallada por Capa

### 1. Domain Layer (`internal/domain/`)

La capa más interna. Define **qué** hace el sistema sin saber **cómo**.

```
internal/domain/
└── repository/                    # Interfaces (contratos)
    ├── school_repository.go       # Interface SchoolRepository
    ├── academic_unit_repository.go # Interface AcademicUnitRepository
    ├── user_repository.go         # Interface UserRepository
    ├── unit_membership_repository.go
    ├── subject_repository.go
    ├── material_repository.go
    ├── stats_repository.go
    └── guardian_repository.go
```

**Características:**
- Solo define **interfaces** (contratos)
- No tiene dependencias externas
- Las entidades vienen de `edugo-infrastructure/postgres/entities`

**Ejemplo de Interface:**
```go
type SchoolRepository interface {
    Create(ctx context.Context, school *entities.School) error
    FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error)
    FindByCode(ctx context.Context, code string) (*entities.School, error)
    Update(ctx context.Context, school *entities.School) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, filters ListFilters) ([]*entities.School, error)
}
```

---

### 2. Application Layer (`internal/application/`)

Orquesta la lógica de negocio. Conoce el dominio pero no la infraestructura.

```
internal/application/
├── dto/                           # Data Transfer Objects
│   ├── school_dto.go              # CreateSchoolRequest, SchoolResponse
│   ├── academic_unit_dto.go       # DTOs para unidades
│   ├── user_dto.go                # CreateUserRequest, UserResponse
│   ├── unit_membership_dto.go     # DTOs para membresías
│   ├── subject_dto.go
│   ├── guardian_dto.go
│   └── stats_dto.go
│
└── service/                       # Servicios de aplicación
    ├── school_service.go          # SchoolService (CRUD escuelas)
    ├── academic_unit_service.go   # AcademicUnitService (unidades)
    ├── user_service.go            # UserService (usuarios)
    ├── unit_membership_service.go # UnitMembershipService
    ├── hierarchy_service.go       # HierarchyService (árbol)
    ├── subject_service.go
    ├── material_service.go
    ├── stats_service.go
    ├── guardian_service.go
    └── *_test.go                  # Tests unitarios
```

**Características:**
- Implementa casos de uso
- Transforma DTOs ↔ Entities
- Valida reglas de negocio
- No conoce HTTP, SQL, etc.

**Ejemplo de Service:**
```go
type SchoolService interface {
    CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error)
    GetSchool(ctx context.Context, id uuid.UUID) (*dto.SchoolResponse, error)
    ListSchools(ctx context.Context, filters repository.ListFilters) ([]dto.SchoolResponse, error)
    UpdateSchool(ctx context.Context, id uuid.UUID, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error)
    DeleteSchool(ctx context.Context, id uuid.UUID) error
}
```

---

### 3. Infrastructure Layer (`internal/infrastructure/`)

Implementa los detalles técnicos: HTTP, bases de datos, servicios externos.

```
internal/infrastructure/
├── http/
│   ├── handler/                   # Handlers HTTP (Gin)
│   │   ├── school_handler.go      # SchoolHandler
│   │   ├── academic_unit_handler.go
│   │   ├── unit_membership_handler.go
│   │   ├── user_handler.go
│   │   ├── subject_handler.go
│   │   ├── material_handler.go
│   │   ├── stats_handler.go
│   │   ├── guardian_handler.go
│   │   └── *_test.go
│   ├── dto/                       # DTOs específicos HTTP
│   └── router/                    # Configuración de rutas
│
└── persistence/
    ├── postgres/
    │   └── repository/            # Implementaciones PostgreSQL
    │       ├── school_repository.go
    │       ├── academic_unit_repository.go
    │       ├── user_repository.go
    │       └── ...
    └── mock/                      # Mocks para testing
        ├── school_repository_mock.go
        └── ...
```

**Características:**
- Implementa interfaces del dominio
- Conoce frameworks específicos (Gin, GORM)
- Traduce errores de DB a errores de dominio

---

## 🔌 Dependency Injection (`internal/container/`)

El contenedor centraliza la creación e inyección de dependencias:

```go
type Container struct {
    // Infrastructure
    DB         *sql.DB
    Logger     logger.Logger
    JWTManager *auth.JWTManager

    // Repositories
    SchoolRepository         repository.SchoolRepository
    AcademicUnitRepository   repository.AcademicUnitRepository
    UserRepository           repository.UserRepository
    UnitMembershipRepository repository.UnitMembershipRepository

    // Services
    SchoolService         service.SchoolService
    AcademicUnitService   service.AcademicUnitService
    UserService           service.UserService

    // Handlers
    SchoolHandler         *handler.SchoolHandler
    AcademicUnitHandler   *handler.AcademicUnitHandler
    UserHandler           *handler.UserHandler
}
```

**Factory Pattern para Repositorios:**
```go
// Decidir entre Mock o PostgreSQL según configuración
if cfg.Database.UseMockRepositories {
    repositoryFactory = factory.NewMockRepositoryFactory()
} else {
    repositoryFactory = factory.NewPostgresRepositoryFactory(db)
}

c.SchoolRepository = repositoryFactory.CreateSchoolRepository()
```

---

## 🔄 Flujo de Request

```
                                    Request HTTP
                                         │
                                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      GIN ROUTER                              │
│  • Parsea URL, headers, body                                │
│  • Ejecuta middlewares (JWT, logging, CORS)                 │
└─────────────────────────────────────────────────────────────┘
                                         │
                                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    HTTP HANDLER                              │
│  • Extrae parámetros y body                                 │
│  • Valida input (binding)                                   │
│  • Llama al Service                                         │
│  • Formatea respuesta HTTP                                  │
└─────────────────────────────────────────────────────────────┘
                                         │
                                         ▼
┌─────────────────────────────────────────────────────────────┐
│                 APPLICATION SERVICE                          │
│  • Ejecuta lógica de negocio                                │
│  • Valida reglas de negocio                                 │
│  • Orquesta repositorios                                    │
│  • Transforma DTOs ↔ Entities                               │
└─────────────────────────────────────────────────────────────┘
                                         │
                                         ▼
┌─────────────────────────────────────────────────────────────┐
│                     REPOSITORY                               │
│  • Ejecuta queries SQL                                      │
│  • Mapea resultados a Entities                              │
│  • Maneja errores de DB                                     │
└─────────────────────────────────────────────────────────────┘
                                         │
                                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    DATABASE                                  │
│  PostgreSQL / MongoDB                                        │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔐 Módulo de Autenticación (`internal/auth/`)

Módulo independiente para autenticación centralizada:

```
internal/auth/
├── dto/
│   ├── login_dto.go         # LoginRequest, LoginResponse
│   ├── verify_dto.go        # VerifyTokenRequest, VerifyTokenResponse
│   └── error_dto.go         # ErrorResponse
├── handler/
│   ├── auth_handler.go      # Login, Logout, Refresh
│   └── verify_handler.go    # Verificación para otros servicios
├── service/
│   ├── auth_service.go      # AuthService (login/logout)
│   └── token_service.go     # TokenService (validación JWT)
├── middleware/
│   └── jwt_middleware.go    # Middleware de autenticación
└── repository/
    └── blacklist_repository.go  # Tokens revocados
```

---

## 🏭 Patrones de Diseño Utilizados

| Patrón | Uso en el Proyecto |
|--------|-------------------|
| **Repository** | Abstracción de acceso a datos |
| **Service/Use Case** | Lógica de negocio encapsulada |
| **Factory** | Creación de repositorios (mock/real) |
| **Dependency Injection** | Container centralizado |
| **DTO** | Transferencia de datos entre capas |
| **Middleware** | Auth, logging, error handling |

---

## 📊 Diagrama de Componentes

```
┌─────────────────────────────────────────────────────────────────────┐
│                         cmd/main.go                                  │
│  • Carga configuración                                              │
│  • Inicializa bootstrap                                             │
│  • Crea Container                                                   │
│  • Configura Router                                                 │
│  • Inicia servidor HTTP                                             │
└─────────────────────────────────────────────────────────────────────┘
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   config/   │      │  container/ │      │  bootstrap/ │
│   Config    │─────▶│  Container  │◀─────│  Resources  │
│   Loader    │      │             │      │             │
└─────────────┘      └─────────────┘      └─────────────┘
                            │
           ┌────────────────┼────────────────┐
           ▼                ▼                ▼
    ┌────────────┐   ┌────────────┐   ┌────────────┐
    │  Handlers  │   │  Services  │   │   Repos    │
    │  (HTTP)    │──▶│  (Logic)   │──▶│  (Data)    │
    └────────────┘   └────────────┘   └────────────┘
```

---

## 🧪 Estrategia de Testing

| Nivel | Ubicación | Herramientas |
|-------|-----------|--------------|
| **Unit Tests** | `internal/**/service/*_test.go` | testify, mocks |
| **Integration Tests** | `test/integration/` | testcontainers |
| **Handler Tests** | `internal/**/handler/*_test.go` | httptest, gin |

**Mocks:**
- Los mocks están en `internal/infrastructure/persistence/mock/`
- Se pueden activar con `USE_MOCK_REPOSITORIES=true`
- Útiles para tests unitarios y desarrollo local

---

## 🔗 Dependencias Externas

```go
// Paquetes internos EduGo
github.com/EduGoGroup/edugo-infrastructure/postgres      // Entidades, DB
github.com/EduGoGroup/edugo-shared/auth                 // JWT compartido
github.com/EduGoGroup/edugo-shared/bootstrap            // Inicialización
github.com/EduGoGroup/edugo-shared/common               // Tipos comunes
github.com/EduGoGroup/edugo-shared/logger               // Logging
github.com/EduGoGroup/edugo-shared/middleware/gin       // Middlewares

// Frameworks principales
github.com/gin-gonic/gin           // HTTP framework
gorm.io/gorm                       // ORM
github.com/golang-jwt/jwt/v5       // JWT
go.mongodb.org/mongo-driver        // MongoDB driver
```
