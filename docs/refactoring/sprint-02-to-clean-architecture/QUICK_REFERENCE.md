# Referencia Rápida - Refactor Clean Architecture

**Guía de consulta rápida para la implementación**

---

## 🚀 Quick Start

### 1. Crear Branch
```bash
git checkout main
git pull origin main
git checkout -b refactor/clean-architecture-domain-services
```

### 2. Ejecutar Fase por Fase
```bash
# Ver plan detallado
cat docs/refactoring/sprint-02-to-clean-architecture/WORK_PLAN.md

# Fase 1: Domain Services
mkdir -p internal/domain/service
# ... seguir WORK_PLAN.md FASE 1

# Validar después de cada fase
make test-unit
make lint
```

---

## 📋 Comandos Útiles

### Tests y Coverage
```bash
# Tests unitarios
go test ./internal/domain/entity/... -v -cover

# Tests de services
go test ./internal/domain/service/... -v -cover

# Coverage completo
make coverage-report

# Ver coverage en browser
open coverage/coverage.html

# Verificar umbral
./scripts/check-coverage.sh coverage/coverage-filtered.out 33
```

### Build y Lint
```bash
# Compilar
make build

# Lint
make lint

# Vet
go vet ./...

# Format
gofmt -w .
```

---

## 🔧 Patterns de Migración

### Pattern 1: Método Simple

**Antes (Entity):**
```go
func (au *AcademicUnit) UpdateDisplayName(name string) error {
    if name == "" {
        return errors.NewValidationError("display_name is required")
    }
    au.displayName = name
    au.updatedAt = time.Now()
    return nil
}
```

**Después (Service):**
```go
func (s *AcademicUnitDomainService) UpdateDisplayName(
    unit *entity.AcademicUnit,
    name string,
) error {
    if name == "" {
        return errors.NewValidationError("display_name is required")
    }
    unit.SetDisplayName(name)
    unit.SetUpdatedAt(time.Now())
    return nil
}
```

**Y en Entity:**
```go
func (au *AcademicUnit) SetDisplayName(name string) {
    au.displayName = name
}
```

---

### Pattern 2: Método con Acceso a Otros Campos

**Antes (Entity):**
```go
func (au *AcademicUnit) SetParent(parentID, parentType) error {
    // Valida contra au.id, au.unitType
    if au.id.Equals(parentID) {
        return errors.New("cannot be own parent")
    }
    au.parentUnitID = &parentID
    return nil
}
```

**Después (Service):**
```go
func (s *AcademicUnitDomainService) SetParent(
    unit *entity.AcademicUnit,
    parentID valueobject.UnitID,
    parentType valueobject.UnitType,
) error {
    // Usa getters públicos
    if unit.ID().Equals(parentID) {
        return errors.New("cannot be own parent")
    }
    unit.SetParentID(parentID)
    return nil
}
```

---

### Pattern 3: Método Recursivo

**Antes (Entity):**
```go
func (au *AcademicUnit) GetAllDescendants() []*AcademicUnit {
    descendants := make([]*AcademicUnit, 0)
    for _, child := range au.children {
        descendants = append(descendants, child)
        // Recursión en entity
        descendants = append(descendants, child.GetAllDescendants()...)
    }
    return descendants
}
```

**Después (Service):**
```go
func (s *AcademicUnitDomainService) GetAllDescendants(
    unit *entity.AcademicUnit,
) []*entity.AcademicUnit {
    descendants := make([]*entity.AcademicUnit, 0)
    for _, child := range unit.Children() {
        descendants = append(descendants, child)
        // Recursión en service
        descendants = append(descendants, s.GetAllDescendants(child)...)
    }
    return descendants
}
```

---

### Pattern 4: Migrar Test

**Antes (Entity Test):**
```go
func TestAcademicUnit_SetParent(t *testing.T) {
    unit := entity.NewAcademicUnit(...)
    err := unit.SetParent(parentID, parentType)
    assert.NoError(t, err)
}
```

**Después (Service Test):**
```go
func TestAcademicUnitService_SetParent(t *testing.T) {
    service := service.NewAcademicUnitDomainService()
    unit := entity.NewAcademicUnit(...)
    err := service.SetParent(unit, parentID, parentType)
    assert.NoError(t, err)
}
```

---

## ⚡ Checklist Rápido por Fase

### FASE 1
```
□ mkdir internal/domain/service
□ Copiar lógica de entity → service  
□ Cambiar firmas: (au *Entity) → (s *Service, au *Entity)
□ Tests básicos
□ make test-unit ✅
```

### FASE 2  
```
□ Marcar entity methods como @deprecated
□ Agregar setters básicos
□ Entity methods delegan a service (temporal)
□ make test-unit ✅
```

### FASE 3
```
□ Copiar entity_test.go → service_test.go
□ Adaptar tests para usar service
□ Reducir entity_test.go a getters/setters
□ make test-unit ✅
□ Coverage >= baseline
```

### FASE 4
```
□ Inyectar domainService en appService
□ Cambiar unit.Method() → service.Method(unit)
□ Actualizar container
□ make test-integration ✅
```

### FASE 5
```
□ Validación completa
□ Eliminar @deprecated
□ Actualizar .coverignore
□ Docs
□ PR ready ✅
```

---

## 🐛 Troubleshooting

### "Tests no compilan después de cambiar entity"
```bash
# Ver qué tests usan el método eliminado
grep -r "\.SetParent" internal/application/
grep -r "\.AddChild" internal/application/

# Actualizar a usar service
```

### "Coverage bajó después del refactor"
```bash
# Ver qué NO está cubierto
make coverage-report
open coverage/coverage.html

# Agregar tests faltantes a service_test.go
```

### "Import cycle detected"
```
❌ domain/service imports domain/entity  
❌ domain/entity imports domain/service

✅ Correcto:
- domain/service → domain/entity (OK)
- domain/entity → NADA (OK)
```

---

## 📚 Ejemplos de Código

### Crear Service Instance
```go
// En application service
type AppService struct {
    domainService *service.AcademicUnitDomainService
    repo          repository.AcademicUnitRepository
}

func NewAppService(
    domainService *service.AcademicUnitDomainService,
    repo repository.AcademicUnitRepository,
) *AppService {
    return &AppService{
        domainService: domainService,
        repo: repo,
    }
}
```

### Usar Service
```go
func (s *AppService) CreateWithParent(...) error {
    // 1. Crear entity
    unit := entity.NewAcademicUnit(...)
    
    // 2. Aplicar lógica con service
    err := s.domainService.SetParent(unit, parentID, parentType)
    if err != nil {
        return err
    }
    
    // 3. Persistir
    return s.repo.Create(ctx, unit)
}
```

---

## ⏱️ Time Tracking Template

```markdown
### YYYY-MM-DD - Sesión de Trabajo

**Inicio:** HH:MM  
**Fin:** HH:MM  
**Duración:** Xh

**Tareas Completadas:**
- [x] Tarea 1
- [x] Tarea 2

**Progreso:**
- FASE X: Y% → Z%

**Blockers:**
- Ninguno / Descripción

**Próxima Sesión:**
- Continuar con Tarea 3
```

---

## 🎯 Definition of Done (DoD)

Una tarea está **DONE** cuando:

- ✅ Código implementado
- ✅ Tests escritos y pasando
- ✅ Coverage mantenida/mejorada
- ✅ Lint sin errores
- ✅ Commit realizado
- ✅ Documentado en PROGRESS_TRACKING.md
- ✅ Peer review (si aplica)

---

**Para más detalles, ver:**
- [WORK_PLAN.md](WORK_PLAN.md) - Plan detallado
- [VALIDATION_CHECKLIST.md](VALIDATION_CHECKLIST.md) - Validaciones
- [PROGRESS_TRACKING.md](PROGRESS_TRACKING.md) - Este documento
