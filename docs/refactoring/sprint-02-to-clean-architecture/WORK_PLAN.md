# Plan de Trabajo - Refactor a Clean Architecture

**Duración Total:** 25-30 horas  
**Fases:** 5  
**Metodología:** Incremental (cada fase es deployable)

---

## 📅 Cronograma General

```
┌──────────────────────────────────────────────────────────────┐
│ FASE 1: Setup & Domain Services Base          │ 6-8h │ Día 1-2  │
├──────────────────────────────────────────────────────────────┤
│ FASE 2: Migrar Entities a Anemic             │ 4-6h │ Día 3-4  │
├──────────────────────────────────────────────────────────────┤
│ FASE 3: Migrar Tests                          │ 6-8h │ Día 5-6  │
├──────────────────────────────────────────────────────────────┤
│ FASE 4: Actualizar Application Layer          │ 4-5h │ Día 7-8  │
├──────────────────────────────────────────────────────────────┤
│ FASE 5: Validación y Limpieza                 │ 3-4h │ Día 9-10 │
└──────────────────────────────────────────────────────────────┘
```

---

## 🎯 FASE 1: Setup & Domain Services Base (6-8h)

### Objetivos
- Crear estructura de domain services
- Implementar `AcademicUnitDomainService` completo
- Implementar `MembershipDomainService` completo
- Tests básicos de services

### Tareas

#### 1.1 Crear estructura de directorios
```bash
mkdir -p internal/domain/service
touch internal/domain/service/academic_unit_service.go
touch internal/domain/service/membership_service.go
touch internal/domain/service/academic_unit_service_test.go
touch internal/domain/service/membership_service_test.go
```
**Tiempo:** 15 min

#### 1.2 Implementar `AcademicUnitDomainService`
**Archivo:** `internal/domain/service/academic_unit_service.go`

**Métodos a implementar:**
```go
type AcademicUnitDomainService struct {}

// Métodos migrados de entity.AcademicUnit
- SetParent(unit, parentID, parentType) error
- RemoveParent(unit) error  
- AddChild(parent, child) error
- RemoveChild(parent, childID) error
- GetAllDescendants(unit) []*entity.AcademicUnit
- GetDepth(unit) int
- UpdateInfo(unit, name, desc) error
- UpdateDisplayName(unit, name) error
- CanHaveChildren(unit) bool
- HasChildren(unit) bool
- IsRoot(unit) bool
- SoftDelete(unit) error
- Restore(unit) error
- Validate(unit) error
```

**Estrategia:**
1. Copiar código de entity
2. Cambiar firma: `func (au *AcademicUnit) Method()` → `func (s *Service) Method(au *AcademicUnit)`
3. Cambiar acceso a campos: `au.field` → `au.field` (getters si privado)

**Tiempo:** 3-4h

#### 1.3 Implementar `MembershipDomainService`
**Archivo:** `internal/domain/service/membership_service.go`

**Métodos:**
```go
type MembershipDomainService struct {}

- IsActive(membership) bool
- IsActiveAt(membership, time) bool
- SetValidUntil(membership, time) error
- ExtendIndefinitely(membership)
- Expire(membership) error
- ChangeRole(membership, role) error
- HasPermission(membership, permission) bool
- Validate(membership) error
```

**Tiempo:** 2-3h

#### 1.4 Tests básicos de services
- Copiar tests de entity_test.go
- Adaptar para usar services
- Verificar que pasan

**Tiempo:** 1-2h

### Entregables Fase 1
- [ ] `academic_unit_service.go` completo
- [ ] `membership_service.go` completo
- [ ] Tests básicos pasando
- [ ] Commit: `feat(domain): add domain services for business logic`

---

## 🎯 FASE 2: Migrar Entities a Anemic (4-6h)

### Objetivos
- Simplificar entities a solo datos + getters/setters
- Mantener compatibilidad temporal (deprecated methods)
- Preparar para migración gradual

### Tareas

#### 2.1 Simplificar `AcademicUnit`

**Estrategia: Deprecation Gradual**

1. **Marcar métodos como deprecated:**
```go
// Deprecated: Use AcademicUnitDomainService.SetParent instead
func (au *AcademicUnit) SetParent(parentID, parentType) error {
    // Delegar al service por ahora
    service := service.NewAcademicUnitDomainService()
    return service.SetParent(au, parentID, parentType)
}
```

2. **Agregar getters/setters básicos:**
```go
func (au *AcademicUnit) SetParentID(id valueobject.UnitID) {
    au.parentUnitID = &id
    au.updatedAt = time.Now()
}

func (au *AcademicUnit) SetDisplayName(name string) {
    au.displayName = name
    au.updatedAt = time.Now()
}
```

**Tiempo:** 2-3h

#### 2.2 Simplificar `UnitMembership`

Similar a AcademicUnit:
- Deprecar métodos con lógica
- Agregar setters básicos
- Delegar a service

**Tiempo:** 1-2h

#### 2.3 Actualizar constructores

Mantener constructores con validaciones básicas:
```go
func NewAcademicUnit(...) (*AcademicUnit, error) {
    // Validaciones mínimas
    if schoolID.IsZero() {
        return nil, errors.New("school_id required")
    }
    // ... crear entity simple
}
```

**Tiempo:** 1h

### Entregables Fase 2
- [ ] Entities simplificadas con deprecated methods
- [ ] Getters/setters agregados
- [ ] Constructores actualizados
- [ ] Todos los tests siguen pasando
- [ ] Commit: `refactor(domain): simplify entities to anemic model`

---

## 🎯 FASE 3: Migrar Tests (6-8h)

### Objetivos
- Migrar tests de entity a service
- Reducir tests de entity a solo getters/setters
- Mantener misma cobertura

### Tareas

#### 3.1 Migrar tests de `AcademicUnit`

**De:**
- `internal/domain/entity/academic_unit_test.go` (656 líneas)

**A:**
- `internal/domain/entity/academic_unit_test.go` (150 líneas) - getters/setters
- `internal/domain/service/academic_unit_service_test.go` (500 líneas) - lógica

**Proceso:**
1. Copiar tests de lógica de negocio
2. Adaptar para usar service:
```go
// Antes
unit.SetParent(parentID, parentType)

// Después  
service := service.NewAcademicUnitDomainService()
service.SetParent(unit, parentID, parentType)
```
3. Validar cobertura equivalente

**Tiempo:** 3-4h

#### 3.2 Migrar tests de `UnitMembership`

Similar a 3.1

**Tiempo:** 2-3h

#### 3.3 Validar cobertura

```bash
# Antes
go test ./internal/domain/entity/... -cover
# Coverage: 48.2%

# Después
go test ./internal/domain/service/... -cover
# Coverage: ~85%
```

**Tiempo:** 1h

### Entregables Fase 3
- [ ] Tests migrados a service
- [ ] Cobertura mantenida o mejorada
- [ ] Todos los tests pasando
- [ ] Commit: `test(domain): migrate tests to domain services`

---

## 🎯 FASE 4: Actualizar Application Layer (4-5h)

### Objetivos
- Inyectar domain services en application services
- Cambiar llamadas de entity methods a service methods
- Actualizar dependency injection

### Tareas

#### 4.1 Actualizar Application Services

**Archivos a modificar:**
- `internal/application/service/academic_unit_application_service.go`
- `internal/application/service/membership_application_service.go`

**Cambios:**

```go
// Antes
type AcademicUnitApplicationService struct {
    repo repository.AcademicUnitRepository
}

// Después
type AcademicUnitApplicationService struct {
    repo          repository.AcademicUnitRepository
    domainService *service.AcademicUnitDomainService  // ✅ Nuevo
}

func NewAcademicUnitApplicationService(
    repo repository.AcademicUnitRepository,
    domainService *service.AcademicUnitDomainService,  // ✅ Inyectar
) *AcademicUnitApplicationService {
    return &AcademicUnitApplicationService{
        repo: repo,
        domainService: domainService,
    }
}

// Cambiar llamadas
func (s *AcademicUnitApplicationService) CreateUnit(...) error {
    unit := entity.NewAcademicUnit(...)
    
    if parentID != nil {
        // Antes: unit.SetParent(parentID, parentType)
        // Después:
        err := s.domainService.SetParent(unit, parentID, parentType)
        if err != nil {
            return err
        }
    }
    
    return s.repo.Create(ctx, unit)
}
```

**Tiempo:** 3-4h

#### 4.2 Actualizar Dependency Injection Container

**Archivo:** `internal/container/container.go`

```go
// Registrar domain services
container.Register("academicUnitDomainService", func() interface{} {
    return service.NewAcademicUnitDomainService()
})

// Inyectar en application services
container.Register("academicUnitApplicationService", func() interface{} {
    return application.NewAcademicUnitApplicationService(
        container.Get("academicUnitRepository"),
        container.Get("academicUnitDomainService"),  // ✅
    )
})
```

**Tiempo:** 1h

### Entregables Fase 4
- [ ] Application services actualizados
- [ ] Dependency injection configurado
- [ ] Integration tests pasando
- [ ] Commit: `refactor(app): use domain services in application layer`

---

## 🎯 FASE 5: Validación y Limpieza (3-4h)

### Objetivos
- Validar que todo funciona end-to-end
- Eliminar código deprecated
- Actualizar documentación
- Preparar PR

### Tareas

#### 5.1 Validación Completa

```bash
# Unit tests
make test-unit

# Integration tests
make test-integration

# Coverage
make coverage-report

# Lint
make lint

# Build
make build
```

**Criterios de aceptación:**
- [ ] Todos los tests pasando
- [ ] Coverage >= 35%
- [ ] Lint sin errores
- [ ] Build exitoso

**Tiempo:** 1h

#### 5.2 Eliminar Deprecated Methods

Una vez validado todo, remover los métodos marcados como `@Deprecated`:

```go
// ❌ Eliminar esto:
// Deprecated: Use AcademicUnitDomainService.SetParent instead
func (au *AcademicUnit) SetParent(...) error { ... }
```

**Tiempo:** 1h

#### 5.3 Actualizar `.coverignore`

```diff
- # Entities de dominio (solo structs, sin lógica)
- # NOTA: Si se agrega lógica de negocio a entities, remover esta exclusión
- internal/domain/entity/

+ # Entities de dominio (anemic - solo datos + getters/setters)
+ # La lógica de negocio está en internal/domain/service/
+ # internal/domain/entity/  ← Ya no se excluye
```

**Tiempo:** 5 min

#### 5.4 Actualizar Documentación

**Archivos a actualizar:**
- `README.md` - Explicar nueva arquitectura
- `docs/ARCHITECTURE.md` - Diagramas actualizados
- `CONTRIBUTING.md` - Guidelines para developers

**Contenido:**
```markdown
## Domain Layer Architecture

### Entities (Anemic Model)
Entities contienen solo datos y getters/setters básicos:
- `internal/domain/entity/academic_unit.go`
- `internal/domain/entity/unit_membership.go`

### Domain Services  
La lógica de negocio está en domain services:
- `internal/domain/service/academic_unit_service.go`
- `internal/domain/service/membership_service.go`

### Usage Example
```go
// ✅ Correcto
service := service.NewAcademicUnitDomainService()
err := service.SetParent(unit, parentID, parentType)

// ❌ Incorrecto (setters directos rompen invariantes)
unit.SetParentID(parentID)  // Sin validación!
```
```

**Tiempo:** 1h

#### 5.5 Preparar PR

**PR Title:** `refactor: migrate from rich DDD to clean architecture with domain services`

**PR Description:**
```markdown
## 🎯 Objetivo
Migrar de DDD Rico (lógica en entities) a Clean Architecture (lógica en services)

## 📊 Cambios
- Entities simplificadas a modelo anémico
- Nueva capa de Domain Services con lógica de negocio
- Tests migrados y cobertura mejorada

## 📈 Métricas
- Coverage: 13.2% → 35%+
- Archivos modificados: 35
- Tests: Todos pasando

## ✅ Checklist
- [x] Unit tests pasando
- [x] Integration tests pasando
- [x] Coverage >= 35%
- [x] Lint sin errores
- [x] Documentación actualizada
```

**Tiempo:** 30min

### Entregables Fase 5
- [ ] Validación completa exitosa
- [ ] Deprecated code eliminado
- [ ] Documentación actualizada
- [ ] PR creado y listo para review
- [ ] Commit: `docs: update documentation for clean architecture`

---

## 📊 Tracking de Progreso

Ver [PROGRESS_TRACKING.md](PROGRESS_TRACKING.md) para seguimiento detallado.

---

## 🚨 Plan de Rollback

Si algo sale mal en cualquier fase:

1. **Git revert** al commit anterior a la fase
2. **Documentar** qué falló
3. **Ajustar** plan
4. **Reintentar** o **abortar**

---

**Próximo documento:** [TARGET_ARCHITECTURE.md](TARGET_ARCHITECTURE.md)
