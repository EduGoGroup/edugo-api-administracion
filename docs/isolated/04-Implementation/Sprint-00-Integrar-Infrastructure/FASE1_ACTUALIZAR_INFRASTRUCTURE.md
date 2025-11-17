# 🔷 FASE 1: Actualizar Infrastructure (Enfoque Correcto)

**Duración:** 2-3 horas  
**Ubicación:** `edugo-infrastructure`  
**Enfoque:** Modificar migraciones existentes (002, 003, 004) en vez de crear 012  
**Versión objetivo:** v0.7.0

---

## 🎯 Filosofía: Una Sola Migración Completa

> **"Plomo al ampa"** - No parcheamos, hacemos bien desde el inicio

### Principio
- ✅ Modificamos migraciones existentes para que queden **completas**
- ✅ Cualquier proyecto nuevo obtiene **TODO de una vez**
- ✅ **UNA SOLA VERDAD** en infrastructure (no fragmentada)
- ❌ NO creamos migración 012 (sería parche)

### Justificación
- Estamos en **desarrollo** (podemos recrear BD sin problema)
- api-mobile **NO se rompe** (campos nuevos son opcionales con defaults)
- Futuro worker obtiene **schema completo** desde el inicio
- Mantenimiento **más simple** (menos migraciones)

---

## 📋 ¿Qué Vamos a Modificar?

### Migración 002: schools
**Agregar:**
- `metadata JSONB DEFAULT '{}'::jsonb`

### Migración 003: academic_units
**Agregar:**
- `parent_unit_id UUID REFERENCES academic_units(id)` (jerarquía)
- `description TEXT`
- `metadata JSONB DEFAULT '{}'::jsonb`
- Función `prevent_academic_unit_cycles()` y trigger
- Vista `v_academic_unit_tree` (CTE recursivo)

**Modificar:**
- `type` CHECK: agregar `'school', 'club', 'department'`
- `academic_year`: quitar NOT NULL, agregar DEFAULT 0

### Migración 004: memberships
**Agregar:**
- `metadata JSONB DEFAULT '{}'::jsonb`

**Modificar:**
- `role` CHECK: agregar `'coordinator', 'admin', 'assistant'`

---

## 📁 Archivos Preparados

Todos los archivos ya están listos en:
```
api-admin/docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/
├── migrations/
│   ├── 002_create_schools.up.sql          # ✅ VERSIÓN COMPLETA
│   ├── 002_create_schools.down.sql
│   ├── 003_create_academic_units.up.sql   # ✅ VERSIÓN COMPLETA CON JERARQUÍA
│   ├── 003_create_academic_units.down.sql
│   ├── 004_create_memberships.up.sql      # ✅ VERSIÓN COMPLETA
│   ├── 004_create_memberships.down.sql
│   └── 012_extend_for_admin_api.*.sql     # ❌ IGNORAR (enfoque viejo)
└── seeds/
    ├── academic_units.sql                  # ✅ SEEDS CON JERARQUÍA
    └── memberships.sql                     # ✅ SEEDS CON TODOS LOS ROLES
```

---

## 🚀 Plan de Ejecución

### TASK 1.1: Backup de Migraciones Actuales (5 min)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Crear backup
mkdir -p postgres/migrations/.backup-$(date +%Y%m%d)
cp postgres/migrations/002_create_schools.* postgres/migrations/.backup-$(date +%Y%m%d)/
cp postgres/migrations/003_create_academic_units.* postgres/migrations/.backup-$(date +%Y%m%d)/
cp postgres/migrations/004_create_memberships.* postgres/migrations/.backup-$(date +%Y%m%d)/
```

---

### TASK 1.2: Copiar Migraciones Completas (5 min)

```bash
# Origen: api-admin
SRC="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/migrations"

# Destino: infrastructure
DEST="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/postgres/migrations"

# Copiar (sobreescribir)
cp $SRC/002_create_schools.up.sql $DEST/
cp $SRC/002_create_schools.down.sql $DEST/
cp $SRC/003_create_academic_units.up.sql $DEST/
cp $SRC/003_create_academic_units.down.sql $DEST/
cp $SRC/004_create_memberships.up.sql $DEST/
cp $SRC/004_create_memberships.down.sql $DEST/

echo "✅ Migraciones copiadas"
```

---

### TASK 1.3: Copiar Seeds (5 min)

```bash
# Origen: api-admin
SRC="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/docs/isolated/04-Implementation/Sprint-00-Integrar-Infrastructure/seeds"

# Destino: infrastructure
DEST="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/postgres/seeds"

# Copiar
cp $SRC/academic_units.sql $DEST/
cp $SRC/memberships.sql $DEST/

echo "✅ Seeds copiados"
```

---

### TASK 1.4: Testing Completo (30 min)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# 1. Levantar PostgreSQL limpio
make dev-down  # Bajar si estaba corriendo
make dev-up-core  # Levantar PostgreSQL fresco

# Esperar 10 segundos
sleep 10

# 2. Ejecutar migraciones (001-011, luego las modificadas 002-004)
cd postgres
make migrate-up

# 3. Validar que las nuevas columnas existen
psql -U edugo -d edugo_dev -c "\d schools" | grep metadata
# Debe mostrar: metadata | jsonb

psql -U edugo -d edugo_dev -c "\d academic_units" | grep parent_unit_id
# Debe mostrar: parent_unit_id | uuid

psql -U edugo -d edugo_dev -c "\d academic_units" | grep metadata
# Debe mostrar: metadata | jsonb

psql -U edugo -d edugo_dev -c "\d academic_units" | grep description
# Debe mostrar: description | text

# 4. Validar CHECK constraints extendidos
psql -U edugo -d edugo_dev -c "
SELECT constraint_name, check_clause 
FROM information_schema.check_constraints 
WHERE constraint_name LIKE '%type_check%';
"
# Debe incluir: school, club, department

psql -U edugo -d edugo_dev -c "
SELECT constraint_name, check_clause 
FROM information_schema.check_constraints 
WHERE constraint_name LIKE '%role_check%';
"
# Debe incluir: coordinator, admin, assistant

# 5. Validar función y vista
psql -U edugo -d edugo_dev -c "\df prevent_academic_unit_cycles"
# Debe mostrar la función

psql -U edugo -d edugo_dev -c "\dv v_academic_unit_tree"
# Debe mostrar la vista

# 6. Testing de rollback (DOWN)
psql -U edugo -d edugo_dev -f migrations/004_create_memberships.down.sql
psql -U edugo -d edugo_dev -f migrations/003_create_academic_units.down.sql
psql -U edugo -d edugo_dev -f migrations/002_create_schools.down.sql

# Validar que se eliminaron
psql -U edugo -d edugo_dev -c "\dt schools"
# Debe decir: no se encontró

# 7. Re-aplicar (UP) para dejar listo
psql -U edugo -d edugo_dev -f migrations/002_create_schools.up.sql
psql -U edugo -d edugo_dev -f migrations/003_create_academic_units.up.sql
psql -U edugo -d edugo_dev -f migrations/004_create_memberships.up.sql

# 8. Ejecutar seeds
psql -U edugo -d edugo_dev -f seeds/users.sql
psql -U edugo -d edugo_dev -f seeds/schools.sql
psql -U edugo -d edugo_dev -f seeds/academic_units.sql
psql -U edugo -d edugo_dev -f seeds/memberships.sql

# 9. Validar jerarquía
psql -U edugo -d edugo_dev -c "SELECT * FROM v_academic_unit_tree LIMIT 5;"
# Debe mostrar árbol jerárquico

# 10. Validar que no hay ciclos (debe fallar)
psql -U edugo -d edugo_dev -c "
UPDATE academic_units 
SET parent_unit_id = id 
WHERE id = 'a1000000-0000-0000-0000-000000000001';
"
# Debe dar error: "Ciclo detectado en jerarquía"

echo "✅ Testing completo exitoso"
```

---

### TASK 1.5: Actualizar CHANGELOG.md (10 min)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure
```

Agregar al inicio de `CHANGELOG.md`:

```markdown
## [0.7.0] - 2025-11-17

### Added (postgres)

**Soporte completo de jerarquía académica:**
- Campo `parent_unit_id` en `academic_units` para estructura jerárquica
- Campos `metadata` JSONB en `schools`, `academic_units`, `memberships`
- Campo `description` TEXT en `academic_units`
- Función `prevent_academic_unit_cycles()` con trigger
- Vista `v_academic_unit_tree` (CTE recursivo) para consultar jerarquía

**Extensiones de tipos y roles:**
- Tipos de `academic_units`: `school`, `club`, `department` (además de `grade`, `class`, `section`)
- Roles de `memberships`: `coordinator`, `admin`, `assistant` (además de `teacher`, `student`, `guardian`)

### Changed (postgres)
- `academic_units.academic_year` ahora es nullable con DEFAULT 0 (0 = sin año específico)
- Constraints CHECK extendidos (backward compatible)

### Seeds
- `academic_units.sql`: Ejemplos de jerarquía (Escuela → Grado → Sección, Escuela → Departamento → Clase)
- `memberships.sql`: Ejemplos de todos los roles

### Migration Path
Para proyectos existentes con v0.6.0:
- Ejecutar `make dev-down && make dev-up-core` (recrear BD)
- Ejecutar `make migrate-up`
- Los nuevos campos tienen defaults seguros (backward compatible)

### Breaking Changes
- Si tienes datos existentes en v0.6.0, necesitas DROP DATABASE (estamos en desarrollo)
- `academic_year` ahora acepta NULL y 0 (antes era NOT NULL)

### Notes
- Modificadas migraciones 002, 003, 004 (enfoque: una sola migración completa)
- NO se creó migración 012 (enfoque de parche descartado)
- Cambios son backward compatible con api-mobile
```

---

### TASK 1.6: Commit y Push (10 min)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Verificar cambios
git status

# Agregar
git add postgres/migrations/002_create_schools.*
git add postgres/migrations/003_create_academic_units.*
git add postgres/migrations/004_create_memberships.*
git add postgres/seeds/academic_units.sql
git add postgres/seeds/memberships.sql
git add CHANGELOG.md

# Commit
git commit -m "feat(postgres): extender schema con jerarquía y metadata para soporte completo

## Cambios en Migraciones

### 002_create_schools
- Add: metadata JSONB (extensibilidad)

### 003_create_academic_units
- Add: parent_unit_id UUID (jerarquía)
- Add: description TEXT
- Add: metadata JSONB
- Add: prevent_academic_unit_cycles() función y trigger
- Add: v_academic_unit_tree vista (CTE recursivo)
- Change: type CHECK extendido (school, club, department)
- Change: academic_year nullable DEFAULT 0

### 004_create_memberships
- Add: metadata JSONB
- Change: role CHECK extendido (coordinator, admin, assistant)

## Seeds

- academic_units.sql: Ejemplos con jerarquía
- memberships.sql: Ejemplos de todos los roles

## Compatibilidad

- ✅ Backward compatible con api-mobile (campos opcionales con defaults)
- ✅ api-admin puede usar jerarquía
- ✅ Futuro worker obtiene schema completo

## Breaking Changes

- Requiere DROP DATABASE si tienes v0.6.0 (estamos en desarrollo)

BREAKING CHANGE: academic_year ahora nullable
Co-authored-by: api-admin team"

# Push
git push origin dev
```

---

### TASK 1.7: Crear Tag v0.7.0 (5 min)

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Asegurarse de estar en dev (o main si es tu flujo)
git checkout dev

# Crear tag anotado
git tag -a v0.7.0 -m "Release v0.7.0: Schema completo con jerarquía

## Highlights

- Soporte completo de jerarquía en academic_units
- Metadata JSONB en todas las tablas core
- Tipos y roles extendidos
- Vista recursiva para consultar árboles
- Función anti-ciclos

## Modified Migrations

- 002_create_schools (+ metadata)
- 003_create_academic_units (+ jerarquía completa)
- 004_create_memberships (+ roles administrativos)

## New Seeds

- academic_units.sql (jerarquía demo)
- memberships.sql (todos los roles)

## Compatible With

- api-mobile ✅ (sin cambios requeridos)
- api-admin ✅ (puede usar jerarquía)
- worker ✅ (futuro)

## Breaking

- DROP DATABASE requerido desde v0.6.0 (dev only)"

# Push tag
git push origin v0.7.0

# Verificar
git tag -l | grep v0.7.0
```

---

### TASK 1.8: Validar Disponibilidad (5 min)

```bash
# Esperar 1-2 minutos para que GitHub procese

# Intentar descargar
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.0

# Si funciona:
echo "✅ FASE 1 COMPLETADA - infrastructure v0.7.0 disponible"

# Verificar en GitHub:
# https://github.com/EduGoGroup/edugo-infrastructure/releases/tag/v0.7.0
```

---

## ✅ Checklist FASE 1

- [ ] Backup de migraciones actuales creado
- [ ] Migraciones 002, 003, 004 copiadas
- [ ] Seeds academic_units y memberships copiados
- [ ] Testing UP exitoso (campos agregados)
- [ ] Testing DOWN exitoso (rollback funciona)
- [ ] Función prevent_academic_unit_cycles existe
- [ ] Vista v_academic_unit_tree existe
- [ ] Seeds ejecutados correctamente
- [ ] Jerarquía funciona (vista muestra árbol)
- [ ] Prevención de ciclos funciona (error al crear ciclo)
- [ ] CHANGELOG.md actualizado
- [ ] Commit creado
- [ ] Push a dev exitoso
- [ ] Tag v0.7.0 creado
- [ ] Tag v0.7.0 pusheado
- [ ] `go get ...@v0.7.0` funciona

---

## 🎯 Output Esperado

Al finalizar FASE 1:
- ✅ `edugo-infrastructure@v0.7.0` disponible en GitHub
- ✅ Migraciones completas (no fragmentadas)
- ✅ Schema soporta jerarquía + metadata
- ✅ api-mobile sigue funcionando sin cambios
- ✅ api-admin listo para migrar (FASE 2)

---

## 🔄 Diferencia con Enfoque Anterior

| Aspecto | Enfoque Viejo (012) | Enfoque Nuevo (Modificar) |
|---------|---------------------|---------------------------|
| **Migraciones** | 001-011 + **012 ALTER** | 001-011 (002,003,004 completas) |
| **Resultado** | Schema parcheado | Schema completo desde inicio |
| **Futuro worker** | Ejecuta 12 migraciones | Ejecuta 11 migraciones |
| **Mantenimiento** | Más complejo | Más simple |
| **Filosofía** | Parche incremental | Una sola verdad |

---

**Próximo paso:** Después de completar FASE 1 → Ejecutar FASE 2 (migrar api-admin)

---

**Fecha:** 17 de Noviembre, 2025  
**Versión:** 2.0 (Enfoque correcto - Modificar en vez de parche)  
**Filosofía:** Plomo al ampa 🚀
