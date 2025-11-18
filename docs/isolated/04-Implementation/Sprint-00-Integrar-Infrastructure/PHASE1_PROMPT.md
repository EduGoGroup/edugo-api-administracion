# PROMPT FASE 1 - SPRINT-00 INFRASTRUCTURE v0.7.0

**Proyecto:** edugo-api-administracion (API de Jerarquía Académica)  
**Ejecutor:** Claude Code Web (Desatendido)  
**Duración estimada:** 2-3 horas  
**Branch destino:** `feature/sprint-00-infrastructure-v0.7.0`

---

## 🎯 CONTEXTO DEL SPRINT

Eres Claude Code Web ejecutando la **Fase 1 del Sprint-00: Integrar Infrastructure v0.7.0**.

### Estado Actual del Proyecto

**Proyecto:** edugo-api-administracion  
**Versión actual:** v0.3.0 (Release completado)  
**Estado funcional:** ✅ COMPLETO (v0.2.0) - API funcionando  
**Estado técnico:** ⚠️ REQUIERE MIGRACIÓN a infrastructure v0.7.0

**Sprints completados:**
- ✅ Sprint-01 a Sprint-06: Funcionalidad completa (Schema, Dominio, Repos, Services, Testing, CI/CD)
- ✅ Sprint-07: CI/CD Fixes (v0.3.0)
- 🎯 **Sprint-00 FASE 1: Actualizar Infrastructure** ← **TU MISIÓN**

### ¿Por Qué Este Sprint?

El proyecto api-administracion está **funcionalmente completo** pero usa:
- ❌ Migraciones locales en `scripts/postgresql/` (desactualizadas)
- ❌ edugo-shared v0.7.0 (correcto pero infrastructure desactualizado)
- ❌ infrastructure v0.6.0 (NO tiene soporte de jerarquía académica)

**Objetivo:** Migrar infrastructure a v0.7.0 con soporte completo de jerarquía académica.

---

## 📂 UBICACIÓN DE LA DOCUMENTACIÓN

**Toda la documentación de este sprint está en:**
```
/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/
```

**Archivos clave:**
- `README.md` - Overview del sprint
- `TASKS.md` - Lista de tareas
- `FASE1_ACTUALIZAR_INFRASTRUCTURE.md` - **TU GUÍA PRINCIPAL**
- `migrations/*.sql` - Migraciones completas ya preparadas
- `seeds/*.sql` - Seeds con jerarquía ya preparados

---

## 🚀 TU MISIÓN (FASE 1)

### Objetivo Principal

Modificar el repositorio **edugo-infrastructure** para publicar la versión **v0.7.0** con soporte completo de:
1. **Jerarquía académica** (parent_unit_id en academic_units)
2. **Metadata JSONB** en schools, academic_units, memberships
3. **Función anti-ciclos** (prevent_academic_unit_cycles)
4. **Vista recursiva** (v_academic_unit_tree)
5. **Tipos y roles extendidos** (school, club, department, coordinator, admin, assistant)

### Filosofía: "Plomo al Ampa"

> **NO parcheamos, hacemos bien desde el inicio**

- ✅ Modificar migraciones existentes (002, 003, 004) para que queden **completas**
- ✅ Cualquier proyecto nuevo obtiene **TODO de una vez**
- ✅ **UNA SOLA VERDAD** en infrastructure
- ❌ NO crear migración 012 (sería un parche)

### Justificación

- Estamos en **desarrollo** (podemos recrear BD sin problema)
- api-mobile **NO se rompe** (campos nuevos son opcionales con defaults)
- Futuro worker obtiene **schema completo** desde el inicio
- Mantenimiento **más simple** (menos migraciones)

---

## 📋 PLAN DE EJECUCIÓN DETALLADO

### CONFIGURACIÓN INICIAL

**NO necesitas crear branch en api-administracion aún.** Tu trabajo es en **edugo-infrastructure**.

```bash
# Navegar a infrastructure
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Verificar rama actual
git status

# Asegurarte de estar en dev y actualizado
git checkout dev
git pull origin dev
```

---

### TASK 1.1: Backup de Migraciones Actuales (5 min)

**Objetivo:** Preservar versiones anteriores por seguridad.

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Crear directorio de backup con timestamp
mkdir -p postgres/migrations/.backup-$(date +%Y%m%d-%H%M%S)

# Copiar migraciones actuales
cp postgres/migrations/002_create_schools.up.sql postgres/migrations/.backup-$(date +%Y%m%d-%H%M%S)/
cp postgres/migrations/002_create_schools.down.sql postgres/migrations/.backup-$(date +%Y%m%d-%H%M%S)/
cp postgres/migrations/003_create_academic_units.up.sql postgres/migrations/.backup-$(date +%Y%m%d-%H%M%S)/
cp postgres/migrations/003_create_academic_units.down.sql postgres/migrations/.backup-$(date +%Y%m%d-%H%M%S)/
cp postgres/migrations/004_create_memberships.up.sql postgres/migrations/.backup-$(date +%Y%m%d-%H%M%S)/
cp postgres/migrations/004_create_memberships.down.sql postgres/migrations/.backup-$(date +%Y%m%d-%H%M%S)/

echo "✅ Backup creado en postgres/migrations/.backup-$(date +%Y%m%d-%H%M%S)/"
```

**Validación:**
```bash
ls -la postgres/migrations/.backup-*/
# Debe mostrar 6 archivos (.up.sql y .down.sql para 002, 003, 004)
```

---

### TASK 1.2: Copiar Migraciones Completas (10 min)

**Objetivo:** Reemplazar migraciones 002, 003, 004 con versiones completas que incluyen jerarquía.

**Origen:** Las migraciones ya están preparadas en api-administracion  
**Destino:** edugo-infrastructure

```bash
# Definir rutas (CAMBIAR SI ES NECESARIO)
SRC="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/migrations"
DEST="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/postgres/migrations"

# Verificar que origen existe
ls -la $SRC/002_create_schools.up.sql
ls -la $SRC/003_create_academic_units.up.sql
ls -la $SRC/004_create_memberships.up.sql

# Si los archivos existen, copiar (SOBRESCRIBIR)
cp $SRC/002_create_schools.up.sql $DEST/
cp $SRC/002_create_schools.down.sql $DEST/
cp $SRC/003_create_academic_units.up.sql $DEST/
cp $SRC/003_create_academic_units.down.sql $DEST/
cp $SRC/004_create_memberships.up.sql $DEST/
cp $SRC/004_create_memberships.down.sql $DEST/

echo "✅ Migraciones copiadas y sobrescritas"
```

**Validación:**
```bash
# Verificar que las migraciones tienen el contenido nuevo
grep -q "metadata JSONB" $DEST/002_create_schools.up.sql && echo "✅ 002: metadata agregado"
grep -q "parent_unit_id" $DEST/003_create_academic_units.up.sql && echo "✅ 003: parent_unit_id agregado"
grep -q "prevent_academic_unit_cycles" $DEST/003_create_academic_units.up.sql && echo "✅ 003: función anti-ciclos agregada"
grep -q "v_academic_unit_tree" $DEST/003_create_academic_units.up.sql && echo "✅ 003: vista recursiva agregada"
```

---

### TASK 1.3: Copiar Seeds Actualizados (5 min)

**Objetivo:** Agregar seeds con ejemplos de jerarquía académica.

```bash
SRC="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/seeds"
DEST="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/postgres/seeds"

# Verificar que origen existe
ls -la $SRC/academic_units.sql
ls -la $SRC/memberships.sql

# Copiar (SOBRESCRIBIR si existen)
cp $SRC/academic_units.sql $DEST/
cp $SRC/memberships.sql $DEST/

echo "✅ Seeds copiados"
```

**Validación:**
```bash
grep -q "parent_unit_id" $DEST/academic_units.sql && echo "✅ Seeds tienen jerarquía"
grep -q "coordinator" $DEST/memberships.sql && echo "✅ Seeds tienen roles extendidos"
```

---

### TASK 1.4: Actualizar CHANGELOG.md (15 min)

**Objetivo:** Documentar los cambios en infrastructure v0.7.0.

**Archivo:** `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/CHANGELOG.md`

**Acción:** Agregar al **inicio** del archivo (antes de cualquier otra versión):

```markdown
## [0.7.0] - 2025-11-17

### Added (postgres)

**Soporte completo de jerarquía académica:**
- Campo `parent_unit_id UUID` en `academic_units` para estructura jerárquica (self-referencing FK)
- Campos `metadata JSONB DEFAULT '{}'::jsonb` en `schools`, `academic_units`, `memberships` (extensibilidad)
- Campo `description TEXT` en `academic_units` (descripciones detalladas)
- Función `prevent_academic_unit_cycles()` con trigger `check_academic_unit_cycle` (prevención de ciclos en árbol)
- Vista `v_academic_unit_tree` usando CTE recursivo (consultas de jerarquía completa)

**Extensiones de tipos y roles:**
- Tipos de `academic_units.type`: Agregados `'school'`, `'club'`, `'department'` (además de `'grade'`, `'class'`, `'section'`)
- Roles de `memberships.role`: Agregados `'coordinator'`, `'admin'`, `'assistant'` (además de `'teacher'`, `'student'`, `'guardian'`)

### Changed (postgres)

- **`academic_units.academic_year`**: Ahora es `nullable` con `DEFAULT 0` (antes era `NOT NULL`)
  - Razón: 0 = "sin año académico específico" (ej: departamentos, clubes no tienen año)
- **Constraints CHECK extendidos**: Backward compatible (valores anteriores siguen siendo válidos)

### Seeds

- `academic_units.sql`: Ejemplos de jerarquía completa
  - Escuela → Grado → Sección
  - Escuela → Departamento → Clase
  - Escuela → Club
- `memberships.sql`: Ejemplos de todos los roles (teacher, student, guardian, **coordinator**, **admin**, **assistant**)

### Migration Path

**Para proyectos existentes con infrastructure v0.6.0:**
1. ⚠️ **DROP DATABASE requerido** (estamos en desarrollo, no hay migración incremental)
2. Ejecutar `make dev-down && make dev-up-core` (recrear contenedores)
3. Ejecutar `make migrate-up` (aplicar migraciones 001-011, incluyendo 002-004 modificadas)
4. Los nuevos campos tienen defaults seguros → **Backward compatible** con api-mobile

**Para proyectos nuevos:**
- Ejecutar `make migrate-up` obtiene schema completo de una vez

### Breaking Changes

- **`academic_units.academic_year`**: Ahora acepta `NULL` y `0` (antes era `NOT NULL`)
  - Impacto: Si api-mobile asume `NOT NULL`, necesita actualizar queries
  - Mitigación: Usar `COALESCE(academic_year, 0)` en queries existentes
- **Recreación de BD requerida** desde v0.6.0 (DROP DATABASE + make dev-up-core)

### Compatibility

- ✅ **api-mobile**: Sin cambios requeridos (campos nuevos son opcionales)
- ✅ **api-admin**: Puede usar jerarquía completa (parent_unit_id, vista recursiva)
- ✅ **worker** (futuro): Obtiene schema completo desde inicio

### Notes

- **Enfoque adoptado:** Modificar migraciones existentes (002, 003, 004) en vez de crear migración 012
- **Razón:** "Plomo al ampa" - Una sola migración completa, no parcheada
- **NO se creó migración 012:** Enfoque de parche descartado en favor de completitud
- **Filosofía:** Una sola verdad en infrastructure, no fragmentada

### Technical Details

**Función prevent_academic_unit_cycles:**
```sql
CREATE OR REPLACE FUNCTION prevent_academic_unit_cycles()
RETURNS TRIGGER AS $$
BEGIN
  -- Detecta ciclos usando CTE recursivo
  -- Lanza excepción si se detecta ciclo
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

**Vista v_academic_unit_tree:**
```sql
CREATE VIEW v_academic_unit_tree AS
WITH RECURSIVE tree AS (
  -- Nodos raíz (parent_unit_id IS NULL)
  SELECT id, parent_unit_id, name, type, 1 as level, ARRAY[id] as path
  FROM academic_units WHERE parent_unit_id IS NULL
  UNION ALL
  -- Nodos hijos (recursión)
  SELECT au.id, au.parent_unit_id, au.name, au.type, t.level + 1, t.path || au.id
  FROM academic_units au
  INNER JOIN tree t ON au.parent_unit_id = t.id
)
SELECT * FROM tree ORDER BY path;
```

---

**Autores:** edugo-api-administracion team  
**Reviewed by:** infrastructure team  
**Breaking Change Approved:** Yes (desarrollo solamente)
```

**Validación:**
```bash
head -100 /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/CHANGELOG.md | grep "0.7.0"
# Debe mostrar la nueva versión al inicio
```

---

### TASK 1.5: Commit y Push (10 min)

**Objetivo:** Publicar cambios en branch `dev` de infrastructure.

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Verificar cambios
git status
# Debe mostrar:
# - postgres/migrations/002_*, 003_*, 004_* (modified)
# - postgres/seeds/academic_units.sql, memberships.sql (modified o new)
# - CHANGELOG.md (modified)
# - postgres/migrations/.backup-* (untracked - NO agregar)

# Agregar solo los archivos relevantes
git add postgres/migrations/002_create_schools.up.sql
git add postgres/migrations/002_create_schools.down.sql
git add postgres/migrations/003_create_academic_units.up.sql
git add postgres/migrations/003_create_academic_units.down.sql
git add postgres/migrations/004_create_memberships.up.sql
git add postgres/migrations/004_create_memberships.down.sql
git add postgres/seeds/academic_units.sql
git add postgres/seeds/memberships.sql
git add CHANGELOG.md

# Verificar que NO se agregó .backup
git status | grep -v backup

# Commit con mensaje semántico detallado
git commit -m "feat(postgres): extender schema con jerarquía académica completa para v0.7.0

## Cambios en Migraciones

### 002_create_schools
- Add: metadata JSONB DEFAULT '{}'::jsonb (extensibilidad para datos custom)

### 003_create_academic_units  
- Add: parent_unit_id UUID REFERENCES academic_units(id) (jerarquía self-referencing)
- Add: description TEXT (descripciones detalladas de unidades)
- Add: metadata JSONB DEFAULT '{}'::jsonb (extensibilidad)
- Add: prevent_academic_unit_cycles() función PL/pgSQL (validación anti-ciclos)
- Add: check_academic_unit_cycle TRIGGER (ejecuta función en INSERT/UPDATE)
- Add: v_academic_unit_tree VIEW (CTE recursivo para consultar árbol completo)
- Change: type CHECK constraint extendido ('school', 'club', 'department' agregados)
- Change: academic_year ahora nullable con DEFAULT 0 (0 = sin año específico)

### 004_create_memberships
- Add: metadata JSONB DEFAULT '{}'::jsonb (extensibilidad)
- Change: role CHECK constraint extendido ('coordinator', 'admin', 'assistant' agregados)

## Seeds Actualizados

- academic_units.sql: Ejemplos con jerarquía (Escuela → Grado, Escuela → Depto, Escuela → Club)
- memberships.sql: Ejemplos de todos los roles (incluyendo coordinator, admin, assistant)

## Compatibilidad

✅ Backward compatible con api-mobile v0.2.0
  - Campos nuevos son opcionales con defaults seguros
  - Constraints CHECK extendidos (no restrictivos)
  - academic_year acepta NULL (api-mobile usa NOT NULL, pero no se rompe)

✅ api-admin v0.3.0 puede usar jerarquía completa
  - parent_unit_id disponible
  - Vista v_academic_unit_tree lista para usar
  - Función anti-ciclos protege integridad

✅ Futuro worker obtiene schema completo desde inicio

## Breaking Changes

⚠️ REQUIERE DROP DATABASE si vienes de v0.6.0
  - Razón: Modificamos migraciones existentes, no creamos migración incremental
  - Impacto: Solo desarrollo (sin datos de producción aún)
  - Solución: make dev-down && make dev-up-core && make migrate-up

⚠️ academic_year ahora nullable
  - Antes: NOT NULL
  - Ahora: NULL permitido, DEFAULT 0
  - Impacto: Queries que asumen NOT NULL deben usar COALESCE(academic_year, 0)
  - Proyectos afectados: api-mobile (bajo impacto, solo queries de filtrado)

## Filosofía de Diseño

\"Plomo al ampa\" - Hacemos bien desde el inicio
- ✅ Una sola migración completa (no parcheada)
- ✅ Una sola verdad en infrastructure
- ❌ NO se creó migración 012 (enfoque de parche descartado)
- ✅ Futuro mantenimiento más simple

## Testing Requerido en Fase 2

- [ ] Levantar PostgreSQL con make dev-up-core
- [ ] Ejecutar migraciones: make migrate-up
- [ ] Validar columnas nuevas existen (metadata, parent_unit_id, description)
- [ ] Validar función prevent_academic_unit_cycles existe
- [ ] Validar vista v_academic_unit_tree existe
- [ ] Test anti-ciclos: Intentar UPDATE para crear ciclo (debe fallar)
- [ ] Test vista recursiva: SELECT * FROM v_academic_unit_tree (debe mostrar árbol)
- [ ] Test rollback: Ejecutar .down.sql y verificar tablas eliminadas
- [ ] Test re-aplicar: Ejecutar .up.sql y verificar tablas creadas
- [ ] Ejecutar seeds y validar datos

BREAKING CHANGE: academic_year ahora nullable, requiere DROP DATABASE desde v0.6.0

Co-authored-by: edugo-api-administracion team <noreply@edugo.com>
Co-Authored-By: Claude Code <noreply@anthropic.com>"

# Push a dev
git push origin dev

echo "✅ Cambios pusheados a dev"
```

**Validación:**
```bash
# Verificar que commit existe en remoto
git log origin/dev -1 --oneline | grep "feat(postgres)"
```

---

### TASK 1.6: Crear Tag v0.7.0 (10 min)

**Objetivo:** Publicar versión v0.7.0 de infrastructure para que api-administracion pueda consumirla.

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Asegurarse de estar en dev actualizado
git checkout dev
git pull origin dev

# Crear tag anotado con descripción completa
git tag -a v0.7.0 -m "Release v0.7.0: Schema PostgreSQL con jerarquía académica completa

## 🎯 Highlights

- ✅ Soporte completo de jerarquía en academic_units (parent_unit_id + vista recursiva)
- ✅ Metadata JSONB en todas las tablas core (schools, academic_units, memberships)
- ✅ Tipos y roles extendidos (school, club, department, coordinator, admin, assistant)
- ✅ Función anti-ciclos con trigger (prevent_academic_unit_cycles)
- ✅ Vista recursiva para consultar árboles (v_academic_unit_tree)

## 📝 Modified Migrations

- 002_create_schools: + metadata JSONB
- 003_create_academic_units: + parent_unit_id, description, metadata, función anti-ciclos, vista recursiva, tipos extendidos, academic_year nullable
- 004_create_memberships: + metadata JSONB, roles extendidos

## 🌱 New Seeds

- academic_units.sql: Jerarquía demo (Escuela → Grado → Sección, Escuela → Depto, Escuela → Club)
- memberships.sql: Todos los roles (teacher, student, guardian, coordinator, admin, assistant)

## ✅ Compatible With

- api-mobile v0.2.0 ✅ (sin cambios requeridos, campos opcionales)
- api-admin v0.3.0 ✅ (puede usar jerarquía completa)
- worker (futuro) ✅ (obtiene schema completo)

## ⚠️ Breaking Changes

- REQUIERE DROP DATABASE desde v0.6.0 (desarrollo solamente)
- academic_year ahora nullable (antes NOT NULL)

## 🔧 Migration Path

Desde v0.6.0:
1. make dev-down
2. make dev-up-core
3. make migrate-up
4. make seed (opcional)

Proyecto nuevo:
1. make dev-up-core
2. make migrate-up
3. make seed (opcional)

## 📚 Documentation

Ver CHANGELOG.md para detalles completos.

## 🤝 Contributors

- edugo-api-administracion team
- infrastructure team

---
Generated with workflow orchestration v2.0
Philosophy: Plomo al ampa - One complete migration, not patched"

# Push tag
git push origin v0.7.0

echo "✅ Tag v0.7.0 creado y pusheado"
```

**Validación:**
```bash
# Verificar tag en remoto
git ls-remote --tags origin | grep v0.7.0
# Debe mostrar: refs/tags/v0.7.0

# Verificar que GitHub procesó el tag (esperar 30 segundos)
sleep 30
```

---

### TASK 1.7: Validar Disponibilidad de v0.7.0 (10 min)

**Objetivo:** Confirmar que infrastructure v0.7.0 está disponible para consumo vía `go get`.

```bash
# Crear directorio temporal para testing
cd /tmp
mkdir -p test-infrastructure-v0.7.0
cd test-infrastructure-v0.7.0

# Inicializar módulo Go temporal
go mod init test-infra

# Intentar obtener infrastructure v0.7.0
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.0

# Validar que se descargó correctamente
if [ $? -eq 0 ]; then
  echo "✅ infrastructure v0.7.0 disponible y descargable"
  
  # Verificar contenido
  cat go.mod | grep "edugo-infrastructure/postgres v0.7.0"
  
  # Limpiar
  cd /tmp
  rm -rf test-infrastructure-v0.7.0
  
  echo "✅ TASK 1.7 COMPLETADA"
else
  echo "❌ ERROR: infrastructure v0.7.0 NO disponible"
  echo "Posibles causas:"
  echo "- Tag no pusheado correctamente"
  echo "- GitHub aún no procesó el tag (esperar 1-2 minutos)"
  echo "- Problemas de conectividad"
  exit 1
fi
```

**Si falla:** Esperar 1-2 minutos adicionales y reintentar. GitHub a veces tarda en indexar tags nuevos.

---

### TASK 1.8: Generar PHASE2_BRIDGE.md (20 min)

**Objetivo:** Documentar TODO lo necesario para que Claude Code Local ejecute Fase 2.

**Ubicación:** `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/PHASE2_BRIDGE.md`

**Contenido:**

```markdown
# PHASE 2 BRIDGE - Sprint-00 Infrastructure v0.7.0

**Fecha Fase 1:** [TIMESTAMP_ACTUAL]  
**Ejecutor Fase 1:** Claude Code Web  
**Estado Fase 1:** ✅ COMPLETADA  
**Versión Infrastructure:** v0.7.0  

---

## 📊 RESUMEN DE FASE 1

### ✅ Completado al 100%

- [x] **TASK 1.1**: Backup de migraciones actuales en infrastructure
- [x] **TASK 1.2**: Migraciones 002, 003, 004 copiadas y sobrescritas
- [x] **TASK 1.3**: Seeds academic_units.sql y memberships.sql copiados
- [x] **TASK 1.4**: CHANGELOG.md actualizado con v0.7.0
- [x] **TASK 1.5**: Commit y push a dev en infrastructure
- [x] **TASK 1.6**: Tag v0.7.0 creado y pusheado
- [x] **TASK 1.7**: Validado `go get ...@v0.7.0` funciona

### 📦 Publicado

- ✅ **infrastructure v0.7.0** disponible en GitHub
- ✅ Migraciones 002, 003, 004 con soporte de jerarquía
- ✅ Seeds con ejemplos de jerarquía

### 🔄 Cambios Realizados en Infrastructure

**Repositorio:** edugo-infrastructure  
**Branch:** dev  
**Tag:** v0.7.0

**Archivos modificados:**
- `postgres/migrations/002_create_schools.up.sql` - Agregado metadata JSONB
- `postgres/migrations/002_create_schools.down.sql` - Rollback actualizado
- `postgres/migrations/003_create_academic_units.up.sql` - Agregado parent_unit_id, metadata, función, vista, tipos extendidos
- `postgres/migrations/003_create_academic_units.down.sql` - Rollback actualizado
- `postgres/migrations/004_create_memberships.up.sql` - Agregado metadata, roles extendidos
- `postgres/migrations/004_create_memberships.down.sql` - Rollback actualizado
- `postgres/seeds/academic_units.sql` - Seeds con jerarquía
- `postgres/seeds/memberships.sql` - Seeds con roles extendidos
- `CHANGELOG.md` - Documentación de v0.7.0

---

## 🎯 FASE 2: MIGRAR API-ADMINISTRACION

### Objetivo

Actualizar **edugo-api-administracion** para usar **infrastructure v0.7.0** y eliminar migraciones locales.

### Prerequisitos

✅ infrastructure v0.7.0 publicado (completado en Fase 1)  
⏳ Docker Desktop corriendo (requerido en Fase 2)  
⏳ PostgreSQL 15+ accesible (requerido en Fase 2)

---

## 📋 TAREAS PENDIENTES PARA FASE 2

### TASK 2.1: Actualizar go.mod en api-administracion

**Archivo:** `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/go.mod`

**Acción:** Actualizar dependencia de infrastructure a v0.7.0

**Código actual (estimado):**
```go
require (
    github.com/EduGoGroup/edugo-infrastructure/postgres v0.6.0
)
```

**Código requerido:**
```go
require (
    github.com/EduGoGroup/edugo-infrastructure/postgres v0.7.0
)
```

**Comando:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.0
go mod tidy
```

**Validación:**
```bash
cat go.mod | grep "edugo-infrastructure/postgres v0.7.0"
# Debe mostrar la línea con v0.7.0
```

---

### TASK 2.2: Eliminar Migraciones Locales (OPCIONAL)

**Archivos a eliminar:**
```
/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/scripts/postgresql/
├── 06_schools.sql
├── 07_academic_units.sql
├── 08_unit_memberships.sql
└── ... (cualquier otro script de migración local)
```

**Razón:** infrastructure v0.7.0 ya tiene estas migraciones completas.

**Comando:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Verificar qué archivos existen
ls -la scripts/postgresql/

# Eliminar (SOLO si existen y son redundantes)
rm -f scripts/postgresql/06_schools.sql
rm -f scripts/postgresql/07_academic_units.sql
rm -f scripts/postgresql/08_unit_memberships.sql

# Si el directorio queda vacío, eliminarlo
rmdir scripts/postgresql/ 2>/dev/null || true

echo "✅ Migraciones locales eliminadas (ahora usamos infrastructure)"
```

**Nota:** Si hay scripts de migración que NO están en infrastructure, **NO eliminarlos**. Documentar en EXECUTION_REPORT.md.

---

### TASK 2.3: Actualizar Configuración de Base de Datos (SI APLICA)

**Verificar:** Si api-administracion tiene configuración custom de conexión a PostgreSQL.

**Archivos a revisar:**
- `internal/config/database.go`
- `.env` o `.env.local`
- `config/*.yaml`

**Acción:** Asegurar que pool de conexiones está configurado correctamente.

**Configuración recomendada:**
```go
// internal/config/database.go
type DatabaseConfig struct {
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func DefaultDatabaseConfig() *DatabaseConfig {
    return &DatabaseConfig{
        MaxOpenConns:    25,
        MaxIdleConns:    5,
        ConnMaxLifetime: 5 * time.Minute,
        ConnMaxIdleTime: 10 * time.Minute,
    }
}
```

**Validación:**
```bash
grep -n "MaxOpenConns" internal/config/database.go
# Si existe, verificar valores son razonables
```

---

### TASK 2.4: Recrear Base de Datos Local con Infrastructure v0.7.0

**⚠️ REQUIERE DOCKER Y POSTGRESQL CORRIENDO**

**Prerequisito:**
```bash
# Verificar Docker está corriendo
docker ps
# Si falla, iniciar Docker Desktop
```

**Acción:** Recrear base de datos para aplicar migraciones v0.7.0.

**Opción A: Usar Makefile de infrastructure (RECOMENDADO)**

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Bajar containers actuales (DESTRUYE DATOS - desarrollo solamente)
make dev-down

# Levantar PostgreSQL fresco
make dev-up-core

# Esperar 10 segundos
sleep 10

# Aplicar migraciones (001-011, incluyendo 002-004 modificadas)
cd postgres
make migrate-up

# Validar que las nuevas columnas existen
psql -U edugo -d edugo_dev -c "\d schools" | grep metadata
# Debe mostrar: metadata | jsonb | | | 

psql -U edugo -d edugo_dev -c "\d academic_units" | grep parent_unit_id
# Debe mostrar: parent_unit_id | uuid | | |

psql -U edugo -d edugo_dev -c "\d academic_units" | grep description
# Debe mostrar: description | text | | |

# Validar función anti-ciclos existe
psql -U edugo -d edugo_dev -c "\df prevent_academic_unit_cycles"
# Debe mostrar la función

# Validar vista recursiva existe
psql -U edugo -d edugo_dev -c "\dv v_academic_unit_tree"
# Debe mostrar la vista

echo "✅ Base de datos recreada con infrastructure v0.7.0"
```

**Opción B: Usar Docker Compose directamente**

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Bajar y eliminar volúmenes
docker compose down -v

# Levantar solo PostgreSQL
docker compose up -d postgres

# Esperar
sleep 10

# Ejecutar migraciones manualmente
psql -U edugo -h localhost -p 5432 -d edugo_dev < postgres/migrations/001_*.up.sql
psql -U edugo -h localhost -p 5432 -d edugo_dev < postgres/migrations/002_*.up.sql
# ... (hasta 011)

# Validar (igual que Opción A)
```

---

### TASK 2.5: Ejecutar Tests de Integración

**⚠️ REQUIERE POSTGRESQL CORRIENDO**

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Ejecutar tests de integración
make test-integration

# O manualmente:
go test -v -tags=integration ./test/integration/...

# Validar que pasan
# Si fallan:
# - Revisar logs
# - Verificar que PostgreSQL tiene las nuevas columnas
# - Verificar go.mod usa infrastructure v0.7.0
```

**Tests esperados:**
- ✅ Tests de repositorios (SchoolRepository, AcademicUnitRepository, MembershipRepository)
- ✅ Tests de consultas recursivas (jerarquía)
- ✅ Tests de función anti-ciclos

**Si algún test falla:**
1. Analizar error
2. Corregir (max 3 intentos por error único)
3. Si error se repite 3 veces: **DETENER e INFORMAR**

---

### TASK 2.6: Ejecutar Tests Unitarios

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Ejecutar tests unitarios
make test

# O manualmente:
go test -v ./...

# Validar que pasan (deben pasar todos porque no cambiamos lógica de negocio)
```

---

### TASK 2.7: Ejecutar Build Completo

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Build completo
make build

# O manualmente:
go build -v ./...

# Debe compilar sin errores
```

---

### TASK 2.8: Verificar Jerarquía Funciona (NUEVO)

**Objetivo:** Validar que la API puede usar las nuevas capacidades de jerarquía.

**Test manual con PostgreSQL:**

```bash
# Insertar unidad académica con parent
psql -U edugo -d edugo_dev -c "
INSERT INTO academic_units (id, school_id, name, type, parent_unit_id)
VALUES (
  gen_random_uuid(),
  (SELECT id FROM schools LIMIT 1),
  'Departamento de Matemáticas',
  'department',
  NULL
)
RETURNING id;
"

# Guardar el ID devuelto como PARENT_ID

# Insertar unidad hija
psql -U edugo -d edugo_dev -c "
INSERT INTO academic_units (id, school_id, name, type, parent_unit_id)
VALUES (
  gen_random_uuid(),
  (SELECT id FROM schools LIMIT 1),
  'Cálculo I',
  'class',
  'PARENT_ID_AQUI'
);
"

# Consultar vista recursiva
psql -U edugo -d edugo_dev -c "SELECT * FROM v_academic_unit_tree LIMIT 10;"

# Debe mostrar jerarquía completa
```

**Test de ciclos (debe FALLAR):**

```bash
# Intentar crear ciclo (debe dar error)
psql -U edugo -d edugo_dev -c "
UPDATE academic_units
SET parent_unit_id = id
WHERE id = (SELECT id FROM academic_units LIMIT 1);
"

# Debe dar error: "Ciclo detectado en jerarquía de unidades académicas"
# Si NO da error: REPORTAR PROBLEMA
```

---

### TASK 2.9: Commit y Push (Fase 2)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Crear branch
git checkout dev
git pull origin dev
git checkout -b feature/sprint-00-infrastructure-v0.7.0

# Agregar cambios
git add go.mod go.sum
git add scripts/postgresql/  # Solo si se eliminaron archivos

# Commit
git commit -m "feat(infra): migrar a infrastructure v0.7.0 con soporte de jerarquía

## Cambios

- Update: edugo-infrastructure/postgres v0.6.0 → v0.7.0
- Remove: Migraciones locales (ahora en infrastructure)

## Nuevas Capacidades (Disponibles)

- ✅ Jerarquía académica (parent_unit_id en academic_units)
- ✅ Metadata JSONB en schools, academic_units, memberships
- ✅ Función anti-ciclos (prevent_academic_unit_cycles)
- ✅ Vista recursiva (v_academic_unit_tree)
- ✅ Tipos extendidos (school, club, department)
- ✅ Roles extendidos (coordinator, admin, assistant)

## Testing

- ✅ Tests unitarios: Pasando
- ✅ Tests integración: Pasando
- ✅ Build: Exitoso
- ✅ Jerarquía validada: Funcional
- ✅ Anti-ciclos validado: Funcional

## Breaking Changes

- Requiere recrear BD local (make dev-down && make dev-up-core)
- academic_year ahora nullable (bajo impacto)

BREAKING CHANGE: infrastructure v0.6.0 → v0.7.0 requiere DROP DATABASE

Co-authored-by: infrastructure team <noreply@edugo.com>
Co-Authored-By: Claude Code <noreply@anthropic.com>"

# Push
git push origin feature/sprint-00-infrastructure-v0.7.0

echo "✅ Branch pusheado"
```

---

### TASK 2.10: Crear Pull Request a dev

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Crear PR con GitHub CLI
gh pr create \
  --title "feat: Migrar a infrastructure v0.7.0 con jerarquía académica" \
  --body "## Sprint-00 FASE 2: Infrastructure v0.7.0

### 🎯 Objetivo

Migrar api-administracion de infrastructure v0.6.0 a v0.7.0 para usar schema completo con jerarquía académica.

### ✅ Cambios Realizados

- **go.mod**: Actualizado infrastructure/postgres v0.6.0 → v0.7.0
- **scripts/postgresql/**: Eliminadas migraciones locales (ahora en infrastructure)

### 🚀 Nuevas Capacidades Disponibles

- ✅ **Jerarquía académica**: Campo \`parent_unit_id\` en \`academic_units\`
- ✅ **Metadata JSONB**: En \`schools\`, \`academic_units\`, \`memberships\`
- ✅ **Función anti-ciclos**: \`prevent_academic_unit_cycles()\` con trigger
- ✅ **Vista recursiva**: \`v_academic_unit_tree\` para consultas de árbol
- ✅ **Tipos extendidos**: \`school\`, \`club\`, \`department\` en \`academic_units.type\`
- ✅ **Roles extendidos**: \`coordinator\`, \`admin\`, \`assistant\` en \`memberships.role\`

### 🧪 Testing

- ✅ Tests unitarios: **PASSING**
- ✅ Tests integración: **PASSING**
- ✅ Build completo: **SUCCESS**
- ✅ Validación jerarquía: **FUNCTIONAL**
- ✅ Validación anti-ciclos: **FUNCTIONAL**

### ⚠️ Breaking Changes

- **Requiere recrear BD local**: \`make dev-down && make dev-up-core\`
- **academic_year ahora nullable**: Antes NOT NULL, ahora acepta NULL y DEFAULT 0
  - Impacto: Bajo (queries existentes siguen funcionando)

### 📝 Notas

- Infrastructure v0.7.0 publicado en Fase 1 (ver edugo-infrastructure)
- Enfoque \"plomo al ampa\": Una sola migración completa, no parcheada
- Backward compatible con api-mobile

### 🔗 Referencias

- Infrastructure CHANGELOG: https://github.com/EduGoGroup/edugo-infrastructure/blob/dev/CHANGELOG.md
- Sprint-00 docs: \`docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/\`

---
**Sprint:** 00/07 (Infrastructure)  
**Fase:** 2/2 (Fase 1: infrastructure publicado, Fase 2: api-admin migrado)  
**Workflow:** Orchestration 2.0

Co-authored-by: infrastructure team" \
  --base dev \
  --head feature/sprint-00-infrastructure-v0.7.0

# Obtener número de PR
PR_NUMBER=$(gh pr view --json number -q '.number')
echo "✅ PR creado: #$PR_NUMBER"
```

---

### TASK 2.11: Monitorear CI/CD (5 min MAX)

**⏱️ REGLAS ESTRICTAS:**
- Tiempo máximo: **5 minutos**
- Revisar cada: **1 minuto**
- Si pasan 5 min sin completar: **DETENER e INFORMAR**
- Si falla: Corregir (max 3 intentos por error único)
- Si error se repite 3x: **DETENER e INFORMAR**

```bash
#!/bin/bash
MAX_WAIT=300  # 5 minutos
CHECK_INTERVAL=60  # 1 minuto
START_TIME=$(date +%s)
declare -A ERROR_COUNTS

echo "⏳ Iniciando monitoreo CI/CD (max 5 minutos)..."

while true; do
  CURRENT_TIME=$(date +%s)
  ELAPSED=$((CURRENT_TIME - START_TIME))

  if [ $ELAPSED -gt $MAX_WAIT ]; then
    echo "❌ TIMEOUT: CI/CD no completó en 5 minutos"
    gh pr checks
    echo ""
    echo "ACCIÓN REQUERIDA:"
    echo "1. Revisar workflows manualmente en GitHub"
    echo "2. Verificar que no haya tests colgados"
    echo "3. Considerar cancelar y reiniciar workflows"
    exit 1
  fi

  # Obtener estado de checks
  CHECKS_JSON=$(gh pr checks --json state,name,conclusion)
  IN_PROGRESS=$(echo "$CHECKS_JSON" | jq '[.[] | select(.state=="IN_PROGRESS")] | length')

  if [ "$IN_PROGRESS" -eq 0 ]; then
    FAILED=$(echo "$CHECKS_JSON" | jq '[.[] | select(.conclusion=="FAILURE")] | length')

    if [ "$FAILED" -eq 0 ]; then
      echo "✅ CI/CD completado exitosamente en $ELAPSED segundos"
      gh pr checks
      break
    else
      echo "❌ CI/CD falló. Analizando errores..."
      gh pr checks

      # [Aquí iría lógica de análisis y corrección]
      # Por ahora, INFORMAR
      echo ""
      echo "ACCIÓN REQUERIDA:"
      echo "1. Revisar logs de workflows fallidos"
      echo "2. Corregir errores identificados"
      echo "3. Push correcciones y reiniciar monitoreo"
      exit 1
    fi
  else
    echo "⏳ CI/CD en progreso... ($ELAPSED/$MAX_WAIT segundos)"
    echo "Checks activos: $IN_PROGRESS"
  fi

  sleep $CHECK_INTERVAL
done

echo "✅ TASK 2.11 COMPLETADA"
```

---

### TASK 2.12: Atender Code Review (GitHub Copilot)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Obtener comentarios de Copilot en el PR
gh pr view --json comments | jq '[.comments[] | select(.author.login=="github-copilot[bot]")]'

# Para CADA comentario:
# 1. Leer y analizar
# 2. Si es PROCEDENTE (seguridad, bugs, performance críticos): CORREGIR
# 3. Si NO es procedente (cosmético, opinión): DOCUMENTAR en EXECUTION_REPORT.md
# 4. Si requiere decisión arquitectónica: DETENER e INFORMAR

# Si hay correcciones:
git add .
git commit -m "fix: atender comentarios de code review

- [Descripción de correcciones]

Co-Authored-By: GitHub Copilot <noreply@github.com>"
git push origin feature/sprint-00-infrastructure-v0.7.0

# Reiniciar monitoreo CI/CD
```

---

### TASK 2.13: Merge a dev

**⚠️ SOLO SI:**
- ✅ CI/CD pasando
- ✅ Comentarios Copilot atendidos o documentados
- ✅ No hay errores repetidos 3x
- ✅ Usuario aprueba (si hay decisiones pendientes)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Merge con squash
gh pr merge --squash --delete-branch

echo "✅ PR mergeado a dev"
echo "✅ Branch feature/sprint-00-infrastructure-v0.7.0 eliminado"
```

---

### TASK 2.14: Generar EXECUTION_REPORT.md

**Ubicación:** `docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/EXECUTION_REPORT.md`

**Contenido:**

```markdown
# Sprint-00 Infrastructure v0.7.0 - Execution Report

**Fecha inicio Fase 1:** [TIMESTAMP]  
**Fecha fin Fase 2:** [TIMESTAMP]  
**Duración total:** [X horas Y minutos]

---

## ✅ Objetivos Completados

- [x] Publicar infrastructure v0.7.0 con jerarquía académica
- [x] Migrar api-administracion a infrastructure v0.7.0
- [x] Validar jerarquía funciona correctamente
- [x] Validar función anti-ciclos funciona
- [x] Pasar todos los tests (unitarios + integración)

---

## 📊 Métricas

### Fase 1 (Claude Code Web)
- **Duración:** [X min]
- **Archivos modificados:** 9 (migraciones, seeds, CHANGELOG)
- **Líneas agregadas:** ~XXX
- **Líneas eliminadas:** ~XXX

### Fase 2 (Claude Code Local)
- **Duración:** [X min]
- **Archivos modificados:** 2-5 (go.mod, go.sum, scripts eliminados)
- **Tests ejecutados:** ~77
- **Tests pasando:** 77/77 ✅
- **Coverage:** ~XX%

### CI/CD
- **Workflows ejecutados:** X
- **Tiempo promedio:** X min
- **Intentos hasta éxito:** 1-3

---

## 🔧 Problemas Encontrados

### Problema 1: [Descripción]
**Solución:** [Cómo se resolvió]  
**Lección aprendida:** [Qué aprendimos]

[Repetir para cada problema]

---

## 💬 Code Review

### Comentarios de GitHub Copilot
- **Total:** X comentarios
- **Atendidos:** X
- **No procedentes:** X (documentados abajo)

### Comentarios No Atendidos

#### Comentario 1: [Resumen]
**Razón:** [Por qué no se atendió]  
**Impacto:** Bajo/Medio/Alto  
**Deuda técnica:** Sí/No

---

## 📝 Notas

- [Cualquier nota adicional]

---

**Generado por:** Claude Code Local (Fase 2)  
**Workflow:** Orchestration 2.0
```

---

### TASK 2.15: Actualizar PROGRESS.json (SI EXISTE)

**Archivo:** `docs/isolated/PROGRESS.json`

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Si existe PROGRESS.json, actualizar
if [ -f docs/isolated/PROGRESS.json ]; then
  jq '(.sprints[] | select(.id=="Sprint-00-Integrar-Infrastructure")) |= {
    id: "Sprint-00-Integrar-Infrastructure",
    status: "completed",
    phase1_completed_at: "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    phase2_completed_at: "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    merged_at: "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    pr_number: '$PR_NUMBER',
    branch: null,
    infrastructure_version: "v0.7.0"
  }' docs/isolated/PROGRESS.json > tmp.json && mv tmp.json docs/isolated/PROGRESS.json

  git checkout dev
  git pull origin dev
  git add docs/isolated/PROGRESS.json
  git commit -m "chore: actualizar estado Sprint-00 a completed (infrastructure v0.7.0)"
  git push origin dev

  echo "✅ PROGRESS.json actualizado"
else
  echo "⚠️ PROGRESS.json no existe, saltando actualización"
fi
```

---

### TASK 2.16: Informar al Usuario

**Generar resumen final:**

```
✅ SPRINT-00 COMPLETADO (2 FASES)

## FASE 1 (Infrastructure)
- Repositorio: edugo-infrastructure
- Branch: dev
- Tag: v0.7.0 ✅ PUBLICADO
- Disponible: go get ...@v0.7.0 ✅

Cambios:
- Migraciones 002, 003, 004 extendidas con jerarquía
- Seeds actualizados con ejemplos de jerarquía
- CHANGELOG.md actualizado

## FASE 2 (API-Administracion)
- Repositorio: edugo-api-administracion
- Branch: feature/sprint-00-infrastructure-v0.7.0 ✅ MERGEADO
- PR: #[PR_NUMBER] ✅ MERGED
- Duración total: [X horas Y minutos]

Cambios:
- go.mod: infrastructure v0.6.0 → v0.7.0
- Migraciones locales eliminadas
- Tests: 77/77 pasando ✅
- CI/CD: Pasando ✅

## NUEVAS CAPACIDADES DISPONIBLES

- ✅ Jerarquía académica (parent_unit_id)
- ✅ Metadata JSONB en todas las tablas
- ✅ Función anti-ciclos con trigger
- ✅ Vista recursiva (v_academic_unit_tree)
- ✅ Tipos extendidos (school, club, department)
- ✅ Roles extendidos (coordinator, admin, assistant)

## PRÓXIMOS PASOS

✅ Sprint-00 completado
✅ Infrastructure v0.7.0 disponible
⏭️ Próximo sprint: [Determinar según roadmap]

Ver reporte completo:
docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/EXECUTION_REPORT.md
```

---

## 🛑 REGLAS DE DETENCIÓN

**DEBES DETENERTE si:**
- ⏱️ Monitoreo CI/CD excede 5 minutos sin completar
- 🔁 Un error se repite 3 veces (mismo tipo)
- 🐳 Docker no está corriendo y no se puede levantar
- 🗄️ PostgreSQL no levanta tras 3 intentos
- 🔧 Comentarios de Copilot requieren decisión arquitectónica
- 📦 infrastructure v0.7.0 no está disponible tras 5 minutos

**Al detenerte:**
1. Generar informe detallado de estado actual
2. Documentar problema exacto con logs
3. Listar pasos completados vs pendientes
4. Informar al usuario con contexto completo

---

## 📁 ARCHIVOS CREADOS EN FASE 1

```
edugo-infrastructure/
├── postgres/migrations/
│   ├── .backup-YYYYMMDD-HHMMSS/
│   │   ├── 002_create_schools.up.sql           (backup)
│   │   ├── 002_create_schools.down.sql         (backup)
│   │   ├── 003_create_academic_units.up.sql    (backup)
│   │   ├── 003_create_academic_units.down.sql  (backup)
│   │   ├── 004_create_memberships.up.sql       (backup)
│   │   └── 004_create_memberships.down.sql     (backup)
│   ├── 002_create_schools.up.sql               (✅ MODIFICADO - con metadata)
│   ├── 002_create_schools.down.sql             (✅ MODIFICADO)
│   ├── 003_create_academic_units.up.sql        (✅ MODIFICADO - con jerarquía completa)
│   ├── 003_create_academic_units.down.sql      (✅ MODIFICADO)
│   ├── 004_create_memberships.up.sql           (✅ MODIFICADO - con metadata y roles)
│   └── 004_create_memberships.down.sql         (✅ MODIFICADO)
├── postgres/seeds/
│   ├── academic_units.sql                      (✅ MODIFICADO - con jerarquía)
│   └── memberships.sql                         (✅ MODIFICADO - con roles extendidos)
└── CHANGELOG.md                                (✅ MODIFICADO - v0.7.0 agregado)
```

---

## 📦 TAG PUBLICADO

- **Tag:** v0.7.0
- **Repositorio:** github.com/EduGoGroup/edugo-infrastructure
- **Branch:** dev
- **Disponible vía:** `go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.0`

---

## ✅ CHECKLIST FASE 2

- [ ] Checkout dev y pull latest
- [ ] Crear branch feature/sprint-00-infrastructure-v0.7.0
- [ ] Actualizar go.mod a infrastructure v0.7.0
- [ ] Ejecutar go mod tidy
- [ ] Eliminar migraciones locales (si aplica)
- [ ] Recrear BD local con infrastructure v0.7.0
- [ ] Validar nuevas columnas existen (metadata, parent_unit_id, description)
- [ ] Validar función prevent_academic_unit_cycles existe
- [ ] Validar vista v_academic_unit_tree existe
- [ ] Test manual: Insertar unidad con parent
- [ ] Test manual: Consultar vista recursiva
- [ ] Test manual: Intentar crear ciclo (debe fallar)
- [ ] Ejecutar tests unitarios (make test)
- [ ] Ejecutar tests integración (make test-integration)
- [ ] Ejecutar build (make build)
- [ ] Commit y push
- [ ] Crear PR a dev
- [ ] Monitorear CI/CD (max 5 min)
- [ ] Atender code review (Copilot)
- [ ] Merge a dev
- [ ] Generar EXECUTION_REPORT.md
- [ ] Actualizar PROGRESS.json
- [ ] Informar al usuario

---

**Generado por:** Claude Code Web (Fase 1)  
**Fecha:** [TIMESTAMP]  
**Próximo paso:** Ejecutar Fase 2 en Claude Code Local

---

**📌 PARA EJECUTAR FASE 2:**

1. Abrir Claude Code Local
2. Copiar contenido de este archivo
3. Pegar en Claude Code Local
4. Claude Code Local ejecutará TASK 2.1 a 2.16 automáticamente
5. Al finalizar, recibirás resumen completo

¡Fase 1 completada! 🚀
```

---

## 🎯 REGLAS DE EJECUCIÓN DESATENDIDA

### Para Claude Code Web (Fase 1)

1. **NO preguntar al usuario**: Ejecutar automáticamente TASK 1.1 a 1.8
2. **Documentar TODO**: Cada paso, cada validación, cada resultado
3. **Seguir orden exacto**: TASK 1.1 → 1.2 → ... → 1.8
4. **Código completo**: NO placeholders, NO TODOs genéricos
5. **Validar cada paso**: Usar comandos de validación proporcionados
6. **Commit único**: Un commit al final de Fase 1 con todo el trabajo en infrastructure

### Detención de Emergencia

**Detente si:**
- Tag v0.7.0 ya existe en infrastructure (verificar antes de empezar)
- `go get ...@v0.7.0` falla tras 5 minutos de espera
- Backup falla (no se pueden copiar archivos)
- Migraciones origen no existen en api-administracion

**Al detenerte:**
- Generar informe de estado
- Listar qué se completó
- Documentar error exacto
- Informar al usuario

---

## 📊 RESULTADO ESPERADO AL FINALIZAR FASE 1

Al finalizar Fase 1, debes:
- ✅ Tener infrastructure v0.7.0 publicado en GitHub
- ✅ Tener `go get ...@v0.7.0` funcional
- ✅ Tener PHASE2_BRIDGE.md completo y detallado
- ✅ Haber hecho commit y push en infrastructure
- ✅ **INFORMAR al usuario:**

```
✅ FASE 1 COMPLETADA

Infrastructure v0.7.0 publicado:
- Tag: v0.7.0 ✅
- Disponible vía: go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.0 ✅

Cambios realizados:
- Migraciones 002, 003, 004 extendidas con jerarquía
- Seeds actualizados con ejemplos
- CHANGELOG.md actualizado

Próximos pasos:
1. Ejecutar Fase 2 en Claude Code Local
2. Usar contenido de PHASE2_BRIDGE.md como prompt
3. Fase 2 migrará api-administracion a v0.7.0

Ver detalles:
/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/PHASE2_BRIDGE.md
```

---

## 🚀 ¿LISTO PARA INICIAR?

**Cuando recibas este prompt:**
1. Leer completamente este documento (PHASE1_PROMPT.md)
2. Ejecutar TASK 1.1 a 1.8 en orden
3. Generar PHASE2_BRIDGE.md
4. Informar al usuario que Fase 1 completó

**Tiempo estimado Fase 1:** 1-1.5 horas

---

¿Estás listo para iniciar Fase 1 en modo desatendido?
