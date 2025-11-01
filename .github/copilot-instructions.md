# GitHub Copilot - Instrucciones Personalizadas: EduGo API Administración

## 🌍 IDIOMA / LANGUAGE

**IMPORTANTE**: Todos los comentarios, sugerencias, code reviews y respuestas en chat deben estar **SIEMPRE EN ESPAÑOL**.

- ✅ Comentarios en Pull Requests: **español**
- ✅ Sugerencias de código: **español**
- ✅ Explicaciones en chat: **español**
- ✅ Mensajes de error: **español**

---

## 🏗️ Arquitectura del Proyecto

Este proyecto implementa **Clean Architecture (Hexagonal)** con Go 1.25:

```
internal/
├── domain/              # Entidades, Value Objects, Interfaces
├── application/         # Servicios, DTOs, Casos de uso
├── infrastructure/      # Implementaciones concretas
│   ├── http/           # Handlers, Middleware
│   └── persistence/    # Repositorios (PostgreSQL, MongoDB)
├── container/          # Inyección de Dependencias
└── config/             # Configuración con Viper
```

### Principios Arquitectónicos
- **Dependency Inversion**: El dominio NO depende de infraestructura
- **Separation of Concerns**: Cada capa tiene responsabilidades claras
- **Dependency Injection**: Usar container/container.go para DI
- **Interface Segregation**: Interfaces pequeñas y específicas

---

## 📦 Dependencia Compartida: edugo-shared

Usamos el módulo `github.com/EduGoGroup/edugo-shared` para funcionalidad compartida:

### Paquetes Disponibles
- **logger**: Logger Zap estructurado (`edugo-shared/logger`)
- **common/errors**: Tipos de error de aplicación (`edugo-shared/common/errors`)

### ⚠️ REGLA CRÍTICA: NO Reimplementar Funcionalidad

```go
// ❌ INCORRECTO: Reimplementar funcionalidad existente
type MyLogger struct { ... }
func (l *MyLogger) Info(msg string) { ... }

// ✅ CORRECTO: Usar edugo-shared
import "github.com/EduGoGroup/edugo-shared/logger"
logger.Info(ctx, "mensaje de log", zap.String("key", "value"))
```

---

## 🎯 Convenciones de Código

### Naming Conventions

```go
// DTOs
type UserDTO struct { ... }          // ✅ Termina en DTO
type CreateInstitutionDTO struct { ... }  // ✅ Termina en DTO

// Servicios
type UserService struct { ... }      // ✅ Termina en Service
type InstitutionService struct { ... }  // ✅ Termina en Service

// Repositorios
type UserRepository interface { ... } // ✅ Termina en Repository
type PostgresUserRepository struct { ... } // ✅ Implementación específica

// Handlers
type UserHandler struct { ... }      // ✅ Termina en Handler
```

### Manejo de Errores

```go
// ✅ CORRECTO: Usar tipos de error de edugo-shared
import "github.com/EduGoGroup/edugo-shared/common/errors"

func (s *UserService) GetUser(ctx context.Context, id string) (*UserDTO, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.IsNotFound(err) {
            return nil, errors.NewNotFoundError("user", id)
        }
        return nil, errors.NewInternalError("failed to get user", err)
    }
    return user, nil
}

// ❌ INCORRECTO: NO usar fmt.Errorf directamente
return nil, fmt.Errorf("user not found: %s", id)

// ❌ INCORRECTO: NO usar errors.New
return nil, errors.New("user not found")
```

### Context en Todas las Funciones

```go
// ✅ CORRECTO: Siempre recibir context.Context como primer parámetro
func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDTO) (*UserDTO, error)
func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error
func (h *UserHandler) CreateUser(c *gin.Context)  // Gin ya provee context

// ❌ INCORRECTO: Métodos sin context
func (s *UserService) CreateUser(dto CreateUserDTO) (*UserDTO, error)
```

### Logging Estructurado

```go
// ✅ CORRECTO: Usar logger de edugo-shared con campos estructurados
import (
    "github.com/EduGoGroup/edugo-shared/logger"
    "go.uber.org/zap"
)

func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDTO) (*UserDTO, error) {
    logger.Info(ctx, "creating user",
        zap.String("email", dto.Email),
        zap.String("role", dto.Role),
    )

    // ... lógica ...

    if err != nil {
        logger.Error(ctx, "failed to create user",
            zap.Error(err),
            zap.String("email", dto.Email),
        )
        return nil, err
    }

    logger.Info(ctx, "user created successfully", zap.String("user_id", user.ID))
    return user, nil
}

// ❌ INCORRECTO: NO usar log estándar
log.Println("user created:", userID)
log.Printf("error: %v", err)

// ❌ INCORRECTO: NO usar fmt.Println
fmt.Println("creating user...")
```

---

## 🗄️ Bases de Datos

### PostgreSQL (Datos Relacionales)

```go
// ✅ Usar lib/pq para queries
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    query := `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`
    err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, errors.NewNotFoundError("user", id)
    }
    return &user, err
}
```

---

## ✅ Testing

### Principios de Testing

```go
// ✅ Tests de integración con testcontainers
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestUserRepository_Integration(t *testing.T) {
    // Setup: Levantar PostgreSQL container
    ctx := context.Background()
    container, err := postgres.RunContainer(ctx, ...)
    require.NoError(t, err)
    defer container.Terminate(ctx)

    // Test: Usar repositorio real
    repo := NewPostgresUserRepository(db)
    // ...

    // Cleanup: Automático con defer
}

// ✅ Tests unitarios con mocks para dependencias externas
type MockUserRepository struct {
    mock.Mock
}

// ✅ Tests deben ser independientes y ejecutarse en paralelo
func TestUserService_CreateUser(t *testing.T) {
    t.Parallel()  // ✅ Permite ejecución paralela
    // ...
}
```

### Cobertura de Tests

- **Objetivo**: >70% de cobertura
- **Prioridad**: Servicios de aplicación y repositorios

---

## 🛠️ Tecnologías y Stack

### Framework y Bibliotecas Core
- **Framework Web**: Gin Gonic
- **Config Management**: Viper
- **Logging**: Zap (via edugo-shared)
- **Database Drivers**:
  - PostgreSQL: `lib/pq`

### Testing
- **Framework**: Testing estándar de Go
- **Containers**: Testcontainers
- **Mocking**: Testify/mock

### DevOps
- **Containerización**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Registry**: GitHub Container Registry (ghcr.io)

---

## 📚 Documentación API

### Swagger/OpenAPI

```go
// ✅ CORRECTO: Agregar anotaciones Swagger en handlers
// @Summary Crear nuevo usuario
// @Description Crea un usuario en el sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserDTO true "Datos del usuario"
// @Success 201 {object} UserDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ...
}
```

### Generar Documentación

```bash
# Regenerar docs Swagger
swag init -g cmd/main.go --output docs

# Acceder a Swagger UI
# http://localhost:8080/swagger/index.html
```

---

## 🌐 Variables de Entorno

### Variables Requeridas

```bash
# Base de datos
POSTGRES_PASSWORD=<contraseña>

# Ambiente
APP_ENV=local|dev|qa|prod
```

### NO Hardcodear Secrets

```go
// ❌ INCORRECTO: Secrets hardcodeados
const dbPassword = "postgres123"

// ✅ CORRECTO: Leer de variables de entorno
dbPassword := viper.GetString("database.password")
```

---

## 🎨 Estilo de Código

### Formato

```bash
# ✅ SIEMPRE formatear con gofmt antes de commit
gofmt -w .

# ✅ Verificar con linter
golangci-lint run
```

### Comentarios

```go
// ✅ CORRECTO: Comentarios en español, explicativos
// CreateUser crea un nuevo usuario en el sistema y envía un email de bienvenida.
// Valida que el email sea único antes de crear el registro.
func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDTO) (*UserDTO, error)

// ❌ INCORRECTO: Comentarios obvios o redundantes
// CreateUser crea un usuario
func (s *UserService) CreateUser(...)
```

### Imports

```go
// ✅ CORRECTO: Agrupar imports
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // Third party
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    // Internal - edugo-shared
    "github.com/EduGoGroup/edugo-shared/logger"
    "github.com/EduGoGroup/edugo-shared/common/errors"

    // Internal - este proyecto
    "github.com/EduGoGroup/edugo-api-administracion/internal/domain"
    "github.com/EduGoGroup/edugo-api-administracion/internal/application"
)
```

---

## ⚡ Mejores Prácticas Adicionales

### 1. Inyección de Dependencias

```go
// ✅ CORRECTO: Constructor con dependencias explícitas
func NewUserService(
    repo UserRepository,
    logger logger.Logger,
) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}

// ❌ INCORRECTO: Dependencias globales o singleton
var globalDB *sql.DB  // ❌ Evitar
```

### 2. Validación de DTOs

```go
// ✅ CORRECTO: Usar validaciones explícitas
import "github.com/go-playground/validator/v10"

type CreateUserDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required,min=2"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var dto CreateUserDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(400, gin.H{"error": "invalid request body"})
        return
    }

    if err := validate.Struct(dto); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // ...
}
```

### 3. Transacciones de Base de Datos

```go
// ✅ CORRECTO: Usar transacciones para operaciones múltiples
func (s *UserService) CreateUserWithProfile(ctx context.Context, dto CreateUserDTO) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()  // Rollback automático si no hay commit

    // Operación 1
    user, err := s.userRepo.SaveTx(ctx, tx, user)
    if err != nil {
        return err
    }

    // Operación 2
    err = s.profileRepo.SaveTx(ctx, tx, profile)
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

---

## 🎓 Recursos de Referencia

- **Workflows CI/CD**: [.github/workflows/README.md](workflows/README.md)
- **CHANGELOG**: [CHANGELOG.md](../CHANGELOG.md)

---

## 📝 Notas Finales para Copilot

### Al Revisar Pull Requests

1. ✅ Verificar que se usen tipos de error de `edugo-shared`
2. ✅ Confirmar que todos los métodos reciben `context.Context`
3. ✅ Validar que se use logging estructurado
4. ✅ Señalar TODOs o funcionalidad incompleta
5. ✅ Verificar que no se reimplemente funcionalidad de `edugo-shared`

### Al Sugerir Código

1. ✅ Seguir Clean Architecture (no mezclar capas)
2. ✅ Usar dependencias de `edugo-shared` cuando corresponda
3. ✅ Incluir logging adecuado
4. ✅ Manejar errores con tipos apropiados
5. ✅ Agregar validaciones necesarias
6. ✅ Escribir código testeable

### Recordatorio de Idioma

🌍 **TODOS los comentarios, sugerencias y explicaciones deben estar en ESPAÑOL.**

---

**Última actualización**: 2025-11-01
**Versión del proyecto**: v0.1.0 (en desarrollo)
**Go Version**: 1.25.3
**edugo-shared Version**: Usar tags cuando estén disponibles
