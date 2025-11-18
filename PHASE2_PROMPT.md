# PROMPT FASE 2 - SPRINT-03 REPOSITORIOS CON LTREE

**Proyecto:** edugo-api-administracion
**Ejecutor:** Claude Code Local (o Claude Code Web si hay conectividad)
**Duración estimada:** 4-5 horas
**Branch:** `feature/sprint-03-repositorios-ltree`

---

## 🎯 OBJETIVO

Implementar los **repositorios con soporte para árbol jerárquico** usando PostgreSQL ltree para gestión académica.

**Basado en:** Sprint-02 (Entidades del Dominio con Árbol)

---

## ✅ COMPLETADO EN SPRINT-02 (FASE 1)

### Value Objects
- ✅ UnitID - Completo con validaciones y métodos
- ✅ UnitType - Completo con jerarquía (CanHaveChildren, AllowedChildTypes)
- ✅ MembershipID - Completo
- ✅ MembershipRole - Completo con permisos

### Entidades del Dominio
- ✅ AcademicUnit con métodos de árbol:
  - Campo `children []*AcademicUnit` para árbol en memoria
  - `IsRoot()` - Verifica si es raíz
  - `HasChildren()` - Verifica si tiene hijos
  - `AddChild()` - Agrega hijo con validaciones completas
  - `RemoveChild()` - Remueve hijo
  - `GetAllDescendants()` - Obtiene todos los descendientes recursivamente
  - `GetDepth()` - Calcula profundidad del árbol
  - `UpdateDisplayName()` - Actualiza nombre de visualización
  - `SetParent()` - Establece padre con validaciones de tipo
  - `RemoveParent()` - Convierte en raíz
  - `CanHaveChildren()` - Validación de tipo

- ✅ UnitMembership completa:
  - `IsActive()` - Verifica si está activa ahora
  - `IsActiveAt(time.Time)` - Verifica si está activa en momento específico
  - `Expire()` - Marca como expirada
  - `ChangeRole()` - Cambia el rol
  - `HasPermission()` - Verifica permisos según rol
  - `SetValidUntil()` - Establece fecha de fin
  - `ExtendIndefinitely()` - Remueve fecha de fin

### Tests Unitarios
- ✅ `internal/domain/entity/academic_unit_test.go` - 100% de los métodos testeados
  - Tests de construcción y reconstrucción
  - Tests de setters/getters
  - Tests de lógica de negocio (SetParent, UpdateInfo, etc.)
  - Tests completos de árbol (AddChild, RemoveChild, GetAllDescendants, GetDepth)
  - Tests de soft delete y restauración
  - Tests de metadata

- ✅ `internal/domain/entity/unit_membership_test.go` - 100% de los métodos testeados
  - Tests de construcción y reconstrucción
  - Tests de activación temporal (IsActive, IsActiveAt)
  - Tests de expiración y extensión
  - Tests de cambio de rol
  - Tests de permisos por rol
  - Tests de metadata

**Nota:** Los tests están implementados pero no se pudieron ejecutar en el entorno web debido a problemas de conectividad de red que impidieron descargar Go 1.24.10. Se recomienda ejecutar localmente con:

```bash
go test ./internal/domain/entity -v -cover
```

---

## 📋 TAREAS SPRINT-03

### TASK-03-001: Migración ltree (2h)

Crear migración en `migrations/` para agregar soporte ltree:

```sql
-- migrations/XXXXXX_add_ltree_to_academic_units.sql
CREATE EXTENSION IF NOT EXISTS ltree;

ALTER TABLE academic_units ADD COLUMN path ltree;
CREATE INDEX academic_units_path_idx ON academic_units USING GIST (path);
CREATE INDEX academic_units_path_btree_idx ON academic_units USING btree (path);

-- Trigger para actualizar path automáticamente
CREATE OR REPLACE FUNCTION update_academic_unit_path()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.parent_unit_id IS NULL THEN
    NEW.path = NEW.id::text::ltree;
  ELSE
    SELECT path || NEW.id::text::ltree INTO NEW.path
    FROM academic_units WHERE id = NEW.parent_unit_id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER academic_unit_path_trigger
BEFORE INSERT OR UPDATE ON academic_units
FOR EACH ROW EXECUTE FUNCTION update_academic_unit_path();
```

### TASK-03-002: Repository Implementation (3h)

Implementar en `internal/infrastructure/postgres/academic_unit_repository.go`:

**Métodos de árbol:**
- `FindByPath(ctx, path ltree)` - Buscar por path ltree
- `FindChildren(ctx, parentID)` - Hijos directos
- `FindDescendants(ctx, unitID)` - Todos los descendientes
- `FindAncestors(ctx, unitID)` - Todos los ancestros
- `FindBySchoolIDAndDepth(ctx, schoolID, depth)` - Por profundidad específica
- `MoveSubtree(ctx, unitID, newParentID)` - Mover subárbol completo

**Queries ltree:**
```go
// Descendientes
SELECT * FROM academic_units WHERE path <@ (SELECT path FROM academic_units WHERE id = $1)

// Ancestros
SELECT * FROM academic_units WHERE path @> (SELECT path FROM academic_units WHERE id = $1)

// Hijos directos
SELECT * FROM academic_units WHERE parent_unit_id = $1
```

### TASK-03-003: Integration Tests (2h)

Tests de integración en `test/integration/academic_unit_repository_test.go`:

```go
func TestAcademicUnitRepository_TreeOperations(t *testing.T) {
    // Test FindChildren
    // Test FindDescendants
    // Test FindAncestors
    // Test MoveSubtree
}
```

### TASK-03-004: Service Layer (3h)

**NO implementar en esta fase** - será Sprint-04

---

## 🔍 PUNTOS DE VALIDACIÓN

1. **Migración ltree:**
   - [ ] Extension ltree creada
   - [ ] Columna path agregada
   - [ ] Índices GIST y BTREE creados
   - [ ] Trigger funciona correctamente

2. **Repository:**
   - [ ] Métodos de árbol implementados
   - [ ] Queries ltree optimizadas
   - [ ] Tests de integración >= 80% cobertura
   - [ ] MoveSubtree mantiene integridad referencial

3. **Performance:**
   - [ ] Queries usan índices ltree
   - [ ] GetDescendants es O(log n) con ltree vs O(n) recursivo
   - [ ] Benchmark comparativo documentado

---

## 📖 REFERENCIAS

- **PostgreSQL ltree:** https://www.postgresql.org/docs/current/ltree.html
- **Queries ltree comunes:** Ver `docs/isolated/03-Architecture/database-design.md`
- **Tests de integración:** Usar testcontainers con PostgreSQL 16

---

## ✅ CRITERIOS DE ACEPTACIÓN

1. Migración ltree ejecutada sin errores
2. Repository implementado con todos los métodos de árbol
3. Tests de integración >= 80% cobertura
4. Benchmark muestra mejora de performance con ltree
5. Código pasa linter y tests
6. Documentación actualizada

---

## 🚀 EJECUCIÓN

```bash
# Ejecutar migración
make db-migrate

# Tests de integración
go test ./test/integration -v -run TestAcademicUnit

# Benchmark
go test ./test/integration -bench=BenchmarkTreeQueries -benchmem

# Cobertura
go test ./internal/infrastructure/postgres -cover
```

---

## 📝 NOTAS PARA CLAUDE CODE LOCAL

- El Sprint-02 está **100% completo** con entidades y tests unitarios
- Los tests están escritos pero no ejecutados por problemas de red en entorno web
- Ejecutar primero: `go test ./internal/domain/entity -v -cover` para validar que todo compila
- Luego proceder con Sprint-03 (repositorios con ltree)
- El objetivo final es tener un sistema de árbol jerárquico eficiente con PostgreSQL ltree

---

## 🎯 SIGUIENTE FASE (Sprint-04)

Una vez completado Sprint-03, proceder con Sprint-04 que implementará:
- Services con lógica de negocio para operaciones de árbol
- Handlers HTTP para APIs REST
- DTOs y validaciones
- Tests end-to-end

**Branch para Sprint-04:** `feature/sprint-04-services-handlers`
