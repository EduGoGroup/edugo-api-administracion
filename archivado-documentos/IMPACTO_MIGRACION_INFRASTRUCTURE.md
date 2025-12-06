# 🔍 Análisis de Impacto: Migración a Infrastructure

**Fecha:** 17 de Noviembre, 2025  
**Proyecto:** edugo-api-administracion  
**Objetivo:** Migrar de tablas locales a `edugo-infrastructure` v0.6.0

---

## 📊 Comparativa de Tablas

### Tabla 1: Schools

| Aspecto | api-admin (LOCAL) | infrastructure (INFRA) | Tipo de Cambio |
|---------|-------------------|------------------------|----------------|
| **Nombre de tabla** | `school` (singular) | `schools` (plural) | 🔴 CRÍTICO |
| **Campo: id** | `UUID` | `UUID` | ✅ Compatible |
| **Campo: name** | `VARCHAR(255)` | `VARCHAR(255)` | ✅ Compatible |
| **Campo: code** | `VARCHAR(50) UNIQUE` | `VARCHAR(50) UNIQUE` | ✅ Compatible |
| **Campo: address** | `TEXT` | `TEXT` | ✅ Compatible |
| **Campo: contact_email** | `VARCHAR(255)` | `email VARCHAR(255)` | 🟡 Renombrar |
| **Campo: contact_phone** | `VARCHAR(50)` | `phone VARCHAR(50)` | 🟡 Renombrar |
| **Campo: metadata** | `JSONB` | ❌ NO existe | 🟡 Campo extra local |
| **Campo: city** | ❌ NO existe | `VARCHAR(100)` | 🟢 Nuevo en infra |
| **Campo: country** | ❌ NO existe | `VARCHAR(100) DEFAULT 'Chile'` | 🟢 Nuevo en infra |
| **Campo: is_active** | ❌ NO existe | `BOOLEAN DEFAULT true` | 🟢 Nuevo en infra |
| **Campo: subscription_tier** | ❌ NO existe | `VARCHAR(50) DEFAULT 'free'` | 🟢 Nuevo en infra |
| **Campo: max_teachers** | ❌ NO existe | `INTEGER DEFAULT 10` | 🟢 Nuevo en infra |
| **Campo: max_students** | ❌ NO existe | `INTEGER DEFAULT 100` | 🟢 Nuevo en infra |
| **Campo: created_at** | `TIMESTAMP` | `TIMESTAMP WITH TIME ZONE` | 🟡 Cambio tipo |
| **Campo: updated_at** | `TIMESTAMP` | `TIMESTAMP WITH TIME ZONE` | 🟡 Cambio tipo |
| **Campo: deleted_at** | ❌ NO existe | `TIMESTAMP WITH TIME ZONE` | 🟢 Nuevo en infra |

**Resumen:**
- 🔴 **Cambio crítico:** `school` → `schools` (requiere refactoring)
- 🟡 **3 campos renombrados:** `contact_email` → `email`, `contact_phone` → `phone`, timestamps con TZ
- 🟡 **1 campo local no en infra:** `metadata` (perder funcionalidad o agregar a infra)
- 🟢 **7 campos nuevos en infra:** `city`, `country`, `is_active`, `subscription_tier`, `max_teachers`, `max_students`, `deleted_at`

---

### Tabla 2: Academic Units

| Aspecto | api-admin (LOCAL) | infrastructure (INFRA) | Tipo de Cambio |
|---------|-------------------|------------------------|----------------|
| **Nombre de tabla** | `academic_unit` (singular) | `academic_units` (plural) | 🔴 CRÍTICO |
| **Campo: id** | `UUID` | `UUID` | ✅ Compatible |
| **Campo: parent_unit_id** | `UUID` (jerárquico) | ❌ NO existe | 🔴 **INCOMPATIBLE** |
| **Campo: school_id** | `UUID REFERENCES school` | `UUID REFERENCES schools` | 🔴 FK diferente |
| **Campo: unit_type** | `VARCHAR(50)` CHECK (school, grade, section, club, department) | `type VARCHAR(50)` CHECK (grade, class, section) | 🔴 **INCOMPATIBLE** |
| **Campo: display_name** | `VARCHAR(255)` | `name VARCHAR(255)` | 🟡 Renombrar |
| **Campo: code** | `VARCHAR(50)` | `VARCHAR(50)` | ✅ Compatible |
| **Campo: description** | `TEXT` | ❌ NO existe | 🟡 Campo extra local |
| **Campo: metadata** | `JSONB` | ❌ NO existe | 🟡 Campo extra local |
| **Campo: level** | ❌ NO existe | `VARCHAR(50)` | 🟢 Nuevo en infra |
| **Campo: academic_year** | ❌ NO existe | `INTEGER NOT NULL` | 🟢 Nuevo en infra |
| **Campo: is_active** | ❌ NO existe | `BOOLEAN DEFAULT true` | 🟢 Nuevo en infra |
| **Campo: created_at** | `TIMESTAMP` | `TIMESTAMP WITH TIME ZONE` | 🟡 Cambio tipo |
| **Campo: updated_at** | `TIMESTAMP` | `TIMESTAMP WITH TIME ZONE` | 🟡 Cambio tipo |
| **Campo: deleted_at** | `TIMESTAMP` (soft delete) | `TIMESTAMP WITH TIME ZONE` | 🟡 Cambio tipo |
| **Constraint UNIQUE** | `(school_id, code)` | `(school_id, code, academic_year)` | 🔴 **INCOMPATIBLE** |

**Resumen:**
- 🔴 **INCOMPATIBILIDAD GRAVE:** Infrastructure NO tiene estructura jerárquica (`parent_unit_id`)
- 🔴 **INCOMPATIBILIDAD GRAVE:** `unit_type` permite valores diferentes (infra no tiene 'school', 'club', 'department')
- 🔴 **Cambio crítico:** `academic_unit` → `academic_units`
- 🟡 **2 campos renombrados:** `display_name` → `name`, `unit_type` → `type`
- 🟡 **2 campos locales no en infra:** `description`, `metadata`, `parent_unit_id`
- 🟢 **3 campos nuevos en infra:** `level`, `academic_year`, `is_active`
- 🔴 **Constraint diferente:** infra requiere `academic_year` en UNIQUE

**⚠️ BLOQUEANTE CRÍTICO:** Infrastructure NO soporta jerarquía de academic_units

---

### Tabla 3: Memberships

| Aspecto | api-admin (LOCAL) | infrastructure (INFRA) | Tipo de Cambio |
|---------|-------------------|------------------------|----------------|
| **Nombre de tabla** | `unit_membership` | `memberships` | 🔴 CRÍTICO |
| **Campo: id** | `UUID` | `UUID` | ✅ Compatible |
| **Campo: unit_id** | `UUID REFERENCES academic_unit` | `academic_unit_id UUID REFERENCES academic_units` | 🟡 Renombrar + FK |
| **Campo: user_id** | `UUID` (sin FK) | `UUID REFERENCES users` | 🟢 Mejor en infra (con FK) |
| **Campo: role** | `VARCHAR(50)` CHECK (student, teacher, coordinator, admin, assistant) | `VARCHAR(50)` CHECK (teacher, student, guardian) | 🔴 **INCOMPATIBLE** |
| **Campo: valid_from** | `TIMESTAMP DEFAULT NOW()` | ❌ NO existe | 🟡 Campo extra local |
| **Campo: valid_until** | `TIMESTAMP` (nullable) | ❌ NO existe | 🟡 Campo extra local |
| **Campo: metadata** | `JSONB` | ❌ NO existe | 🟡 Campo extra local |
| **Campo: school_id** | ❌ NO existe | `UUID REFERENCES schools NOT NULL` | 🟢 Nuevo en infra |
| **Campo: is_active** | ❌ NO existe | `BOOLEAN DEFAULT true` | 🟢 Nuevo en infra |
| **Campo: enrolled_at** | ❌ NO existe | `TIMESTAMP WITH TIME ZONE DEFAULT NOW()` | 🟢 Nuevo en infra |
| **Campo: withdrawn_at** | ❌ NO existe | `TIMESTAMP WITH TIME ZONE` | 🟢 Nuevo en infra |
| **Campo: created_at** | `TIMESTAMP` | `TIMESTAMP WITH TIME ZONE` | 🟡 Cambio tipo |
| **Campo: updated_at** | `TIMESTAMP` | `TIMESTAMP WITH TIME ZONE` | 🟡 Cambio tipo |
| **Constraint UNIQUE** | `(unit_id, user_id, valid_from)` | `(user_id, school_id, academic_unit_id, role)` | 🔴 **INCOMPATIBLE** |

**Resumen:**
- 🔴 **INCOMPATIBILIDAD GRAVE:** `role` permite valores diferentes (infra no tiene 'coordinator', 'admin', 'assistant')
- 🔴 **Cambio crítico:** `unit_membership` → `memberships`
- 🟡 **1 campo renombrado:** `unit_id` → `academic_unit_id`
- 🟡 **3 campos locales no en infra:** `valid_from`, `valid_until`, `metadata` (sistema de vigencia temporal)
- 🟢 **4 campos nuevos en infra:** `school_id`, `is_active`, `enrolled_at`, `withdrawn_at`
- 🔴 **Constraint diferente:** infra requiere `school_id` en UNIQUE

---

## 🔴 BLOQUEANTES CRÍTICOS IDENTIFICADOS

### Bloqueante 1: Estructura Jerárquica No Soportada

**Problema:**
- `api-admin` implementa árbol jerárquico con `parent_unit_id` (Facultad → Departamento → Carrera)
- `infrastructure` NO tiene campo `parent_unit_id` (estructura plana)

**Funcionalidades afectadas:**
- ✅ Consultas recursivas (CTE) para obtener ancestros/descendientes
- ✅ Vista `v_unit_tree` con jerarquía completa
- ✅ Función `prevent_academic_unit_cycles()` (prevenir ciclos)
- ✅ Queries en repositorio que usan `parent_unit_id`

**Opciones:**

#### Opción 1.A: Agregar `parent_unit_id` a Infrastructure (RECOMENDADA)
```sql
-- Migración nueva en infrastructure: 012_add_hierarchy_to_academic_units.up.sql
ALTER TABLE academic_units 
ADD COLUMN parent_unit_id UUID REFERENCES academic_units(id) ON DELETE SET NULL;

CREATE INDEX idx_academic_units_parent ON academic_units(parent_unit_id);
```

**Pros:**
- ✅ Mantiene funcionalidad existente de api-admin
- ✅ Otros proyectos (api-mobile, worker) pueden usar jerarquía si la necesitan
- ✅ No requiere refactoring de lógica de negocio

**Contras:**
- ❌ Requiere nuevo release de infrastructure (v0.7.0)
- ❌ Es bloqueante (no se puede migrar hasta tener nueva versión)

**Duración estimada:** 2 horas (crear migración + release + testing)

---

#### Opción 1.B: Mantener Jerarquía en Tabla Separada Local
```sql
-- Nueva tabla en api-admin: academic_unit_hierarchy
CREATE TABLE academic_unit_hierarchy (
    child_unit_id UUID REFERENCES academic_units(id) ON DELETE CASCADE,
    parent_unit_id UUID REFERENCES academic_units(id) ON DELETE CASCADE,
    PRIMARY KEY (child_unit_id)
);
```

**Pros:**
- ✅ No depende de infrastructure
- ✅ Implementación rápida

**Contras:**
- ❌ Tabla adicional local (contra el principio de infrastructure como verdad)
- ❌ Queries más complejas (JOIN extra)
- ❌ Lógica de jerarquía fragmentada

**Duración estimada:** 3 horas (migración + refactor queries)

---

#### Opción 1.C: Usar `metadata` JSONB para Jerarquía
```sql
-- En academic_units, agregar campo metadata:
metadata JSONB DEFAULT '{}'::jsonb

-- Ejemplo:
{ "parent_unit_id": "uuid-del-padre", "hierarchy_path": ["uuid1", "uuid2"] }
```

**Pros:**
- ✅ No requiere cambios en infrastructure

**Contras:**
- ❌ Performance degradado (no se puede indexar JSONB eficientemente para jerarquía)
- ❌ Pierde integridad referencial (no puede ser FK)
- ❌ CTEs recursivos no funcionan
- ❌ Anti-patrón (datos relacionales en JSON)

**Duración estimada:** 4-5 horas (refactor completo de queries)

---

### Bloqueante 2: Valores de `unit_type` Incompatibles

**Problema:**
- `api-admin` usa: `school`, `grade`, `section`, `club`, `department`
- `infrastructure` solo permite: `grade`, `class`, `section`

**Impacto:**
- ❌ NO se pueden migrar unidades de tipo `school`, `club`, `department`
- ❌ Seeds actuales fallarán (tienen tipo `school`)

**Opciones:**

#### Opción 2.A: Extender CHECK constraint en Infrastructure (RECOMENDADA)
```sql
-- Migración en infrastructure: 012_extend_academic_unit_types.up.sql
ALTER TABLE academic_units 
DROP CONSTRAINT IF EXISTS academic_units_type_check;

ALTER TABLE academic_units 
ADD CONSTRAINT academic_units_type_check 
CHECK (type IN ('school', 'grade', 'class', 'section', 'club', 'department'));
```

**Pros:**
- ✅ Compatible con api-admin
- ✅ Otros proyectos pueden usar tipos adicionales
- ✅ No requiere refactoring de código

**Contras:**
- ❌ Requiere nuevo release de infrastructure
- ❌ Es bloqueante

**Duración estimada:** 30 minutos (modificar constraint + release)

---

#### Opción 2.B: Mapear Tipos de api-admin a Infrastructure
```go
// Mapeo en código Go
func mapLocalTypeToInfra(localType string) string {
    switch localType {
    case "school": return "grade"      // Raíz como "grado"
    case "club": return "section"       // Clubs como "secciones"
    case "department": return "class"   // Departamentos como "clases"
    default: return localType
    }
}
```

**Pros:**
- ✅ No requiere cambios en infrastructure

**Contras:**
- ❌ Pérdida de semántica (un "club" no es una "section")
- ❌ Confusión en reportes y UI
- ❌ Dificulta debugging

**Duración estimada:** 2 horas (implementar mapeo + testing)

---

### Bloqueante 3: Valores de `role` en Memberships Incompatibles

**Problema:**
- `api-admin` usa: `student`, `teacher`, `coordinator`, `admin`, `assistant`
- `infrastructure` solo permite: `teacher`, `student`, `guardian`

**Impacto:**
- ❌ NO se pueden migrar memberships con rol `coordinator`, `admin`, `assistant`

**Opciones:**

#### Opción 3.A: Extender CHECK constraint en Infrastructure (RECOMENDADA)
```sql
-- Migración en infrastructure: 012_extend_membership_roles.up.sql
ALTER TABLE memberships 
DROP CONSTRAINT IF EXISTS memberships_role_check;

ALTER TABLE memberships 
ADD CONSTRAINT memberships_role_check 
CHECK (role IN ('teacher', 'student', 'guardian', 'coordinator', 'admin', 'assistant'));
```

**Pros:**
- ✅ Compatible con api-admin
- ✅ Roles administrativos disponibles para todos

**Contras:**
- ❌ Requiere nuevo release de infrastructure
- ❌ Es bloqueante

**Duración estimada:** 30 minutos

---

#### Opción 3.B: Mapear Roles a Existentes
```go
func mapLocalRoleToInfra(localRole string) string {
    switch localRole {
    case "coordinator", "admin", "assistant": return "teacher"  // Admin como teacher
    default: return localRole
    }
}
```

**Pros:**
- ✅ No requiere cambios en infrastructure

**Contras:**
- ❌ Pérdida de información (no distingue entre teacher y admin)
- ❌ Lógica de permisos se complica

**Duración estimada:** 2 horas

---

### Bloqueante 4: Campo `academic_year` Requerido en Infrastructure

**Problema:**
- `infrastructure` requiere `academic_year INTEGER NOT NULL`
- `api-admin` NO tiene este campo

**Impacto:**
- ❌ Migraciones fallarán (constraint NOT NULL)
- ❌ Unique constraint diferente: `(school_id, code, academic_year)`

**Opciones:**

#### Opción 4.A: Agregar `academic_year` a Lógica de api-admin
```go
// En entity.AcademicUnit agregar:
type AcademicUnit struct {
    // ... campos existentes
    academicYear int  // Nuevo campo
}

// Valor default: año actual
academicYear := time.Now().Year()
```

**Pros:**
- ✅ Compatible con infrastructure
- ✅ Funcionalidad útil (unidades por año escolar)

**Contras:**
- ❌ Refactoring en dominio, DTOs, repositorios
- ❌ Seeds deben especificar año
- ❌ Cambio en lógica de negocio

**Duración estimada:** 3 horas

---

#### Opción 4.B: Hacer `academic_year` Nullable en Infrastructure
```sql
-- Migración en infrastructure: 012_make_academic_year_nullable.up.sql
ALTER TABLE academic_units ALTER COLUMN academic_year DROP NOT NULL;
ALTER TABLE academic_units ALTER COLUMN academic_year SET DEFAULT 0;  -- 0 = sin año
```

**Pros:**
- ✅ api-admin no requiere cambios

**Contras:**
- ❌ Requiere release de infrastructure
- ❌ Pierde semántica (año 0 no tiene sentido)

**Duración estimada:** 30 minutos (+ release)

---

## 📋 Campos Extra de api-admin No en Infrastructure

### Campos que se Perderían al Migrar

1. **`school.metadata` (JSONB):** Metadata adicional de escuelas
2. **`academic_unit.description` (TEXT):** Descripción de unidades
3. **`academic_unit.metadata` (JSONB):** Metadata adicional de unidades
4. **`unit_membership.valid_from` (TIMESTAMP):** Fecha inicio membresía
5. **`unit_membership.valid_until` (TIMESTAMP):** Fecha fin membresía
6. **`unit_membership.metadata` (JSONB):** Metadata adicional de membresías

**Opciones:**

#### Opción A: Agregar Campos a Infrastructure (RECOMENDADA para `metadata`)
```sql
-- Migración en infrastructure: 012_add_metadata_fields.up.sql
ALTER TABLE schools ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
ALTER TABLE academic_units ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
ALTER TABLE academic_units ADD COLUMN description TEXT;
ALTER TABLE memberships ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
```

**Justificación:**
- `metadata` es un patrón común en todas las tablas (extensibilidad)
- `description` es útil para documentar unidades académicas
- Otros proyectos pueden beneficiarse

**Duración estimada:** 1 hora (+ release)

---

#### Opción B: Mantener en Tablas Locales Separadas
```sql
-- api-admin local
CREATE TABLE school_metadata (school_id UUID PRIMARY KEY, metadata JSONB);
CREATE TABLE academic_unit_metadata (unit_id UUID PRIMARY KEY, metadata JSONB, description TEXT);
```

**Pros:**
- ✅ No depende de infrastructure

**Contras:**
- ❌ Queries más complejas (JOINs extra)
- ❌ Contra principio de infrastructure como verdad

---

## 🎯 Plan de Acción Recomendado

### FASE 1: Actualizar Infrastructure (BLOQUEANTE)

**Duración:** 3-4 horas  
**Responsable:** Equipo infrastructure  
**Bloqueante para:** Migración de api-admin

**Tareas:**

1. **Crear migración `012_extend_for_admin_api.up.sql`:**
   ```sql
   -- 1. Agregar jerarquía a academic_units
   ALTER TABLE academic_units 
   ADD COLUMN parent_unit_id UUID REFERENCES academic_units(id) ON DELETE SET NULL;
   
   CREATE INDEX idx_academic_units_parent ON academic_units(parent_unit_id);
   
   -- 2. Extender tipos de academic_units
   ALTER TABLE academic_units DROP CONSTRAINT IF EXISTS academic_units_type_check;
   ALTER TABLE academic_units ADD CONSTRAINT academic_units_type_check 
   CHECK (type IN ('school', 'grade', 'class', 'section', 'club', 'department'));
   
   -- 3. Extender roles de memberships
   ALTER TABLE memberships DROP CONSTRAINT IF EXISTS memberships_role_check;
   ALTER TABLE memberships ADD CONSTRAINT memberships_role_check 
   CHECK (role IN ('teacher', 'student', 'guardian', 'coordinator', 'admin', 'assistant'));
   
   -- 4. Agregar metadata y description
   ALTER TABLE schools ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
   ALTER TABLE academic_units ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
   ALTER TABLE academic_units ADD COLUMN description TEXT;
   ALTER TABLE memberships ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
   
   -- 5. Hacer academic_year nullable (opcional)
   ALTER TABLE academic_units ALTER COLUMN academic_year DROP NOT NULL;
   ALTER TABLE academic_units ALTER COLUMN academic_year SET DEFAULT 0;
   ```

2. **Crear migración down correspondiente**

3. **Testing de migraciones:**
   ```bash
   cd edugo-infrastructure
   make test-migrations
   ```

4. **Crear release v0.7.0:**
   ```bash
   git tag -a v0.7.0 -m "feat: extend schema for api-admin compatibility"
   git push origin v0.7.0
   ```

5. **Actualizar documentación:**
   - `TABLE_OWNERSHIP.md` (documentar nuevos campos)
   - `CHANGELOG.md` (v0.7.0)

**Output esperado:**
- ✅ `edugo-infrastructure@v0.7.0` publicado
- ✅ Tablas compatibles con api-admin

---

### FASE 2: Actualizar api-admin (Después de FASE 1)

**Duración:** 4-5 horas  
**Dependencia:** Requiere infrastructure v0.7.0

**Sprint-00 Actualizado:**

1. **Actualizar go.mod:**
   ```bash
   go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.0
   go get github.com/EduGoGroup/edugo-infrastructure/migrations@v0.7.0
   go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0
   go mod tidy
   ```

2. **Refactoring de código (cambios de nombres):**
   - `school` → `schools` (en queries SQL)
   - `academic_unit` → `academic_units`
   - `unit_membership` → `memberships`
   - `contact_email` → `email`
   - `contact_phone` → `phone`
   - `display_name` → `name`
   - `unit_id` → `academic_unit_id`
   - `unit_type` → `type`

3. **Agregar campo `academic_year`:**
   - Actualizar entity `AcademicUnit`
   - Actualizar DTOs
   - Actualizar repositorios
   - Valor default: `time.Now().Year()`

4. **Actualizar repositorios (archivos afectados):**
   - `school_repository_impl.go` (~20 cambios)
   - `academic_unit_repository_impl.go` (~30 cambios)
   - `unit_membership_repository_impl.go` (~25 cambios)

5. **Actualizar tests:**
   - Fixtures con nuevos nombres de tablas
   - Seeds con `academic_year`

6. **Eliminar migraciones locales:**
   ```bash
   rm -rf scripts/postgresql/
   ```

7. **Validación:**
   ```bash
   go test ./... -v
   go build ./...
   ```

**Output esperado:**
- ✅ Código usa tablas de infrastructure
- ✅ Tests pasan 100%
- ✅ No hay migraciones locales

---

## ⏱️ Estimación Total

| Fase | Duración | Bloqueante | Responsable |
|------|----------|------------|-------------|
| **FASE 1:** Actualizar infrastructure | 3-4 horas | SÍ (bloquea FASE 2) | Equipo infra |
| **FASE 2:** Migrar api-admin | 4-5 horas | NO | Equipo api-admin |
| **TOTAL** | **7-9 horas** | - | - |

---

## ✅ Checklist de Pre-Migración

Antes de comenzar FASE 1:

- [ ] Revisar y aprobar migración `012_extend_for_admin_api.up.sql`
- [ ] Validar que api-mobile NO se rompe con nuevos campos (backward compatible)
- [ ] Backup de BD de desarrollo
- [ ] Tests de infrastructure pasan 100%

Antes de comenzar FASE 2:

- [ ] infrastructure v0.7.0 publicado y disponible
- [ ] `go get` puede descargar v0.7.0
- [ ] Rama `feature/migrate-to-infrastructure` creada en api-admin
- [ ] Backup de código actual

---

## 🎯 Recomendación Final

**Ejecutar FASE 1 (actualizar infrastructure) PRIMERO** por las siguientes razones:

1. ✅ **Infrastructure es la verdad:** Debe soportar todos los casos de uso
2. ✅ **Reutilización:** Jerarquía y metadata pueden ser útiles para api-mobile en el futuro
3. ✅ **Semántica correcta:** Roles y tipos específicos tienen significado de negocio
4. ✅ **Backward compatible:** Nuevos campos son opcionales (no rompen api-mobile)

**NO recomendamos workarounds** (mapeos, tablas locales) porque:
- ❌ Crean deuda técnica
- ❌ Violan principio de infrastructure como fuente de verdad
- ❌ Dificultan mantenimiento futuro

---

**Documento creado:** 17 de Noviembre, 2025  
**Próximo paso:** Aprobar migración de infrastructure y crear release v0.7.0
