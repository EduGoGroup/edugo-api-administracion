# 📋 Sprint-00: Integración Completa con Infrastructure

**Fecha:** 17 de Noviembre, 2025  
**Versión:** 2.0.0  
**Duración Total:** 7-9 horas  
**Prioridad:** CRÍTICA (debe ejecutarse PRIMERO)

---

## 🎯 Objetivo

Migrar `edugo-api-administracion` para usar `edugo-infrastructure` como fuente de verdad de la base de datos, siguiendo el mismo patrón que `edugo-api-mobile`.

**Regla de oro:** Infrastructure tiene la última palabra sobre el schema de BD.

---

## 📊 Resumen de Cambios

| Aspecto | Antes | Después |
|---------|-------|---------|
| **Migraciones** | Local (scripts/postgresql/) | Infrastructure v0.7.0 |
| **Dependencias** | shared v0.5.0 | shared v0.7.0 + infra v0.7.0 |
| **Tablas** | Singular (school) | Plural (schools) |
| **Jerarquía** | ✅ Soportada | ✅ Soportada (agregada a infra) |
| **LOC eliminadas** | - | ~500 líneas (scripts SQL) |

---

## 🚧 IMPORTANTE: 2 FASES

Este Sprint se divide en **2 fases secuenciales**:

### FASE 1: Actualizar Infrastructure (BLOQUEANTE)
- **Duración:** 3-4 horas
- **Ubicación:** `edugo-infrastructure` (proyecto hermano)
- **Output:** Release v0.7.0

### FASE 2: Migrar api-admin
- **Duración:** 4-5 horas
- **Ubicación:** `edugo-api-administracion` (este proyecto)
- **Dependencia:** Requiere infrastructure v0.7.0

**⚠️ NO se puede ejecutar FASE 2 sin completar FASE 1**

---

# 🔷 FASE 1: Actualizar Infrastructure (3-4 horas)

## 📍 Ubicación de Trabajo

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure
```

---

## TASK 1.1: Copiar Migración 012 a Infrastructure

**Duración:** 10 minutos

```bash
# Desde edugo-api-administracion
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Copiar migraciones
cp docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/migrations/012_extend_for_admin_api.up.sql \
   /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/postgres/migrations/

cp docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/migrations/012_extend_for_admin_api.down.sql \
   /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/postgres/migrations/
```

**Validación:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/postgres/migrations
ls -la 012_extend_for_admin_api.*
# Debe mostrar:
# 012_extend_for_admin_api.up.sql
# 012_extend_for_admin_api.down.sql
```

---

## TASK 1.2: Validar Sintaxis de Migraciones

**Duración:** 10 minutos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/postgres

# Validar sintaxis SQL
psql -U postgres -d postgres -f migrations/012_extend_for_admin_api.up.sql --dry-run 2>&1 | head -20

# O usar herramienta de linting SQL si tienes
sqlfluff lint migrations/012_extend_for_admin_api.up.sql
```

**Si tienes errores de sintaxis:**
- Corregir en el archivo
- Re-validar

---

## TASK 1.3: Testing de Migraciones (UP y DOWN)

**Duración:** 30 minutos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# 1. Levantar PostgreSQL de testing
make dev-up-core
# Esperar 10 segundos para que PostgreSQL arranque

# 2. Ejecutar TODAS las migraciones (001 a 011)
cd postgres
make migrate-up

# 3. Ejecutar migración 012 (la nueva)
psql -U edugo -d edugo_dev -f migrations/012_extend_for_admin_api.up.sql

# 4. Validar que se aplicó correctamente
psql -U edugo -d edugo_dev -c "\d academic_units" | grep parent_unit_id
# Debe mostrar: parent_unit_id | uuid

psql -U edugo -d edugo_dev -c "\d+ academic_units" | grep metadata
# Debe mostrar: metadata | jsonb

# 5. Testing de rollback (DOWN)
psql -U edugo -d edugo_dev -f migrations/012_extend_for_admin_api.down.sql

# 6. Validar que se revirtió
psql -U edugo -d edugo_dev -c "\d academic_units" | grep parent_unit_id
# NO debe mostrar nada

# 7. Re-aplicar (UP) para dejar lista
psql -U edugo -d edugo_dev -f migrations/012_extend_for_admin_api.up.sql
```

**Validación final:**
```bash
# Verificar que existen:
psql -U edugo -d edugo_dev -c "
SELECT 
    column_name, 
    data_type 
FROM information_schema.columns 
WHERE table_name = 'academic_units' 
  AND column_name IN ('parent_unit_id', 'metadata', 'description')
ORDER BY column_name;
"

# Debe mostrar:
#  column_name    | data_type
# ----------------+-----------
#  description    | text
#  metadata       | jsonb
#  parent_unit_id | uuid
```

---

## TASK 1.4: Actualizar CHANGELOG.md de Infrastructure

**Duración:** 10 minutos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure
```

Agregar al inicio de `CHANGELOG.md`:

```markdown
## [0.7.0] - 2025-11-17

### Added (postgres)
- **Migración 012:** Soporte de jerarquía en `academic_units` con `parent_unit_id`
- Campo `metadata` JSONB en `schools`, `academic_units`, `memberships`
- Campo `description` TEXT en `academic_units`
- Tipos extendidos en `academic_units`: `school`, `club`, `department` (además de `grade`, `class`, `section`)
- Roles extendidos en `memberships`: `coordinator`, `admin`, `assistant` (además de `teacher`, `student`, `guardian`)
- Función `prevent_academic_unit_cycles()` para validar jerarquía
- Vista `v_academic_unit_tree` con CTE recursivo para árbol jerárquico
- Vista mejorada `v_active_memberships` con información completa

### Changed (postgres)
- `academic_units.academic_year` ahora es nullable (default: 0)
- Constraint `academic_units_type_check` extendido
- Constraint `memberships_role_check` extendido

### Migration Path
Para proyectos existentes:
- `academic_year` con valor 0 indica "sin año específico"
- `parent_unit_id` NULL indica unidad raíz (sin padre)
- Campos `metadata` tienen default `{}` (vacío)
```

---

## TASK 1.5: Commit y Push a Infrastructure

**Duración:** 5 minutos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# 1. Agregar archivos
git add postgres/migrations/012_extend_for_admin_api.up.sql
git add postgres/migrations/012_extend_for_admin_api.down.sql
git add CHANGELOG.md

# 2. Commit
git commit -m "feat(postgres): add hierarchy support and extend schema for api-admin

- Add parent_unit_id to academic_units for hierarchical structure
- Add metadata JSONB to schools, academic_units, memberships
- Add description TEXT to academic_units
- Extend academic_units types: school, club, department
- Extend memberships roles: coordinator, admin, assistant
- Add prevent_academic_unit_cycles() function and trigger
- Add v_academic_unit_tree view (recursive CTE)
- Improve v_active_memberships view
- Make academic_units.academic_year nullable (default: 0)

BREAKING CHANGE: academic_year is now nullable
Migration: 012_extend_for_admin_api.up.sql"

# 3. Push a rama actual (probablemente dev)
git push origin dev

# O si quieres crear una rama específica:
git checkout -b feature/add-hierarchy-support
git push origin feature/add-hierarchy-support
```

---

## TASK 1.6: Crear Release v0.7.0 (Tag)

**Duración:** 10 minutos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Asegurarse de estar en la rama correcta (dev o main)
git checkout dev  # o main

# Merge de la feature branch si creaste una
git merge feature/add-hierarchy-support  # si usaste branch

# Crear tag
git tag -a v0.7.0 -m "Release v0.7.0: Hierarchy support for academic_units

## Added
- Hierarchical structure with parent_unit_id
- Metadata fields (JSONB) across tables
- Extended types and roles
- Tree views with recursive CTEs

## Changed  
- academic_year is now nullable

## Migration
- 012_extend_for_admin_api.up.sql
- Full backward compatibility maintained"

# Push tag
git push origin v0.7.0

# Push rama
git push origin dev  # o main
```

**Validación:**
```bash
# Ver tags
git tag -l | grep v0.7.0

# Ver en GitHub
# https://github.com/EduGoGroup/edugo-infrastructure/releases/tag/v0.7.0
```

---

## TASK 1.7: Validar que v0.7.0 está Disponible

**Duración:** 5 minutos

```bash
# Esperar 1-2 minutos para que GitHub procese el tag

# Intentar descargar
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.0

# O verificar en GitHub:
# https://github.com/EduGoGroup/edugo-infrastructure/tags
```

**Si el tag NO aparece:**
- Verificar que hiciste `git push origin v0.7.0`
- Verificar permisos en GitHub
- Esperar 2-3 minutos más

---

## ✅ Checklist FASE 1 Completada

Antes de pasar a FASE 2, verificar:

- [ ] Migraciones 012 copiadas a infrastructure
- [ ] Testing de UP exitoso (campos agregados)
- [ ] Testing de DOWN exitoso (rollback funciona)
- [ ] CHANGELOG.md actualizado
- [ ] Commit creado en infrastructure
- [ ] Push a GitHub exitoso
- [ ] Tag v0.7.0 creado
- [ ] Tag v0.7.0 visible en GitHub
- [ ] `go get ...@v0.7.0` funciona

**Tiempo estimado FASE 1:** 1.5 - 2 horas (si no hay problemas)

---

# 🔷 FASE 2: Migrar api-admin (4-5 horas)

## 📍 Ubicación de Trabajo

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion
```

---

## TASK 2.1: Crear Rama de Trabajo

**Duración:** 2 minutos

```bash
git checkout dev  # o tu rama base
git pull origin dev

git checkout -b feature/migrate-to-infrastructure-v0.7.0
```

---

## TASK 2.2: Actualizar go.mod con Infrastructure v0.7.0

**Duración:** 10 minutos

```bash
# 1. Agregar módulos de infrastructure
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.0
go get github.com/EduGoGroup/edugo-infrastructure/migrations@v0.7.0

# 2. Actualizar módulos de shared
go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0
go get github.com/EduGoGroup/edugo-shared/common@v0.7.0
go get github.com/EduGoGroup/edugo-shared/logger@v0.7.0
go get github.com/EduGoGroup/edugo-shared/lifecycle@v0.7.0
go get github.com/EduGoGroup/edugo-shared/bootstrap@v0.5.0  # mantener 0.5.0

# 3. Limpiar
go mod tidy
```

**Validación:**
```bash
cat go.mod | grep -E "(infrastructure|shared)" | head -10

# Debe mostrar versiones:
# github.com/EduGoGroup/edugo-infrastructure/postgres v0.7.0
# github.com/EduGoGroup/edugo-shared/auth v0.7.0
```

---

## TASK 2.3: Refactoring de Repositorios (Cambio de Nombres)

**Duración:** 2 horas

### Archivos a Modificar

1. **`internal/infrastructure/persistence/postgres/repository/school_repository_impl.go`**

Cambios:
```go
// ANTES:
query := `INSERT INTO school (id, name, code, ...`
query := `SELECT ... FROM school WHERE id = $1`
query := `UPDATE school SET ...`
query := `DELETE FROM school WHERE id = $1`

// DESPUÉS:
query := `INSERT INTO schools (id, name, code, ...`
query := `SELECT ... FROM schools WHERE id = $1`
query := `UPDATE schools SET ...`
query := `DELETE FROM schools WHERE id = $1`

// Renombrar campos:
// contact_email → email
// contact_phone → phone
```

2. **`internal/infrastructure/persistence/postgres/repository/academic_unit_repository_impl.go`**

Cambios:
```go
// ANTES:
query := `INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, ...`
query := `SELECT ... FROM academic_unit WHERE id = $1`
query := `UPDATE academic_unit SET ...`
query := `DELETE FROM academic_unit WHERE id = $1`

// DESPUÉS:
query := `INSERT INTO academic_units (id, parent_unit_id, school_id, type, name, ...`
query := `SELECT ... FROM academic_units WHERE id = $1`
query := `UPDATE academic_units SET ...`
query := `DELETE FROM academic_units WHERE id = $1`

// Renombrar campos:
// unit_type → type
// display_name → name
// Agregar: academic_year (usar 0 como default)
```

3. **`internal/infrastructure/persistence/postgres/repository/unit_membership_repository_impl.go`**

Cambios:
```go
// ANTES:
query := `INSERT INTO unit_membership (id, unit_id, user_id, role, ...`
query := `SELECT ... FROM unit_membership WHERE id = $1`
query := `UPDATE unit_membership SET ...`
query := `DELETE FROM unit_membership WHERE id = $1`

// DESPUÉS:
query := `INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, ...`
query := `SELECT ... FROM memberships WHERE id = $1`
query := `UPDATE memberships SET ...`
query := `DELETE FROM memberships WHERE id = $1`

// Renombrar campos:
// unit_id → academic_unit_id
// Agregar: school_id (obtener de academic_unit)
```

**Script auxiliar para refactoring:**

```bash
# Crear backup
cp -r internal internal.backup

# Reemplazos automáticos (revisar después)
find internal -name "*.go" -type f -exec sed -i '' 's/FROM school/FROM schools/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/INTO school/INTO schools/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/UPDATE school/UPDATE schools/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/JOIN school/JOIN schools/g' {} +

find internal -name "*.go" -type f -exec sed -i '' 's/FROM academic_unit/FROM academic_units/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/INTO academic_unit/INTO academic_units/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/UPDATE academic_unit/UPDATE academic_units/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/JOIN academic_unit/JOIN academic_units/g' {} +

find internal -name "*.go" -type f -exec sed -i '' 's/FROM unit_membership/FROM memberships/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/INTO unit_membership/INTO memberships/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/UPDATE unit_membership/UPDATE memberships/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/JOIN unit_membership/JOIN memberships/g' {} +

# Renombrar campos (más delicado, revisar manualmente)
find internal -name "*.go" -type f -exec sed -i '' 's/contact_email/email/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/contact_phone/phone/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/display_name/name/g' {} +
find internal -name "*.go" -type f -exec sed -i '' 's/unit_type/type/g' {} +
```

**⚠️ IMPORTANTE:** Revisar MANUALMENTE cada cambio antes de continuar

---

## TASK 2.4: Agregar Campo `academic_year` a Entity

**Duración:** 30 minutos

1. **Modificar `internal/domain/entity/academic_unit.go`:**

```go
type AcademicUnit struct {
    id           valueobject.UnitID
    parentUnitID *valueobject.UnitID
    schoolID     valueobject.SchoolID
    unitType     valueobject.UnitType
    displayName  string  // RENOMBRAR a "name" si quieres
    code         string
    description  string
    academicYear int      // NUEVO: año académico (0 = sin año)
    metadata     map[string]interface{}
    createdAt    time.Time
    updatedAt    time.Time
    deletedAt    *time.Time
}

// Agregar getter/setter
func (a *AcademicUnit) AcademicYear() int {
    return a.academicYear
}

func (a *AcademicUnit) SetAcademicYear(year int) {
    a.academicYear = year
}
```

2. **Actualizar constructores para incluir `academicYear`:**

```go
func NewAcademicUnit(
    id valueobject.UnitID,
    parentUnitID *valueobject.UnitID,
    schoolID valueobject.SchoolID,
    unitType valueobject.UnitType,
    displayName string,
    code string,
    description string,
    academicYear int,  // NUEVO
    metadata map[string]interface{},
) (*AcademicUnit, error) {
    // ... validaciones

    return &AcademicUnit{
        id:           id,
        parentUnitID: parentUnitID,
        schoolID:     schoolID,
        unitType:     unitType,
        displayName:  displayName,
        code:         code,
        description:  description,
        academicYear: academicYear,  // NUEVO
        metadata:     metadata,
        createdAt:    time.Now(),
        updatedAt:    time.Now(),
    }, nil
}
```

3. **Actualizar DTOs en `internal/application/dto/`:**

Agregar `AcademicYear` en los DTOs de request/response

---

## TASK 2.5: Actualizar Repositorios (Queries SQL)

**Duración:** 1.5 horas

Para cada repositorio modificado, actualizar:

### Ejemplo: `academic_unit_repository_impl.go`

```go
func (r *postgresAcademicUnitRepository) Create(ctx context.Context, unit *entity.AcademicUnit) error {
    query := `
        INSERT INTO academic_units (
            id, parent_unit_id, school_id, type, name, code, description, 
            academic_year, metadata, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
    
    // ... preparar valores
    
    academicYear := unit.AcademicYear()
    if academicYear == 0 {
        academicYear = time.Now().Year()  // Default: año actual
    }
    
    _, err := r.db.ExecContext(ctx, query,
        unit.ID().String(),
        parentID,
        unit.SchoolID().String(),
        unit.UnitType().String(),
        unit.DisplayName(),  // o unit.Name() si renombraste
        unit.Code(),
        unit.Description(),
        academicYear,  // NUEVO
        metadataJSON,
        unit.CreatedAt(),
        unit.UpdatedAt(),
    )
    
    return err
}
```

**Archivos a actualizar:**
- `school_repository_impl.go` (~15 queries)
- `academic_unit_repository_impl.go` (~20 queries)
- `unit_membership_repository_impl.go` (~15 queries)
- `stats_repository_impl.go` (si usa las tablas)

---

## TASK 2.6: Actualizar Tests

**Duración:** 1 hora

1. **Actualizar fixtures:**

```go
// ANTES:
db.Exec(`INSERT INTO school (id, name, code) VALUES (...)`)
db.Exec(`INSERT INTO academic_unit (id, school_id, unit_type, display_name) VALUES (...)`)

// DESPUÉS:
db.Exec(`INSERT INTO schools (id, name, code) VALUES (...)`)
db.Exec(`INSERT INTO academic_units (id, school_id, type, name, academic_year) VALUES (..., 2025)`)
```

2. **Actualizar seeds:**

```bash
# Modificar: scripts/postgresql/02_seeds_hierarchy.sql (si aún existe)
# O crear: internal/infrastructure/persistence/postgres/seeds/academic_hierarchy.go
```

3. **Ejecutar tests:**

```bash
go test ./... -v

# Si hay errores, revisar:
# - Nombres de tablas
# - Nombres de campos
# - Valores de academic_year
```

---

## TASK 2.7: Eliminar Migraciones Locales

**Duración:** 5 minutos

```bash
# Backup (por si acaso)
mv scripts/postgresql scripts/postgresql.backup

# O eliminar directamente
rm -rf scripts/postgresql/

# Actualizar .gitignore
echo "scripts/postgresql.backup/" >> .gitignore
```

---

## TASK 2.8: Actualizar Documentación

**Duración:** 20 minutos

### 1. Actualizar `README.md`:

```markdown
## Base de Datos

**IMPORTANTE:** Este proyecto usa migraciones centralizadas de `edugo-infrastructure`.

### Setup de Base de Datos

\`\`\`bash
# 1. Clonar infrastructure
git clone git@github.com:EduGoGroup/edugo-infrastructure.git
cd edugo-infrastructure

# 2. Levantar PostgreSQL
make dev-up-core

# 3. Ejecutar migraciones
cd postgres
make migrate-up

# ✅ Listo - BD con schema v0.7.0
\`\`\`

### Tablas Usadas

- `schools` - Escuelas/instituciones
- `academic_units` - Unidades académicas (con jerarquía)
- `memberships` - Asignaciones usuario-escuela-unidad

Ver: [TABLE_OWNERSHIP.md](https://github.com/EduGoGroup/edugo-infrastructure/blob/dev/database/TABLE_OWNERSHIP.md)
```

### 2. Actualizar `docs/DEVELOPMENT.md`:

Agregar prerequisito de infrastructure v0.7.0

---

## TASK 2.9: Build y Validación Final

**Duración:** 15 minutos

```bash
# 1. Build
go build ./...

# 2. Tests completos
go test ./... -v -race

# 3. Linting
golangci-lint run

# 4. Verificar imports
go mod verify
go mod tidy
```

---

## TASK 2.10: Commit y Push

**Duración:** 5 minutos

```bash
git add .
git commit -m "feat: migrate to infrastructure v0.7.0

- Update go.mod to infrastructure/postgres@v0.7.0
- Rename tables: school→schools, academic_unit→academic_units, unit_membership→memberships
- Rename fields: contact_email→email, display_name→name, unit_type→type
- Add academic_year field to AcademicUnit entity
- Update all repositories with new table/field names
- Remove local migrations (now in infrastructure)
- Update tests and fixtures
- Update documentation

BREAKING CHANGE: Requires infrastructure v0.7.0
Closes #XX"

git push origin feature/migrate-to-infrastructure-v0.7.0
```

---

## TASK 2.11: Crear Pull Request

**Duración:** 10 minutos

```markdown
# PR Title
feat: Migrate to infrastructure v0.7.0

## Descripción
Migración completa a `edugo-infrastructure` v0.7.0, eliminando migraciones locales y siguiendo el patrón de `api-mobile`.

## Cambios Principales
- ✅ Dependencia de infrastructure v0.7.0
- ✅ Actualización de shared a v0.7.0
- ✅ Renombrado de tablas (singular → plural)
- ✅ Renombrado de campos según estándar de infra
- ✅ Campo `academic_year` agregado
- ✅ Migraciones locales eliminadas
- ✅ Tests actualizados y pasando

## Testing
- [ ] Tests unitarios pasan (100%)
- [ ] Tests de integración pasan
- [ ] Build exitoso
- [ ] Linting sin errores

## Dependencias
- Requiere: `edugo-infrastructure@v0.7.0`

## Checklist
- [ ] Código revisado
- [ ] Tests pasan
- [ ] Documentación actualizada
- [ ] CHANGELOG.md actualizado (si aplica)

## Relacionado
- Sprint-00-Integrar-Infrastructure
- edugo-infrastructure#XX (PR de migración 012)
```

---

## ✅ Checklist FASE 2 Completada

- [ ] Rama de trabajo creada
- [ ] go.mod actualizado (infra v0.7.0, shared v0.7.0)
- [ ] Repositorios refactorizados (nombres de tablas y campos)
- [ ] Campo `academic_year` agregado a entity
- [ ] Queries SQL actualizadas (~50 cambios)
- [ ] Tests actualizados y pasando
- [ ] Migraciones locales eliminadas
- [ ] Documentación actualizada
- [ ] Build exitoso
- [ ] Commit y push realizados
- [ ] Pull Request creado

**Tiempo estimado FASE 2:** 4-5 horas

---

# 📊 Resumen Final

| Fase | Duración | Bloqueante | Output |
|------|----------|------------|--------|
| **FASE 1** | 3-4h | SÍ | infrastructure v0.7.0 |
| **FASE 2** | 4-5h | NO (requiere FASE 1) | api-admin migrado |
| **TOTAL** | **7-9h** | - | Migración completa |

---

# 🎯 Criterios de Éxito

Sprint-00 está **COMPLETADO** cuando:

1. ✅ Infrastructure v0.7.0 publicado en GitHub
2. ✅ `go get ...@v0.7.0` funciona
3. ✅ api-admin usa tablas de infrastructure
4. ✅ NO existen migraciones locales en `scripts/`
5. ✅ Tests pasan 100%
6. ✅ PR mergeado a `dev`

---

**Documento creado:** 17 de Noviembre, 2025  
**Versión:** 2.0.0 (Plan Completo)  
**Próximo paso:** Ejecutar FASE 1
