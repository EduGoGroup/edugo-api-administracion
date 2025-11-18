# Análisis de Impacto - Refactor a Clean Architecture

**Fecha:** 2025-11-17  
**Versión:** 1.0  
**Estado:** 🟡 En Revisión

---

## 📊 Resumen Ejecutivo

Este documento analiza el impacto de migrar de **DDD Rico** (lógica en entities) a **Clean Architecture Estricta** (lógica en domain services).

### Impacto General

| Categoría | Nivel de Impacto | Archivos Afectados | Horas Estimadas |
|-----------|------------------|--------------------|--------------------|
| Domain Entities | 🔴 **Alto** | 2 archivos | 4-6h |
| Domain Services | 🟢 **Nuevo** | 5 archivos nuevos | 6-8h |
| Tests Unitarios | 🔴 **Alto** | 2 archivos | 6-8h |
| Repositorios | 🟡 **Medio** | 3-5 archivos | 3-4h |
| Application Services | 🟡 **Medio** | 2-3 archivos | 2-3h |
| Tests Integración | 🟢 **Bajo** | 3 archivos | 1-2h |
| **TOTAL** | - | **~35 archivos** | **25-30h** |

---

## 🎯 Alcance del Cambio

### 1. Entities Afectadas

#### 1.1 `AcademicUnit` Entity

**Antes (DDD Rico):**
```go
// internal/domain/entity/academic_unit.go (400+ líneas)
type AcademicUnit struct {
    id, parentUnitID, schoolID, unitType, displayName, code, description
    metadata, children, createdAt, updatedAt, deletedAt
}

// Métodos con lógica de negocio (18 métodos)
func (au *AcademicUnit) SetParent(parentID, parentType) error { /* 30 líneas validación */ }
func (au *AcademicUnit) AddChild(child) error { /* 40 líneas validación */ }
func (au *AcademicUnit) RemoveChild(childID) error { /* 15 líneas */ }
func (au *AcademicUnit) GetAllDescendants() []*AcademicUnit { /* recursivo */ }
func (au *AcademicUnit) GetDepth() int { /* recursivo */ }
func (au *AcademicUnit) UpdateInfo(name, desc) error { /* validación */ }
// ... 12 métodos más
```

**Después (Clean Architecture):**
```go
// internal/domain/entity/academic_unit.go (~150 líneas)
type AcademicUnit struct {
    id, parentUnitID, schoolID, unitType, displayName, code, description
    metadata, children, createdAt, updatedAt, deletedAt
}

// Solo getters/setters básicos (10 métodos simples)
func (au *AcademicUnit) ID() valueobject.UnitID { return au.id }
func (au *AcademicUnit) SetParentID(id valueobject.UnitID) { au.parentUnitID = &id }
// ... getters/setters básicos
```

**Cambios:**
- ❌ Eliminar: 18 métodos con lógica
- ✅ Mantener: Campos + getters/setters
- ➡️ Mover a: `AcademicUnitDomainService`

#### 1.2 `UnitMembership` Entity

**Antes:**
```go
// internal/domain/entity/unit_membership.go (300+ líneas)
func (um *UnitMembership) IsActive() bool { /* lógica temporal */ }
func (um *UnitMembership) IsActiveAt(t time.Time) bool { /* validación */ }
func (um *UnitMembership) Expire() error { /* cambio estado */ }
func (um *UnitMembership) ChangeRole(role) error { /* validación */ }
func (um *UnitMembership) HasPermission(perm) bool { /* lógica permisos */ }
// ... 8 métodos más
```

**Después:**
```go
// internal/domain/entity/unit_membership.go (~100 líneas)
// Solo getters/setters
func (um *UnitMembership) Role() valueobject.MembershipRole { return um.role }
func (um *UnitMembership) SetRole(role) { um.role = role }
```

**Cambios:**
- ❌ Eliminar: 13 métodos con lógica
- ✅ Mantener: Campos + getters/setters
- ➡️ Mover a: `MembershipDomainService`

---

### 2. Domain Services (Nuevos)

#### 2.1 `AcademicUnitDomainService`

**Archivo:** `internal/domain/service/academic_unit_service.go`

```go
package service

type AcademicUnitDomainService struct {
    // Posibles dependencias si las necesita
}

func NewAcademicUnitDomainService() *AcademicUnitDomainService {
    return &AcademicUnitDomainService{}
}

// Métodos migrados de entity
func (s *AcademicUnitDomainService) SetParent(
    unit *entity.AcademicUnit,
    parentID valueobject.UnitID,
    parentType valueobject.UnitType,
) error {
    // Toda la lógica de validación que estaba en entity.SetParent()
}

func (s *AcademicUnitDomainService) AddChild(
    parent *entity.AcademicUnit,
    child *entity.AcademicUnit,
) error {
    // Toda la lógica que estaba en entity.AddChild()
}

func (s *AcademicUnitDomainService) GetAllDescendants(
    unit *entity.AcademicUnit,
) []*entity.AcademicUnit {
    // Lógica recursiva que estaba en entity.GetAllDescendants()
}

func (s *AcademicUnitDomainService) GetDepth(
    unit *entity.AcademicUnit,
) int {
    // Lógica que estaba en entity.GetDepth()
}

// ... resto de métodos
```

**Tamaño estimado:** ~300-350 líneas

#### 2.2 `MembershipDomainService`

**Archivo:** `internal/domain/service/membership_service.go`

```go
package service

type MembershipDomainService struct {}

func NewMembershipDomainService() *MembershipDomainService {
    return &MembershipDomainService{}
}

func (s *MembershipDomainService) IsActive(
    membership *entity.UnitMembership,
) bool {
    // Lógica que estaba en entity.IsActive()
}

func (s *MembershipDomainService) IsActiveAt(
    membership *entity.UnitMembership,
    t time.Time,
) bool {
    // Lógica que estaba en entity.IsActiveAt()
}

func (s *MembershipDomainService) HasPermission(
    membership *entity.UnitMembership,
    permission string,
) bool {
    // Lógica que estaba en entity.HasPermission()
}

func (s *MembershipDomainService) Expire(
    membership *entity.UnitMembership,
) error {
    // Lógica que estaba en entity.Expire()
}

// ... resto de métodos
```

**Tamaño estimado:** ~250 líneas

---

### 3. Tests Impactados

#### 3.1 Tests de Entity (Refactorizar)

**Antes:**
- `internal/domain/entity/academic_unit_test.go` (656 líneas)
- `internal/domain/entity/unit_membership_test.go` (490 líneas)

**Después:**
- `internal/domain/entity/academic_unit_test.go` (~150 líneas)
  - Solo tests de getters/setters
  - Tests de constructor y reconstruct
  
- `internal/domain/entity/unit_membership_test.go` (~100 líneas)
  - Solo tests básicos

#### 3.2 Tests de Service (Crear)

**Nuevos archivos:**
- `internal/domain/service/academic_unit_service_test.go` (~500 líneas)
  - Migrar todos los tests de lógica de negocio
  - Tests de SetParent, AddChild, GetDescendants, GetDepth
  
- `internal/domain/service/membership_service_test.go` (~400 líneas)
  - Tests de IsActive, HasPermission, Expire, etc.

**Trabajo:**
1. Copiar tests de entity_test.go
2. Adaptar para usar services
3. Validar misma cobertura

---

### 4. Repositorios Impactados

#### Archivos a Modificar:

**1. `internal/infrastructure/persistence/postgres/repository/academic_unit_repository_impl.go`**

**Antes:**
```go
func (r *repository) Create(ctx context.Context, unit *entity.AcademicUnit) error {
    // Usa unit.ID(), unit.ParentUnitID(), etc.
}
```

**Después:**
```go
// Sin cambios significativos - los getters siguen existiendo
func (r *repository) Create(ctx context.Context, unit *entity.AcademicUnit) error {
    // Mismo código
}
```

**Impacto:** 🟢 **Mínimo** - Los getters no cambian

**2. `internal/infrastructure/persistence/postgres/repository/unit_membership_repository_impl.go`**

**Impacto:** 🟢 **Mínimo** - Similar a academic_unit

---

### 5. Application Services Impactados

#### Archivos a Modificar:

**1. `internal/application/service/academic_unit_application_service.go`**

**Antes:**
```go
func (s *AcademicUnitApplicationService) CreateUnit(...) error {
    unit, err := entity.NewAcademicUnit(...)
    if parentID != nil {
        err = unit.SetParent(parentID, parentType) // ❌ Ya no existe
    }
    return s.repo.Create(ctx, unit)
}
```

**Después:**
```go
func (s *AcademicUnitApplicationService) CreateUnit(...) error {
    unit, err := entity.NewAcademicUnit(...)
    if parentID != nil {
        // ✅ Usar domain service
        err = s.domainService.SetParent(unit, parentID, parentType)
    }
    return s.repo.Create(ctx, unit)
}
```

**Cambios necesarios:**
1. Inyectar `AcademicUnitDomainService` en constructor
2. Cambiar llamadas de `unit.Method()` a `service.Method(unit, ...)`

**Impacto:** 🟡 **Medio** - Refactorizar ~10-15 llamadas

---

## 📁 Matriz de Archivos Impactados

| Archivo | Tipo de Cambio | Líneas Antes | Líneas Después | Esfuerzo |
|---------|----------------|--------------|----------------|----------|
| `internal/domain/entity/academic_unit.go` | ✏️ Simplificar | 400 | 150 | 4h |
| `internal/domain/entity/unit_membership.go` | ✏️ Simplificar | 300 | 100 | 3h |
| `internal/domain/service/academic_unit_service.go` | ➕ Crear | 0 | 350 | 5h |
| `internal/domain/service/membership_service.go` | ➕ Crear | 0 | 250 | 4h |
| `internal/domain/entity/academic_unit_test.go` | ✏️ Reducir | 656 | 150 | 3h |
| `internal/domain/entity/unit_membership_test.go` | ✏️ Reducir | 490 | 100 | 2h |
| `internal/domain/service/academic_unit_service_test.go` | ➕ Crear | 0 | 500 | 6h |
| `internal/domain/service/membership_service_test.go` | ➕ Crear | 0 | 400 | 5h |
| `internal/application/service/*.go` | ✏️ Refactor | varies | varies | 3h |
| `test/integration/*.go` | ✏️ Menor | varies | varies | 2h |
| **TOTAL** | - | **~2000** | **~2300** | **~30h** |

---

## ⚠️ Riesgos Detallados

### 1. Riesgo: Invariantes Rotas

**Descripción:** Al exponer setters públicos, código externo puede romper invariantes.

**Ejemplo:**
```go
// ❌ PELIGRO - Ahora es posible:
unit.SetParentID(someID)  // Sin validar si el tipo es compatible!
```

**Mitigación:**
1. Documentar claramente que setters NO deben usarse directamente
2. Hacer setters package-private cuando sea posible
3. Forzar uso de domain services en code reviews
4. Tests de integración que validen flujos completos

**Nivel:** 🔴 **Alto**

### 2. Riesgo: Performance

**Descripción:** Llamadas adicionales a services pueden impactar performance.

**Antes:**
```go
unit.AddChild(child)  // 1 llamada
```

**Después:**
```go
service.AddChild(parent, child)  // 1 llamada + paso de parámetros
```

**Mitigación:**
- Benchmarks antes/después
- Go inline optimizations manejan esto bien

**Nivel:** 🟢 **Bajo**

### 3. Riesgo: Complejidad del Código

**Descripción:** Más archivos y niveles de indirección.

**Impacto:**
- +5 archivos nuevos (services + tests)
- Desarrolladores deben conocer ambas capas

**Mitigación:**
- Documentación clara
- Ejemplos de uso en README
- Diagramas de arquitectura

**Nivel:** 🟡 **Medio**

---

## 📊 Impacto en Coverage

### Cobertura Actual

```
internal/domain/entity/           48.2%  (excluido en .coverignore)
Cobertura total proyecto:         13.2%
```

### Cobertura Esperada Post-Refactor

```
internal/domain/entity/           ~90%  (solo getters/setters)
internal/domain/service/          ~85%  (lógica de negocio testeada)
Cobertura total proyecto:         ~35-40%
```

**Acción necesaria:** Actualizar `.coverignore`:

```diff
- # Entities de dominio (solo structs, sin lógica)
- # NOTA: Si se agrega lógica de negocio a entities, remover esta exclusión
- internal/domain/entity/

+ # Entities de dominio (anemic - solo datos)
+ # La lógica está en domain services
```

---

## 🔄 Dependencias del Cambio

### Bloqueadores

- ❌ **Ninguno** - Podemos empezar inmediatamente

### Dependencias

1. **PR #28 actual**: Debe mergearse o cerrarse primero
2. **Branch limpia**: Crear nueva branch desde main

### Recomendación

**Opción A (Recomendada):**
1. Aplicar fix rápido al PR #28 (quitar entity de .coverignore)
2. Merge PR #28
3. Iniciar refactor en nuevo PR

**Opción B:**
1. Cerrar PR #28
2. Hacer refactor completo
3. Crear nuevo PR con refactor

---

## 📈 Métricas de Validación

### Pre-Refactor (Baseline)

- [ ] Todos los tests pasando
- [ ] Coverage: 13.2% total, 48.2% domain/entity
- [ ] 0 archivos en domain/service/
- [ ] Build exitoso
- [ ] Lint sin errores

### Post-Refactor (Objetivo)

- [ ] Todos los tests pasando
- [ ] Coverage: >35% total, ~85% domain/service
- [ ] 5+ archivos en domain/service/
- [ ] Build exitoso
- [ ] Lint sin errores
- [ ] Performance similar (±5%)

---

## 👥 Comunicación

### Stakeholders a Notificar

1. **Tech Lead**: Aprobar inicio de refactor
2. **Equipo Dev**: Comunicar que no toquen entities/services
3. **QA**: Tests de regresión después del merge

### Plan de Comunicación

```
Día 0: Presentar plan a tech lead
Día 1: Slack announcement - "Starting refactor"
Día 5: Status update - "50% complete"
Día 10: Status update - "Ready for review"
Día 12: Demo + QA validation
```

---

**Aprobación:**
- [ ] Revisado por: _________________
- [ ] Aprobado por: _________________
- [ ] Fecha: _________________

---

**Siguiente documento:** [WORK_PLAN.md](WORK_PLAN.md)
