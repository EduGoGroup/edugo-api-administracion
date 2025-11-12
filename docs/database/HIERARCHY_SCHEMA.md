# Schema de Jerarquía Académica

**Proyecto:** edugo-api-administracion  
**Fecha:** 12 de Noviembre, 2025  
**Versión:** 1.0  

---

## 📋 Tabla de Contenidos

- [Resumen](#resumen)
- [Modelo de Datos](#modelo-de-datos)
- [Tablas](#tablas)
- [Funciones](#funciones)
- [Vistas](#vistas)
- [Índices](#índices)
- [Ejemplos de Uso](#ejemplos-de-uso)
- [Migraciones](#migraciones)

---

## Resumen

Este schema implementa la jerarquía académica para EduGo, permitiendo modelar:

- **Escuelas** (instituciones educativas)
- **Unidades Académicas** jerárquicas (grados, secciones, clubes, departamentos)
- **Membresías** de usuarios en unidades (con roles y vigencia temporal)

### Características Principales

✅ Jerarquía multinivel con auto-referencia  
✅ Prevención de ciclos mediante trigger  
✅ Soft deletes en unidades académicas  
✅ Membresías con vigencia temporal  
✅ Vistas optimizadas con CTE recursivo  
✅ Índices de performance  
✅ Metadata extensible (JSONB)  

---

## Modelo de Datos

```
┌─────────────┐
│   school    │
└──────┬──────┘
       │ 1
       │
       │ N
┌──────▼──────────────┐
│  academic_unit      │◄──┐ parent_unit_id (auto-referencia)
│  (self-reference)   │   │
└──────┬──────────────┘   │
       │ 1               │
       │                 │
       │ N               │
┌──────▼────────────┐    │
│ unit_membership   │    │
│                   │    │
└───────────────────┘    │
         │               │
         └───────────────┘
```

---

## Tablas

### 1. `school`

Escuelas del sistema EduGo.

**Columnas:**

| Columna | Tipo | Descripción | Constraints |
|---------|------|-------------|-------------|
| `id` | UUID | Identificador único | PK, DEFAULT uuid_generate_v4() |
| `name` | VARCHAR(255) | Nombre de la escuela | NOT NULL, CHECK no vacío |
| `code` | VARCHAR(50) | Código único | NOT NULL, UNIQUE, CHECK no vacío |
| `address` | TEXT | Dirección física | |
| `contact_email` | VARCHAR(255) | Email de contacto | CHECK formato válido |
| `contact_phone` | VARCHAR(50) | Teléfono de contacto | |
| `metadata` | JSONB | Metadata adicional | DEFAULT '{}' |
| `created_at` | TIMESTAMP | Fecha de creación | NOT NULL, DEFAULT NOW() |
| `updated_at` | TIMESTAMP | Fecha de actualización | NOT NULL, DEFAULT NOW(), AUTO-UPDATE |

**Índices:**
- `idx_school_code` en `code`
- `idx_school_created_at` en `created_at DESC`

**Ejemplo:**
```sql
INSERT INTO school (name, code, address, contact_email) VALUES
  ('Colegio San José', 'ESC-001', 'Calle Principal 123', 'contacto@sanjose.edu');
```

---

### 2. `academic_unit`

Unidades académicas con estructura jerárquica.

**Columnas:**

| Columna | Tipo | Descripción | Constraints |
|---------|------|-------------|-------------|
| `id` | UUID | Identificador único | PK, DEFAULT uuid_generate_v4() |
| `parent_unit_id` | UUID | Unidad padre (jerarquía) | FK academic_unit(id), ON DELETE SET NULL |
| `school_id` | UUID | Escuela propietaria | NOT NULL, FK school(id), ON DELETE CASCADE |
| `unit_type` | VARCHAR(50) | Tipo de unidad | NOT NULL, CHECK tipo válido |
| `display_name` | VARCHAR(255) | Nombre para mostrar | NOT NULL, CHECK no vacío |
| `code` | VARCHAR(50) | Código único en escuela | UNIQUE (school_id, code) |
| `description` | TEXT | Descripción | |
| `metadata` | JSONB | Metadata adicional | DEFAULT '{}' |
| `created_at` | TIMESTAMP | Fecha de creación | NOT NULL, DEFAULT NOW() |
| `updated_at` | TIMESTAMP | Fecha de actualización | NOT NULL, DEFAULT NOW(), AUTO-UPDATE |
| `deleted_at` | TIMESTAMP | Soft delete | |

**Tipos de Unidad (`unit_type`):**
- `school` - Nivel raíz de la escuela
- `grade` - Grado académico (ej: Primer Grado, Segundo Año)
- `section` - Sección de grado (ej: Sección A, Sección B)
- `club` - Club extracurricular (ej: Club de Robótica)
- `department` - Departamento administrativo (ej: Depto. Matemáticas)

**Constraints Especiales:**
- `academic_unit_no_self_reference`: No puede ser su propio padre
- Trigger `prevent_academic_unit_cycles`: Previene ciclos en jerarquía

**Índices:**
- `idx_academic_unit_school_id` en `school_id`
- `idx_academic_unit_parent_id` en `parent_unit_id`
- `idx_academic_unit_type` en `unit_type`
- `idx_academic_unit_deleted_at` en `deleted_at`
- `idx_academic_unit_school_type` en `(school_id, unit_type)` WHERE deleted_at IS NULL

**Ejemplo Jerárquico:**
```sql
-- Nivel 1: Escuela (raíz)
INSERT INTO academic_unit (school_id, unit_type, display_name, code) VALUES
  ('...', 'school', 'Colegio San José', 'SJ-ROOT');

-- Nivel 2: Grados
INSERT INTO academic_unit (parent_unit_id, school_id, unit_type, display_name, code) VALUES
  ('id-school', '...', 'grade', 'Primer Grado', 'SJ-G1');

-- Nivel 3: Secciones
INSERT INTO academic_unit (parent_unit_id, school_id, unit_type, display_name, code) VALUES
  ('id-grado', '...', 'section', 'Primer Grado - Sección A', 'SJ-G1-A');
```

---

### 3. `unit_membership`

Relación usuarios-unidades académicas con roles y vigencia temporal.

**Columnas:**

| Columna | Tipo | Descripción | Constraints |
|---------|------|-------------|-------------|
| `id` | UUID | Identificador único | PK, DEFAULT uuid_generate_v4() |
| `unit_id` | UUID | Unidad académica | NOT NULL, FK academic_unit(id), ON DELETE CASCADE |
| `user_id` | UUID | Usuario | NOT NULL |
| `role` | VARCHAR(50) | Rol en la unidad | NOT NULL, CHECK rol válido |
| `valid_from` | TIMESTAMP | Inicio de vigencia | NOT NULL, DEFAULT NOW() |
| `valid_until` | TIMESTAMP | Fin de vigencia | CHECK valid_until > valid_from |
| `metadata` | JSONB | Metadata adicional | DEFAULT '{}' |
| `created_at` | TIMESTAMP | Fecha de creación | NOT NULL, DEFAULT NOW() |
| `updated_at` | TIMESTAMP | Fecha de actualización | NOT NULL, DEFAULT NOW(), AUTO-UPDATE |

**Roles Válidos:**
- `student` - Estudiante
- `teacher` - Profesor
- `coordinator` - Coordinador
- `admin` - Administrador
- `assistant` - Asistente

**Constraints:**
- `unit_membership_unique`: UNIQUE (unit_id, user_id, valid_from)
- `unit_membership_dates_valid`: valid_until IS NULL OR valid_until > valid_from

**Índices:**
- `idx_unit_membership_unit_id` en `unit_id`
- `idx_unit_membership_user_id` en `user_id`
- `idx_unit_membership_role` en `role`
- `idx_unit_membership_valid_dates` en `(valid_from, valid_until)`

**Ejemplo:**
```sql
-- Estudiante activo
INSERT INTO unit_membership (unit_id, user_id, role, valid_from) VALUES
  ('id-seccion', 'id-usuario', 'student', '2025-01-15');

-- Profesor con vigencia definida
INSERT INTO unit_membership (unit_id, user_id, role, valid_from, valid_until) VALUES
  ('id-seccion', 'id-profesor', 'teacher', '2025-01-10', '2025-12-20');
```

---

## Funciones

### `prevent_academic_unit_cycles()`

Función trigger que previene ciclos en la jerarquía de unidades académicas.

**Trigger:** `trigger_prevent_academic_unit_cycles`  
**Eventos:** BEFORE INSERT OR UPDATE OF parent_unit_id ON academic_unit  

**Comportamiento:**
1. Si `parent_unit_id` es NULL, permite la operación (nodo raíz)
2. Recorre hacia arriba la jerarquía siguiendo parent_unit_id
3. Si detecta que un ancestro es igual al nodo actual → **RAISE EXCEPTION**
4. Límite de profundidad: 50 niveles

**Ejemplo de Prevención:**
```sql
-- Esto fallará con: "Ciclo detectado en jerarquía..."
UPDATE academic_unit 
SET parent_unit_id = 'id-hijo'
WHERE id = 'id-padre';
```

---

## Vistas

### `v_unit_tree`

Vista con árbol jerárquico completo usando CTE recursivo.

**Columnas:**
- `id`, `parent_unit_id`, `school_id`, `unit_type`, `display_name`, `code`, `description`
- `depth` - Profundidad en el árbol (1 = raíz)
- `path` - Array de UUIDs desde raíz hasta nodo actual
- `full_path` - Path textual legible (ej: "Escuela > Grado > Sección")
- `school_name`, `school_code` - Datos de la escuela

**Uso:**
```sql
-- Ver todo el árbol
SELECT depth, full_path, unit_type FROM v_unit_tree ORDER BY path;

-- Buscar path completo de una unidad
SELECT full_path FROM v_unit_tree WHERE id = 'id-unidad';

-- Obtener todos los hijos de un nodo (por path)
SELECT * FROM v_unit_tree WHERE path @> ARRAY['id-padre'::UUID];
```

---

### `v_active_memberships`

Vista de membresías activas con información denormalizada.

**Columnas:**
- `id`, `unit_id`, `user_id`, `role`, `valid_from`, `valid_until`, `metadata`
- `unit_name`, `unit_type`, `school_id`, `school_name`

**Filtros Aplicados:**
- `valid_until IS NULL OR valid_until > NOW()` - Solo membresías vigentes
- `deleted_at IS NULL` - Solo unidades activas

**Uso:**
```sql
-- Todos los estudiantes activos de una escuela
SELECT * FROM v_active_memberships 
WHERE school_name = 'Colegio San José' AND role = 'student';

-- Membresías de un usuario
SELECT * FROM v_active_memberships WHERE user_id = 'id-usuario';
```

---

## Índices

### Índices de Performance

**school:**
- `idx_school_code` - Búsquedas por código de escuela
- `idx_school_created_at DESC` - Ordenar por fecha de creación

**academic_unit:**
- `idx_academic_unit_school_id` - Filtrar por escuela
- `idx_academic_unit_parent_id` - Recorrer jerarquía
- `idx_academic_unit_type` - Filtrar por tipo
- `idx_academic_unit_deleted_at` - Soft deletes
- `idx_academic_unit_school_type` (parcial) - Filtros combinados en unidades activas

**unit_membership:**
- `idx_unit_membership_unit_id` - Membresías por unidad
- `idx_unit_membership_user_id` - Membresías por usuario
- `idx_unit_membership_role` - Filtrar por rol
- `idx_unit_membership_valid_dates` - Rango de vigencia

---

## Ejemplos de Uso

### Crear Jerarquía Completa

```sql
-- 1. Crear escuela
INSERT INTO school (name, code) VALUES ('Mi Escuela', 'ESC-004')
RETURNING id INTO school_id;

-- 2. Crear unidad raíz
INSERT INTO academic_unit (school_id, unit_type, display_name, code) VALUES
  (school_id, 'school', 'Mi Escuela', 'ESC-004-ROOT')
RETURNING id INTO root_id;

-- 3. Crear grado
INSERT INTO academic_unit (parent_unit_id, school_id, unit_type, display_name, code) VALUES
  (root_id, school_id, 'grade', 'Quinto Grado', 'ESC-004-G5')
RETURNING id INTO grade_id;

-- 4. Crear sección
INSERT INTO academic_unit (parent_unit_id, school_id, unit_type, display_name, code) VALUES
  (grade_id, school_id, 'section', 'Quinto Grado - Sección C', 'ESC-004-G5-C')
RETURNING id INTO section_id;
```

### Consultar Árbol Jerárquico

```sql
-- Árbol completo de una escuela
SELECT 
    REPEAT('  ', depth - 1) || display_name AS jerarquia,
    unit_type,
    code
FROM v_unit_tree
WHERE school_name = 'Colegio San José'
ORDER BY path;
```

### Asignar Membresías

```sql
-- Asignar estudiantes a una sección
INSERT INTO unit_membership (unit_id, user_id, role, valid_from) VALUES
  ('section-id', 'user-1-id', 'student', '2025-01-15'),
  ('section-id', 'user-2-id', 'student', '2025-01-15');

-- Asignar profesor con vigencia de 1 año
INSERT INTO unit_membership (unit_id, user_id, role, valid_from, valid_until) VALUES
  ('section-id', 'teacher-id', 'teacher', '2025-01-10', '2026-01-10');
```

### Soft Delete de Unidad

```sql
-- Marcar unidad como eliminada
UPDATE academic_unit SET deleted_at = NOW() WHERE id = 'unit-id';

-- Las vistas automáticamente la excluyen
SELECT * FROM v_unit_tree; -- No aparece la unidad eliminada
```

---

## Migraciones

### Aplicar Schema

```bash
# Schema completo
psql -U edugo -d edugo -f scripts/postgresql/01_academic_hierarchy.sql

# Seeds de prueba
psql -U edugo -d edugo -f scripts/postgresql/02_seeds_hierarchy.sql
```

### Rollback

```sql
-- Eliminar en orden inverso (por FKs)
DROP VIEW IF EXISTS v_active_memberships CASCADE;
DROP VIEW IF EXISTS v_unit_tree CASCADE;
DROP TABLE IF EXISTS unit_membership CASCADE;
DROP TABLE IF EXISTS academic_unit CASCADE;
DROP TABLE IF EXISTS school CASCADE;
DROP FUNCTION IF EXISTS prevent_academic_unit_cycles() CASCADE;
DROP FUNCTION IF EXISTS update_school_updated_at() CASCADE;
DROP FUNCTION IF EXISTS update_academic_unit_updated_at() CASCADE;
DROP FUNCTION IF EXISTS update_unit_membership_updated_at() CASCADE;
```

---

## Notas de Implementación

### Consideraciones de Performance

1. **Índices Parciales:** `idx_academic_unit_school_type` solo indexa unidades activas (deleted_at IS NULL)
2. **CTE Recursivo:** v_unit_tree es eficiente para árboles de hasta ~1000 nodos por escuela
3. **Soft Deletes:** Permite auditoría sin pérdida de datos, filtrar con `deleted_at IS NULL`

### Limitaciones

- **Profundidad máxima:** 50 niveles en jerarquía (configurable en `prevent_academic_unit_cycles`)
- **Sin índice NOW():** Los índices con funciones no inmutables fueron removidos por restricciones de PostgreSQL

### Futuras Mejoras

- [ ] Agregar tabla `users` y FK en `unit_membership.user_id`
- [ ] Materializar vista `v_unit_tree` para escuelas grandes (>5000 unidades)
- [ ] Agregar auditoría completa (created_by, updated_by)
- [ ] Implementar RLS (Row Level Security) por escuela

---

**Documentado por:** Claude Code  
**Fecha:** 12 de Noviembre, 2025
