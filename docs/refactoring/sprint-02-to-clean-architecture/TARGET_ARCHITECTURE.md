# Arquitectura Objetivo - Clean Architecture

**Versión:** 1.0  
**Fecha:** 2025-11-17

---

## 🏛️ Diagrama de Capas

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Handlers Layer                       │
│              (gin, middleware, DTOs request/response)        │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                 Application Services Layer                   │
│           (orchestration, use cases, transactions)           │
│                                                              │
│  - AcademicUnitApplicationService                           │
│  - MembershipApplicationService                             │
│                                                              │
│  Responsibilities:                                           │
│  • Orchestrate domain services                              │
│  • Handle transactions                                       │
│  • Map DTOs to/from domain                                  │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Domain Services Layer ⭐ NEW              │
│                  (business logic, validation)                │
│                                                              │
│  - AcademicUnitDomainService                                │
│  - MembershipDomainService                                  │
│                                                              │
│  Responsibilities:                                           │
│  • Business rules validation                                 │
│  • Complex domain logic                                      │
│  • Tree operations                                           │
│  • Permission validation                                     │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Domain Entities Layer                     │
│                  (data + simple getters/setters)             │
│                                                              │
│  - AcademicUnit (anemic)                                    │
│  - UnitMembership (anemic)                                  │
│  - Value Objects                                             │
│                                                              │
│  Responsibilities:                                           │
│  • Hold domain data                                          │
│  • Simple getters/setters                                    │
│  • No business logic                                         │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                     Repository Layer                         │
│                  (persistence abstraction)                   │
│                                                              │
│  - AcademicUnitRepository                                   │
│  - MembershipRepository                                     │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Infrastructure Layer                       │
│              (PostgreSQL, MongoDB, external APIs)            │
└─────────────────────────────────────────────────────────────┘
```

---

## 📁 Estructura de Directorios

```
internal/
├── domain/
│   ├── entity/                      ← Anemic models
│   │   ├── academic_unit.go         (150 líneas - solo datos)
│   │   └── unit_membership.go       (100 líneas - solo datos)
│   ├── service/                     ⭐ NEW
│   │   ├── academic_unit_service.go (350 líneas - lógica)
│   │   └── membership_service.go    (250 líneas - lógica)
│   ├── valueobject/
│   └── repository/                  (interfaces)
├── application/
│   └── service/
│       ├── academic_unit_application_service.go  (usa domain service)
│       └── membership_application_service.go     (usa domain service)
└── infrastructure/
    └── persistence/
        └── postgres/
            └── repository/
                ├── academic_unit_repository_impl.go
                └── membership_repository_impl.go
```

---

## 🔄 Flujo de Ejemplo: Crear Unidad con Padre

### Antes (DDD Rico)

```go
// Handler
func (h *Handler) CreateUnit(c *gin.Context) {
    // 1. Parse request
    var req CreateUnitRequest
    c.BindJSON(&req)
    
    // 2. Call application service
    unit, err := h.appService.CreateUnit(req)
}

// Application Service
func (s *AppService) CreateUnit(req) (*entity.AcademicUnit, error) {
    // 3. Create entity
    unit := entity.NewAcademicUnit(...)
    
    // 4. Set parent (lógica en entity) ❌
    if req.ParentID != nil {
        err := unit.SetParent(req.ParentID, parentType)
    }
    
    // 5. Save
    return s.repo.Create(ctx, unit)
}

// Entity (con lógica)
func (au *AcademicUnit) SetParent(parentID, parentType) error {
    // ❌ 30 líneas de validación aquí
    if !parentType.CanHaveChildren() { ... }
    if au.id.Equals(parentID) { ... }
    // ...
}
```

### Después (Clean Architecture)

```go
// Handler (sin cambios)
func (h *Handler) CreateUnit(c *gin.Context) {
    var req CreateUnitRequest
    c.BindJSON(&req)
    unit, err := h.appService.CreateUnit(req)
}

// Application Service
func (s *AppService) CreateUnit(req) (*entity.AcademicUnit, error) {
    // 3. Create entity (simple, sin lógica)
    unit := entity.NewAcademicUnit(...)
    
    // 4. Use domain service para lógica ✅
    if req.ParentID != nil {
        err := s.domainService.SetParent(unit, req.ParentID, parentType)
    }
    
    // 5. Save
    return s.repo.Create(ctx, unit)
}

// Domain Service (con toda la lógica)
func (s *DomainService) SetParent(unit, parentID, parentType) error {
    // ✅ 30 líneas de validación aquí
    if !parentType.CanHaveChildren() { ... }
    if unit.ID().Equals(parentID) { ... }
    
    // Modificar entity
    unit.SetParentID(parentID)
    unit.SetUpdatedAt(time.Now())
}

// Entity (anemic, sin lógica)
type AcademicUnit struct {
    id, parentUnitID, ...
}

func (au *AcademicUnit) SetParentID(id UnitID) {
    au.parentUnitID = &id
}
```

---

## ✅ Principios Aplicados

### 1. Separación de Responsabilidades (SRP)
- **Entity**: Solo datos
- **Domain Service**: Solo lógica de negocio
- **Application Service**: Solo orquestación

### 2. Dependency Inversion (DIP)
```go
// Application Service depende de abstracción
type ApplicationService struct {
    domainService DomainServiceInterface  // abstracción
    repo          RepositoryInterface     // abstracción
}
```

### 3. Open/Closed
- Fácil agregar nuevos services sin modificar entities
- Entities cerradas a modificación, abiertas a extensión

---

## 📊 Comparación

| Aspecto | DDD Rico | Clean Architecture |
|---------|----------|-------------------|
| **Lógica en Entity** | ✅ Sí | ❌ No |
| **Domain Services** | ❌ No | ✅ Sí |
| **Testabilidad** | 🟡 Media | ✅ Alta |
| **Complejidad Entity** | 🔴 Alta (400 LOC) | 🟢 Baja (150 LOC) |
| **Archivos** | 🟢 Menos | 🟡 Más |
| **Invariantes** | ✅ Protegidas | ⚠️ Manual |
| **Uncle Bob Approved** | ❌ No | ✅ Sí |

---

**Ver también:** [WORK_PLAN.md](WORK_PLAN.md) para implementación
