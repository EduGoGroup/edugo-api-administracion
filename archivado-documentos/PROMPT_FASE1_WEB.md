# PROMPT FASE 1 - CLAUDE CODE WEB

**Proyecto:** edugo-api-administracion  
**Sprint:** Sprint-03 - Repositorios con ltree  
**Ejecutor:** Claude Code Web  
**Duración estimada:** 3-4 horas  
**Branch:** `feature/sprint-03-repositorios-ltree`

---

## 🎯 TU OBJETIVO (Fase 1 - Sin Docker)

Implementar la **capa de repositorios con soporte ltree** haciendo lo máximo posible **SIN** requerir Docker/PostgreSQL.

### ✅ Lo que SÍ puedes hacer:
- Escribir código de repositorios
- Escribir migraciones SQL
- Escribir tests de integración (estructura)
- Usar **stubs/mocks** donde necesites PostgreSQL
- Documentar TODO lo que dejas como stub

### ❌ Lo que NO puedes hacer:
- Ejecutar tests de integración (requieren Docker)
- Conectar a PostgreSQL real
- Ejecutar migraciones

---

## 📋 TAREAS FASE 1

### TASK-01: Crear Migración ltree (30min) ✅ PUEDES HACERLO

**Archivo:** `migrations/005_add_ltree_to_academic_units.up.sql`

```sql
-- Habilitar extensión ltree
CREATE EXTENSION IF NOT EXISTS ltree;

-- Agregar columna path
ALTER TABLE academic_units ADD COLUMN path ltree;

-- Índices para performance
CREATE INDEX academic_units_path_idx ON academic_units USING GIST (path);
CREATE INDEX academic_units_path_btree_idx ON academic_units USING btree (path);

-- Trigger para actualizar path automáticamente
CREATE OR REPLACE FUNCTION update_academic_unit_path()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.parent_unit_id IS NULL THEN
    -- Si es raíz, path es solo su ID
    NEW.path = NEW.id::text::ltree;
  ELSE
    -- Si tiene padre, concatenar path del padre + su ID
    SELECT path || NEW.id::text::ltree INTO NEW.path
    FROM academic_units WHERE id = NEW.parent_unit_id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER academic_unit_path_trigger
BEFORE INSERT OR UPDATE ON academic_units
FOR EACH ROW EXECUTE FUNCTION update_academic_unit_path();

-- Función para prevenir ciclos usando ltree
CREATE OR REPLACE FUNCTION prevent_academic_unit_cycles()
RETURNS TRIGGER AS $$
BEGIN
  -- Verificar que el nuevo padre no es descendiente de esta unidad
  IF NEW.parent_unit_id IS NOT NULL THEN
    IF EXISTS (
      SELECT 1 FROM academic_units 
      WHERE id = NEW.parent_unit_id 
      AND path <@ (SELECT path FROM academic_units WHERE id = NEW.id)
    ) THEN
      RAISE EXCEPTION 'Cannot set parent: would create a cycle in the hierarchy';
    END IF;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_cycles_trigger
BEFORE UPDATE ON academic_units
FOR EACH ROW EXECUTE FUNCTION prevent_academic_unit_cycles();
```

**Down migration:** `migrations/005_add_ltree_to_academic_units.down.sql`

```sql
DROP TRIGGER IF EXISTS prevent_cycles_trigger ON academic_units;
DROP TRIGGER IF EXISTS academic_unit_path_trigger ON academic_units;
DROP FUNCTION IF EXISTS prevent_academic_unit_cycles();
DROP FUNCTION IF EXISTS update_academic_unit_path();
DROP INDEX IF EXISTS academic_units_path_btree_idx;
DROP INDEX IF EXISTS academic_units_path_idx;
ALTER TABLE academic_units DROP COLUMN IF EXISTS path;
DROP EXTENSION IF EXISTS ltree;
```

**✅ Acción:** Crear ambos archivos.

---

### TASK-02: Implementar Repository con ltree (2h) ✅ PUEDES HACERLO

**Archivo:** `internal/infrastructure/persistence/postgres/repository/academic_unit_repository_impl.go`

**Métodos a agregar:**

```go
// FindByPath busca una unidad por su path ltree
func (r *postgresAcademicUnitRepository) FindByPath(
    ctx context.Context, 
    path string,
) (*entity.AcademicUnit, error) {
    query := `
        SELECT id, parent_unit_id, school_id, unit_type, display_name, 
               code, description, metadata, created_at, updated_at, deleted_at
        FROM academic_units
        WHERE path = $1 AND deleted_at IS NULL
    `
    // Implementar escaneo y mapeo
}

// FindChildren retorna los hijos directos de una unidad
func (r *postgresAcademicUnitRepository) FindChildren(
    ctx context.Context,
    parentID valueobject.UnitID,
) ([]*entity.AcademicUnit, error) {
    query := `
        SELECT id, parent_unit_id, school_id, unit_type, display_name,
               code, description, metadata, created_at, updated_at, deleted_at
        FROM academic_units
        WHERE parent_unit_id = $1 AND deleted_at IS NULL
        ORDER BY display_name
    `
    // Implementar escaneo y mapeo
}

// FindDescendants retorna TODOS los descendientes usando ltree
func (r *postgresAcademicUnitRepository) FindDescendants(
    ctx context.Context,
    unitID valueobject.UnitID,
) ([]*entity.AcademicUnit, error) {
    query := `
        SELECT u.id, u.parent_unit_id, u.school_id, u.unit_type, 
               u.display_name, u.code, u.description, u.metadata,
               u.created_at, u.updated_at, u.deleted_at
        FROM academic_units u
        WHERE u.path <@ (
            SELECT path FROM academic_units WHERE id = $1
        )
        AND u.id != $1
        AND u.deleted_at IS NULL
        ORDER BY u.path
    `
    // Implementar escaneo y mapeo
}

// FindAncestors retorna TODOS los ancestros usando ltree
func (r *postgresAcademicUnitRepository) FindAncestors(
    ctx context.Context,
    unitID valueobject.UnitID,
) ([]*entity.AcademicUnit, error) {
    query := `
        SELECT u.id, u.parent_unit_id, u.school_id, u.unit_type,
               u.display_name, u.code, u.description, u.metadata,
               u.created_at, u.updated_at, u.deleted_at
        FROM academic_units u
        WHERE u.path @> (
            SELECT path FROM academic_units WHERE id = $1
        )
        AND u.id != $1
        AND u.deleted_at IS NULL
        ORDER BY u.path
    `
    // Implementar escaneo y mapeo
}

// FindBySchoolIDAndDepth retorna unidades de una escuela a una profundidad específica
func (r *postgresAcademicUnitRepository) FindBySchoolIDAndDepth(
    ctx context.Context,
    schoolID valueobject.SchoolID,
    depth int,
) ([]*entity.AcademicUnit, error) {
    query := `
        SELECT u.id, u.parent_unit_id, u.school_id, u.unit_type,
               u.display_name, u.code, u.description, u.metadata,
               u.created_at, u.updated_at, u.deleted_at
        FROM academic_units u
        WHERE u.school_id = $1
        AND nlevel(u.path) = $2
        AND u.deleted_at IS NULL
        ORDER BY u.path
    `
    // Implementar escaneo y mapeo
}

// MoveSubtree mueve un subárbol completo a un nuevo padre
func (r *postgresAcademicUnitRepository) MoveSubtree(
    ctx context.Context,
    unitID valueobject.UnitID,
    newParentID *valueobject.UnitID,
) error {
    // Esta operación requiere transacción
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. Obtener path actual
    // 2. Actualizar parent_unit_id
    // 3. El trigger actualizará el path automáticamente
    query := `
        UPDATE academic_units
        SET parent_unit_id = $1, updated_at = NOW()
        WHERE id = $2
    `
    
    // Implementar lógica de actualización
    
    return tx.Commit()
}
```

**✅ Acción:** Implementa estos métodos en el repository existente.

**📝 Nota importante:** El código compilará sin PostgreSQL. Solo escribe la lógica.

---

### TASK-03: Escribir Tests de Integración (1h) ⚠️ USAR STUBS

**Archivo:** `test/integration/academic_unit_ltree_test.go`

Como **NO tienes Docker**, usa esta estrategia:

```go
//go:build integration

package integration

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

// STUB: Este test requiere PostgreSQL con ltree
// Claude Code Local ejecutará esto con testcontainers reales
func TestAcademicUnitRepository_FindChildren(t *testing.T) {
    t.Skip("STUB: Requiere PostgreSQL con ltree - Completar en FASE 2 (Claude Local)")
    
    // TODO FASE 2: Descomentar y ejecutar
    // db, cleanup := getTestDB(t)
    // defer cleanup()
    // 
    // repo := repository.NewPostgresAcademicUnitRepository(db)
    // 
    // // Setup: Crear jerarquía
    // parent := createTestUnit(t, repo, "Grade 1")
    // child1 := createTestUnit(t, repo, "Section A", parent.ID())
    // child2 := createTestUnit(t, repo, "Section B", parent.ID())
    // 
    // // Test
    // children, err := repo.FindChildren(ctx, parent.ID())
    // assert.NoError(t, err)
    // assert.Len(t, children, 2)
}

func TestAcademicUnitRepository_FindDescendants(t *testing.T) {
    t.Skip("STUB: Requiere PostgreSQL con ltree - Completar en FASE 2")
    // Similar estructura...
}

// Continuar con todos los métodos...
```

**✅ Acción:** 
1. Escribe la ESTRUCTURA de todos los tests
2. Marca con `t.Skip("STUB: ...")` 
3. Documenta en comentarios lo que debe hacer
4. Deja el código comentado para Fase 2

---

### TASK-04: Documentar Handoff (30min) ✅ CRÍTICO

**Archivo:** `HANDOFF_FASE1_A_FASE2.md`

Crea este archivo documentando:

```markdown
# Handoff: Fase 1 (Web) → Fase 2 (Local)

## ✅ Completado en Fase 1 (Claude Web)

### Migraciones SQL
- [x] `migrations/005_add_ltree_to_academic_units.up.sql`
- [x] `migrations/005_add_ltree_to_academic_units.down.sql`

### Repository Implementation
- [x] `FindByPath()` - implementado
- [x] `FindChildren()` - implementado  
- [x] `FindDescendants()` - implementado
- [x] `FindAncestors()` - implementado
- [x] `FindBySchoolIDAndDepth()` - implementado
- [x] `MoveSubtree()` - implementado

### Tests (Estructura)
- [x] `test/integration/academic_unit_ltree_test.go` - estructura creada
- ⚠️ Tests marcados con `t.Skip()` - **REQUIEREN FASE 2**

## ⏸️ PENDIENTE para Fase 2 (Claude Local)

### Tests de Integración - DESCOMENTAR Y EJECUTAR
Archivo: `test/integration/academic_unit_ltree_test.go`

**Razón:** Requieren Docker + PostgreSQL con extensión ltree

**Tareas:**
1. Quitar todos los `t.Skip()` 
2. Descomentar código de tests
3. Ejecutar con testcontainers
4. Validar que pasan

### Migraciones - EJECUTAR
Archivo: `migrations/005_add_ltree_to_academic_units.up.sql`

**Razón:** Requieren PostgreSQL corriendo

**Tareas:**
1. Levantar PostgreSQL local/testcontainer
2. Ejecutar migración
3. Validar que ltree funciona
4. Validar que trigger actualiza path correctamente

### Validaciones con Base de Datos Real

**Tests específicos que requieren ejecución:**
- [ ] `TestFindChildren` - Validar query ltree
- [ ] `TestFindDescendants` - Validar path <@
- [ ] `TestFindAncestors` - Validar path @>
- [ ] `TestMoveSubtree` - Validar transacción
- [ ] Benchmark: ltree vs recursión

## 🔍 Stubs/Mocks Usados

**Ninguno requerido** - El código de repository no usa mocks, solo queries.
Los tests están estructurados pero skipeados.

## 📊 Cobertura Esperada Post-Fase 2

- Repository ltree methods: >= 80%
- Integration tests: Todos pasando

## 🚀 Comando para Fase 2

```bash
# 1. Ejecutar migraciones (si es necesario)
# make db-migrate

# 2. Ejecutar tests de integración
go test -tags=integration ./test/integration/... -v -run TestAcademicUnit

# 3. Benchmark
go test -tags=integration ./test/integration/... -bench=BenchmarkTreeQueries

# 4. Cobertura
go test ./internal/infrastructure/postgres/... -cover
```

## ⚠️ Notas Importantes

1. **ltree extension**: Debe estar disponible en PostgreSQL
2. **Trigger validation**: Validar que path se actualiza automáticamente
3. **Performance**: Comparar ltree vs recursión en memoria
4. **Edge cases**: Validar ciclos, path nulos, jerarquías profundas
```

**✅ Acción:** Crea este archivo al finalizar tu trabajo.

---

## 🎯 CHECKLIST DE ENTREGA (Fase 1)

Antes de marcar como "listo para Fase 2", verifica:

### Código
- [ ] Migración ltree creada (up + down)
- [ ] Repository con 6 métodos ltree implementados
- [ ] Código compila sin errores (`go build ./...`)
- [ ] Lint pasa sin errores (`make lint`)

### Tests
- [ ] Tests de integración estructurados
- [ ] Todos los tests con `t.Skip("STUB: ...")`
- [ ] Código de test comentado pero completo
- [ ] Casos de prueba documentados

### Documentación
- [ ] `HANDOFF_FASE1_A_FASE2.md` creado
- [ ] Stubs/skips claramente marcados
- [ ] Comentarios explican qué hacer en Fase 2

### Git
- [ ] Commit: `feat(infrastructure): add ltree repository methods (FASE 1 - stubs)`
- [ ] Branch pusheada
- [ ] **NO crear PR todavía** - será en Fase 2

---

## 📝 TEMPLATE DE STUB

Cuando no puedas ejecutar algo, usa:

```go
// STUB_FASE2: Este código requiere [PostgreSQL/Docker/Internet]
// Completar en Fase 2 con Claude Code Local
func TestSomething(t *testing.T) {
    t.Skip("STUB_FASE2: Requiere PostgreSQL con ltree extension")
    
    // TODO_FASE2: Descomentar y ejecutar
    // db, cleanup := getTestDB(t)
    // defer cleanup()
    // ... resto del test
}
```

---

## 🚨 SI ENCUENTRAS PROBLEMAS

### Problema: "No puedo ejecutar tests de integración"
**Solución:** Usa `t.Skip()` con mensaje claro. Está OK.

### Problema: "No tengo PostgreSQL para validar queries"
**Solución:** Escribe los queries basándote en docs de ltree. Yo validaré en Fase 2.

### Problema: "No puedo instalar dependencias por internet"
**Solución:** Si ya están en go.mod, no hay problema. Si necesitas nuevas, documenta en HANDOFF.

---

## 📚 Referencias para ltree

**Operators ltree:**
- `@>` - contiene (ancestros)
- `<@` - está contenido (descendientes)  
- `~` - match pattern
- `||` - concatenación

**Funciones ltree:**
- `nlevel(path)` - profundidad
- `subpath(path, start, end)` - subpath
- `index(path, label)` - buscar label

**Docs:** https://www.postgresql.org/docs/current/ltree.html

---

## ✅ ENTREGABLES FASE 1

Al finalizar, debes tener:

```
migrations/
  ├── 005_add_ltree_to_academic_units.up.sql    ✅ Creado
  └── 005_add_ltree_to_academic_units.down.sql  ✅ Creado

internal/infrastructure/persistence/postgres/repository/
  └── academic_unit_repository_impl.go          ✅ 6 métodos agregados

test/integration/
  └── academic_unit_ltree_test.go               ⚠️ Con t.Skip()

HANDOFF_FASE1_A_FASE2.md                        ✅ Documentado

.git/
  └── branch: feature/sprint-03-repositorios-ltree  ✅ Pusheada
```

---

## 🎯 CRITERIO DE ÉXITO

**Tu trabajo está COMPLETO cuando:**

1. ✅ Código compila (`go build ./...`)
2. ✅ Lint pasa (`make lint`)
3. ✅ Tests unitarios pasan (los que no requieren DB)
4. ✅ HANDOFF document completo
5. ✅ Branch pusheada
6. ⚠️ Tests de integración SKIPEADOS (OK para Fase 1)

**NO es tu responsabilidad:**
- ❌ Ejecutar tests de integración
- ❌ Conectar a PostgreSQL
- ❌ Ejecutar migraciones
- ❌ Crear PR

Eso lo hará Claude Local en Fase 2.

---

## 🚀 COMANDO PARA INICIAR

```bash
git checkout main
git pull origin main
git checkout -b feature/sprint-03-repositorios-ltree
# ... hacer tu trabajo
git add -A
git commit -m "feat(infrastructure): add ltree repository methods (FASE 1 - stubs)"
git push origin feature/sprint-03-repositorios-ltree
```

---

**Al finalizar, deja el mensaje:**

```
✅ FASE 1 COMPLETA
Branch: feature/sprint-03-repositorios-ltree pusheada
Handoff: Ver HANDOFF_FASE1_A_FASE2.md
Siguiente: Claude Code Local ejecutará Fase 2
```

---

¡Buena suerte Claude Web! 🚀
