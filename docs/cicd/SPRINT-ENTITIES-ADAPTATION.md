# Sprint: Adaptar api-administracion a Entities Centralizadas

**Proyecto:** edugo-api-administracion  
**Fecha:** 20 de Noviembre, 2025  
**Objetivo:** Migrar de entities locales a entities centralizadas de infrastructure  
**Dependencia:** infrastructure Sprint ENTITIES completado ✅  
**Prioridad:** ALTA - Parte del plan de estandarización

---

## 🎯 Contexto

**Situación actual:**
- api-administracion tiene **7 entities locales** en `internal/domain/entity/`
- Estas entities tienen **mucha lógica de negocio** (validaciones, métodos business logic)
- Se usan en **38 archivos diferentes** del proyecto
- Infrastructure ahora provee entities base (sin lógica) para PostgreSQL

**Objetivo del Sprint:**
- Reemplazar entities locales por entities de infrastructure
- Extraer lógica de negocio a Domain Services (donde aún no esté)
- Mantener backward compatibility durante la transición

---

## 📊 Inventario de Entities Actuales

### Entities Locales en api-administracion

**Ubicación:** `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/internal/domain/entity/`

| # | Archivo | Entity | Lines | Complejidad | Mapea a Infrastructure |
|---|---------|--------|-------|-------------|------------------------|
| 1 | `academic_unit.go` | `AcademicUnit` | 413 | 🔴 ALTA (mucha lógica) | `postgres/entities.AcademicUnit` |
| 2 | `guardian_relation.go` | `GuardianRelation` | 133 | 🟡 MEDIA | `postgres/entities.GuardianRelation` |
| 3 | `school.go` | `School` | 166 | 🟡 MEDIA | `postgres/entities.School` |
| 4 | `subject.go` | `Subject` | 71 | 🟢 BAJA | `postgres/entities.Subject` |
| 5 | `unit.go` | `Unit` | 138 | 🟡 MEDIA | `postgres/entities.Unit` |
| 6 | `unit_membership.go` | `UnitMembership` | 230 | 🟡 MEDIA | `postgres/entities.Membership` |
| 7 | `user.go` | `User` | 165 | 🟡 MEDIA | `postgres/entities.User` |

**Total:** 7 entities, ~1,316 líneas de código

### Características de Entities Actuales

**🔴 Problema Principal:** Entities tienen **MUCHA LÓGICA DE NEGOCIO**

Ejemplos:
- `AcademicUnit`:
  - ✅ Constructor con validaciones (`NewAcademicUnit`)
  - ✅ Métodos de negocio (`SetParent`, `AddChild`, `RemoveChild`)
  - ✅ Navegación de árbol (`GetAllDescendants`, `GetDepth`)
  - ✅ Soft delete (`SoftDelete`, `Restore`)
  - ⚠️ **Deprecated notice** - Ya identificaron que debe migrar a domain services

- `UnitMembership`:
  - Validaciones de vigencia (`IsActive`, `IsActiveAt`)
  - Permisos (`HasPermission`)
  - Cambio de rol (`ChangeRole`)

**Esta lógica debe migrar a Domain Services**, no ir a infrastructure.

---

## 🗺️ Estrategia de Migración

### Principio Fundamental

> **Infrastructure entities = Solo estructura de BD**  
> **Domain entities (local) = Lógica de negocio**

### Dos Opciones de Arquitectura

#### Opción A: Entities de Infrastructure como Base (RECOMENDADO)

```go
// Domain layer usa entities de infrastructure directamente
import pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"

// En repositories: Usar pgentities.User directamente
func (r *userRepo) Create(ctx context.Context, user *pgentities.User) error {
    return r.db.Create(user).Error
}

// Lógica de negocio va a Domain Services
type UserDomainService struct {}

func (s *UserDomainService) ValidateNewUser(user *pgentities.User) error {
    if user.FirstName == "" {
        return errors.NewValidationError("first_name is required")
    }
    // ... más validaciones
    return nil
}
```

**Ventajas:**
- ✅ Elimina duplicación de estructuras
- ✅ Cambios en BD solo en infrastructure
- ✅ Fuerza separación clara de responsabilidades

**Desventajas:**
- ⚠️ Expone entities de infrastructure en domain layer
- ⚠️ Pierde encapsulación (campos públicos)

#### Opción B: Mantener Domain Entities + DTOs

```go
// Domain layer mantiene su propia entity (sin tags de BD)
type User struct {
    id        UserID
    email     Email
    firstName string
    // ... campos privados
}

// Mapper entre domain y infrastructure
func (u *User) ToInfrastructure() *pgentities.User {
    return &pgentities.User{
        ID:        u.id.Value(),
        Email:     u.email.Value(),
        FirstName: u.firstName,
        // ...
    }
}
```

**Ventajas:**
- ✅ Mantiene encapsulación
- ✅ Domain desacoplado de infrastructure

**Desventajas:**
- ⚠️ Duplica estructuras
- ⚠️ Requiere mappers en todos lados
- ⚠️ Más código que mantener

### 🎯 Decisión para Este Sprint: **OPCIÓN A**

**Razón:** api-administracion ya tiene Domain Services que pueden absorber la lógica. Priorizar simplicidad y eliminar duplicación.

---

## 📋 Mapeo Detallado: Local → Infrastructure

### 1. User

**Actual:** `internal/domain/entity/user.go` (165 líneas)  
**Target:** `github.com/EduGoGroup/edugo-infrastructure/postgres/entities.User`

**Campos:**
```go
// Actual (privados)
id        valueobject.UserID
email     valueobject.Email
firstName string
lastName  string
role      enum.SystemRole
isActive  bool
createdAt time.Time
updatedAt time.Time

// Infrastructure (públicos, tags db)
ID        uuid.UUID `db:"id"`
Email     string    `db:"email"`
FirstName string    `db:"first_name"`
LastName  string    `db:"last_name"`
Role      string    `db:"role"`
IsActive  bool      `db:"is_active"`
CreatedAt time.Time `db:"created_at"`
UpdatedAt time.Time `db:"updated_at"`
```

**Lógica a migrar:**
- ✅ `NewUser()` → `UserDomainService.CreateUser()`
- ✅ `Deactivate()` → `UserDomainService.Deactivate()`
- ✅ `UpdateName()` → `UserDomainService.UpdateName()`
- ✅ `ChangeRole()` → `UserDomainService.ChangeRole()`
- ✅ `IsTeacher()`, `IsStudent()`, `IsGuardian()` → Helpers en service

**Archivos que usan `entity.User`:** 8 archivos
- `internal/domain/repository/user_repository.go`
- `internal/application/service/user_service.go`
- `internal/infrastructure/persistence/postgres/repository/user_repository_impl.go`
- `internal/application/dto/user_dto.go`
- Tests

### 2. School

**Actual:** `internal/domain/entity/school.go` (166 líneas)  
**Target:** `github.com/EduGoGroup/edugo-infrastructure/postgres/entities.School`

**Lógica a migrar:**
- `NewSchool()` → `SchoolDomainService.CreateSchool()`
- `UpdateInfo()` → `SchoolDomainService.UpdateInfo()`
- `UpdateContactInfo()` → `SchoolDomainService.UpdateContactInfo()`
- `SetMetadata()`, `GetMetadata()` → Métodos helper en service

**Archivos afectados:** 6 archivos

### 3. AcademicUnit

**Actual:** `internal/domain/entity/academic_unit.go` (413 líneas) 🔴 **MÁS COMPLEJO**  
**Target:** `github.com/EduGoGroup/edugo-infrastructure/postgres/entities.AcademicUnit`

**Lógica a migrar:**
- ✅ Ya existe `AcademicUnitDomainService` (domain service creado previamente)
- ✅ Entity ya tiene notice de deprecation para métodos
- Métodos como `SetParent()`, `AddChild()`, `RemoveChild()` ya están en domain service

**Campo especial:**
```go
// Actual: tiene slice de children para árbol en memoria
children []*AcademicUnit

// Infrastructure: NO debería tener esto (solo BD)
// Solución: Usar AcademicUnitTreeService para construir árbol
```

**Archivos afectados:** 10 archivos

### 4. GuardianRelation

**Actual:** `internal/domain/entity/guardian_relation.go` (133 líneas)  
**Target:** `github.com/EduGoGroup/edugo-infrastructure/postgres/entities.GuardianRelation`

**Lógica a migrar:**
- `NewGuardianRelation()` → `GuardianDomainService.CreateRelation()`
- `Deactivate()`, `Activate()` → Service methods
- `ChangeRelationshipType()` → Service method

**Archivos afectados:** 5 archivos

### 5. Subject

**Actual:** `internal/domain/entity/subject.go` (71 líneas) 🟢 **MÁS SIMPLE**  
**Target:** `github.com/EduGoGroup/edugo-infrastructure/postgres/entities.Subject`

**Lógica a migrar:**
- `NewSubject()` → `SubjectDomainService.CreateSubject()`
- `UpdateInfo()` → Service method

**Archivos afectados:** 5 archivos

### 6. Unit

**Actual:** `internal/domain/entity/unit.go` (138 líneas)  
**Target:** `github.com/EduGoGroup/edugo-infrastructure/postgres/entities.Unit`

⚠️ **NOTA:** Infrastructure tiene **dos entities** que pueden mapear:
- `Unit` (tabla simple)
- `AcademicUnit` (jerarquía compleja)

**Decisión:** Depende de qué tabla usa actualmente api-administracion. Verificar migraciones.

**Archivos afectados:** 6 archivos

### 7. UnitMembership

**Actual:** `internal/domain/entity/unit_membership.go` (230 líneas)  
**Target:** `github.com/EduGoGroup/edugo-infrastructure/postgres/entities.Membership`

**Lógica a migrar:**
- Ya existe `MembershipDomainService` (creado previamente)
- `IsActive()`, `IsActiveAt()` → Service methods
- `HasPermission()` → Mover a AuthorizationService
- `ChangeRole()`, `SetValidUntil()`, `Expire()` → Service methods

**Archivos afectados:** 8 archivos

---

## 📋 Tareas del Sprint

### Fase 0: Pre-requisitos

#### Tarea 0.1: Verificar infrastructure entities disponibles

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Verificar que Sprint ENTITIES esté completo
ls -la postgres/entities/
ls -la mongodb/entities/

# Verificar que tenga releases
git tag | grep "postgres/entities"
```

**Criterio de éxito:** 
- ✅ 14 entities PostgreSQL existen
- ✅ Tag de release disponible

---

### Fase 1: Preparar Domain Services

**Objetivo:** Crear/actualizar domain services para absorber lógica de negocio

#### Tarea 1.1: Crear UserDomainService

**Ubicación:** `internal/domain/service/user_domain_service.go`

```go
package service

import (
    pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
    "github.com/EduGoGroup/edugo-shared/common/errors"
    "github.com/google/uuid"
)

// UserDomainService contiene lógica de negocio de usuarios
type UserDomainService struct{}

// NewUserDomainService crea una nueva instancia
func NewUserDomainService() *UserDomainService {
    return &UserDomainService{}
}

// CreateUser crea un nuevo usuario con validaciones de negocio
func (s *UserDomainService) CreateUser(
    email, firstName, lastName, role string,
) (*pgentities.User, error) {
    // Validaciones de negocio (antes en entity.NewUser)
    if email == "" {
        return nil, errors.NewValidationError("email is required")
    }
    if firstName == "" {
        return nil, errors.NewValidationError("first_name is required")
    }
    if lastName == "" {
        return nil, errors.NewValidationError("last_name is required")
    }
    if role != "teacher" && role != "student" && role != "guardian" {
        return nil, errors.NewValidationError("invalid role")
    }
    if role == "admin" {
        return nil, errors.NewBusinessRuleError("cannot create admin users")
    }

    now := time.Now()
    return &pgentities.User{
        ID:        uuid.New(),
        Email:     email,
        FirstName: firstName,
        LastName:  lastName,
        Role:      role,
        IsActive:  true,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// ValidateUserUpdate valida actualizaciones de usuario
func (s *UserDomainService) ValidateUserUpdate(user *pgentities.User, firstName, lastName string) error {
    if firstName == "" || lastName == "" {
        return errors.NewValidationError("first_name and last_name are required")
    }
    return nil
}

// CanDeactivate valida si un usuario puede ser desactivado
func (s *UserDomainService) CanDeactivate(user *pgentities.User) error {
    if !user.IsActive {
        return errors.NewBusinessRuleError("user is already inactive")
    }
    return nil
}

// CanChangeRole valida cambio de rol
func (s *UserDomainService) CanChangeRole(user *pgentities.User, newRole string) error {
    if user.Role == newRole {
        return errors.NewBusinessRuleError("new role is the same as current")
    }
    if newRole == "admin" {
        return errors.NewBusinessRuleError("cannot promote to admin role")
    }
    return nil
}
```

**Archivos a crear:**
- [ ] `internal/domain/service/user_domain_service.go`
- [ ] `internal/domain/service/user_domain_service_test.go`

#### Tarea 1.2: Crear SchoolDomainService

Similar a User, crear service para School.

**Archivos a crear:**
- [ ] `internal/domain/service/school_domain_service.go`
- [ ] `internal/domain/service/school_domain_service_test.go`

#### Tarea 1.3: Crear SubjectDomainService

**Archivos a crear:**
- [ ] `internal/domain/service/subject_domain_service.go`
- [ ] `internal/domain/service/subject_domain_service_test.go`

#### Tarea 1.4: Crear GuardianDomainService

**Archivos a crear:**
- [ ] `internal/domain/service/guardian_domain_service.go`
- [ ] `internal/domain/service/guardian_domain_service_test.go`

#### Tarea 1.5: Actualizar AcademicUnitDomainService

⚠️ **Ya existe**, solo necesita actualizar para usar entities de infrastructure.

**Archivo a modificar:**
- [ ] `internal/domain/service/academic_unit_service.go`

#### Tarea 1.6: Actualizar MembershipDomainService

⚠️ **Ya existe**, solo actualizar.

**Archivo a modificar:**
- [ ] `internal/domain/service/membership_service.go`

**Criterio de éxito Fase 1:** 
- ✅ 4 nuevos domain services creados
- ✅ 2 domain services actualizados
- ✅ Tests unitarios pasan

---

### Fase 2: Actualizar go.mod

#### Tarea 2.1: Agregar infrastructure como dependencia

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Agregar infrastructure entities
go get github.com/EduGoGroup/edugo-infrastructure/postgres/entities@postgres/entities/v0.1.0

# Verificar
go mod tidy
go list -m github.com/EduGoGroup/edugo-infrastructure
```

**Criterio de éxito:** 
- ✅ `go.mod` tiene infrastructure
- ✅ `go mod tidy` sin errores

---

### Fase 3: Actualizar Repositories (Capa de Persistencia)

**Objetivo:** Repositories usan infrastructure entities directamente

#### Tarea 3.1: Actualizar UserRepository

**Archivo:** `internal/domain/repository/user_repository.go`

```go
// ANTES
import "github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id UserID) (*entity.User, error)
    // ...
}

// DESPUÉS
import pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"

type UserRepository interface {
    Create(ctx context.Context, user *pgentities.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*pgentities.User, error)
    // ...
}
```

**Archivos a modificar:**
- [ ] `internal/domain/repository/user_repository.go`
- [ ] `internal/infrastructure/persistence/postgres/repository/user_repository_impl.go`

#### Tarea 3.2: Actualizar SchoolRepository

**Archivos a modificar:**
- [ ] `internal/domain/repository/school_repository.go`
- [ ] `internal/infrastructure/persistence/postgres/repository/school_repository_impl.go`

#### Tarea 3.3: Actualizar SubjectRepository

**Archivos a modificar:**
- [ ] `internal/domain/repository/subject_repository.go`
- [ ] `internal/infrastructure/persistence/postgres/repository/subject_repository_impl.go`

#### Tarea 3.4: Actualizar AcademicUnitRepository

**Archivos a modificar:**
- [ ] `internal/domain/repository/academic_unit_repository.go`
- [ ] `internal/infrastructure/persistence/postgres/repository/academic_unit_repository_impl.go`

#### Tarea 3.5: Actualizar GuardianRepository

**Archivos a modificar:**
- [ ] `internal/domain/repository/guardian_repository.go`
- [ ] `internal/infrastructure/persistence/postgres/repository/guardian_repository_impl.go`

#### Tarea 3.6: Actualizar UnitMembershipRepository

**Archivos a modificar:**
- [ ] `internal/domain/repository/unit_membership_repository.go`
- [ ] `internal/infrastructure/persistence/postgres/repository/unit_membership_repository_impl.go`

#### Tarea 3.7: Actualizar UnitRepository (si existe)

**Archivos a modificar:**
- [ ] `internal/domain/repository/unit_repository.go`
- [ ] `internal/infrastructure/persistence/postgres/repository/unit_repository_impl.go`

**Criterio de éxito Fase 3:** 
- ✅ Todos los repositories actualizados
- ✅ Compila sin errores

---

### Fase 4: Actualizar Application Services

**Objetivo:** Application services usan domain services para lógica de negocio

#### Tarea 4.1: Actualizar UserService

**Archivo:** `internal/application/service/user_service.go`

**Cambios:**
```go
// ANTES
import "github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) error {
    user, err := entity.NewUser(req.Email, req.FirstName, req.LastName, req.Role)
    if err != nil {
        return err
    }
    return s.userRepo.Create(ctx, user)
}

// DESPUÉS
import (
    pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
    "github.com/EduGoGroup/edugo-api-administracion/internal/domain/service"
)

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) error {
    // Usar domain service para lógica de negocio
    user, err := s.userDomainService.CreateUser(req.Email, req.FirstName, req.LastName, req.Role)
    if err != nil {
        return err
    }
    return s.userRepo.Create(ctx, user)
}
```

**Archivos a modificar:**
- [ ] `internal/application/service/user_service.go`

#### Tarea 4.2: Actualizar SchoolService

**Archivos a modificar:**
- [ ] `internal/application/service/school_service.go`

#### Tarea 4.3: Actualizar SubjectService

**Archivos a modificar:**
- [ ] `internal/application/service/subject_service.go`

#### Tarea 4.4: Actualizar GuardianService

**Archivos a modificar:**
- [ ] `internal/application/service/guardian_service.go`

#### Tarea 4.5: Actualizar AcademicUnitService

**Archivos a modificar:**
- [ ] `internal/application/service/academic_unit_service.go`

#### Tarea 4.6: Actualizar UnitMembershipService

**Archivos a modificar:**
- [ ] `internal/application/service/unit_membership_service.go`

#### Tarea 4.7: Actualizar HierarchyService

**Archivos a modificar:**
- [ ] `internal/application/service/hierarchy_service.go`

#### Tarea 4.8: Actualizar UnitService (si existe)

**Archivos a modificar:**
- [ ] `internal/application/service/unit_service.go`

**Criterio de éxito Fase 4:** 
- ✅ Application services actualizados
- ✅ Usan domain services para validaciones
- ✅ Compila sin errores

---

### Fase 5: Actualizar DTOs

**Objetivo:** DTOs mapean desde/hacia infrastructure entities

#### Tarea 5.1: Actualizar UserDTO

**Archivo:** `internal/application/dto/user_dto.go`

```go
// ANTES
import "github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"

func ToUserResponse(user *entity.User) UserResponse {
    return UserResponse{
        ID:        user.ID().String(),
        Email:     user.Email().Value(),
        FirstName: user.FirstName(),
        // ...
    }
}

// DESPUÉS
import pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"

func ToUserResponse(user *pgentities.User) UserResponse {
    return UserResponse{
        ID:        user.ID.String(),
        Email:     user.Email,
        FirstName: user.FirstName,
        // ...
    }
}
```

**Archivos a modificar:**
- [ ] `internal/application/dto/user_dto.go`
- [ ] `internal/application/dto/school_dto.go`
- [ ] `internal/application/dto/subject_dto.go`
- [ ] `internal/application/dto/guardian_dto.go`
- [ ] `internal/application/dto/academic_unit_dto.go`
- [ ] `internal/application/dto/unit_membership_dto.go`
- [ ] `internal/application/dto/unit_dto.go`
- [ ] `internal/infrastructure/http/dto/school_dto.go`
- [ ] `internal/infrastructure/http/dto/academic_unit_dto.go`

**Criterio de éxito Fase 5:** 
- ✅ DTOs actualizados
- ✅ Mappers funcionan correctamente
- ✅ Compila sin errores

---

### Fase 6: Actualizar Tests

#### Tarea 6.1: Actualizar integration tests

**Archivos a modificar:**
- [ ] `test/integration/academic_unit_ltree_test.go`
- [ ] `test/integration/integration_flows_test.go`

#### Tarea 6.2: Actualizar service tests

**Archivos a modificar:**
- [ ] `internal/application/service/hierarchy_service_test.go`
- [ ] `internal/domain/service/academic_unit_service_test.go`
- [ ] `internal/domain/service/membership_service_test.go`

#### Tarea 6.3: Actualizar entity tests (ahora van a domain service tests)

**Archivos:**
- [ ] Eliminar `internal/domain/entity/academic_unit_test.go`
- [ ] Eliminar `internal/domain/entity/unit_membership_test.go`
- [ ] Mover tests de lógica a domain service tests

**Criterio de éxito Fase 6:** 
- ✅ Todos los tests actualizados
- ✅ `go test ./...` pasa

---

### Fase 7: Eliminar Entities Locales

⚠️ **SOLO después de que todo compile y tests pasen**

#### Tarea 7.1: Verificar que no quedan imports

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Buscar imports antiguos
grep -r "internal/domain/entity" --include="*.go" | grep -v "^Binary"

# Debería retornar 0 resultados
```

#### Tarea 7.2: Eliminar carpeta de entities

```bash
# Backup primero
git checkout -b backup/pre-entity-cleanup

# Eliminar entities locales
rm -rf internal/domain/entity/

# Verificar que compila
go build ./...
go test ./...
```

#### Tarea 7.3: Commit y push

```bash
git checkout dev
git add .
git commit -m "refactor: migrar a infrastructure entities

- Reemplazar entities locales por infrastructure entities
- Mover lógica de negocio a domain services
- Actualizar repositories, services, DTOs
- Eliminar entities duplicadas

BREAKING CHANGE: Domain entities ahora usan infrastructure"

git push origin dev
```

**Criterio de éxito Fase 7:** 
- ✅ Carpeta `internal/domain/entity/` eliminada
- ✅ Compila sin errores
- ✅ Tests pasan
- ✅ Commit creado

---

### Fase 8: Validación Final

#### Tarea 8.1: Validar build

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Build completo
go build -v ./...

# Tests completos
go test -v ./...

# Test de cobertura
go test -cover ./...
```

#### Tarea 8.2: Validar que API funciona

```bash
# Levantar API local
make run

# En otra terminal, ejecutar health check
curl http://localhost:8081/health

# Ejecutar integration tests
make test-integration
```

#### Tarea 8.3: Validar CI/CD

```bash
# Push a branch
git push origin dev

# Verificar que GitHub Actions/GitLab CI pase
# Revisar logs de pipeline
```

**Criterio de éxito Fase 8:** 
- ✅ Build exitoso
- ✅ Tests pasan (unit + integration)
- ✅ API levanta correctamente
- ✅ CI/CD pipeline verde

---

## 📊 Estimación de Esfuerzo

| Fase | Tareas | Complejidad | Tiempo Estimado |
|------|--------|-------------|-----------------|
| Fase 0: Pre-requisitos | 1 tarea | 🟢 Baja | 10 min |
| Fase 1: Domain Services | 6 servicios | 🟡 Media | 4-5 horas |
| Fase 2: go.mod | 1 tarea | 🟢 Baja | 10 min |
| Fase 3: Repositories | 7 repos | 🟡 Media | 2-3 horas |
| Fase 4: Application Services | 8 services | 🔴 Alta | 3-4 horas |
| Fase 5: DTOs | 9 archivos | 🟡 Media | 2 horas |
| Fase 6: Tests | ~15 archivos | 🔴 Alta | 3-4 horas |
| Fase 7: Cleanup | 3 tareas | 🟢 Baja | 30 min |
| Fase 8: Validación | 3 tareas | 🟡 Media | 1 hora |
| **TOTAL** | **48 archivos** | | **16-20 horas** |

**Estimación realista:** 2-3 días de trabajo (si se hace en sesiones de 6-8 horas)

---

## 🔗 Dependencias

**Antes de este sprint:**
- ✅ Infrastructure Sprint ENTITIES completado
- ✅ Releases de infrastructure disponibles

**Después de este sprint:**
- ➡️ api-administracion usa entities centralizadas
- ➡️ Se puede replicar proceso en api-mobile y worker

---

## ⚠️ Riesgos y Mitigaciones

### Riesgo 1: Lógica de negocio compleja en AcademicUnit

**Impacto:** Alto  
**Probabilidad:** Alta

**Mitigación:**
- ✅ Ya existe `AcademicUnitDomainService` creado previamente
- ✅ Entity tiene notice de deprecation en métodos
- Plan: Migrar métodos de árbol (`AddChild`, `GetAllDescendants`) a domain service

### Riesgo 2: Breaking changes en application services

**Impacto:** Alto  
**Probabilidad:** Media

**Mitigación:**
- Hacer cambios en branch separado
- Ejecutar tests después de cada fase
- No mergear a `main` hasta que todo funcione

### Riesgo 3: Value objects vs tipos primitivos

**Impacto:** Medio  
**Probabilidad:** Alta

**Contexto:**
- Entities actuales usan `valueobject.UserID`, `valueobject.Email`, etc.
- Infrastructure entities usan `uuid.UUID`, `string`, etc.

**Mitigación:**
- Mantener value objects en domain layer
- Convertir en bordes (repositories, DTOs)
- Ejemplo:
  ```go
  // En repository
  func (r *userRepo) FindByID(ctx context.Context, id UserID) (*pgentities.User, error) {
      var user pgentities.User
      err := r.db.Where("id = ?", id.Value()).First(&user).Error
      return &user, err
  }
  ```

### Riesgo 4: Tests flaky después de cambios

**Impacto:** Medio  
**Probabilidad:** Media

**Mitigación:**
- Actualizar mocks para usar infrastructure entities
- Usar testify/assert para comparaciones
- Revisar tests de integración con BD real

---

## 📝 Checklist de Validación Final

Antes de dar el sprint por completado:

### Build y Tests
- [ ] `go build ./...` exitoso
- [ ] `go test ./...` exitoso (100% de tests pasan)
- [ ] `make test-integration` exitoso
- [ ] No quedan imports de `internal/domain/entity`

### Código
- [ ] Carpeta `internal/domain/entity/` eliminada
- [ ] 4 nuevos domain services creados
- [ ] 7 repositories actualizados
- [ ] 8 application services actualizados
- [ ] 9 DTOs actualizados
- [ ] Tests actualizados o migrados

### Funcionalidad
- [ ] API levanta sin errores
- [ ] Endpoints de salud responden correctamente
- [ ] Tests de integración pasan
- [ ] No hay regresiones funcionales

### CI/CD
- [ ] Pipeline de GitHub Actions pasa
- [ ] Pipeline de GitLab CI pasa (si aplica)
- [ ] No hay warnings de linter

### Documentación
- [ ] README actualizado si es necesario
- [ ] Comentarios en código explicando cambios importantes
- [ ] CHANGELOG.md actualizado

---

## 📈 Métricas de Éxito

**Métricas cuantitativas:**
- ✅ 0 imports de `internal/domain/entity`
- ✅ 100% de tests pasan
- ✅ Reducción de ~1,316 líneas de código de entities duplicadas
- ✅ 1 única fuente de verdad para estructuras de BD

**Métricas cualitativas:**
- ✅ Código más mantenible
- ✅ Separación clara: estructura (infrastructure) vs lógica (domain services)
- ✅ Cambios en BD solo requieren actualizar infrastructure

---

## 🔄 Siguiente Paso

Una vez completado este sprint:

1. **Replicar en api-mobile:**
   - Crear `SPRINT-ENTITIES-ADAPTATION.md` para api-mobile
   - Seguir mismo patrón
   - api-mobile tiene más entities (Material, Assessment, Progress)

2. **Replicar en worker:**
   - Migrar entities de MongoDB
   - Worker usa `MaterialAssessment`, `MaterialSummary`, `MaterialEvent`

3. **Consolidar:**
   - Crear guía general de migración
   - Documentar lecciones aprendidas
   - Agregar ejemplos a docs/

---

## 📚 Referencias

- **Sprint de Infrastructure:** `/Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/02-infrastructure/SPRINT-ENTITIES.md`
- **Migraciones PostgreSQL:** `edugo-infrastructure/postgres/migrations/`
- **Código actual de api-administracion:** `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/`

---

## 🛠️ Comandos Útiles

### Búsqueda y Reemplazo Masivo

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Buscar todos los archivos que importan entity
grep -rl "internal/domain/entity" --include="*.go" .

# Reemplazar imports (usar con precaución - mejor hacerlo por fases)
find . -name "*.go" -exec sed -i '' 's|github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity|github.com/EduGoGroup/edugo-infrastructure/postgres/entities|g' {} +

# Verificar cambios antes de commit
git diff
```

### Verificar Dependencias

```bash
# Ver qué usa infrastructure
go list -m -json github.com/EduGoGroup/edugo-infrastructure | jq

# Ver dependencias de un paquete específico
go list -f '{{.Deps}}' ./internal/domain/repository
```

### Testing Específico

```bash
# Test de un paquete específico
go test -v ./internal/domain/repository/...

# Test con cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Test con race detector
go test -race ./...
```

---

**Generado por:** Claude Code  
**Fecha:** 20 de Noviembre, 2025  
**Sesión:** Análisis de api-administracion entities
