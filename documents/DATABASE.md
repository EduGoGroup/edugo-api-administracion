# 🗄️ Base de Datos

> Modelo de datos, entidades, relaciones y diagramas ER

## 📊 Visión General

EduGo API Administración utiliza dos bases de datos:

| Base de Datos | Propósito | Puerto Default |
|---------------|-----------|----------------|
| **PostgreSQL 15** | Datos principales (escuelas, usuarios, membresías) | 5432 |
| **MongoDB 7.0** | Logs, eventos, auditoría | 27017 |

---

## 🏛️ Diagrama Entidad-Relación (ERD)

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  ┌─────────────┐         ┌─────────────────────┐                        │
│  │   SCHOOL    │         │   ACADEMIC_UNIT     │                        │
│  ├─────────────┤         ├─────────────────────┤                        │
│  │ id (PK)     │◄───────┐│ id (PK)             │                        │
│  │ name        │        ││ school_id (FK)      │────┐                   │
│  │ code        │        ││ parent_unit_id (FK) │────┤ (self-ref)        │
│  │ address     │        ││ type                │    │                   │
│  │ email       │        ││ name                │◄───┘                   │
│  │ phone       │        ││ code                │                        │
│  │ metadata    │        ││ description         │                        │
│  │ created_at  │        ││ metadata            │                        │
│  │ updated_at  │        ││ created_at          │                        │
│  │ deleted_at  │        ││ updated_at          │                        │
│  └─────────────┘        ││ deleted_at          │                        │
│        │                │└─────────────────────┘                        │
│        │                │          │                                    │
│        │                │          │                                    │
│        │                │          ▼                                    │
│        │                │  ┌─────────────────────┐   ┌───────────────┐ │
│        │                │  │    MEMBERSHIP       │   │     USER      │ │
│        │                │  ├─────────────────────┤   ├───────────────┤ │
│        │                └─▶│ id (PK)             │   │ id (PK)       │ │
│        │                   │ academic_unit_id(FK)│   │ email         │ │
│        │                   │ user_id (FK)        │◀──│ first_name    │ │
│        │                   │ role                │   │ last_name     │ │
│        │                   │ enrolled_at         │   │ password_hash │ │
│        │                   │ withdrawn_at        │   │ role          │ │
│        │                   │ is_active           │   │ is_active     │ │
│        │                   │ created_at          │   │ created_at    │ │
│        │                   │ updated_at          │   │ updated_at    │ │
│        │                   └─────────────────────┘   │ deleted_at    │ │
│        │                                             └───────────────┘ │
│        │                                                     │         │
│        │                                                     │         │
│        │                   ┌─────────────────────┐           │         │
│        │                   │  STUDENT_GUARDIAN   │           │         │
│        │                   ├─────────────────────┤           │         │
│        └──────────────────▶│ id (PK)             │           │         │
│                            │ student_id (FK)     │◀──────────┤         │
│                            │ guardian_id (FK)    │◀──────────┘         │
│                            │ relationship        │                     │
│                            │ is_primary          │                     │
│                            │ created_at          │                     │
│                            │ updated_at          │                     │
│                            └─────────────────────┘                     │
│                                                                        │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## 📋 Entidades Detalladas

### 1. School (Escuela)

Representa una institución educativa.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | No | Primary Key |
| `name` | VARCHAR(255) | No | Nombre de la escuela |
| `code` | VARCHAR(50) | No | Código único (slug) |
| `address` | TEXT | Sí | Dirección física |
| `email` | VARCHAR(100) | Sí | Email de contacto |
| `phone` | VARCHAR(20) | Sí | Teléfono de contacto |
| `metadata` | JSONB | Sí | Datos adicionales flexibles |
| `created_at` | TIMESTAMP | No | Fecha de creación |
| `updated_at` | TIMESTAMP | No | Última actualización |
| `deleted_at` | TIMESTAMP | Sí | Soft delete |

**Índices:**
- `PRIMARY KEY (id)`
- `UNIQUE (code)`
- `INDEX (name)`
- `INDEX (deleted_at)` -- Para filtrar soft-deleted

---

### 2. Academic Unit (Unidad Académica)

Representa cualquier nivel jerárquico: escuela, grado, sección, departamento, club.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | No | Primary Key |
| `school_id` | UUID | No | FK → School |
| `parent_unit_id` | UUID | Sí | FK → Academic Unit (auto-referencia) |
| `type` | VARCHAR(50) | No | Tipo: school, grade, section, club, department |
| `name` | VARCHAR(255) | No | Nombre para mostrar |
| `code` | VARCHAR(50) | No | Código único dentro de la escuela |
| `description` | TEXT | Sí | Descripción |
| `metadata` | JSONB | Sí | Datos adicionales |
| `created_at` | TIMESTAMP | No | Fecha de creación |
| `updated_at` | TIMESTAMP | No | Última actualización |
| `deleted_at` | TIMESTAMP | Sí | Soft delete |

**Tipos de Unidad (`type`):**
- `school` - Nivel escuela (raíz)
- `grade` - Grado (1°, 2°, etc.)
- `section` - Sección (A, B, C)
- `club` - Club extracurricular
- `department` - Departamento académico

**Índices:**
- `PRIMARY KEY (id)`
- `FOREIGN KEY (school_id) REFERENCES school(id)`
- `FOREIGN KEY (parent_unit_id) REFERENCES academic_unit(id)`
- `UNIQUE (school_id, code)`
- `INDEX (type)`
- `INDEX (parent_unit_id)`

---

### 3. User (Usuario)

Usuarios del sistema (administradores, profesores, estudiantes, padres).

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | No | Primary Key |
| `email` | VARCHAR(100) | No | Email único |
| `password_hash` | VARCHAR(255) | No | Hash bcrypt de la contraseña |
| `first_name` | VARCHAR(50) | No | Nombre |
| `last_name` | VARCHAR(50) | No | Apellido |
| `role` | VARCHAR(20) | No | Rol del sistema |
| `is_active` | BOOLEAN | No | Estado activo/inactivo |
| `created_at` | TIMESTAMP | No | Fecha de creación |
| `updated_at` | TIMESTAMP | No | Última actualización |
| `deleted_at` | TIMESTAMP | Sí | Soft delete |

**Roles del Sistema (`role`):**
- `super_admin` - Administrador global
- `school_admin` - Administrador de escuela
- `teacher` - Profesor
- `student` - Estudiante
- `guardian` - Padre/Tutor

**Índices:**
- `PRIMARY KEY (id)`
- `UNIQUE (email)`
- `INDEX (role)`
- `INDEX (is_active)`

---

### 4. Membership (Membresía)

Relaciona usuarios con unidades académicas con un rol específico.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | No | Primary Key |
| `academic_unit_id` | UUID | No | FK → Academic Unit |
| `user_id` | UUID | No | FK → User |
| `role` | VARCHAR(50) | No | Rol dentro de la unidad |
| `enrolled_at` | TIMESTAMP | No | Fecha de inscripción |
| `withdrawn_at` | TIMESTAMP | Sí | Fecha de baja |
| `is_active` | BOOLEAN | No | Membresía activa |
| `created_at` | TIMESTAMP | No | Fecha de creación |
| `updated_at` | TIMESTAMP | No | Última actualización |

**Roles en Unidad (`role`):**
- `director` - Director
- `coordinator` - Coordinador
- `teacher` - Profesor
- `assistant` - Asistente
- `student` - Estudiante
- `observer` - Observador

**Índices:**
- `PRIMARY KEY (id)`
- `FOREIGN KEY (academic_unit_id) REFERENCES academic_unit(id)`
- `FOREIGN KEY (user_id) REFERENCES user(id)`
- `UNIQUE (academic_unit_id, user_id, role)` -- Un usuario puede tener múltiples roles
- `INDEX (user_id)`
- `INDEX (is_active)`

---

### 5. Student Guardian (Relación Estudiante-Tutor)

Relaciona estudiantes con sus tutores/padres.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | No | Primary Key |
| `student_id` | UUID | No | FK → User (rol student) |
| `guardian_id` | UUID | No | FK → User (rol guardian) |
| `relationship` | VARCHAR(50) | No | Tipo de relación |
| `is_primary` | BOOLEAN | No | Es el tutor principal |
| `created_at` | TIMESTAMP | No | Fecha de creación |
| `updated_at` | TIMESTAMP | No | Última actualización |

**Tipos de Relación (`relationship`):**
- `father` - Padre
- `mother` - Madre
- `guardian` - Tutor legal
- `other` - Otro familiar

---

## 🌳 Jerarquía de Unidades Académicas

```
School (Colegio ABC)
├── Grade (1° Primaria)
│   ├── Section (Sección A)
│   └── Section (Sección B)
├── Grade (2° Primaria)
│   ├── Section (Sección A)
│   └── Section (Sección B)
├── Department (Matemáticas)
│   └── Club (Club de Olimpiadas)
└── Club (Club de Ajedrez)
```

**Consulta de Jerarquía (CTE Recursiva):**
```sql
WITH RECURSIVE hierarchy AS (
    -- Caso base: unidad raíz
    SELECT id, parent_unit_id, name, type, 1 AS depth
    FROM academic_unit
    WHERE id = :unit_id
    
    UNION ALL
    
    -- Caso recursivo: subir hasta la raíz
    SELECT au.id, au.parent_unit_id, au.name, au.type, h.depth + 1
    FROM academic_unit au
    JOIN hierarchy h ON au.id = h.parent_unit_id
)
SELECT * FROM hierarchy ORDER BY depth DESC;
```

---

## 📈 Modelo MongoDB (Logs/Eventos)

MongoDB almacena datos no relacionales:

### audit_log Collection
```json
{
  "_id": ObjectId,
  "timestamp": ISODate,
  "action": "CREATE|UPDATE|DELETE",
  "entity_type": "school|user|membership|...",
  "entity_id": "uuid",
  "user_id": "uuid",
  "changes": {
    "before": {...},
    "after": {...}
  },
  "ip_address": "string",
  "user_agent": "string"
}
```

### event_log Collection
```json
{
  "_id": ObjectId,
  "timestamp": ISODate,
  "event_type": "LOGIN|LOGOUT|TOKEN_REFRESH|...",
  "user_id": "uuid",
  "metadata": {...},
  "success": true|false
}
```

---

## 🔧 Configuración de Conexión

### PostgreSQL
```yaml
database:
  postgres:
    host: "localhost"
    port: 5432
    database: "edugo"
    user: "edugo_user"
    password: "${POSTGRES_PASSWORD}"
    max_connections: 25
    ssl_mode: "disable"
```

### MongoDB
```yaml
database:
  mongodb:
    uri: "${MONGODB_URI}"
    database: "edugo"
    timeout: 10s
```

---

## 🔄 Migraciones

Las migraciones se manejan desde `edugo-infrastructure/postgres`:

```bash
# Las entidades están definidas en:
github.com/EduGoGroup/edugo-infrastructure/postgres/entities

# Incluye:
# - entities/school.go
# - entities/academic_unit.go
# - entities/user.go
# - entities/membership.go
# - etc.
```

---

## 📊 Índices Recomendados

```sql
-- Performance para búsquedas frecuentes
CREATE INDEX idx_academic_unit_school_type 
ON academic_unit(school_id, type) 
WHERE deleted_at IS NULL;

CREATE INDEX idx_membership_user_active 
ON membership(user_id, is_active) 
WHERE is_active = true;

CREATE INDEX idx_user_email_active 
ON "user"(email, is_active) 
WHERE is_active = true;

-- Full-text search (opcional)
CREATE INDEX idx_school_name_fulltext 
ON school USING gin(to_tsvector('spanish', name));
```

---

## 🛡️ Soft Delete

Todas las entidades principales implementan soft delete:

```go
type School struct {
    // ...
    DeletedAt *time.Time `gorm:"index"`
}

// Al eliminar:
db.Delete(&school) // Sets deleted_at = NOW()

// Al consultar (GORM excluye automáticamente):
db.Find(&schools) // WHERE deleted_at IS NULL

// Para incluir eliminados:
db.Unscoped().Find(&schools)
```
