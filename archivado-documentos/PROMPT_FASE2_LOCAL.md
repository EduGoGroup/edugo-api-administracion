# PROMPT FASE 2 - CLAUDE CODE LOCAL

**Proyecto:** edugo-api-administracion  
**Sprint:** Sprint-03 - Repositorios con ltree  
**Ejecutor:** Claude Code Local  
**Duración estimada:** 1-2 horas  
**Branch:** `feature/sprint-03-repositorios-ltree` (continuación)

---

## 🎯 TU OBJETIVO (Fase 2 - Con Docker)

Completar el trabajo que Claude Code Web dejó en Fase 1, específicamente:
1. **Quitar stubs/mocks**
2. **Ejecutar migraciones**
3. **Ejecutar y validar tests de integración**
4. **Crear PR a dev**
5. **Monitorear CI/CD**
6. **Hacer merge**

---

## 📋 PREREQUISITOS

Antes de empezar, lee:
1. `HANDOFF_FASE1_A_FASE2.md` - Qué hizo Claude Web
2. `PROMPT_FASE1_WEB.md` - Contexto de Fase 1

---

## 📋 TAREAS FASE 2

### TASK-01: Revisar Trabajo de Fase 1 (15min)

```bash
# 1. Checkout de la branch
git checkout feature/sprint-03-repositorios-ltree
git pull origin feature/sprint-03-repositorios-ltree

# 2. Verificar que compila
go build ./...

# 3. Leer handoff
cat HANDOFF_FASE1_A_FASE2.md
```

**Validar que existe:**
- ✅ Migraciones ltree (up + down)
- ✅ Repository con métodos ltree
- ✅ Tests estructurados (con t.Skip)

---

### TASK-02: Ejecutar Migraciones (15min)

**Opción A: Con testcontainer (recomendado)**
```bash
# Las migraciones se ejecutan automáticamente en setupTestDB()
# Solo verifica que test/integration/main_test.go incluye la migración
```

**Opción B: Manualmente**
```bash
# Levantar PostgreSQL
docker run -d --name edugo-test -e POSTGRES_PASSWORD=test -p 5432:5432 postgres:15-alpine

# Conectar y crear extension
psql -h localhost -U postgres -d postgres -c "CREATE EXTENSION ltree;"

# Ejecutar migración
psql -h localhost -U postgres -d postgres < migrations/005_add_ltree_to_academic_units.up.sql
```

**✅ Validación:**
```sql
-- Verificar que path existe
SELECT column_name, data_type FROM information_schema.columns 
WHERE table_name = 'academic_units' AND column_name = 'path';

-- Verificar que trigger funciona
INSERT INTO academic_units (...) VALUES (...);
SELECT id, path FROM academic_units;
```

---

### TASK-03: Descomentar y Ejecutar Tests (45min)

**Archivo:** `test/integration/academic_unit_ltree_test.go`

**Para cada test:**

1. **Quitar `t.Skip()`**
2. **Descomentar código**
3. **Ejecutar**

```go
// ANTES (Fase 1):
func TestFindChildren(t *testing.T) {
    t.Skip("STUB_FASE2: Requiere PostgreSQL")
    // db, cleanup := getTestDB(t)
    // ...
}

// DESPUÉS (Fase 2):
func TestFindChildren(t *testing.T) {
    db, cleanup := getTestDB(t)
    defer cleanup()
    
    ctx := context.Background()
    repo := repository.NewPostgresAcademicUnitRepository(db)
    
    // Test implementación...
}
```

**Ejecutar:**
```bash
go test -tags=integration ./test/integration/... -v -run TestAcademicUnit
```

**Criterio de éxito:**
- ✅ Todos los tests pasan
- ✅ Sin `t.Skip()` en el código
- ✅ Coverage repository >= 80%

---

### TASK-04: Benchmarks de Performance (20min)

**Archivo:** `test/integration/academic_unit_benchmark_test.go`

```go
//go:build integration

package integration

import (
    "testing"
)

func BenchmarkFindDescendants_LTree(b *testing.B) {
    db, cleanup := getTestDB(&testing.T{})
    defer cleanup()
    
    repo := repository.NewPostgresAcademicUnitRepository(db)
    
    // Setup: Crear jerarquía de 100 nodos
    // ...
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        repo.FindDescendants(ctx, rootID)
    }
}

func BenchmarkGetAllDescendants_InMemory(b *testing.B) {
    // Comparar con método recursivo en memoria
    service := service.NewAcademicUnitDomainService()
    
    // Setup: Jerarquía en memoria
    // ...
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.GetAllDescendants(root)
    }
}
```

**Ejecutar:**
```bash
go test -tags=integration ./test/integration/... -bench=. -benchmem
```

**Documentar resultados** en el PR.

---

### TASK-05: Actualizar Repository Interface (10min)

**Archivo:** `internal/domain/repository/academic_unit_repository.go`

Agregar interfaces de los nuevos métodos:

```go
type AcademicUnitRepository interface {
    // Métodos existentes...
    
    // Métodos ltree (nuevos)
    FindByPath(ctx context.Context, path string) (*entity.AcademicUnit, error)
    FindChildren(ctx context.Context, parentID valueobject.UnitID) ([]*entity.AcademicUnit, error)
    FindDescendants(ctx context.Context, unitID valueobject.UnitID) ([]*entity.AcademicUnit, error)
    FindAncestors(ctx context.Context, unitID valueobject.UnitID) ([]*entity.AcademicUnit, error)
    FindBySchoolIDAndDepth(ctx context.Context, schoolID valueobject.SchoolID, depth int) ([]*entity.AcademicUnit, error)
    MoveSubtree(ctx context.Context, unitID valueobject.UnitID, newParentID *valueobject.UnitID) error
}
```

---

### TASK-06: Validación Completa (15min)

```bash
# Compilación
go build ./...

# Tests unitarios (deben seguir pasando)
go test ./internal/domain/... -v

# Tests de integración (ahora DEBEN pasar)
go test -tags=integration ./test/integration/... -v

# Coverage
make coverage-report

# Lint
make lint
```

**Criterios:**
- ✅ Compilación OK
- ✅ Tests unitarios: 100% pasando
- ✅ Tests integración: 100% pasando
- ✅ Coverage >= 15% (threshold actual)
- ✅ Lint sin errores

---

### TASK-07: Crear PR a dev (10min)

**Título:** `feat(infrastructure): implement ltree-based tree operations in repositories`

**Descripción:**

```markdown
## 🎯 Sprint-03: Repositorios con ltree

Implementación de operaciones de árbol optimizadas usando PostgreSQL ltree.

### ✅ Completado

**Migraciones:**
- Extension ltree habilitada
- Columna path agregada
- Índices GIST y BTREE creados
- Triggers automáticos para path

**Repository Methods:**
- FindByPath, FindChildren, FindDescendants
- FindAncestors, FindBySchoolIDAndDepth
- MoveSubtree con transacciones

**Tests:**
- Integration tests >= 80% coverage
- Benchmarks ltree vs recursión
- Todos los tests pasando ✅

### 📊 Performance

Benchmark results:
- ltree FindDescendants: XXX ns/op
- In-memory recursive: XXX ns/op
- **Improvement:** XXXx faster

### 🔄 Proceso

- FASE 1 (Web): Estructura y código
- FASE 2 (Local): Tests y validación ✅
```

**Crear PR:**
```bash
# Ya lo harás con la tool de GitHub
```

---

### TASK-08: Monitorear CI/CD (Variable)

Una vez creado el PR:

1. **Esperar a que corra el pipeline**
2. **Verificar jobs:**
   - Unit Tests
   - Integration Tests  
   - Lint
   - Security

3. **Si falla alguno:**
   - Analizar logs
   - Corregir error
   - Push de fix
   - Volver a monitorear

---

### TASK-09: Resolver Comentarios de Revisores (Variable)

**Si Copilot comenta:**
1. Leer comentario
2. Evaluar si es válido
3. Aplicar sugerencia si corresponde
4. Commit y push

**Si hay comentarios humanos:**
1. Discutir si no estás de acuerdo
2. Aplicar cambios si son razonables
3. Commit y push

---

### TASK-10: Merge a dev (5min)

**Cuando:**
- ✅ Pipeline verde
- ✅ Sin comentarios pendientes
- ✅ Aprobado (si requiere)

**Método:** Squash merge

```bash
# Lo harás con la tool de GitHub
# merge_method: squash
```

---

## 📊 CHECKLIST FINAL FASE 2

Antes de mergear, verifica:

### Código
- [ ] Sin `t.Skip()` en tests
- [ ] Sin comentarios `TODO_FASE2`
- [ ] Sin stubs/mocks
- [ ] Compilación OK

### Tests
- [ ] Unit tests: ✅
- [ ] Integration tests: ✅
- [ ] Coverage >= 15%
- [ ] Benchmarks ejecutados

### CI/CD
- [ ] Pipeline verde
- [ ] Comentarios resueltos
- [ ] PR aprobado

### Documentación
- [ ] HANDOFF marcado como completado
- [ ] Benchmark results documentados
- [ ] CHANGELOG actualizado (opcional)

---

## 🎯 COMANDOS RÁPIDOS

```bash
# Checkout
git checkout feature/sprint-03-repositorios-ltree

# Quitar skip de un test (ejemplo)
sed -i '' '/t.Skip("STUB_FASE2/d' test/integration/academic_unit_ltree_test.go

# Ejecutar tests de integración
go test -tags=integration ./test/integration/... -v

# Ver coverage
make coverage-report
open coverage/coverage.html

# Crear PR (hacer con tool)
# Monitorear (hacer con tool)
# Merge (hacer con tool)
```

---

## 📝 NOTAS

- Este es trabajo de **completar**, no de crear desde cero
- Claude Web ya hizo el trabajo pesado
- Tu rol es **quitar training wheels** (stubs) y **validar con DB real**
- **Duración:** Debería ser rápido (1-2h) vs los 4-5h originales

---

**¡Éxito en Fase 2!** 🚀
